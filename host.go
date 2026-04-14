package main

// Host functions provided by the ApiQube engine to this plugin.
//
// When built as a WASM module, these are linked at load time by wazero.
// The //go:wasmimport directive only takes effect under TinyGo WASI build,
// so a regular `go build` compiles the stubs below cleanly.

// httpResponse mirrors what the host returns from host_http_request.
type httpResponse struct {
	Status     int               `json:"status"`
	Headers    map[string]string `json:"headers"`
	Body       []byte            `json:"body"`
	DurationMs int64             `json:"duration_ms"`
	Error      string            `json:"error,omitempty"`
}

// httpRequest is the Go-friendly wrapper around host_http_request.
func httpRequest(method, url string, headers map[string]string, body []byte) (*httpResponse, error) {
	// TODO: implementation
	// 1. Marshal headers to JSON
	// 2. Pack method/url/headers/body as ptr+len pairs
	// 3. Call hostHTTPRequestImport (see below)
	// 4. Read response bytes, unmarshal into httpResponse
	_ = method
	_ = url
	_ = headers
	_ = body
	return nil, nil
}

// logMessage sends a log line to the host for display/storage.
func logMessage(level, msg string) {
	// TODO: implementation via host_log import
	_ = level
	_ = msg
}

// now returns the current Unix timestamp in milliseconds from the host.
// (Plugins don't have direct access to system clock in WASI.)
func now() int64 {
	// TODO: implementation via host_now import
	return 0
}
