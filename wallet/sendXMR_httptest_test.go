package wallet

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

// withMoneroWalletEnv points MONERO_WALLET at url for the duration of the test and
// restores the previous value afterward.
func withMoneroWalletEnv(t *testing.T, url string) {
	t.Helper()
	prev, had := os.LookupEnv("MONERO_WALLET")
	if err := os.Setenv("MONERO_WALLET", url); err != nil {
		t.Fatalf("failed to set MONERO_WALLET: %v", err)
	}
	t.Cleanup(func() {
		if had {
			os.Setenv("MONERO_WALLET", prev)
		} else {
			os.Unsetenv("MONERO_WALLET")
		}
	})
}

// TestSendXMR_MockServerSuccess covers a mock server returning a valid 200 JSON-RPC
// transfer response, asserting fields decode correctly and that PaymentID round-trips
// through the request body the mock server observes.
func TestSendXMR_MockServerSuccess(t *testing.T) {
	var observedBody xmrWalletTransferShell

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("mock server: failed to read request body: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err := json.Unmarshal(body, &observedBody); err != nil {
			t.Errorf("mock server: failed to unmarshal request body: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"id": "0",
			"jsonrpc": "2.0",
			"result": {
				"amount": 100000000000,
				"fee": 86897600000,
				"multisig_txset": "",
				"tx_blob": "",
				"tx_hash": "7663438de4f72b25a0e395b770ea9ecf7108cd2f0c4b75be0b14a103d3362be9",
				"tx_key": "somekey",
				"tx_metadata": "",
				"unsigned_txset": ""
			}
		}`))
	}))
	defer server.Close()

	withMoneroWalletEnv(t, server.URL)

	paymentID := "0123456789abcdef"
	receipt, err := SendXMR(context.Background(), XMRWalletTransfer{
		Destinations: []XMRPayment{
			{Amount: 100000000000, Address: "someaddress"},
		},
		Priority:  1,
		RingSize:  16,
		GetTxKey:  true,
		PaymentID: &paymentID,
	})
	if err != nil {
		t.Fatalf("SendXMR() unexpected error: %v", err)
	}
	if receipt.Result.TxHash != "7663438de4f72b25a0e395b770ea9ecf7108cd2f0c4b75be0b14a103d3362be9" {
		t.Errorf("unexpected TxHash: %v", receipt.Result.TxHash)
	}
	if receipt.Result.Amount != 100000000000 {
		t.Errorf("unexpected Amount: %v", receipt.Result.Amount)
	}
	if receipt.Result.Fee != 86897600000 {
		t.Errorf("unexpected Fee: %v", receipt.Result.Fee)
	}

	if observedBody.Params.PaymentID == nil || *observedBody.Params.PaymentID != paymentID {
		t.Errorf("PaymentID did not round-trip through the request body observed by the mock server: got %v, want %v", observedBody.Params.PaymentID, paymentID)
	}
}

// TestSendXMR_ContextTimeoutIsAmbiguous covers SendXMR's context/timeout behavior: a
// mock server that sleeps past a short context deadline must return an error that is
// errors.Is(err, context.DeadlineExceeded) AND errors.Is(err, ErrAmbiguousBroadcast),
// because we can't tell whether the wallet processed the transfer before we gave up
// waiting on the response.
func TestSendXMR_ContextTimeoutIsAmbiguous(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id":"0","jsonrpc":"2.0","result":{}}`))
	}))
	defer server.Close()

	withMoneroWalletEnv(t, server.URL)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	_, err := SendXMR(ctx, XMRWalletTransfer{
		Destinations: []XMRPayment{{Amount: 100, Address: "someaddress"}},
		Priority:     1,
		RingSize:     16,
	})
	if err == nil {
		t.Fatal("SendXMR() expected an error due to context timeout, got nil")
	}
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("SendXMR() error = %v, want errors.Is(err, context.DeadlineExceeded) to be true", err)
	}
	if !errors.Is(err, ErrAmbiguousBroadcast) {
		t.Errorf("SendXMR() error = %v, want errors.Is(err, ErrAmbiguousBroadcast) to be true", err)
	}
}

