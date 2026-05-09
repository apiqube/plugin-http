package main

import (
	"encoding/json"
)

// WASM-exported entry points. Each takes a packed (ptr<<32)|len uint64
// pointing into the plugin's linear memory and returns the same form.
//
// The //export directive is honored by TinyGo's WASI target. Under a regular
// `go build` it is treated as a comment, so this file compiles for unit
// tests and linting without producing exports.

// plugin_info returns the plugin metadata as JSON bytes.
//
//export plugin_info
func pluginInfoExport() uint64 {
	data, err := json.Marshal(Info())
	if err != nil {
		return 0
	}
	return writeBytes(data)
}

// plugin_init runs one-time initialization with the host-supplied config.
// For v1.0 this is a no-op; future versions may consume timeout/retry defaults.
//
//export plugin_init
func pluginInitExport(_ uint64) uint64 {
	return 0
}

// validate checks a TestInput without performing the request.
//
//export validate
func validateExport(inputPtr uint64) uint64 {
	input, ok := decodeTestInput(inputPtr)
	if !ok {
		return writeJSON([]FieldError{{Message: "invalid input JSON"}})
	}
	errs := validateInput(input)
	if len(errs) == 0 {
		return 0
	}
	return writeJSON(errs)
}

// execute performs one HTTP request and returns the TestOutput.
//
//export execute
func executeExport(inputPtr uint64) uint64 {
	input, ok := decodeTestInput(inputPtr)
	if !ok {
		return writeJSON(TestOutput{Error: "invalid input JSON"})
	}

	req, err := BuildHostRequest(input)
	if err != nil {
		return writeJSON(TestOutput{Error: err.Error()})
	}

	resp, err := callHostHTTPRequest(req)
	if err != nil {
		return writeJSON(TestOutput{Error: err.Error()})
	}

	return writeJSON(BuildOutput(resp))
}

// plugin_destroy releases pinned allocations.
//
//export plugin_destroy
func pluginDestroyExport() {
	resetAllocations()
}

// validateInput applies the v1.0 rules: target+resource OR fields.url must
// produce a URL, and method must be present (defaulting to GET if empty is
// also acceptable, but we surface absence as a hint for stricter callers).
func validateInput(input TestInput) []FieldError {
	var errs []FieldError
	if _, err := composeURL(input); err != nil {
		errs = append(errs, FieldError{Field: "target", Message: err.Error()})
	}
	return errs
}

// decodeTestInput reads packed bytes and JSON-decodes into a TestInput.
func decodeTestInput(packed uint64) (TestInput, bool) {
	raw := readBytes(packed)
	if len(raw) == 0 {
		return TestInput{}, true
	}
	var input TestInput
	if err := json.Unmarshal(raw, &input); err != nil {
		return TestInput{}, false
	}
	return input, true
}

// writeJSON marshals v and returns its packed pointer. Marshal failure
// returns 0 (host treats as empty result).
func writeJSON(v any) uint64 {
	data, err := json.Marshal(v)
	if err != nil {
		return 0
	}
	return writeBytes(data)
}
