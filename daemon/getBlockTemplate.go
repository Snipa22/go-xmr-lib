package daemon

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Snipa22/go-xmr-lib/support"
	"net/http"
	"os"
)

func GetBlockTemplate(address string, reserve int) (BlockTemplate, error) {
	newBt := BlockTemplate{}
	requestStruct := GenericParamsRequest{
		Jsonrpc: "2.0",
		ID:      "0",
		Method:  "get_block_template",
		Params:  fmt.Sprintf(`{"wallet_address":"%v","reserve_size":%v}`, address, reserve),
	}
	daemonBTResponse := BlockTemplateResponse{}
	if reqJson, err := json.Marshal(requestStruct); err != nil {
		return BlockTemplate{}, err
	} else {
		resp, err := http.Post(os.Getenv("MONERO_DAEMON"), "application/json", bytes.NewReader(reqJson))
		if err != nil {
			return BlockTemplate{}, err
		}

		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&daemonBTResponse)
		if err != nil {
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
	}
	return newBt, nil
}
