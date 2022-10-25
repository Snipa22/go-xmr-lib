package wallet

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os"
)

func SendXMR(inXfer XMRWalletTransfer) (XMRTransferReceipt, error) {
	receipt := XMRTransferReceipt{}
	reqData := xmrWalletTransferShell{
		Jsonrpc: "2.0",
		Id:      "0",
		Method:  "transfer",
		Params:  inXfer,
	}
	if reqJson, err := json.Marshal(reqData); err != nil {
		return XMRTransferReceipt{}, err
	} else {
		resp, err := http.Post(os.Getenv("MONERO_WALLET"), "application/json", bytes.NewReader(reqJson))
		if err != nil {
			return XMRTransferReceipt{}, err
		}

		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&receipt)
		if err != nil {
			return XMRTransferReceipt{}, err
		}
		if receipt.Error != nil {
			return receipt, errors.New(receipt.Error.Message)
		}
		return receipt, nil
	}
}
