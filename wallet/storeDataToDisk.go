package wallet

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Snipa22/go-xmr-lib/support/structs"
	"io"
	"net/http"
	"os"
)

// SaveWalletState calls monero-wallet-rpc's "store" method to flush wallet state to disk.
func SaveWalletState(ctx context.Context) error {
	requestStruct := structs.GenericRPC{
		Jsonrpc: "2.0",
		ID:      "0",
		Method:  "store",
	}

	reqJson, err := json.Marshal(requestStruct)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, os.Getenv("MONERO_WALLET"), bytes.NewReader(reqJson))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := WalletHTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return fmt.Errorf("xmr wallet rpc: store returned HTTP %d: %s", resp.StatusCode, body)
	}

	// The "store" response has no fields we care about beyond confirming success; drain
	// it so the connection can be reused, but there's nothing useful to decode into.
	_, err = io.Copy(io.Discard, resp.Body)
	return err
}
