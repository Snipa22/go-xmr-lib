package wallet

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestGetTransferByTxid_MockServer covers GetTransferByTxid against a mock JSON-RPC
// response shaped like the real monero-wallet-rpc docs for get_transfer_by_txid.
func TestGetTransferByTxid_MockServer(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"id": "0",
			"jsonrpc": "2.0",
			"result": {
				"transfer": {
					"address": "53zii2WaqQwZU4oUsCUcrHgaSv2CrUGCSFJLdQnkLPyH7ZLPYHjtoHhi14dqjF6jywNRknYLwbate2eGv8TuZcS7GuR7wMY",
					"amount": 100000000000,
					"amounts": [100000000000],
					"confirmations": 19,
					"double_spend_seen": false,
					"fee": 53840000,
					"height": 1140109,
					"locked": false,
					"note": "",
					"payment_id": "0000000000000000",
					"subaddr_index": {"major": 0, "minor": 0},
					"suggested_confirmations_threshold": 1,
					"timestamp": 1658360753,
					"txid": "765f7124d01bd2eb2d4e7e59aa44a28c24339a41e4009f463955b087017b0ca3",
					"type": "in",
					"unlock_time": 0
				},
				"transfers": [{
					"address": "53zii2WaqQwZU4oUsCUcrHgaSv2CrUGCSFJLdQnkLPyH7ZLPYHjtoHhi14dqjF6jywNRknYLwbate2eGv8TuZcS7GuR7wMY",
					"amount": 100000000000,
					"amounts": [100000000000],
					"confirmations": 19,
					"double_spend_seen": false,
					"fee": 53840000,
					"height": 1140109,
					"locked": false,
					"note": "",
					"payment_id": "0000000000000000",
					"subaddr_index": {"major": 0, "minor": 0},
					"suggested_confirmations_threshold": 1,
					"timestamp": 1658360753,
					"txid": "765f7124d01bd2eb2d4e7e59aa44a28c24339a41e4009f463955b087017b0ca3",
					"type": "in",
					"unlock_time": 0
				}]
			}
		}`))
	}))
	defer server.Close()

	withMoneroWalletEnv(t, server.URL)

	resp, err := GetTransferByTxid(context.Background(), "765f7124d01bd2eb2d4e7e59aa44a28c24339a41e4009f463955b087017b0ca3", nil)
	if err != nil {
		t.Fatalf("GetTransferByTxid() unexpected error: %v", err)
	}
	if resp.Result.Transfer.Txid != "765f7124d01bd2eb2d4e7e59aa44a28c24339a41e4009f463955b087017b0ca3" {
		t.Errorf("unexpected Transfer.Txid: %v", resp.Result.Transfer.Txid)
	}
	if resp.Result.Transfer.Amount != 100000000000 {
		t.Errorf("unexpected Transfer.Amount: %v", resp.Result.Transfer.Amount)
	}
	if resp.Result.Transfer.Confirmations != 19 {
		t.Errorf("unexpected Transfer.Confirmations: %v", resp.Result.Transfer.Confirmations)
	}
	if resp.Result.Transfer.PaymentID != "0000000000000000" {
		t.Errorf("unexpected Transfer.PaymentID: %v", resp.Result.Transfer.PaymentID)
	}
	if resp.Result.Transfer.SubaddrIndex != (XMRSubaddrIndex{Major: 0, Minor: 0}) {
		t.Errorf("unexpected Transfer.SubaddrIndex: %v", resp.Result.Transfer.SubaddrIndex)
	}
	if resp.Result.Transfer.Type != "in" {
		t.Errorf("unexpected Transfer.Type: %v", resp.Result.Transfer.Type)
	}
	if len(resp.Result.Transfers) != 1 {
		t.Fatalf("unexpected Transfers length: %v", len(resp.Result.Transfers))
	}
}

// TestGetTransfers_MockServer covers GetTransfers against a mock JSON-RPC response
// shaped like the real monero-wallet-rpc docs for get_transfers.
func TestGetTransfers_MockServer(t *testing.T) {
	var observedBody getTransfersShell

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, readErr := io.ReadAll(r.Body)
		if readErr != nil {
			t.Errorf("mock server: failed to read request body: %v", readErr)
		} else if unmarshalErr := json.Unmarshal(body, &observedBody); unmarshalErr != nil {
			t.Errorf("mock server: failed to unmarshal request body: %v", unmarshalErr)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"id": "0",
			"jsonrpc": "2.0",
			"result": {
				"in": [{
					"address": "77Vx9cs1VPicFndSVgYUvTdLCJEZw9h81hXLMYsjBCXSJfUehLa9TDW3Ffh45SQa7xb6dUs18mpNxfUhQGqfwXPSMrvKhVp",
					"amount": 200000000000,
					"amounts": [200000000000],
					"confirmations": 1,
					"double_spend_seen": false,
					"fee": 21650200000,
					"height": 153624,
					"locked": false,
					"note": "",
					"payment_id": "0000000000000000",
					"subaddr_index": {"major": 1, "minor": 0},
					"suggested_confirmations_threshold": 1,
					"timestamp": 1535918400,
					"txid": "c36258a276018c3a4bc1f195a7fb530f50cd63a4fa765fb7c6f7f49fc051762a",
					"type": "in",
					"unlock_time": 0
				}]
			}
		}`))
	}))
	defer server.Close()

	withMoneroWalletEnv(t, server.URL)

	accountIndex := 1
	resp, err := GetTransfers(context.Background(), GetTransfersParams{
		In:           true,
		AccountIndex: &accountIndex,
	})
	if err != nil {
		t.Fatalf("GetTransfers() unexpected error: %v", err)
	}
	if len(resp.Result.In) != 1 {
		t.Fatalf("unexpected In length: %v", len(resp.Result.In))
	}
	if resp.Result.In[0].Txid != "c36258a276018c3a4bc1f195a7fb530f50cd63a4fa765fb7c6f7f49fc051762a" {
		t.Errorf("unexpected In[0].Txid: %v", resp.Result.In[0].Txid)
	}
	if resp.Result.In[0].Amount != 200000000000 {
		t.Errorf("unexpected In[0].Amount: %v", resp.Result.In[0].Amount)
	}
	if len(resp.Result.Out) != 0 {
		t.Errorf("unexpected Out length: %v", len(resp.Result.Out))
	}

	if observedBody.Params.In != true {
		t.Errorf("In did not round-trip through the request body observed by the mock server: got %v", observedBody.Params.In)
	}
	if observedBody.Params.AccountIndex == nil || *observedBody.Params.AccountIndex != accountIndex {
		t.Errorf("AccountIndex did not round-trip through the request body observed by the mock server: got %v, want %v", observedBody.Params.AccountIndex, accountIndex)
	}
}
