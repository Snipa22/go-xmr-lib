package structs

import "github.com/Snipa22/go-xmr-lib/support/serialization"

type GenericRPC struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      string `json:"id"`
	Method  string `json:"method"`
}

type BlockTemplateData struct {
	BlockTemplate serialization.Block
	Difficulty    int
	Height        int
	ReservedSize  int
	PreviousHash  string
	SeedHash      string
}
