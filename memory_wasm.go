//go:build wasip1

package main

import "unsafe"

// allocate is the host-callable export used to reserve memory inside the
// plugin's linear address space. The host calls it before writing data the
// plugin must read.
//
//export allocate
func allocate(size uint32) uint64 {
	if size == 0 {
		return 0
	}
	buf := make([]byte, size)
	allocations = append(allocations, buf)
	return packPointer(buf)
}

func packPointer(buf []byte) uint64 {
	if len(buf) == 0 {
		return 0
	}
	ptr := unsafe.Pointer(unsafe.SliceData(buf))
	return (uint64(uintptr(ptr)) << 32) | uint64(len(buf))
}

func readBytes(packed uint64) []byte {
	ptr := uint32(packed >> 32)
	length := uint32(packed)
	if length == 0 {
		return nil
	}
	return unsafe.Slice((*byte)(unsafe.Pointer(uintptr(ptr))), length)
}

func writeBytes(data []byte) uint64 {
	if len(data) == 0 {
		return 0
	}
	buf := make([]byte, len(data))
	copy(buf, data)
	allocations = append(allocations, buf)
	return packPointer(buf)
}

// resetMemoryState is a no-op on WASM — `allocations` already holds the
// only pinning state, and `resetAllocations` clears it directly.
func resetMemoryState() {}
