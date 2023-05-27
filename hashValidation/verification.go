package hashValidation

type Validator interface {
	Hash(input []byte, seed []byte) ([]byte, error)
	NewSeed(input []byte) error
	SetCurrentSeed(input []byte)
	Info() (ValidatorInfo, error)
}

type ValidatorInfo struct {
	RandomxService string `json:"randomx_service"`
	Algorithm      string `json:"algorithm"`
	Threads        int    `json:"threads"`
	Seed           string `json:"seed"`
	Hashes         int    `json:"hashes"`
}
