package main

// These types mirror github.com/apiqube/engine/internal/plugin.
// They are duplicated here (not imported) because plugins are compiled
// to WASM independently and communicate with the host via JSON bytes.

// PluginInfo is the metadata returned by plugin_info().
type PluginInfo struct {
	Name        string               `json:"name"`
	Version     string               `json:"version"`
	Description string               `json:"description"`
	Protocols   []string             `json:"protocols"`
	Fields      map[string]FieldSpec `json:"fields"`
}

// FieldSpec declares one manifest field this plugin adds.
type FieldSpec struct {
	Type        string `json:"type"` // string, number, bool, object, array, map, any
	Required    bool   `json:"required"`
	Description string `json:"description"`
}

// TestInput is what the host passes to execute().
type TestInput struct {
	Target  string            `json:"target"`
	Headers map[string]string `json:"headers,omitempty"`
	Timeout string            `json:"timeout,omitempty"`
	Fields  map[string]any    `json:"fields"`
}

// TestOutput is what execute() returns to the host.
type TestOutput struct {
	Status     any               `json:"status"`
	Headers    map[string]string `json:"headers,omitempty"`
	Body       any               `json:"body,omitempty"`
	DurationMs int64             `json:"duration_ms"`
	Error      string            `json:"error,omitempty"`
	Metadata   map[string]any    `json:"metadata,omitempty"`
}

// FieldError reports a validation problem on a specific field.
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}
