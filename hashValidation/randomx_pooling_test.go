package hashValidation

import (
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// countingListener wraps a net.Listener and counts every Accept that returns a
// connection, i.e. every distinct TCP connection the client actually opened.
// This is the ground truth for "did the client reuse connections or churn
// them": the HTTP layer cannot fake it, and unlike inspecting transport
// internals it measures the thing that actually costs wall clock in
// production.
type countingListener struct {
	net.Listener
	accepts int64
}

func (l *countingListener) Accept() (net.Conn, error) {
	conn, err := l.Listener.Accept()
	if err != nil {
		return nil, err
	}
	atomic.AddInt64(&l.accepts, 1)
	return conn, nil
}

func (l *countingListener) count() int64 {
	return atomic.LoadInt64(&l.accepts)
}

// newConnCountingDaemon starts a stub randomx-service whose /hash handler
// sleeps for delay before answering, standing in for real RandomX compute
// time, and reports how many TCP connections were accepted.
//
// The delay matters: it is what forces genuine request overlap, which is the
// condition that exposes an idle-connection ceiling below the caller's real
// concurrency. With an instant handler the client can serialize through a
// couple of sockets and the bug hides.
func newConnCountingDaemon(t *testing.T, delay time.Duration) (*httptest.Server, *countingListener) {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc("/hash", func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.Copy(io.Discard, r.Body)
		time.Sleep(delay)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(stubHashHex))
	})
	mux.HandleFunc("/seed", func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.Copy(io.Discard, r.Body)
		w.WriteHeader(http.StatusNoContent)
	})
	srv := httptest.NewUnstartedServer(mux)
	counter := &countingListener{Listener: srv.Listener}
	srv.Listener = counter
	srv.Start()
	t.Cleanup(srv.Close)
	return srv, counter
}

// hashStorm fires workers concurrent Hash calls per round, for rounds rounds,
// waiting for every call in a round to finish before starting the next. The
// barrier is the point: each round is a genuine burst of `workers` simultaneous
// requests, exactly how leaf-direct's bounded randomxPool dispatches, and the
// gap between rounds is where connection reuse either happens or doesn't.
func hashStorm(t *testing.T, v *RXVerifier, seed []byte, workers, rounds int) {
	t.Helper()
	var (
		errMu sync.Mutex
		errs  []error
	)
	for round := 0; round < rounds; round++ {
		var wg sync.WaitGroup
		start := make(chan struct{})
		for w := 0; w < workers; w++ {
			wg.Add(1)
			go func(round, w int) {
				defer wg.Done()
				<-start
				h, err := v.Hash([]byte(fmt.Sprintf("input-%d-%d", round, w)), seed)
				if err != nil {
					errMu.Lock()
					if len(errs) < 8 {
						errs = append(errs, fmt.Errorf("round %d worker %d: %w", round, w, err))
					}
					errMu.Unlock()
					return
				}
				if hex.EncodeToString(h) != stubHashHex {
					errMu.Lock()
					if len(errs) < 8 {
						errs = append(errs, fmt.Errorf("round %d worker %d: Hash = %x, want %v", round, w, h, stubHashHex))
					}
					errMu.Unlock()
				}
			}(round, w)
		}
		close(start)
		wg.Wait()
	}
	errMu.Lock()
	defer errMu.Unlock()
	for _, err := range errs {
		t.Errorf("hash storm error: %v", err)
	}
}

