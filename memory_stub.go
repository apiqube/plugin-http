//go:build !wasip1

package main

// Stub memory bridge for the regular-Go (non-WASM) build used by unit tests.
//
// The wasi build packs a real linear-memory pointer into the upper 32 bits
// of `packed`. On a 64-bit host that pointer doesn't fit, so we instead use
// an opaque integer handle and a process-local map. The wire format stays
// the same — the host never sees these handles in the test build.

var (
	stubMemory  = map[uint32][]byte{}
	stubCounter uint32
)

func resetMemoryState() {
	stubMemory = map[uint32][]byte{}
	stubCounter = 0
}

// allocate reserves a fresh slice of the given size and returns its handle.
//
//export allocate
func allocate(size uint32) uint64 {
	if size == 0 {
		return 0
	}
	buf := make([]byte, size)
	allocations = append(allocations, buf)
	return registerBuffer(buf)
}

func packPointer(buf []byte) uint64 {
	if len(buf) == 0 {
		return 0
	}
	return registerBuffer(buf)
}

func readBytes(packed uint64) []byte {
	if packed == 0 {
		return nil
	}
	idx := uint32(packed >> 32)
	return stubMemory[idx]
}

func writeBytes(data []byte) uint64 {
	if len(data) == 0 {
		return 0
	}
	buf := make([]byte, len(data))
	copy(buf, data)
	allocations = append(allocations, buf)
	return registerBuffer(buf)
}

func registerBuffer(buf []byte) uint64 {
	stubCounter++
	idx := stubCounter
	stubMemory[idx] = buf
	return (uint64(idx) << 32) | uint64(len(buf))
}
