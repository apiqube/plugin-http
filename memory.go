package main

// Memory bridge state shared across both build variants. The actual
// allocate/pack/read/write logic lives in memory_wasm.go (TinyGo wasi target)
// and memory_stub.go (regular Go for tests) — they pick a representation
// suited for their address space.

// allocations holds every byte slice we hand to the host so the GC keeps them
// alive until plugin_destroy releases them.
var allocations [][]byte

// resetAllocations releases every pinned buffer. Called from plugin_destroy
// and from tests between runs.
func resetAllocations() {
	allocations = nil
	resetMemoryState()
}
