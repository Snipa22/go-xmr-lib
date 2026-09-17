package wallet

type xmrWalletTransferShell struct {
	Jsonrpc string            `json:"jsonrpc"`
	Id      string            `json:"id"`
	Method  string            `json:"method"`
	Params  XMRWalletTransfer `json:"params"`
}

type XMRWalletTransfer struct {
	Destinations []XMRPayment `json:"destinations"`
	Priority     int          `json:"priority"`
	RingSize     int          `json:"ring_size"`
	GetTxKey     bool         `json:"get_tx_key"`

	// AccountIndex selects which subaddress account to send from (default 0 if omitted).
	AccountIndex *int `json:"account_index,omitempty"`
	// SubaddrIndices restricts which subaddresses within the account may be used as
	// transfer sources.
	SubaddrIndices []int `json:"subaddr_indices,omitempty"`
	// PaymentID tags this transfer request so it can be looked up later via
	// GetTransfers/GetTransferByTxid to reconcile an ErrAmbiguousBroadcast outcome from
	// SendXMR.
	PaymentID *string `json:"payment_id,omitempty"`
	// UnlockTime is the number of blocks before the monero can be spent (0/omitted = no
	// additional unlock time beyond the default).
	UnlockTime *int `json:"unlock_time,omitempty"`
	// DoNotRelay, if true, builds but does not relay/broadcast the transaction.
	DoNotRelay *bool `json:"do_not_relay,omitempty"`
}

type XMRPayment struct {
	Amount  int64  `json:"amount"`
	Address string `json:"address"`
}

type XMRTransferReceipt struct {
	Error *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
	Id      string `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
	Result  struct {
		Amount        int64   `json:"amount"`
		Fee           int64   `json:"fee"`
		MultisigTxset *string `json:"multisig_txset"`
		TxBlob        *string `json:"tx_blob"`
		TxHash        string  `json:"tx_hash"`
		TxKey         *string `json:"tx_key"`
		TxMetadata    *string `json:"tx_metadata"`
		UnsignedTxset *string `json:"unsigned_txset"`
	} `json:"result"`
}

// XMRSubaddrIndex identifies a subaddress by its account (major) and address (minor)
// index, as used throughout monero-wallet-rpc's transfer-history methods.
type XMRSubaddrIndex struct {
	Major int `json:"major"`
	Minor int `json:"minor"`
}

// XMRTransferDestination is one destination of an outgoing transfer, as reported by
// get_transfer_by_txid/get_transfers (only populated for outgoing transactions, and only
// if the wallet cache constructed the transaction itself).
type XMRTransferDestination struct {
	Amount  int64  `json:"amount"`
	Address string `json:"address"`
}

// XMRTransferInfo mirrors the "transfer" object documented for monero-wallet-rpc's
// get_transfer_by_txid and get_transfers methods.
type XMRTransferInfo struct {
	Address                         string                   `json:"address"`
	Amount                          int64                    `json:"amount"`
	Amounts                         []int64                  `json:"amounts"`
	Confirmations                   int64                    `json:"confirmations"`
	Destinations                    []XMRTransferDestination `json:"destinations,omitempty"`
	DoubleSpendSeen                 bool                     `json:"double_spend_seen"`
	Fee                             int64                    `json:"fee"`
	Height                          int64                    `json:"height"`
	Locked                          bool                     `json:"locked"`
	Note                            string                   `json:"note"`
	PaymentID                       string                   `json:"payment_id"`
	SubaddrIndex                    XMRSubaddrIndex          `json:"subaddr_index"`
	SubaddrIndices                  []XMRSubaddrIndex        `json:"subaddr_indices,omitempty"`
	SuggestedConfirmationsThreshold int64                    `json:"suggested_confirmations_threshold"`
	Timestamp                       int64                    `json:"timestamp"`
	Txid                            string                   `json:"txid"`
	Type                            string                   `json:"type"`
	UnlockTime                      int64                    `json:"unlock_time"`
}

type getTransferByTxidRequestParams struct {
	Txid         string `json:"txid"`
	AccountIndex *int   `json:"account_index,omitempty"`
}

type getTransferByTxidShell struct {
	Jsonrpc string                         `json:"jsonrpc"`
	Id      string                         `json:"id"`
	Method  string                         `json:"method"`
	Params  getTransferByTxidRequestParams `json:"params"`
}

// XMRGetTransferByTxidResponse is the decoded JSON-RPC response of get_transfer_by_txid.
type XMRGetTransferByTxidResponse struct {
	Error *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
	Id      string `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
	Result  struct {
		Transfer  XMRTransferInfo   `json:"transfer"`
		Transfers []XMRTransferInfo `json:"transfers"`
	} `json:"result"`
}

// GetTransfersParams mirrors the documented input params of monero-wallet-rpc's
// get_transfers method.
type GetTransfersParams struct {
	In             bool  `json:"in,omitempty"`
	Out            bool  `json:"out,omitempty"`
	Pending        bool  `json:"pending,omitempty"`
	Failed         bool  `json:"failed,omitempty"`
	Pool           bool  `json:"pool,omitempty"`
	FilterByHeight *bool `json:"filter_by_height,omitempty"`
	MinHeight      *int  `json:"min_height,omitempty"`
	MaxHeight      *int  `json:"max_height,omitempty"`
	AccountIndex   *int  `json:"account_index,omitempty"`
	SubaddrIndices []int `json:"subaddr_indices,omitempty"`
	AllAccounts    *bool `json:"all_accounts,omitempty"`
}

type getTransfersShell struct {
	Jsonrpc string             `json:"jsonrpc"`
	Id      string             `json:"id"`
	Method  string             `json:"method"`
	Params  GetTransfersParams `json:"params"`
}

// XMRGetTransfersResponse is the decoded JSON-RPC response of get_transfers.
type XMRGetTransfersResponse struct {
	Error *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
	Id      string `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
	Result  struct {
		In      []XMRTransferInfo `json:"in,omitempty"`
		Out     []XMRTransferInfo `json:"out,omitempty"`
		Pending []XMRTransferInfo `json:"pending,omitempty"`
		Failed  []XMRTransferInfo `json:"failed,omitempty"`
		Pool    []XMRTransferInfo `json:"pool,omitempty"`
	} `json:"result"`
}
