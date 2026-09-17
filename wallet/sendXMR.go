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
	"time"
)

// WalletHTTPClient is the *http.Client used by every wallet package RPC call. It is
// intentionally exported so callers can override it (e.g. to change the timeout or
// transport) instead of relying on http.DefaultClient, which has no timeout at all.
//
// The default 30s timeout is deliberately generous: monero-wallet-rpc's "transfer" call
// can legitimately take a while to build/sign/relay a transaction, and an overly
// aggressive client timeout would make ErrAmbiguousBroadcast outcomes (see errors.go)
// far more common than they need to be.
var WalletHTTPClient = &http.Client{Timeout: 30 * time.Second}

// SendXMR calls monero-wallet-rpc's "transfer" method to move funds out of the wallet.
// This is the money-moving call in this library: callers MUST bound it via ctx (and/or
// WalletHTTPClient's Timeout) and MUST handle errors.Is(err, ErrAmbiguousBroadcast)
// specially -- that condition means the transfer may have already been broadcast even
// though this call returned an error, and blindly retrying could double-spend. Set
// inXfer.PaymentID before calling so GetTransfers/GetTransferByTxid can later be used to
// reconcile an ambiguous outcome.
func SendXMR(ctx context.Context, inXfer XMRWalletTransfer) (XMRTransferReceipt, error) {
	receipt := XMRTransferReceipt{}
	reqData := xmrWalletTransferShell{
		Jsonrpc: "2.0",
		Id:      "0",
		Method:  "transfer",
		Params:  inXfer,
	}

	reqJson, err := json.Marshal(reqData)
	if err != nil {
		// Never left the process -- definitely not ambiguous.
		return XMRTransferReceipt{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, os.Getenv("MONERO_WALLET"), bytes.NewReader(reqJson))
	if err != nil {
		// Malformed request/URL -- never sent, definitely not ambiguous.
		return XMRTransferReceipt{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := WalletHTTPClient.Do(req)
	if err != nil {
		// We handed the request to the transport; we can no longer be sure the wallet
		// didn't receive/process it if this looks like a timeout/cancellation/reset.
		return XMRTransferReceipt{}, wrapIfAmbiguous(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return XMRTransferReceipt{}, fmt.Errorf("xmr wallet rpc: transfer returned HTTP %d: %s", resp.StatusCode, body)
	}

	decoder := json.NewDecoder(resp.Body)
	if err = decoder.Decode(&receipt); err != nil {
		// A response was received (2xx) but the body was cut off/reset mid-read, or the
		// connection died before we finished reading it -- ambiguous.
		return XMRTransferReceipt{}, wrapIfAmbiguous(err)
	}
	if receipt.Error != nil {
		// The wallet definitively responded with a JSON-RPC error (e.g. insufficient
		// funds) -- this is a clean rejection, not an ambiguous outcome.
		return receipt, errors.New(receipt.Error.Message)
	}
	return receipt, nil
}