// TestSendXMR_NonSuccessStatusIsDistinctError covers a mock server returning a non-2xx
// status. This must produce a distinct status-code error (not a generic JSON-decode
// error) and must NOT be classified as ErrAmbiguousBroadcast: a clean HTTP-level
// rejection with a readable body is a definitive outcome, not an ambiguous one.
func TestSendXMR_NonSuccessStatusIsDistinctError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("internal server error: wallet is locked"))
	}))
	defer server.Close()

	withMoneroWalletEnv(t, server.URL)

	_, err := SendXMR(context.Background(), XMRWalletTransfer{
		Destinations: []XMRPayment{{Amount: 100, Address: "someaddress"}},
		Priority:     1,
		RingSize:     16,
	})
	if err == nil {
		t.Fatal("SendXMR() expected an error due to non-2xx status, got nil")
	}
	if errors.Is(err, ErrAmbiguousBroadcast) {
		t.Errorf("SendXMR() error = %v, a clean non-2xx status must NOT be classified as ErrAmbiguousBroadcast", err)
	}
	if !strings.Contains(err.Error(), strconv.Itoa(http.StatusInternalServerError)) {
		t.Errorf("SendXMR() error = %v, want it to mention the HTTP status code %d", err, http.StatusInternalServerError)
	}
	if !strings.Contains(err.Error(), "wallet is locked") {
		t.Errorf("SendXMR() error = %v, want it to include a snippet of the response body", err)
	}
}

// TestSendXMR_ConnectionResetIsAmbiguous covers the "connection reset mid-response"
// class of ambiguous failure: the server accepts the connection, then hangs up abruptly
// without ever writing a response. This must be classified as ErrAmbiguousBroadcast.
func TestSendXMR_ConnectionResetIsAmbiguous(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	defer ln.Close()

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			// Accept the connection (and thus the request) then hang up without
			// writing anything back, simulating a mid-flight connection reset.
			conn.Close()
		}
	}()

	withMoneroWalletEnv(t, "http://"+ln.Addr().String())

	_, err = SendXMR(context.Background(), XMRWalletTransfer{
		Destinations: []XMRPayment{{Amount: 100, Address: "someaddress"}},
		Priority:     1,
		RingSize:     16,
	})
	if err == nil {
		t.Fatal("SendXMR() expected an error due to connection reset, got nil")
	}
	if !errors.Is(err, ErrAmbiguousBroadcast) {
		t.Errorf("SendXMR() error = %v, want errors.Is(err, ErrAmbiguousBroadcast) to be true", err)
	}
}

// closeTrackingBody wraps an io.ReadCloser and records whether Close was called on it,
// so tests can assert that response bodies are actually closed by SendXMR.
type closeTrackingBody struct {
	io.ReadCloser
	closed *bool
}

func (c *closeTrackingBody) Close() error {
	*c.closed = true
	return c.ReadCloser.Close()
}

// closeTrackingTransport wraps http.DefaultTransport, replacing every response body
// with a closeTrackingBody so the test can observe whether it was closed.
type closeTrackingTransport struct {
	closed *bool
}

func (t *closeTrackingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return resp, err
	}
	resp.Body = &closeTrackingBody{ReadCloser: resp.Body, closed: t.closed}
	return resp, nil
}

// TestSendXMR_ClosesResponseBody asserts that SendXMR closes the HTTP response body it
// received, so successful RPC calls don't leak connections/file descriptors. This is a
// black-box way to test resp.Body.Close() behavior via a custom http.RoundTripper that
// wraps the body and records whether Close() was invoked.
func TestSendXMR_ClosesResponseBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id":"0","jsonrpc":"2.0","result":{"tx_hash":"deadbeef"}}`))
	}))
	defer server.Close()

	withMoneroWalletEnv(t, server.URL)

	closed := false
	prevClient := WalletHTTPClient
	WalletHTTPClient = &http.Client{
		Timeout:   prevClient.Timeout,
		Transport: &closeTrackingTransport{closed: &closed},
	}
	t.Cleanup(func() { WalletHTTPClient = prevClient })

	if _, err := SendXMR(context.Background(), XMRWalletTransfer{
		Destinations: []XMRPayment{{Amount: 100, Address: "someaddress"}},
		Priority:     1,
		RingSize:     16,
	}); err != nil {
		t.Fatalf("SendXMR() unexpected error: %v", err)
	}
	if !closed {
		t.Error("SendXMR() did not close the HTTP response body")
	}
}
