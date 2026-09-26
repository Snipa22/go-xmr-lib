package hashValidation

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// stubHashHex is the (arbitrary but well-formed) 32 byte hash the stub daemon
// hands back on POST /hash.
const stubHashHex = "0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f20"

// newHangingDaemon returns a server whose every handler blocks until either the
// (very long) sleep elapses or the client hangs up. This is the shape of the
// real failure that hung go-crypto-pool's validation workers forever: the
// daemon accepts the connection and then simply never answers.
//
// The handler drains the request body and then selects on r.Context().Done()
// purely so httptest.Server.Close() doesn't block for the full sleep once the
// client has given up (net/http only starts the background read that detects a
// client hangup once the request body has been consumed).
func newHangingDaemon(t *testing.T, sleep time.Duration) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.Copy(io.Discard, r.Body)
		select {
		case <-time.After(sleep):
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(stubHashHex))
		case <-r.Context().Done():
		}
	}))
	t.Cleanup(srv.Close)
	return srv
}

// newStubDaemon returns a server that answers the three real randomx-service
// endpoints the way a healthy daemon would.
func newStubDaemon(t *testing.T) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc("/hash", func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.Copy(io.Discard, r.Body)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(stubHashHex))
	})
	mux.HandleFunc("/seed", func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.Copy(io.Discard, r.Body)
		w.WriteHeader(http.StatusNoContent)
	})
	mux.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(ValidatorInfo{
			RandomxService: "stub",
			Algorithm:      "rx/0",
			Threads:        1,
		})
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return srv
}

func isTimeoutShaped(err error) bool {
	var ne net.Error
	if errors.As(err, &ne) && ne.Timeout() {
		return true
	}
	return errors.Is(err, context.DeadlineExceeded)
}

// TestHashWithContextRespectsDeadline is THE regression test for the
// production incident: a caller-supplied deadline must actually abort the
// in-flight request. Before the fix Hash took no context at all and the
// verifier's http.Client had no timeout, so this call blocked forever and
// permanently consumed a share-validation worker slot.
func TestHashWithContextRespectsDeadline(t *testing.T) {
	// Sleep far longer than both the caller's deadline and the test's own
	// tolerance, so a pass can only mean the deadline was enforced.
	srv := newHangingDaemon(t, 30*time.Second)

	// Deliberately the real default (10s) client: the 150ms context, not the
	// client timeout, has to be what ends this call.
	v := NewRXVerifier(srv.URL)
	seed := []byte("deadline-test-seed")
	v.SetCurrentSeed(seed) // skip the reseed round trip; exercise /hash itself

	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()

	start := time.Now()
	hash, err := v.HashWithContext(ctx, []byte("input"), seed)
	elapsed := time.Since(start)

	if err == nil {
		t.Fatalf("HashWithContext returned nil error (hash=%x) against a hanging daemon; the deadline was not enforced", hash)
	}
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("HashWithContext error = %v, want it to wrap context.DeadlineExceeded", err)
	}
	if elapsed > 5*time.Second {
		t.Fatalf("HashWithContext took %v to honour a 150ms deadline", elapsed)
	}
	t.Logf("HashWithContext aborted after %v with %v", elapsed, err)
}

