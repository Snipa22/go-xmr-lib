package daemon

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
)

func GetLastBlockHeader() (BlockHeader, error) {
	requestStruct := GenericJSONRPCRequest{
		Jsonrpc: "2.0",
		ID:      "0",
		Method:  "get_last_block_header",
	}
	respHeader := BlockHeaderResponse{}
	if reqJson, err := json.Marshal(requestStruct); err != nil {
		return BlockHeader{}, err
	} else {
		resp, err := http.Post(os.Getenv("MONERO_DAEMON"), "application/json", bytes.NewReader(reqJson))
		if err != nil {
			return BlockHeader{}, err
		}

		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&respHeader)
		if err != nil {
			return BlockHeader{}, err
		}
		return respHeader.Result.BlockHeader, nil
	}
}
