package main

import "testing"

func TestAllocate(t *testing.T) {
	resetAllocations()
	defer resetAllocations()

	packed := allocate(10)
	if packed == 0 {
		t.Fatal("allocate(10) returned 0")
	}
	if len(allocations) != 1 || len(allocations[0]) != 10 {
		t.Errorf("allocations wrong: %v", allocations)
	}
}

func TestAllocate_Zero(t *testing.T) {
	resetAllocations()
	if got := allocate(0); got != 0 {
		t.Errorf("allocate(0) should return 0, got %d", got)
	}
	if len(allocations) != 0 {
		t.Errorf("zero allocate should not pin a buffer")
	}
}

func TestPackPointer_Empty(t *testing.T) {
	if got := packPointer(nil); got != 0 {
		t.Errorf("packPointer(nil) = %d; want 0", got)
	}
	if got := packPointer([]byte{}); got != 0 {
		t.Errorf("packPointer(empty) = %d; want 0", got)
	}
}

func TestWriteAndReadBytes_RoundTrip(t *testing.T) {
	resetAllocations()
	defer resetAllocations()

	original := []byte("hello world")
	packed := writeBytes(original)
	if packed == 0 {
		t.Fatal("writeBytes returned 0")
	}
	got := readBytes(packed)
	if string(got) != "hello world" {
		t.Errorf("round-trip mismatch: %q", got)
	}
}

func TestWriteBytes_Empty(t *testing.T) {
	if got := writeBytes(nil); got != 0 {
		t.Errorf("writeBytes(nil) should be 0, got %d", got)
	}
}

func TestReadBytes_Zero(t *testing.T) {
	if got := readBytes(0); got != nil {
		t.Errorf("readBytes(0) should be nil, got %v", got)
	}
}

func TestResetAllocations(t *testing.T) {
	_ = writeBytes([]byte("a"))
	_ = writeBytes([]byte("b"))
	resetAllocations()
	if len(allocations) != 0 {
		t.Errorf("after reset, allocations should be empty")
	}
}
