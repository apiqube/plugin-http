package main

import (
	"encoding/json"
	"testing"
)

func TestPluginInfoExport(t *testing.T) {
	resetAllocations()
	defer resetAllocations()

	packed := pluginInfoExport()
	if packed == 0 {
		t.Fatal("pluginInfoExport returned 0")
	}
	raw := readBytes(packed)
	var info PluginInfo
	if err := json.Unmarshal(raw, &info); err != nil {
		t.Fatalf("decode info: %v", err)
	}
	if info.Name != "http" {
		t.Errorf("name wrong: %q", info.Name)
	}
	if len(info.Capabilities) != 1 || info.Capabilities[0] != "http" {
		t.Errorf("capabilities wrong: %v", info.Capabilities)
	}
}

func TestPluginInitExport_NoOp(t *testing.T) {
	if got := pluginInitExport(0); got != 0 {
		t.Errorf("plugin_init should return 0, got %d", got)
	}
}

func TestValidateExport_OK(t *testing.T) {
	resetAllocations()
	defer resetAllocations()

	in := TestInput{Method: "GET", Target: "http://api/", Resource: "/x"}
	data, _ := json.Marshal(in)
	packed := writeBytes(data)

	got := validateExport(packed)
	if got != 0 {
		t.Errorf("valid input should return 0, got %d (errors=%s)", got, readBytes(got))
	}
}

func TestValidateExport_MissingURL(t *testing.T) {
	resetAllocations()
	defer resetAllocations()

	in := TestInput{Method: "GET"}
	data, _ := json.Marshal(in)
	packed := writeBytes(data)

	got := validateExport(packed)
	if got == 0 {
		t.Fatal("missing URL should produce errors")
	}
	var errs []FieldError
	if err := json.Unmarshal(readBytes(got), &errs); err != nil {
		t.Fatal(err)
	}
	if len(errs) == 0 || errs[0].Field != "target" {
		t.Errorf("expected target error, got %v", errs)
	}
}

func TestValidateExport_InvalidJSON(t *testing.T) {
	resetAllocations()
	defer resetAllocations()

	packed := writeBytes([]byte("{not json}"))
	got := validateExport(packed)
	if got == 0 {
		t.Fatal("invalid JSON should produce a FieldError list")
	}
	var errs []FieldError
	_ = json.Unmarshal(readBytes(got), &errs)
	if len(errs) != 1 {
		t.Errorf("expected 1 error, got %v", errs)
	}
}

func TestExecuteExport_HostFailureSurfacesError(t *testing.T) {
	// Under unit-test build host_http_request stub returns 0, which
	// callHostHTTPRequest reports as "host returned no response".
	resetAllocations()
	defer resetAllocations()

	in := TestInput{Method: "GET", Target: "http://api"}
	data, _ := json.Marshal(in)
	packed := writeBytes(data)

	got := executeExport(packed)
	if got == 0 {
		t.Fatal("execute should return a TestOutput even on host failure")
	}
	var out TestOutput
	if err := json.Unmarshal(readBytes(got), &out); err != nil {
		t.Fatal(err)
	}
	if out.Error == "" {
		t.Errorf("expected non-empty Error from stub host, got %+v", out)
	}
}

func TestExecuteExport_BadInputJSON(t *testing.T) {
	resetAllocations()
	defer resetAllocations()

	packed := writeBytes([]byte("not json"))
	got := executeExport(packed)
	var out TestOutput
	_ = json.Unmarshal(readBytes(got), &out)
	if out.Error == "" {
		t.Errorf("bad JSON should yield error TestOutput")
	}
}

func TestExecuteExport_BuildFailure(t *testing.T) {
	resetAllocations()
	defer resetAllocations()

	// No target, no resource, no fields.url → BuildHostRequest fails.
	in := TestInput{Method: "GET"}
	data, _ := json.Marshal(in)
	packed := writeBytes(data)

	got := executeExport(packed)
	var out TestOutput
	_ = json.Unmarshal(readBytes(got), &out)
	if out.Error == "" {
		t.Errorf("build failure should produce error TestOutput")
	}
}

func TestPluginDestroyExport_ClearsAllocations(t *testing.T) {
	resetAllocations()
	_ = writeBytes([]byte("x"))
	if len(allocations) == 0 {
		t.Fatal("expected at least one allocation")
	}
	pluginDestroyExport()
	if len(allocations) != 0 {
		t.Errorf("destroy should clear allocations, got %d", len(allocations))
	}
}

func TestDecodeTestInput_Empty(t *testing.T) {
	in, ok := decodeTestInput(0)
	if !ok {
		t.Error("empty input should decode as zero value, ok=true")
	}
	if in.Method != "" {
		t.Errorf("zero value expected, got %+v", in)
	}
}
