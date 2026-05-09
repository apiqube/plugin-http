package main

// Wire-format types for the host-plugin contract.
//
// These mirror engine's internal/wire and internal/plugin types. They are
// duplicated here (not imported) because plugins are compiled to WASM and
// communicate with the host only through JSON bytes — there is no shared Go
// runtime between plugin and host. The JSON tags are the contract; the Go
// types are local conveniences.

// PluginInfo is the metadata returned by plugin_info.
//
// Capabilities lists the host capabilities this plugin requires (e.g. "http").
// At load time the host validates each entry against its supported set; an
// unsupported capability fails the plugin load with a clear error.
type PluginInfo struct {
	Name         string               `json:"name"`
	Version      string               `json:"version"`
	Description  string               `json:"description"`
	Protocols    []string             `json:"protocols"`
	Capabilities []string             `json:"capabilities,omitempty"`
	Fields       map[string]FieldSpec `json:"fields,omitempty"`
	Events       map[string]EventSpec `json:"events,omitempty"`
}

// FieldSpec declares one manifest field this plugin reads, or one field of
// an event payload this plugin emits.
type FieldSpec struct {
	Type        string `json:"type"`
	Required    bool   `json:"required,omitempty"`
	Description string `json:"description,omitempty"`
}

// EventSpec declares one event this plugin can emit at runtime. Frontends
// read these schemas to know what events exist and what fields they carry.
type EventSpec struct {
	Description string               `json:"description"`
	Fields      map[string]FieldSpec `json:"fields,omitempty"`
}

// TestInput is what the host passes to execute().
//
// Method and Resource are core fields readable by every plugin (HTTP
// method/path, gRPC method-name, etc.). Plugin-specific data lives in Fields.
type TestInput struct {
	Method   string            `json:"method,omitempty"`
	Resource string            `json:"resource,omitempty"`
	Target   string            `json:"target"`
	Headers  map[string]string `json:"headers,omitempty"`
	Timeout  string            `json:"timeout,omitempty"`
	Fields   map[string]any    `json:"fields,omitempty"`
}

// TestOutput is what execute() returns.
//
// Events carries plugin events accumulated during a single Execute call.
// Streaming plugins additionally emit via host_emit_event during the call;
// both paths feed the engine's event handlers.
type TestOutput struct {
	Status     any               `json:"status"`
	Headers    map[string]string `json:"headers,omitempty"`
	Body       any               `json:"body,omitempty"`
	DurationMs int64             `json:"duration_ms"`
	Error      string            `json:"error,omitempty"`
	Metadata   map[string]any    `json:"metadata,omitempty"`
	Events     []PluginEvent     `json:"events,omitempty"`
}

// FieldError reports a validation problem on a specific manifest field.
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// PluginEvent is the wire shape for events a plugin emits — either via the
// host_emit_event host import (streaming) or batched in TestOutput.Events.
type PluginEvent struct {
	Plugin string         `json:"plugin"`
	Kind   string         `json:"kind"`
	Data   map[string]any `json:"data,omitempty"`
}
