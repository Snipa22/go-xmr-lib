package daemon

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Snipa22/go-xmr-lib/support"
	"io"
	"net/http"
	"os"
)

func GetBlockTemplate(ctx context.Context, address string, reserve int) (BlockTemplate, error) {
	newBt := BlockTemplate{}
	requestStruct := GenericParamsRequest{
		Jsonrpc: "2.0",
		ID:      "0",
		Method:  "get_block_template",
		Params:  fmt.Sprintf(`{"wallet_address":"%v","reserve_size":%v}`, address, reserve),
	}
	daemonBTResponse := BlockTemplateResponse{}

	reqJson, err := json.Marshal(requestStruct)
	if err != nil {
		return BlockTemplate{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, os.Getenv("MONERO_DAEMON"), bytes.NewReader(reqJson))
	if err != nil {
		return BlockTemplate{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := DaemonHTTPClient.Do(req)
	if err != nil {
		return BlockTemplate{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return BlockTemplate{}, fmt.Errorf("xmr daemon rpc: get_block_template returned HTTP %d: %s", resp.StatusCode, body)
	}

	decoder := json.NewDecoder(resp.Body)
	if err = decoder.Decode(&daemonBTResponse); err != nil {
		return BlockTemplate{}, err
	}
	// Convert the inbound raw BT into a useful stateful object
	newBt.Template, err = support.ParseBlockFromTemplateBlob(daemonBTResponse.Result.BlocktemplateBlob)
	if err != nil {
		return BlockTemplate{}, err
	}
	newBt.SeedHash = daemonBTResponse.Result.SeedHash
	newBt.Difficulty = daemonBTResponse.Result.Difficulty
	newBt.Height = daemonBTResponse.Result.Height
	newBt.ReservedOffset = daemonBTResponse.Result.ReservedOffset
	newBt.RawTemplate = daemonBTResponse.Result.BlocktemplateBlob
	newBt.RawResponse = daemonBTResponse
	return newBt, nil
}
