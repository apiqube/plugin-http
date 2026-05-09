// Package main is the entry point for the HTTP executor plugin.
//
// This is a WebAssembly plugin compiled from Go via TinyGo (target=wasi).
// Engine loads it via wazero and calls the exports declared in this package:
//
//	plugin_info()    — returns plugin metadata + field/event definitions
//	plugin_init()    — one-time initialization (no-op for v1.0)
//	validate()       — checks a TestInput is shaped correctly, no I/O
//	execute()        — performs the HTTP request, returns the result
//	plugin_destroy() — releases pinned allocations
//	allocate()       — host-callable memory reservation
//
// The plugin imports four host functions from the "apiqube" namespace:
// host_log, host_now, host_emit_event, host_http_request.
//
// Build:
//
//	tinygo build -o plugin-http.wasm -target=wasi ./
package main

// main is required by the WASI command runtime. TinyGo's WASI target lets
// exported functions be called after main returns, so this empty body is
// sufficient for our reactor-style plugin pattern.
func main() {}