// TestWithContextVariantsRespectDeadline covers the other two endpoints, plus
// the already-expired-context case (nothing should even hit the wire).
func TestWithContextVariantsRespectDeadline(t *testing.T) {
	srv := newHangingDaemon(t, 30*time.Second)

	t.Run("NewSeedWithContext", func(t *testing.T) {
		v := NewRXVerifier(srv.URL)
		ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
		defer cancel()
		start := time.Now()
		err := v.NewSeedWithContext(ctx, []byte("some-new-seed"))
		if err == nil || !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("NewSeedWithContext error = %v, want context.DeadlineExceeded", err)
		}
		if elapsed := time.Since(start); elapsed > 5*time.Second {
			t.Fatalf("NewSeedWithContext took %v to honour a 150ms deadline", elapsed)
		}
	})

	t.Run("InfoWithContext", func(t *testing.T) {
		v := NewRXVerifier(srv.URL)
		ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
		defer cancel()
		start := time.Now()
		_, err := v.InfoWithContext(ctx)
		if err == nil || !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("InfoWithContext error = %v, want context.DeadlineExceeded", err)
		}
		if elapsed := time.Since(start); elapsed > 5*time.Second {
			t.Fatalf("InfoWithContext took %v to honour a 150ms deadline", elapsed)
		}
	})

	t.Run("HashWithContext/already-expired", func(t *testing.T) {
		v := NewRXVerifier(srv.URL)
		seed := []byte("expired-ctx-seed")
		v.SetCurrentSeed(seed)
		ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
		defer cancel()
		<-ctx.Done()
		start := time.Now()
		if _, err := v.HashWithContext(ctx, []byte("input"), seed); !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("HashWithContext with an expired context: error = %v, want context.DeadlineExceeded", err)
		}
		if elapsed := time.Since(start); elapsed > time.Second {
			t.Fatalf("HashWithContext with an already-expired context took %v", elapsed)
		}
	})

	// A reseed triggered from inside Hash must inherit the caller's deadline
	// too -- that request is on the same worker's critical path.
	t.Run("HashWithContext/reseed-path", func(t *testing.T) {
		v := NewRXVerifier(srv.URL)
		ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
		defer cancel()
		start := time.Now()
		if _, err := v.HashWithContext(ctx, []byte("input"), []byte("unseen-seed")); !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("HashWithContext (reseed path) error = %v, want context.DeadlineExceeded", err)
		}
		if elapsed := time.Since(start); elapsed > 5*time.Second {
			t.Fatalf("HashWithContext (reseed path) took %v to honour a 150ms deadline", elapsed)
		}
	})
}

// TestDefaultClientTimeoutIsBounded pins the default: an RXVerifier must never
// be built on top of a timeout-less http.Client again.
func TestDefaultClientTimeoutIsBounded(t *testing.T) {
	if DefaultRXVerifierTimeout <= 0 {
		t.Fatalf("DefaultRXVerifierTimeout = %v, must be positive", DefaultRXVerifierTimeout)
	}
	for _, server := range []string{"", "http://127.0.0.1:39093"} {
		v := NewRXVerifier(server)
		if v.httpSession.Timeout != DefaultRXVerifierTimeout {
			t.Fatalf("NewRXVerifier(%q).httpSession.Timeout = %v, want %v", server, v.httpSession.Timeout, DefaultRXVerifierTimeout)
		}
	}
	// A non-positive explicit timeout must not be taken literally as "no
	// timeout" -- that is the exact bug being fixed.
	if v := NewRXVerifierWithTimeout("http://127.0.0.1:39093", 0); v.httpSession.Timeout != DefaultRXVerifierTimeout {
		t.Fatalf("NewRXVerifierWithTimeout(..., 0).httpSession.Timeout = %v, want %v", v.httpSession.Timeout, DefaultRXVerifierTimeout)
	}
	if v := NewRXVerifierWithTimeout("http://127.0.0.1:39093", -time.Second); v.httpSession.Timeout != DefaultRXVerifierTimeout {
		t.Fatalf("NewRXVerifierWithTimeout(..., -1s).httpSession.Timeout = %v, want %v", v.httpSession.Timeout, DefaultRXVerifierTimeout)
	}
}

