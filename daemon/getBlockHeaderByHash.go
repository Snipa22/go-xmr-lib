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

func GetBlockHeaderByHash(ctx context.Context, hash string) (BlockHeader, error) {
	requestStruct := GenericParamsRequest{
		Jsonrpc: "2.0",
		ID:      "0",
		Method:  "get_block_header_by_hash",
		Params:  fmt.Sprintf(`{"hash": "%v"}`, hash),
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
		return BlockHeader{}, fmt.Errorf("xmr daemon rpc: get_block_header_by_hash returned HTTP %d: %s", resp.StatusCode, body)
	}

	decoder := json.NewDecoder(resp.Body)
	if err = decoder.Decode(&respHeader); err != nil {
		return BlockHeader{}, err
	}
	return respHeader.Result.BlockHeader, nil
}
