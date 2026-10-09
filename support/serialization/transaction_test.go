package serialization

import (
	"bytes"
	"testing"
)

// TestConstructTXExtraUnrecognizedTagReturnsError confirms that an
// unrecognized tag byte (anything outside 0x00-0x03) returns a clean error
// immediately instead of hanging forever. This test MUST run with a short
// timeout (`go test -timeout 5s`) so that a regression that reintroduces the
// infinite loop fails the suite loudly instead of silently hanging CI.
func TestConstructTXExtraUnrecognizedTagReturnsError(t *testing.T) {
	for _, tag := range []byte{0x7F, 0xFF} {
		extra, err := ConstructTXExtra([]byte{tag})
		if err == nil {
			t.Fatalf("tag 0x%02x: expected an error, got nil (extra=%+v)", tag, extra)
		}
		if err != ErrUnhandledTxExtraTag {
			t.Fatalf("tag 0x%02x: expected ErrUnhandledTxExtraTag, got %v", tag, err)
		}
	}
}

// TestConstructTXExtraTagSweepNeverHangs is the sweep test required by the
// dispatch brief: for every byte value 0x04-0xFF (252 values), confirm
// ConstructTXExtra returns promptly with an error and never hangs. The whole
// sweep must complete well under the 5s -timeout used for this package.
func TestConstructTXExtraTagSweepNeverHangs(t *testing.T) {
	for tag := 0x04; tag <= 0xFF; tag++ {
		b := byte(tag)
		extra, err := ConstructTXExtra([]byte{b})
		if err == nil {
			t.Fatalf("tag 0x%02x: expected an error, got nil (extra=%+v)", b, extra)
		}
		if err != ErrUnhandledTxExtraTag {
			t.Fatalf("tag 0x%02x: expected ErrUnhandledTxExtraTag, got %v", b, err)
		}
	}
}

// TestConstructTXExtraTruncatedPubKeyReturnsError covers tag 0x01 with fewer
// than the required 32 trailing bytes.
func TestConstructTXExtraTruncatedPubKeyReturnsError(t *testing.T) {
	_, err := ConstructTXExtra([]byte{0x01, 0x02, 0x03})
	if err != ErrTxExtraTruncated {
		t.Fatalf("expected ErrTxExtraTruncated, got %v", err)
	}
}

// TestConstructTXExtraTruncatedNonceReturnsError covers tag 0x02 where the
// nonce-length byte claims more bytes than actually remain in the buffer.
func TestConstructTXExtraTruncatedNonceReturnsError(t *testing.T) {
	_, err := ConstructTXExtra([]byte{0x02, 0xFF, 0x01})
	if err != ErrTxExtraTruncated {
		t.Fatalf("expected ErrTxExtraTruncated, got %v", err)
	}

	// Also cover the case where even the length byte itself is missing.
	_, err = ConstructTXExtra([]byte{0x02})
	if err != ErrTxExtraTruncated {
		t.Fatalf("expected ErrTxExtraTruncated for missing length byte, got %v", err)
	}
}

// TestConstructTXExtraTruncatedMergeMiningReturnsError covers tag 0x03 where
// the merge-mining-length byte claims more bytes than actually remain in the
// buffer.
func TestConstructTXExtraTruncatedMergeMiningReturnsError(t *testing.T) {
	_, err := ConstructTXExtra([]byte{0x03, 0xFF, 0x01})
	if err != ErrTxExtraTruncated {
		t.Fatalf("expected ErrTxExtraTruncated, got %v", err)
	}

	// Also cover the case where even the length byte itself is missing.
	_, err = ConstructTXExtra([]byte{0x03})
	if err != ErrTxExtraTruncated {
		t.Fatalf("expected ErrTxExtraTruncated for missing length byte, got %v", err)
	}
}

// TestConstructTXExtraValidAllFourTags is the regression test required by
// the dispatch brief: well-formed tx_extra data covering all four existing
// valid tag bytes (0x00/0x01/0x02/0x03) must still parse identically to
// before the fix (no behavior change on the happy path).
func TestConstructTXExtraValidAllFourTags(t *testing.T) {
	pubKey := make([]byte, 32)
	for i := range pubKey {
		pubKey[i] = byte(i + 1)
	}
	nonce := []byte{0xAA, 0xBB, 0xCC}
	mergeMining := []byte{0x11, 0x22}

	var in []byte
	in = append(in, 0x00) // padding
	in = append(in, 0x01) // pubkey tag
	in = append(in, pubKey...)
	in = append(in, 0x02, byte(len(nonce))) // nonce tag
	in = append(in, nonce...)
	in = append(in, 0x03, byte(len(mergeMining))) // merge mining tag
	in = append(in, mergeMining...)

	extra, err := ConstructTXExtra(in)
	if err != nil {
		t.Fatalf("expected no error for well-formed input, got %v", err)
	}
	if extra.Padding == nil || len(*extra.Padding) != 1 || (*extra.Padding)[0] != 0 {
		t.Fatalf("unexpected Padding: %+v", extra.Padding)
	}
	if !bytes.Equal(extra.PubKey, pubKey) {
		t.Fatalf("unexpected PubKey: got %x want %x", extra.PubKey, pubKey)
	}
	if !bytes.Equal(extra.Nonce, nonce) {
		t.Fatalf("unexpected Nonce: got %x want %x", extra.Nonce, nonce)
	}
	if !bytes.Equal(extra.MergeMiningTag, mergeMining) {
		t.Fatalf("unexpected MergeMiningTag: got %x want %x", extra.MergeMiningTag, mergeMining)
	}
}

// TestConstructTXExtraEmptyInput confirms the len(mutable) == 0 base case
// (both for an entirely empty buffer and for a buffer fully consumed by
// valid tags) still returns cleanly with no error.
func TestConstructTXExtraEmptyInput(t *testing.T) {
	extra, err := ConstructTXExtra([]byte{})
	if err != nil {
		t.Fatalf("expected no error for empty input, got %v", err)
	}
	if extra.MergeMiningTag == nil {
		t.Fatalf("expected MergeMiningTag to be initialized to an empty slice, got nil")
	}
}
