package wallet

import (
	"context"
	"errors"
	"io"
	"net"
	"strings"
	"syscall"
)

// ErrAmbiguousBroadcast is returned (wrapped) by SendXMR when the outcome of a transfer
// request is genuinely unknown: the request may have reached monero-wallet-rpc and been
// relayed/broadcast, but we were unable to confirm this because of a timeout, cancellation,
// connection reset, or a response that was cut off mid-read. Callers MUST NOT assume the
// transfer did not happen in this case -- use GetTransferByTxid/GetTransfers (keyed on a
// PaymentID set on the original request, or on a tx_hash captured from a prior successful
// response) to reconcile before retrying, to avoid double-spending funds.
//
// Use errors.Is(err, ErrAmbiguousBroadcast) to detect this condition.
var ErrAmbiguousBroadcast = errors.New("xmr wallet rpc: ambiguous result, transaction may have been broadcast")

// isAmbiguousTransportError reports whether err represents a class of failure where we
// genuinely can't tell if monero-wallet-rpc processed (and possibly relayed) the request
// before the failure occurred: context deadline/cancellation, a net.Error that reports
// Timeout(), or a connection-reset/EOF-mid-response style error. It deliberately does NOT
// match plain JSON-marshal errors (those never leave the process) or clean non-2xx/JSON-RPC
// error responses (the wallet definitively rejected the request and told us so).
func isAmbiguousTransportError(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
		return true
	}
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return true
	}
	// io.EOF/io.ErrUnexpectedEOF here means the connection was closed/cut in the middle of
	// (or before) reading the response, after the request had already gone out over the
	// wire -- we cannot tell whether the wallet processed it.
	if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
		return true
	}
	if errors.Is(err, syscall.ECONNRESET) {
		return true
	}
	// Fallback substring detection: on some platforms/transports a connection reset or
	// broken pipe surfaces only as a formatted *net.OpError string rather than a typed
	// syscall.Errno we can errors.Is against.
	msg := err.Error()
	if strings.Contains(msg, "connection reset by peer") || strings.Contains(msg, "broken pipe") {
		return true
	}
	return false
}

// wrapIfAmbiguous wraps err so that errors.Is(result, ErrAmbiguousBroadcast) is true AND
// errors.Is(result, err) (the original cause) is true, if and only if err looks like one of
// the ambiguous-outcome cases described on isAmbiguousTransportError. Otherwise err is
// returned unchanged.
func wrapIfAmbiguous(err error) error {
	if !isAmbiguousTransportError(err) {
		return err
	}
	return &ambiguousBroadcastError{cause: err}
}

// ambiguousBroadcastError wraps an underlying transport-level error while also exposing
// ErrAmbiguousBroadcast via Unwrap, so callers can use errors.Is against either.
type ambiguousBroadcastError struct {
	cause error
}

func (e *ambiguousBroadcastError) Error() string {
	return ErrAmbiguousBroadcast.Error() + ": " + e.cause.Error()
}

func (e *ambiguousBroadcastError) Unwrap() []error {
	return []error{ErrAmbiguousBroadcast, e.cause}
}
