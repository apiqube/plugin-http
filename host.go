package main

import (
	"encoding/json"
	"errors"
)

// Host imports are split across two build-tagged files:
//
//	host_wasm.go  (build: wasip1) — //go:wasmimport declarations, no bodies
//	host_stub.go  (build: !wasip1) — pure-Go stubs so unit tests compile
//
// The wrapper functions in this file are tag-agnostic and can be called from
// any test or export.

// callHostLog packs a level/message pair into a host_log call.
func callHostLog(level int, message string) {
	if message == "" {
		return
	}
	packed := writeBytes([]byte(message))
	if packed == 0 {
		return
	}
	hostLogImport(uint32(level), uint32(packed>>32), uint32(packed))
}

// callHostNow returns the current host clock in Unix milliseconds.
func callHostNow() int64 {
	return int64(hostNowImport())
}

// EmitEvent emits a plugin event to the host. Failure to encode is silent —
// streaming events are best-effort; missing one is preferable to crashing
// the plugin mid-execute.
func EmitEvent(plugin, kind string, data map[string]any) {
	ev := PluginEvent{Plugin: plugin, Kind: kind, Data: data}
	raw, err := json.Marshal(ev)
	if err != nil {
		return
	}
	packed := writeBytes(raw)
	if packed == 0 {
		return
	}
	hostEmitEventImport(uint32(packed>>32), uint32(packed))
}

// callHostHTTPRequest sends a HostHTTPRequest to the engine and decodes the
// response. The host returns 0 on a write failure; we surface that as an
// explicit error.
func callHostHTTPRequest(req HostHTTPRequest) (HostHTTPResponse, error) {
	raw, err := json.Marshal(req)
	if err != nil {
		return HostHTTPResponse{}, err
	}
	inPacked := writeBytes(raw)
	if inPacked == 0 {
		return HostHTTPResponse{}, errors.New("plugin: failed to allocate request bytes")
	}
	outPacked := hostHTTPRequestImport(uint32(inPacked>>32), uint32(inPacked))
	if outPacked == 0 {
		return HostHTTPResponse{}, errors.New("plugin: host returned no response")
	}
	var resp HostHTTPResponse
	if err := json.Unmarshal(readBytes(outPacked), &resp); err != nil {
		return HostHTTPResponse{}, err
	}
	return resp, nil
}
