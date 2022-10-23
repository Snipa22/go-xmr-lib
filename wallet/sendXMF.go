package wallet

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
)

func SendXMR(inXfer XMRWalletTransfer) (XMRTransferReceipt, error) {
	receipt := XMRTransferReceipt{}
	if reqJson, err := json.Marshal(inXfer); err != nil {
		return XMRTransferReceipt{}, err
	} else {
		resp, err := http.Post(os.Getenv("MONERO_DAEMON"), "application/json", bytes.NewReader(reqJson))
		if err != nil {
			return XMRTransferReceipt{}, err
		}

		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&receipt)
		if err != nil {
			return XMRTransferReceipt{}, err
		}
		return receipt, nil
	}
}
