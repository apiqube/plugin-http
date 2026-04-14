package main

// Exported WASM functions. Each takes a uint64 (packed ptr+len to JSON bytes)
// and returns a uint64 with the result bytes.
//
// The //export and //go:wasmimport directives only take effect when building
// with TinyGo (-target=wasi). Under a regular `go build` these are ignored,
// so the file still compiles cleanly for tests and linting.

import "encoding/json"

// plugin_info returns the PluginInfo marshaled as JSON.
//
//export plugin_info
func pluginInfo() uint64 {
	data, _ := json.Marshal(info())
	return toPointer(data)
}

// plugin_init runs one-time initialization.
//
//export plugin_init
func pluginInit(cfgPtr uint64) uint64 {
	// TODO: implementation
	// 1. readBytes(cfgPtr) → config map
	// 2. Apply plugin-level settings
	// 3. Return 0 on success or packed error bytes
	_ = cfgPtr
	return 0
}

// validate checks a TestInput for required fields.
//
//export validate
func validateInput(inputPtr uint64) uint64 {
	// TODO: implementation
	// 1. readBytes(inputPtr) → TestInput
	// 2. Check method required, valid value
	// 3. Check endpoint XOR url
	// 4. Return JSON-packed []FieldError
	_ = inputPtr
	return 0
}

// execute runs a single test case.
//
//export execute
func executeTest(inputPtr uint64) uint64 {
	// TODO: implementation
	// 1. readBytes(inputPtr) → TestInput
	// 2. Build URL from target + endpoint/url
	// 3. Call httpRequest() which calls the host function
	// 4. Marshal TestOutput, return toPointer(data)
	_ = inputPtr
	return 0
}

// plugin_destroy releases plugin resources.
//
//export plugin_destroy
func pluginDestroy() {
	// TODO: implementation (typically noop for HTTP)
}

// toPointer packs a byte slice into a (ptr << 32) | len uint64.
// The host dereferences this to read the JSON result.
func toPointer(data []byte) uint64 {
	// TODO: when building for WASM, use unsafe.Pointer to get the raw pointer
	// Under regular go build this is a no-op for compilation testing.
	return uint64(len(data))
}

// readBytes is the inverse of toPointer — reads bytes at the given packed location.
func readBytes(packed uint64) []byte {
	// TODO: WASM implementation via unsafe
	_ = packed
	return nil
}
