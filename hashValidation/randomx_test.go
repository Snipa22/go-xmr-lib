package hashValidation

import (
	"os"
	"testing"
)

// TestRXVerifierHashAgainstLiveDaemon proves the Hash() fix (previously
// misused http.NewRequest's argument order and hit the wrong endpoint,
// /seed instead of /hash) actually round-trips a real hash against a real
// running randomx-service instance.
//
// Requires a real randomx-service daemon (https://github.com/tevador/randomx-service)
// reachable at RANDOMX_SERVICE_ADDR (default http://127.0.0.1:39093). Skips
// if unreachable rather than failing, since this is an external process not
// managed by `go test`.
func TestRXVerifierHashAgainstLiveDaemon(t *testing.T) {
	addr := os.Getenv("RANDOMX_SERVICE_ADDR")
	v := NewRXVerifier(addr)

	info, err := v.Info()
	if err != nil {
		t.Skipf("randomx-service not reachable, skipping: %v", err)
	}
	if info.Algorithm != "rx/0" {
		t.Fatalf("Info().Algorithm = %q, want %q", info.Algorithm, "rx/0")
	}

	seed := []byte("integration-test-seed-0001")
	if err := v.NewSeed(seed); err != nil {
		t.Fatalf("NewSeed: %v", err)
	}

	hash1, err := v.Hash([]byte("integration test input alpha"), seed)
	if err != nil {
		t.Fatalf("Hash: %v", err)
	}
	if len(hash1) != 32 {
		t.Fatalf("Hash returned %d bytes, want 32", len(hash1))
	}

	// Same input + same seed must reproduce the exact same hash --
	// RandomX is deterministic. This is the property share_processor's
	// hash validation depends on: the value returned here must be
	// byte-for-byte comparable against what a leaf node's own RandomX
	// hash produced for the same input/seed.
	hash2, err := v.Hash([]byte("integration test input alpha"), seed)
	if err != nil {
		t.Fatalf("Hash (repeat): %v", err)
	}
	if string(hash1) != string(hash2) {
		t.Fatalf("Hash is not deterministic: got %x then %x for identical input+seed", hash1, hash2)
	}

	// A different input with the same seed must produce a different hash.
	hash3, err := v.Hash([]byte("integration test input beta"), seed)
	if err != nil {
		t.Fatalf("Hash (different input): %v", err)
	}
	if string(hash1) == string(hash3) {
		t.Fatalf("Hash returned identical output for different inputs: %x", hash1)
	}
}
