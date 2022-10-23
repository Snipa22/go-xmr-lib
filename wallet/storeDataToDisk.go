package wallet

import (
	"bytes"
	"encoding/json"
	"github.com/Snipa22/go-xmr-lib/support/structs"
	"net/http"
	"os"
)

func SaveWalletState() error {
	requestStruct := structs.GenericRPC{
		Jsonrpc: "2.0",
		ID:      "0",
		Method:  "store",
	}
	if reqJson, err := json.Marshal(requestStruct); err != nil {
		return err
	} else {
		_, err = http.Post(os.Getenv("MONERO_WALLET"), "application/json", bytes.NewReader(reqJson))
		if err != nil {
			return err
		}
		return nil
	}
}
