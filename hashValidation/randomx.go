package hashValidation

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

// DefaultRXVerifierTimeout bounds every HTTP request an RXVerifier issues.
//
// The zero-value http.Client this type previously used had NO timeout at all,
// so a single hung request against randomx-service (GC pause, seed-cache
// thrash under memory pressure, etc.) blocked the calling goroutine forever.
// In go-crypto-pool that goroutine is a slot in a bounded share-validation
// worker pool, so every hang permanently consumed a worker until the pool
// saturated with near-idle CPU.
//
// A RandomX hash round-trip against a healthy daemon completes in
// milliseconds; nothing legitimate here takes ten seconds. If it does, the
// daemon is genuinely stuck and the caller needs the call to fail fast rather
// than lose a worker.
const DefaultRXVerifierTimeout = 10 * time.Second

// DefaultRXVerifierMaxIdleConnsPerHost is the keep-alive pool size this type's
// HTTP transport holds open per backend.
//
// Why this exists at all: an RXVerifier built on a bare &http.Client{} has a
// nil Transport, so net/http silently falls back to http.DefaultTransport,
// whose MaxIdleConnsPerHost is 2. Every RXVerifier only ever talks to a single
// host (one randomx-service daemon), and go-crypto-pool drives it from a
// bounded validation worker pool sized to runtime.NumCPU() -- so with the
// stdlib default, exactly 2 of those workers got to reuse a warm connection
// and every other worker paid a fresh connect()/accept()/teardown on every
// single Hash call, with the just-used connection thrown away instead of
// returned to the pool. That is connection churn masquerading as backend
// latency: the observed production symptom was the RandomX validation queue
// backing up while host CPU sat nearly idle, because most of the wall clock in
// the request path was socket setup rather than RandomX compute.
//
// Why 1024 specifically, rather than "a big number":
//
//   - The pool only ever needs one connection per concurrently in-flight
//     request, and concurrency here is capped by the caller's RandomX worker
//     pool. The largest real worker count deployed across the SXMR fleet today
//     is 128 (sxmr-phx-dump's runtime.NumCPU()), so 128 is the actual number
//     this has to cover.
//   - 1024 is 8x that, which leaves genuine headroom for an operator who
//     explicitly raises -randomx-workers past NumCPU (the flag allows it) and
//     for a process that runs more than one verifier against the same daemon,
//     without being effectively unbounded.
//   - The cost of the ceiling being generous is only the idle sockets actually
//     opened -- idle conns are reaped by IdleConnTimeout (90s, inherited from
//     http.DefaultTransport) and a few KB of kernel/socket state each. On a
//     loopback-only, high-frequency, short-request client that is nothing; an
//     idle-conn cap set below real concurrency, as the stdlib default was, is
//     what actually costs throughput.
//
// This is deliberately a finite, justified number and not an unbounded pool.
const DefaultRXVerifierMaxIdleConnsPerHost = 1024

// DefaultRXVerifierMaxIdleConns is the transport-wide idle connection ceiling.
//
// It is held equal to DefaultRXVerifierMaxIdleConnsPerHost on purpose: the
// global cap applies before the per-host cap, so leaving it at
// http.DefaultTransport's 100 would make 100 -- not 1024 -- the real limit for
// this transport. Since a verifier's transport only ever speaks to one
// randomx-service host, one shared value is the whole story.
const DefaultRXVerifierMaxIdleConns = DefaultRXVerifierMaxIdleConnsPerHost

// newPooledTransport returns the transport RXVerifier uses by default: a clone
// of http.DefaultTransport (so Proxy, DialContext, ForceAttemptHTTP2,
// IdleConnTimeout, TLSHandshakeTimeout and friends all keep their stdlib
// behaviour) with only the idle-connection ceilings raised.
//
// Note what is NOT set here: MaxConnsPerHost is left at 0 (unlimited).
// Bounding total connections per host would reintroduce exactly the failure
// being fixed -- workers blocking in the transport waiting for a connection
// slot. Concurrency is the caller's worker pool's job to bound; this
// transport's job is to make sure those workers get a warm socket.
func newPooledTransport() *http.Transport {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.MaxIdleConnsPerHost = DefaultRXVerifierMaxIdleConnsPerHost
	transport.MaxIdleConns = DefaultRXVerifierMaxIdleConns
	return transport
}

// Compile-time proof that the thin Hash/NewSeed/Info wrappers kept the exact
// signatures the published Validator interface requires, so this change is not
// breaking for existing consumers, and that the new *WithContext variants
// satisfy ValidatorWithContext.
var (
	_ Validator            = (*RXVerifier)(nil)
	_ ValidatorWithContext = (*RXVerifier)(nil)
)