// TestConnectionReuseUnderConcurrentHashing is THE regression test for this
// incident, and it is structured as a head-to-head proof rather than a bare
// threshold check: the identical workload is run twice against identical stub
// daemons, once with the stdlib http.DefaultTransport (what RXVerifier used to
// get by leaving http.Client.Transport nil) and once with the library's own
// pooled transport. The only difference between the two runs is the transport.
//
// With the default transport, MaxIdleConnsPerHost = 2 means at most 2 sockets
// survive the end of each round, so every subsequent round has to dial roughly
// `workers - 2` fresh connections. With the pooled transport the whole burst's
// worth of sockets stays in the keep-alive pool and round 2..N reuse them, so
// the accepted-connection count stays at roughly one burst's worth no matter
// how many rounds run.
func TestConnectionReuseUnderConcurrentHashing(t *testing.T) {
	const (
		// 64 concurrent requests: well above the stdlib per-host idle cap of
		// 2, and in the same ballpark as the real NumCPU()-sized RandomX
		// worker pools this client is driven from in production.
		workers = 64
		// 4 rounds so reuse-vs-churn is visible: churn scales with rounds,
		// reuse does not. Same proof shape used for the go-crypto-pool
		// legacytransport fix.
		rounds = 4
		// Per-request server-side delay standing in for RandomX compute, long
		// enough that all `workers` requests are genuinely in flight together.
		delay = 20 * time.Millisecond
	)
	seed := []byte("pooling-test-seed")

	// Run A: the old behaviour -- nil Transport on the http.Client, i.e.
	// http.DefaultTransport.
	defaultSrv, defaultCounter := newConnCountingDaemon(t, delay)
	defaultVerifier := NewRXVerifierWithTransport(defaultSrv.URL, DefaultRXVerifierTimeout, http.DefaultTransport)
	defaultVerifier.SetCurrentSeed(seed)
	defaultStart := time.Now()
	hashStorm(t, defaultVerifier, seed, workers, rounds)
	defaultElapsed := time.Since(defaultStart)
	defaultConns := defaultCounter.count()

	// Run B: the fix -- the pooled transport the constructors now install.
	pooledSrv, pooledCounter := newConnCountingDaemon(t, delay)
	pooledVerifier := NewRXVerifier(pooledSrv.URL)
	pooledVerifier.SetCurrentSeed(seed)
	pooledStart := time.Now()
	hashStorm(t, pooledVerifier, pooledVerifier.CurrentSeed(), workers, rounds)
	pooledElapsed := time.Since(pooledStart)
	pooledConns := pooledCounter.count()

	totalRequests := int64(workers * rounds)
	t.Logf("%d concurrent Hash calls x %d rounds = %d requests (server delay %v), NumCPU=%d",
		workers, rounds, totalRequests, delay, runtime.NumCPU())
	t.Logf("http.DefaultTransport (MaxIdleConnsPerHost=2): %d TCP connections accepted in %v",
		defaultConns, defaultElapsed)
	t.Logf("pooled transport      (MaxIdleConnsPerHost=%d): %d TCP connections accepted in %v",
		DefaultRXVerifierMaxIdleConnsPerHost, pooledConns, pooledElapsed)

	// The pooled transport must not need materially more than a single
	// burst's worth of sockets for the whole run: one connection per
	// concurrent request, then reuse. A little slack absorbs a connection the
	// server or transport legitimately retires mid-run.
	if maxPooled := int64(workers + workers/8); pooledConns > maxPooled {
		t.Fatalf("pooled transport opened %d connections for %d requests at concurrency %d; want <= %d (one burst, then reuse)",
			pooledConns, totalRequests, workers, maxPooled)
	}

	// ...and it must be dramatically fewer than the default transport needed
	// for the identical workload. If this ever stops holding, either the
	// pooling was lost or the test stopped generating real concurrency; both
	// are failures worth shouting about. The default transport is expected
	// around workers + (rounds-1)*(workers-2) == 250 here; 2x workers is a
	// deliberately conservative floor for that.
	if minDefault := int64(2 * workers); defaultConns < minDefault {
		t.Fatalf("http.DefaultTransport opened only %d connections for %d requests; expected >= %d churned connections, so this test is no longer proving anything",
			defaultConns, totalRequests, minDefault)
	}
	if defaultConns <= pooledConns {
		t.Fatalf("http.DefaultTransport opened %d connections and the pooled transport %d; the pooled transport must open strictly fewer",
			defaultConns, pooledConns)
	}

	t.Logf("connection churn eliminated: %d -> %d connections for the same %d requests (%.1fx fewer)",
		defaultConns, pooledConns, totalRequests, float64(defaultConns)/float64(pooledConns))
}

