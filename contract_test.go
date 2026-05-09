package main

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestPluginInfo_JSONRoundTrip(t *testing.T) {
	in := PluginInfo{
		Name:         "http",
		Version:      "0.1.0",
		Description:  "test",
		Protocols:    []string{"http", "https"},
		Capabilities: []string{"http"},
		Fields: map[string]FieldSpec{
			"body": {Type: "any", Description: "request body"},
		},
		Events: map[string]EventSpec{
			"Sent": {Description: "request sent"},
		},
	}
	data, err := json.Marshal(in)
	if err != nil {
		t.Fatal(err)
	}
	var out PluginInfo
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(in, out) {
		t.Errorf("PluginInfo roundtrip mismatch:\n in: %#v\nout: %#v", in, out)
	}
}

func TestTestInput_JSONShape(t *testing.T) {
	in := TestInput{
		Method:   "POST",
		Resource: "/users/1",
		Target:   "http://api.example.com",
		Headers:  map[string]string{"Authorization": "Bearer t"},
		Fields:   map[string]any{"body": map[string]any{"k": "v"}},
	}
	data, err := json.Marshal(in)
	if err != nil {
		t.Fatal(err)
	}
	// Snake-case JSON wire format check.
	s := string(data)
	for _, want := range []string{`"method":"POST"`, `"resource":"/users/1"`, `"target":"http://api.example.com"`} {
		if !contains(s, want) {
			t.Errorf("expected %q in %s", want, s)
		}
	}
}

func TestTestOutput_JSONShape(t *testing.T) {
	out := TestOutput{
		Status:     201,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       map[string]any{"id": 42},
		DurationMs: 123,
		Events: []PluginEvent{
			{Plugin: "http", Kind: "Sent", Data: map[string]any{"path": "/x"}},
		},
	}
	data, err := json.Marshal(out)
	if err != nil {
		t.Fatal(err)
	}
	s := string(data)
	for _, want := range []string{`"duration_ms":123`, `"events":[`, `"plugin":"http"`} {
		if !contains(s, want) {
			t.Errorf("expected %q in %s", want, s)
		}
	}
}

func TestPluginEvent_RoundTrip(t *testing.T) {
	in := PluginEvent{Plugin: "http", Kind: "Replied", Data: map[string]any{"status": float64(200)}}
	data, err := json.Marshal(in)
	if err != nil {
		t.Fatal(err)
	}
	var out PluginEvent
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatal(err)
	}
	if in.Plugin != out.Plugin || in.Kind != out.Kind {
		t.Errorf("identity fields differ: in=%+v out=%+v", in, out)
	}
}

func TestFieldError_JSON(t *testing.T) {
	in := FieldError{Field: "method", Message: "required"}
	data, _ := json.Marshal(in)
	if string(data) != `{"field":"method","message":"required"}` {
		t.Errorf("unexpected JSON: %s", data)
	}
}

func contains(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
