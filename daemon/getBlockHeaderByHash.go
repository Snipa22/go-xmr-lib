package daemon

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func GetBlockHeaderByHash(hash string) (BlockHeader, error) {
	requestStruct := GenericParamsRequest{
		Jsonrpc: "2.0",
		ID:      "0",
		Method:  "get_block_header_by_hash",
		Params:  fmt.Sprintf(`{"hash": "%v"}`, hash),
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
