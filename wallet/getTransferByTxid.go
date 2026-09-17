package wallet

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

// GetTransferByTxid calls monero-wallet-rpc's "get_transfer_by_txid" method. It is a
// read-only query, so unlike SendXMR its errors are not classified against
// ErrAmbiguousBroadcast -- a failure here just means "we couldn't look it up", not "we
// don't know if a transfer happened."
//
// accountIndex is optional; pass nil to omit it (monero-wallet-rpc defaults to account 0).
func GetTransferByTxid(ctx context.Context, txid string, accountIndex *int) (XMRGetTransferByTxidResponse, error) {
	result := XMRGetTransferByTxidResponse{}
	reqData := getTransferByTxidShell{
		Jsonrpc: "2.0",
		Id:      "0",
		Method:  "get_transfer_by_txid",
		Params: getTransferByTxidRequestParams{
			Txid:         txid,
			AccountIndex: accountIndex,
		},
	}

	reqJson, err := json.Marshal(reqData)
	if err != nil {
		return XMRGetTransferByTxidResponse{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, os.Getenv("MONERO_WALLET"), bytes.NewReader(reqJson))
	if err != nil {
		return XMRGetTransferByTxidResponse{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := WalletHTTPClient.Do(req)
	if err != nil {
		return XMRGetTransferByTxidResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return XMRGetTransferByTxidResponse{}, fmt.Errorf("xmr wallet rpc: get_transfer_by_txid returned HTTP %d: %s", resp.StatusCode, body)
	}

	decoder := json.NewDecoder(resp.Body)
	if err = decoder.Decode(&result); err != nil {
		return XMRGetTransferByTxidResponse{}, err
	}
	if result.Error != nil {
		return result, errors.New(result.Error.Message)
	}
	return result, nil
}
