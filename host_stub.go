//go:build !wasip1

package main

// Pure-Go stubs for the WASM host imports — used when compiling for unit
// tests, vet, and lint on a regular Go target. Under TinyGo's wasi target
// the real imports in host_wasm.go take over.

func hostLogImport(level, ptr, length uint32) {
	_ = level
	_ = ptr
	_ = length
}

func hostNowImport() uint64 { return 0 }

func hostEmitEventImport(ptr, length uint32) {
	_ = ptr
	_ = length
}

func hostHTTPRequestImport(ptr, length uint32) uint64 {
	_ = ptr
	_ = length
	return 0
}
