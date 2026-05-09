//go:build wasip1

package main

// host_log writes a log line at the given level. Level encoding matches
// engine's capabilities.LogLevel (0=debug, 1=info, 2=warn, 3=error).
//
//go:wasmimport apiqube host_log
func hostLogImport(level uint32, ptr uint32, length uint32)

// host_now returns the current Unix milliseconds.
//
//go:wasmimport apiqube host_now
func hostNowImport() uint64

// host_emit_event delivers a JSON-encoded PluginEvent to the engine's
// event handler.
//
//go:wasmimport apiqube host_emit_event
func hostEmitEventImport(ptr uint32, length uint32)

// host_http_request performs an HTTP request via the engine's http.Client.
// Returns a packed (ptr<<32)|len pointing to a JSON-encoded HostHTTPResponse.
//
//go:wasmimport apiqube host_http_request
func hostHTTPRequestImport(ptr uint32, length uint32) uint64