type RXVerifier struct {
	httpSession *http.Client
	uri         string

	// seedMu guards currentSeed. RXVerifier instances are shared across
	// go-crypto-pool's multi-goroutine validation pool, so every read and
	// write of currentSeed has to be synchronized. It is an RWMutex because
	// the common case by far is Hash observing a seed that already matches
	// (a pure read); a plain Mutex would serialize otherwise fully parallel
	// Hash calls on that path.
	seedMu      sync.RWMutex
	currentSeed []byte
}

// NewRXVerifier returns a RXVerifier instance, which is used to wrap a connection to the RandomX verificaton daemon.
// It expects a server argument, which should be the HTTP interface including the port, but this is optional if
// the daemon is running on the default port/localhost bindings.
//
// The returned verifier's HTTP client carries DefaultRXVerifierTimeout and a
// connection-pooled transport (see DefaultRXVerifierMaxIdleConnsPerHost); use
// NewRXVerifierWithTimeout to pick a different bound, or
// NewRXVerifierWithTransport to supply your own transport.
func NewRXVerifier(server string) *RXVerifier {
	return NewRXVerifierWithTimeout(server, DefaultRXVerifierTimeout)
}

// NewRXVerifierWithTimeout behaves like NewRXVerifier but lets the caller set
// the per-request timeout on the verifier's own HTTP client. A timeout of zero
// or less falls back to DefaultRXVerifierTimeout -- an unbounded client is
// never a valid configuration for this type.
//
// The returned verifier's client uses the connection-pooled transport
// described on DefaultRXVerifierMaxIdleConnsPerHost, never the stdlib default
// transport.
func NewRXVerifierWithTimeout(server string, timeout time.Duration) *RXVerifier {
	return NewRXVerifierWithTransport(server, timeout, nil)
}

// NewRXVerifierWithTransport behaves like NewRXVerifierWithTimeout but lets the
// caller supply the http.RoundTripper the verifier's client will use, so a
// deployment with a concurrency profile the defaults here don't suit (or one
// that wants to wrap the transport for metrics/tracing) can tune it without
// waiting on a go-xmr-lib release.
//
// A nil transport means "use the library default", which is the pooled
// transport from newPooledTransport -- so passing nil is identical to calling
// NewRXVerifierWithTimeout. Callers building their own *http.Transport should
// clone http.DefaultTransport and raise MaxIdleConnsPerHost/MaxIdleConns
// rather than starting from a zero-value transport, and must not leave those
// at the stdlib defaults (2 and 100) for a high-concurrency validation pool.
func NewRXVerifierWithTransport(server string, timeout time.Duration, transport http.RoundTripper) *RXVerifier {
	if len(server) == 0 {
		server = "http://127.0.0.1:39093"
	}
	if timeout <= 0 {
		timeout = DefaultRXVerifierTimeout
	}
	if transport == nil {
		transport = newPooledTransport()
	}
	return &RXVerifier{
		httpSession: &http.Client{
			Timeout: timeout,
			// Explicit transport: a nil Transport here falls back to
			// http.DefaultTransport and its MaxIdleConnsPerHost of 2, which
			// throttled concurrent validation to connection-churn speed.
			Transport: transport,
		},
		uri: server,
	}
}

// Hash performs an actual round of hashing against a given input, for RandomX, this sends it off to the verifier.
//
// It is a thin wrapper around HashWithContext using context.Background(); the
// call is still bounded by the verifier's http.Client timeout. Prefer
// HashWithContext when the caller has a deadline of its own to enforce.
func (s *RXVerifier) Hash(input []byte, seed []byte) ([]byte, error) {
	return s.HashWithContext(context.Background(), input, seed)
}

// HashWithContext performs a round of hashing against the given input, sending
// it to the verification daemon, and aborts the request as soon as ctx is done.
func (s *RXVerifier) HashWithContext(ctx context.Context, input []byte, seed []byte) ([]byte, error) {
	// Note: this is deliberately a check-then-reseed, not a reseed held under
	// the write lock for the duration of the HTTP call. The daemon has a
	// single global seed, so two goroutines racing to install *different*
	// seeds can never both be satisfied no matter what this struct locks;
	// the daemon answers that case itself with a 422 (badSeed) because the
	// per-request RandomX-Seed header is validated server side. Holding the
	// write lock across the network call would only convert that into every
	// Hash call blocking behind an unrelated reseed.
	if !s.seedMatches(seed) {
		if err := s.NewSeedWithContext(ctx, seed); err != nil {
			return nil, err
		}
	}
	// Bug fix: this previously called http.NewRequest(url, method, body) --
	// backwards (NewRequest's real signature is (method, url, body)) -- and
	// pointed at /seed instead of /hash, meaning every real call here would
	// either fail to build the request or silently reseed instead of
	// hashing. Confirmed against the real randomx-service HTTP API
	// (POST /hash, see doc/API.md): this now issues the hash request
	// correctly.
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%v/hash", s.uri), bytes.NewReader(input))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x.randomx+bin")
	req.Header.Set("RandomX-Seed", hex.EncodeToString(seed))
	resp, err := s.httpSession.Do(req)
	if err != nil {
		return nil, err
	}
	defer drainAndClose(resp)
	switch resp.StatusCode {
	case 200:
		if b, err := io.ReadAll(resp.Body); err != nil {
			return nil, err
		} else {
			return hex.DecodeString(string(b))
		}
	case 400:
		return nil, invalidBody
	case 403:
		return nil, notInitialized
	case 413:
		return nil, payloadSize
	case 415:
		return nil, badHeader
	case 422:
		return nil, badSeed
	default:
		return nil, nil
	}
}