// TestClientTimeoutBoundsHangWithoutContext proves the no-context path -- the
// plain Hash/NewSeed/Info wrappers every existing caller uses -- is bounded
// too, by the http.Client timeout itself. The verifier is built with a short
// test-only timeout via NewRXVerifierWithTimeout so the test doesn't have to
// sit through the real 10s default.
func TestClientTimeoutBoundsHangWithoutContext(t *testing.T) {
	srv := newHangingDaemon(t, 30*time.Second)
	const testTimeout = 250 * time.Millisecond

	t.Run("Hash", func(t *testing.T) {
		v := NewRXVerifierWithTimeout(srv.URL, testTimeout)
		seed := []byte("client-timeout-seed")
		v.SetCurrentSeed(seed)
		start := time.Now()
		hash, err := v.Hash([]byte("input"), seed)
		elapsed := time.Since(start)
		if err == nil {
			t.Fatalf("Hash returned nil error (hash=%x) against a hanging daemon; the client timeout was not enforced", hash)
		}
		if !isTimeoutShaped(err) {
			t.Fatalf("Hash error = %v, want a timeout-shaped error", err)
		}
		if elapsed > 5*time.Second {
			t.Fatalf("Hash took %v with a %v client timeout", elapsed, testTimeout)
		}
		t.Logf("Hash aborted after %v with %v", elapsed, err)
	})

	t.Run("NewSeed", func(t *testing.T) {
		v := NewRXVerifierWithTimeout(srv.URL, testTimeout)
		start := time.Now()
		err := v.NewSeed([]byte("client-timeout-new-seed"))
		elapsed := time.Since(start)
		if err == nil || !isTimeoutShaped(err) {
			t.Fatalf("NewSeed error = %v, want a timeout-shaped error", err)
		}
		if elapsed > 5*time.Second {
			t.Fatalf("NewSeed took %v with a %v client timeout", elapsed, testTimeout)
		}
		t.Logf("NewSeed aborted after %v with %v", elapsed, err)
	})

	t.Run("Info", func(t *testing.T) {
		v := NewRXVerifierWithTimeout(srv.URL, testTimeout)
		start := time.Now()
		_, err := v.Info()
		elapsed := time.Since(start)
		if err == nil || !isTimeoutShaped(err) {
			t.Fatalf("Info error = %v, want a timeout-shaped error", err)
		}
		if elapsed > 5*time.Second {
			t.Fatalf("Info took %v with a %v client timeout", elapsed, testTimeout)
		}
		t.Logf("Info aborted after %v with %v", elapsed, err)
	})
}

// countingBody counts Close calls on a response body.
type countingBody struct {
	io.ReadCloser
	closes *int64
}

func (c *countingBody) Close() error {
	atomic.AddInt64(c.closes, 1)
	return c.ReadCloser.Close()
}

// countingTransport wraps every response body it sees so the test can assert
// Close was called exactly once per request that made it onto the wire.
type countingTransport struct {
	base     http.RoundTripper
	requests int64
	closes   int64
}

func (t *countingTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	atomic.AddInt64(&t.requests, 1)
	resp, err := t.base.RoundTrip(r)
	if err != nil {
		return nil, err
	}
	resp.Body = &countingBody{ReadCloser: resp.Body, closes: &t.closes}
	return resp, nil
}

func (t *countingTransport) counts() (requests, closes int64) {
	return atomic.LoadInt64(&t.requests), atomic.LoadInt64(&t.closes)
}

// statusDaemon answers with a caller-chosen status code per path, so the
// non-200/204 branches (which previously returned without reading OR closing
// the body) can be exercised.
func statusDaemon(t *testing.T, status int, body string) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.Copy(io.Discard, r.Body)
		w.WriteHeader(status)
		if body != "" {
			_, _ = w.Write([]byte(body))
		}
	}))
	t.Cleanup(srv.Close)
	return srv
}

