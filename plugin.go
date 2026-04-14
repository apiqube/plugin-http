// Package main is the entry point for the HTTP executor plugin.
//
// This is a WebAssembly plugin compiled from Go via TinyGo.
// It implements the ApiQube plugin contract:
//
//	plugin_info()  → returns plugin metadata + field definitions
//	plugin_init()  → one-time initialization
//	execute()      → run one test case, return result
//	validate()     → check if a test case is valid for this plugin
//	plugin_destroy() → cleanup
//
// Build:
//
//	tinygo build -o plugin-http.wasm -target=wasi ./
package main

// main is required but does nothing — the plugin's entry points are
// individual exported functions called by the host.
func main() {}
