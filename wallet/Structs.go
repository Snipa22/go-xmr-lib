package wallet

type XMRWalletTransfer struct {
	Jsonrpc string `json:"jsonrpc"`
	Id      string `json:"id"`
	Method  string `json:"method"`
	Params  struct {
		Destinations []struct {
			Amount  int64  `json:"amount"`
			Address string `json:"address"`
		} `json:"destinations"`
		Priority int  `json:"priority"`
		RingSize int  `json:"ring_size"`
		GetTxKey bool `json:"get_tx_key"`
	} `json:"params"`
}

type XMRTransferReceipt struct {
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
