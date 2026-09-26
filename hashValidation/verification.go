package hashValidation

import "context"

type Validator interface {
	Hash(input []byte, seed []byte) ([]byte, error)
	NewSeed(input []byte) error
	SetCurrentSeed(input []byte)
	Info() (*ValidatorInfo, error)
}

// ValidatorWithContext is the context-aware counterpart of Validator. It is
// additive: Validator itself is deliberately unchanged so existing
// implementations and callers keep compiling, but a caller that needs its own
// deadline actually enforced on the wire (e.g. a bounded share-validation
// worker pool, where a hung request otherwise consumes a worker slot
// indefinitely) should depend on this interface instead.
type ValidatorWithContext interface {
	Validator
	HashWithContext(ctx context.Context, input []byte, seed []byte) ([]byte, error)
	NewSeedWithContext(ctx context.Context, input []byte) error
	InfoWithContext(ctx context.Context) (*ValidatorInfo, error)
}

type ValidatorInfo struct {
	RandomxService string `json:"randomx_service"`
	Algorithm      string `json:"algorithm"`
	Threads        int    `json:"threads"`
	Seed           string `json:"seed"`
	Hashes         int    `json:"hashes"`
}
