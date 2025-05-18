package daemon

import "github.com/Snipa22/go-xmr-lib/support/serialization"

type BlockHeader struct {
	BlockSize                 int    `json:"block_size"`
	BlockWeight               int    `json:"block_weight"`
	CumulativeDifficulty      int64  `json:"cumulative_difficulty"`
	CumulativeDifficultyTop64 int    `json:"cumulative_difficulty_top64"`
	Depth                     int    `json:"depth"`
	Difficulty                int    `json:"difficulty"`
	DifficultyTop64           int    `json:"difficulty_top64"`
	Hash                      string `json:"hash"`
	Height                    int    `json:"height"`
	LongTermWeight            int    `json:"long_term_weight"`
	MajorVersion              int    `json:"major_version"`
	MinerTxHash               string `json:"miner_tx_hash"`
	MinorVersion              int    `json:"minor_version"`
	Nonce                     int    `json:"nonce"`
	NumTxes                   int    `json:"num_txes"`
	OrphanStatus              bool   `json:"orphan_status"`
	PowHash                   string `json:"pow_hash"`
	PrevHash                  string `json:"prev_hash"`
	Reward                    int64  `json:"reward"`
	Timestamp                 int    `json:"timestamp"`
	WideCumulativeDifficulty  string `json:"wide_cumulative_difficulty"`
	WideDifficulty            string `json:"wide_difficulty"`
}

type BlockTemplate struct {
	SeedHash       string
	Template       serialization.Block
	Difficulty     int64
	Height         int
	RawTemplate    string
	ReservedOffset int
	RawResponse    BlockTemplateResponse
}

type BlockHeaderResponse struct {
	ID      string `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
	Result  struct {
		BlockHeader BlockHeader `json:"block_header"`
		Credits     int         `json:"credits"`
		Status      string      `json:"status"`
		TopHash     string      `json:"top_hash"`
		Untrusted   bool        `json:"untrusted"`
	} `json:"result"`
}

type BlockTemplateResponse struct {
	ID      string `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
	Result  struct {
		BlockhashingBlob  string `json:"blockhashing_blob"`
		BlocktemplateBlob string `json:"blocktemplate_blob"`
		Difficulty        int64  `json:"difficulty"`
		DifficultyTop64   int    `json:"difficulty_top64"`
		ExpectedReward    int64  `json:"expected_reward"`
		Height            int    `json:"height"`
		NextSeedHash      string `json:"next_seed_hash"`
		PrevHash          string `json:"prev_hash"`
		ReservedOffset    int    `json:"reserved_offset"`
		SeedHash          string `json:"seed_hash"`
		SeedHeight        int    `json:"seed_height"`
		Status            string `json:"status"`
		Untrusted         bool   `json:"untrusted"`
		WideDifficulty    string `json:"wide_difficulty"`
	} `json:"result"`
}

type GenericParamsRequest struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      string `json:"id"`
	Method  string `json:"method"`
	Params  string `json:"params"`
}

type GenericJSONRPCRequest struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      string `json:"id"`
	Method  string `json:"method"`
}