// TestDefaultTransportIsPooled pins the actual configuration the constructors
// install, so a future refactor cannot silently drop back to
// http.DefaultTransport (nil Transport) or to the stdlib ceilings.
func TestDefaultTransportIsPooled(t *testing.T) {
	for _, tc := range []struct {
		name string
		v    *RXVerifier
	}{
		{"NewRXVerifier", NewRXVerifier("http://127.0.0.1:39093")},
		{"NewRXVerifier/empty-server", NewRXVerifier("")},
		{"NewRXVerifierWithTimeout", NewRXVerifierWithTimeout("http://127.0.0.1:39093", time.Second)},
		{"NewRXVerifierWithTransport/nil", NewRXVerifierWithTransport("http://127.0.0.1:39093", time.Second, nil)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if tc.v.httpSession.Transport == nil {
				t.Fatal("httpSession.Transport is nil, so net/http falls back to http.DefaultTransport and its MaxIdleConnsPerHost of 2")
			}
			transport, ok := tc.v.httpSession.Transport.(*http.Transport)
			if !ok {
				t.Fatalf("httpSession.Transport is %T, want *http.Transport", tc.v.httpSession.Transport)
			}
			if transport == http.DefaultTransport {
				t.Fatal("httpSession.Transport is http.DefaultTransport itself; it must be an owned clone with raised idle-conn ceilings")
			}
			if transport.MaxIdleConnsPerHost != DefaultRXVerifierMaxIdleConnsPerHost {
				t.Errorf("MaxIdleConnsPerHost = %d, want %d", transport.MaxIdleConnsPerHost, DefaultRXVerifierMaxIdleConnsPerHost)
			}
			if transport.MaxIdleConns != DefaultRXVerifierMaxIdleConns {
				t.Errorf("MaxIdleConns = %d, want %d", transport.MaxIdleConns, DefaultRXVerifierMaxIdleConns)
			}
			// The global cap is applied as well as the per-host cap, so it can
			// never be the smaller of the two or it becomes the real limit.
			if transport.MaxIdleConns < transport.MaxIdleConnsPerHost {
				t.Errorf("MaxIdleConns = %d < MaxIdleConnsPerHost = %d; the global cap would become the effective per-host limit",
					transport.MaxIdleConns, transport.MaxIdleConnsPerHost)
			}
			// Bounding total connections per host would put workers back to
			// blocking inside the transport, which is the failure being fixed.
			if transport.MaxConnsPerHost != 0 {
				t.Errorf("MaxConnsPerHost = %d, want 0 (unbounded); concurrency is the caller's worker pool's job to bound", transport.MaxConnsPerHost)
			}
			// Cloned from http.DefaultTransport, so the stdlib's dialer, proxy
			// handling and idle reaping all survive.
			if transport.Proxy == nil {
				t.Error("Proxy is nil; the transport was not cloned from http.DefaultTransport")
			}
			if transport.DialContext == nil {
				t.Error("DialContext is nil; the transport was not cloned from http.DefaultTransport")
			}
			if transport.IdleConnTimeout <= 0 {
				t.Errorf("IdleConnTimeout = %v, want the inherited positive default so idle sockets are still reaped", transport.IdleConnTimeout)
			}
			// The timeout work from v1.0.0 must be untouched by this change.
			if tc.v.httpSession.Timeout <= 0 {
				t.Errorf("httpSession.Timeout = %v, must stay positive", tc.v.httpSession.Timeout)
			}
		})
	}
}

