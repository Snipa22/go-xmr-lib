package daemon

import (
	"net/http"
	"time"
)

// DaemonHTTPClient is the *http.Client used by every daemon package RPC call. It is
// intentionally exported so callers can override it (e.g. to change the timeout or
// transport) instead of relying on http.DefaultClient, which has no timeout at all.
var DaemonHTTPClient = &http.Client{Timeout: 30 * time.Second}
