package daemon

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func GetLastBlockHeader(ctx context.Context) (BlockHeader, error) {
	requestStruct := GenericJSONRPCRequest{
		Jsonrpc: "2.0",
		ID:      "0",
		Method:  "get_last_block_header",
	}
	respHeader := BlockHeaderResponse{}

	reqJson, err := json.Marshal(requestStruct)
	if err != nil {
		return BlockHeader{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, os.Getenv("MONERO_DAEMON"), bytes.NewReader(reqJson))
	if err != nil {
		return BlockHeader{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := DaemonHTTPClient.Do(req)
	if err != nil {
		return BlockHeader{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return BlockHeader{}, fmt.Errorf("xmr daemon rpc: get_last_block_header returned HTTP %d: %s", resp.StatusCode, body)
	}

	decoder := json.NewDecoder(resp.Body)
	if err = decoder.Decode(&respHeader); err != nil {
		return BlockHeader{}, err
	}
	return respHeader.Result.BlockHeader, nil
}