// TestPooledCeilingCoversRealWorkerPools ties the chosen numbers to the reason
// they were chosen: they must cover the largest RandomX worker pool actually
// deployed (128 on sxmr-phx-dump's NumCPU()) with real headroom for an operator
// who raises -randomx-workers past NumCPU, while staying finite.
func TestPooledCeilingCoversRealWorkerPools(t *testing.T) {
	const largestDeployedWorkerPool = 128

	if DefaultRXVerifierMaxIdleConnsPerHost < largestDeployedWorkerPool {
		t.Fatalf("DefaultRXVerifierMaxIdleConnsPerHost = %d, must be at least the largest deployed worker pool (%d)",
			DefaultRXVerifierMaxIdleConnsPerHost, largestDeployedWorkerPool)
	}
	if got := DefaultRXVerifierMaxIdleConnsPerHost / largestDeployedWorkerPool; got < 4 {
		t.Fatalf("DefaultRXVerifierMaxIdleConnsPerHost = %d is only %dx the largest deployed worker pool (%d); want >= 4x headroom for explicit -randomx-workers overrides",
			DefaultRXVerifierMaxIdleConnsPerHost, got, largestDeployedWorkerPool)
	}
	// Finite, not "effectively unlimited": an unbounded idle pool is an
	// explicit non-goal.
	if DefaultRXVerifierMaxIdleConnsPerHost > 8192 {
		t.Fatalf("DefaultRXVerifierMaxIdleConnsPerHost = %d is unreasonably large; the ceiling must stay a justified finite number",
			DefaultRXVerifierMaxIdleConnsPerHost)
	}
	t.Logf("MaxIdleConnsPerHost=%d / MaxIdleConns=%d covers the largest deployed RandomX worker pool (%d) with %dx headroom; local NumCPU=%d",
		DefaultRXVerifierMaxIdleConnsPerHost, DefaultRXVerifierMaxIdleConns, largestDeployedWorkerPool,
		DefaultRXVerifierMaxIdleConnsPerHost/largestDeployedWorkerPool, runtime.NumCPU())
}

// TestNewRXVerifierWithTransportHonoursInjection proves the injection hook
// actually routes traffic through the caller's RoundTripper (so a deployment
// can tune pooling, or wrap it for metrics, without a library release) and
// that it does not disturb the timeout behaviour from v1.0.0.
func TestNewRXVerifierWithTransportHonoursInjection(t *testing.T) {
	srv := newStubDaemon(t)
	ct := &countingTransport{base: newPooledTransport()}

	v := NewRXVerifierWithTransport(srv.URL, 2*time.Second, ct)
	if v.httpSession.Transport != http.RoundTripper(ct) {
		t.Fatalf("httpSession.Transport = %T, want the injected *countingTransport", v.httpSession.Transport)
	}
	if v.httpSession.Timeout != 2*time.Second {
		t.Fatalf("httpSession.Timeout = %v, want 2s", v.httpSession.Timeout)
	}

	seed := []byte("injected-transport-seed")
	v.SetCurrentSeed(seed)
	h, err := v.Hash([]byte("input"), seed)
	if err != nil {
		t.Fatalf("Hash: %v", err)
	}
	if hex.EncodeToString(h) != stubHashHex {
		t.Fatalf("Hash = %x, want %v", h, stubHashHex)
	}
	if requests, closes := ct.counts(); requests != 1 || closes != 1 {
		t.Fatalf("injected transport saw %d request(s) and %d Close() call(s), want 1 and 1", requests, closes)
	}

	// A non-positive timeout still falls back to the default even on this
	// path -- an unbounded client remains an invalid configuration.
	if v := NewRXVerifierWithTransport(srv.URL, 0, ct); v.httpSession.Timeout != DefaultRXVerifierTimeout {
		t.Fatalf("NewRXVerifierWithTransport(..., 0, ct).httpSession.Timeout = %v, want %v", v.httpSession.Timeout, DefaultRXVerifierTimeout)
	}
	// ...and an empty server string still gets the loopback default.
	if v := NewRXVerifierWithTransport("", time.Second, ct); v.uri != "http://127.0.0.1:39093" {
		t.Fatalf("NewRXVerifierWithTransport(\"\", ...).uri = %q, want the loopback default", v.uri)
	}
}
