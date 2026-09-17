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

// GetTransfers calls monero-wallet-rpc's "get_transfers" method. It is a read-only query,
// so unlike SendXMR its errors are not classified against ErrAmbiguousBroadcast -- a
// failure here just means "we couldn't look it up", not "we don't know if a transfer
// happened."
//
// Combined with GetTransferByTxid, this gives a caller a real way to reconcile an
// ErrAmbiguousBroadcast outcome from SendXMR by checking whether a specific tx_hash or
// payment_id from the original request actually landed.
func GetTransfers(ctx context.Context, params GetTransfersParams) (XMRGetTransfersResponse, error) {
	result := XMRGetTransfersResponse{}
	reqData := getTransfersShell{
		Jsonrpc: "2.0",
		Id:      "0",
		Method:  "get_transfers",
		Params:  params,
	}

	reqJson, err := json.Marshal(reqData)
	if err != nil {
		return XMRGetTransfersResponse{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, os.Getenv("MONERO_WALLET"), bytes.NewReader(reqJson))
	if err != nil {
		return XMRGetTransfersResponse{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := WalletHTTPClient.Do(req)
	if err != nil {
		return XMRGetTransfersResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return XMRGetTransfersResponse{}, fmt.Errorf("xmr wallet rpc: get_transfers returned HTTP %d: %s", resp.StatusCode, body)
	}

	decoder := json.NewDecoder(resp.Body)
	if err = decoder.Decode(&result); err != nil {
		return XMRGetTransfersResponse{}, err
	}
	if result.Error != nil {
		return result, errors.New(result.Error.Message)
	}
	return result, nil
}