// NewSeed sets the seed value into the remote daemon and updates the local seed directly.
//
// It is a thin wrapper around NewSeedWithContext using context.Background();
// the call is still bounded by the verifier's http.Client timeout.
func (s *RXVerifier) NewSeed(input []byte) error {
	return s.NewSeedWithContext(context.Background(), input)
}

// NewSeedWithContext sets the seed value into the remote daemon and updates the
// local seed directly, aborting the request as soon as ctx is done.
func (s *RXVerifier) NewSeedWithContext(ctx context.Context, input []byte) error {
	if s.seedMatches(input) {
		return nil
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%v/seed", s.uri), bytes.NewReader(input))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x.randomx+bin")
	resp, err := s.httpSession.Do(req)
	if err != nil {
		return err
	}
	defer drainAndClose(resp)
	switch resp.StatusCode {
	case 204:
		s.SetCurrentSeed(input)
		return nil
	case 400:
		return invalidBody
	case 413:
		return invalidBody
	case 415:
		return badHeader
	default:
		return invalidResponse
	}
}

// Info returns the information available from the /info endpoint in the daemon.
//
// It is a thin wrapper around InfoWithContext using context.Background(); the
// call is still bounded by the verifier's http.Client timeout.
func (s *RXVerifier) Info() (*ValidatorInfo, error) {
	return s.InfoWithContext(context.Background())
}

// InfoWithContext returns the information available from the /info endpoint in
// the daemon, aborting the request as soon as ctx is done.
func (s *RXVerifier) InfoWithContext(ctx context.Context) (*ValidatorInfo, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%v/info", s.uri), nil)
	if err != nil {
		return nil, err
	}
	resp, err := s.httpSession.Do(req)
	if err != nil {
		return nil, err
	}
	defer drainAndClose(resp)
	if b, err := io.ReadAll(resp.Body); err != nil {
		return nil, err
	} else {
		retVal := &ValidatorInfo{}
		if err = json.Unmarshal(b, retVal); err != nil {
			return nil, err
		}
		return retVal, nil
	}
}

// SetCurrentSeed is used to overwrite the current seed in an RXVerifier setup, this is largely for multithreaded
// implementations where multiple clients need to synchronize their seed across instances together
func (s *RXVerifier) SetCurrentSeed(input []byte) {
	s.seedMu.Lock()
	defer s.seedMu.Unlock()
	s.currentSeed = input
}

// CurrentSeed returns a copy of the seed this verifier believes the daemon is
// currently loaded with. It is safe for concurrent use.
func (s *RXVerifier) CurrentSeed() []byte {
	s.seedMu.RLock()
	defer s.seedMu.RUnlock()
	if s.currentSeed == nil {
		return nil
	}
	return bytes.Clone(s.currentSeed)
}

// seedMatches reports whether seed is already the verifier's current seed,
// taking the read lock for the comparison.
func (s *RXVerifier) seedMatches(seed []byte) bool {
	s.seedMu.RLock()
	defer s.seedMu.RUnlock()
	return bytes.Equal(seed, s.currentSeed)
}

// drainAndClose closes resp.Body on every code path, first draining a bounded
// amount of any unread remainder so the underlying connection can go back to
// the transport's keep-alive pool instead of being torn down. Prior to this,
// no code path in this file closed the body at all, leaking a connection/file
// descriptor per call.
func drainAndClose(resp *http.Response) {
	_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 4096))
	_ = resp.Body.Close()
}

var badSeed = errors.New("seed in the hashing daemon does not match provided seed")
var notInitialized = errors.New("hashing daemon not initialized")
var invalidBody = errors.New("invalid body")
var badHeader = errors.New("bad content-type header")
var invalidResponse = errors.New("unhandled status response from verification daemon")
var payloadSize = errors.New("the POST body is too large")