// TestResponseBodyClosedOnEveryPath proves Bug 2 is fixed: every request this
// file issues gets its body closed exactly once, on the success path and on
// every non-200/204 status branch.
func TestResponseBodyClosedOnEveryPath(t *testing.T) {
	cases := []struct {
		name         string
		status       int
		body         string
		wantRequests int64
		call         func(v *RXVerifier) error
	}{
		{
			name:         "Hash/200",
			status:       http.StatusOK,
			body:         stubHashHex,
			wantRequests: 1,
			call: func(v *RXVerifier) error {
				seed := []byte("close-seed")
				v.SetCurrentSeed(seed)
				h, err := v.Hash([]byte("input"), seed)
				if err != nil {
					return err
				}
				if hex.EncodeToString(h) != stubHashHex {
					return fmt.Errorf("Hash = %x, want %v", h, stubHashHex)
				}
				return nil
			},
		},
		{
			name:         "Hash/400",
			status:       http.StatusBadRequest,
			body:         "bad request",
			wantRequests: 1,
			call: func(v *RXVerifier) error {
				seed := []byte("close-seed")
				v.SetCurrentSeed(seed)
				if _, err := v.Hash([]byte("input"), seed); !errors.Is(err, invalidBody) {
					return fmt.Errorf("Hash error = %v, want invalidBody", err)
				}
				return nil
			},
		},
		{
			name:         "Hash/422",
			status:       http.StatusUnprocessableEntity,
			body:         "seed mismatch",
			wantRequests: 1,
			call: func(v *RXVerifier) error {
				seed := []byte("close-seed")
				v.SetCurrentSeed(seed)
				if _, err := v.Hash([]byte("input"), seed); !errors.Is(err, badSeed) {
					return fmt.Errorf("Hash error = %v, want badSeed", err)
				}
				return nil
			},
		},
		{
			name:         "Hash/500-unhandled-status",
			status:       http.StatusInternalServerError,
			body:         "boom",
			wantRequests: 1,
			call: func(v *RXVerifier) error {
				seed := []byte("close-seed")
				v.SetCurrentSeed(seed)
				// The default branch's return values are pre-existing
				// behaviour and deliberately untouched here; what matters for
				// this test is that the body still gets closed.
				_, _ = v.Hash([]byte("input"), seed)
				return nil
			},
		},
		{
			name:         "Hash/reseed-then-hash",
			status:       http.StatusNoContent, // /seed answers 204, /hash then sees 204 -> default branch
			wantRequests: 2,
			call: func(v *RXVerifier) error {
				// Fresh verifier: the seed doesn't match, so this is two
				// requests (POST /seed then POST /hash) and both bodies must
				// be closed.
				_, _ = v.Hash([]byte("input"), []byte("a-brand-new-seed"))
				return nil
			},
		},
		{
			name:         "NewSeed/204",
			status:       http.StatusNoContent,
			wantRequests: 1,
			call: func(v *RXVerifier) error {
				return v.NewSeed([]byte("fresh-seed"))
			},
		},
		{
			name:         "NewSeed/400",
			status:       http.StatusBadRequest,
			body:         "bad request",
			wantRequests: 1,
			call: func(v *RXVerifier) error {
				if err := v.NewSeed([]byte("fresh-seed")); !errors.Is(err, invalidBody) {
					return fmt.Errorf("NewSeed error = %v, want invalidBody", err)
				}
				return nil
			},
		},
		{
			name:         "NewSeed/500-unhandled-status",
			status:       http.StatusInternalServerError,
			body:         "boom",
			wantRequests: 1,
			call: func(v *RXVerifier) error {
				if err := v.NewSeed([]byte("fresh-seed")); !errors.Is(err, invalidResponse) {
					return fmt.Errorf("NewSeed error = %v, want invalidResponse", err)
				}
				return nil
			},
		},
		{
			name:         "Info/200",
			status:       http.StatusOK,
			body:         `{"algorithm":"rx/0","threads":4}`,
			wantRequests: 1,
			call: func(v *RXVerifier) error {
				info, err := v.Info()
				if err != nil {
					return err
				}
				if info.Algorithm != "rx/0" || info.Threads != 4 {
					return fmt.Errorf("Info = %+v, want rx/0 with 4 threads", info)
				}
				return nil
			},
		},
		{
			name:         "Info/500",
			status:       http.StatusInternalServerError,
			body:         "boom",
			wantRequests: 1,
			call: func(v *RXVerifier) error {
				// Non-JSON error body: Info surfaces a decode error, but the
				// body still has to be closed.
				if _, err := v.Info(); err == nil {
					return errors.New("Info returned nil error for a 500/non-JSON body")
				}
				return nil
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			srv := statusDaemon(t, tc.status, tc.body)
			v := NewRXVerifier(srv.URL)
			ct := &countingTransport{base: http.DefaultTransport}
			v.httpSession.Transport = ct

			if err := tc.call(v); err != nil {
				t.Fatalf("call: %v", err)
			}

			requests, closes := ct.counts()
			if requests != tc.wantRequests {
				t.Fatalf("issued %d requests, want %d", requests, tc.wantRequests)
			}
			if closes != requests {
				t.Fatalf("closed %d response bodies for %d requests; every body must be closed exactly once", closes, requests)
			}
			t.Logf("%d request(s), %d Close() call(s)", requests, closes)
		})
	}
}

// TestConcurrentHashSeedAccessIsRaceFree is the Bug 3 regression test: many
// goroutines hammering one shared RXVerifier with a MIX of the already-current
// seed and brand new seeds. Run under `go test -race`, this fails (race
// detected) if the currentSeed mutex is removed, and passes with it.
func TestConcurrentHashSeedAccessIsRaceFree(t *testing.T) {
	srv := newStubDaemon(t)
	v := NewRXVerifier(srv.URL)

	seeds := make([][]byte, 8)
	for i := range seeds {
		seeds[i] = []byte(fmt.Sprintf("concurrent-seed-%02d", i))
	}
	// Pre-install one of them so a good share of the Hash calls take the
	// "seed already matches" read-only fast path.
	v.SetCurrentSeed(seeds[0])

	const (
		workers    = 64
		iterations = 40
	)

	var (
		errMu  sync.Mutex
		errs   []error
		hashes int64
	)
	recordErr := func(err error) {
		errMu.Lock()
		defer errMu.Unlock()
		if len(errs) < 16 {
			errs = append(errs, err)
		}
	}

	start := make(chan struct{})
	var wg sync.WaitGroup
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func(w int) {
			defer wg.Done()
			<-start
			for i := 0; i < iterations; i++ {
				switch (w + i) % 8 {
				case 0, 1, 2, 3:
					// Dominant case: hash against the seed most workers share.
					h, err := v.Hash([]byte(fmt.Sprintf("input-%d-%d", w, i)), seeds[0])
					if err != nil {
						recordErr(fmt.Errorf("worker %d iter %d Hash(shared seed): %w", w, i, err))
						continue
					}
					if hex.EncodeToString(h) != stubHashHex {
						recordErr(fmt.Errorf("worker %d iter %d Hash = %x, want %v", w, i, h, stubHashHex))
						continue
					}
					atomic.AddInt64(&hashes, 1)
				case 4, 5:
					// Concurrent reseeds through Hash's own reseed path.
					if _, err := v.Hash([]byte("input"), seeds[(w+i)%len(seeds)]); err != nil {
						recordErr(fmt.Errorf("worker %d iter %d Hash(rotating seed): %w", w, i, err))
						continue
					}
					atomic.AddInt64(&hashes, 1)
				case 6:
					// Direct writers.
					if err := v.NewSeed(seeds[(w*3+i)%len(seeds)]); err != nil {
						recordErr(fmt.Errorf("worker %d iter %d NewSeed: %w", w, i, err))
					}
					v.SetCurrentSeed(seeds[(w+i)%len(seeds)])
				case 7:
					// Direct readers plus an unrelated endpoint on the same
					// shared http.Client.
					_ = v.CurrentSeed()
					if _, err := v.InfoWithContext(context.Background()); err != nil {
						recordErr(fmt.Errorf("worker %d iter %d Info: %w", w, i, err))
					}
				}
			}
		}(w)
	}
	close(start)
	wg.Wait()

	errMu.Lock()
	defer errMu.Unlock()
	if len(errs) > 0 {
		msgs := make([]string, 0, len(errs))
		for _, err := range errs {
			msgs = append(msgs, err.Error())
		}
		t.Fatalf("%d concurrent operation error(s):\n%v", len(errs), strings.Join(msgs, "\n"))
	}
	t.Logf("%d workers x %d iterations completed, %d successful hashes, seed=%q",
		workers, iterations, atomic.LoadInt64(&hashes), v.CurrentSeed())
}

// TestSeedStateUnderConcurrentWriters keeps the seed accessors themselves hot
// (no HTTP involved) so the race detector has the tightest possible window on
// the currentSeed field.
func TestSeedStateUnderConcurrentWriters(t *testing.T) {
	v := NewRXVerifier("http://127.0.0.1:1")
	var wg sync.WaitGroup
	for w := 0; w < 32; w++ {
		wg.Add(1)
		go func(w int) {
			defer wg.Done()
			seed := []byte(fmt.Sprintf("seed-%02d", w))
			for i := 0; i < 500; i++ {
				v.SetCurrentSeed(seed)
				_ = v.seedMatches(seed)
				_ = v.CurrentSeed()
			}
		}(w)
	}
	wg.Wait()
}
