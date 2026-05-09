package main

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func TestBuildHostRequest_TargetAndResource(t *testing.T) {
	cases := []struct {
		name     string
		target   string
		resource string
		wantURL  string
	}{
		{"plain", "http://api", "/users", "http://api/users"},
		{"trailing target slash", "http://api/", "/users", "http://api/users"},
		{"no leading resource slash", "http://api", "users/1", "http://api/users/1"},
		{"both slashes", "http://api/", "/users/", "http://api/users/"},
		{"empty resource", "http://api", "", "http://api"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req, err := BuildHostRequest(TestInput{Method: "GET", Target: c.target, Resource: c.resource})
			if err != nil {
				t.Fatal(err)
			}
			if req.URL != c.wantURL {
				t.Errorf("URL = %q; want %q", req.URL, c.wantURL)
			}
		})
	}
}

func TestBuildHostRequest_AbsoluteURLOverride(t *testing.T) {
	req, err := BuildHostRequest(TestInput{
		Method: "GET",
		Target: "http://other",
		Fields: map[string]any{"url": "https://override/x"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if req.URL != "https://override/x" {
		t.Errorf("override should win, got %q", req.URL)
	}
}

func TestBuildHostRequest_QueryParams(t *testing.T) {
	req, err := BuildHostRequest(TestInput{
		Method: "GET", Target: "http://api/users",
		Fields: map[string]any{
			"query": map[string]any{"page": 2, "tag": "admin"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// keys are sorted, so page comes before tag.
	if !strings.HasSuffix(req.URL, "?page=2&tag=admin") {
		t.Errorf("query encoding wrong: %q", req.URL)
	}
}

func TestBuildHostRequest_QueryEncodesSpecial(t *testing.T) {
	req, _ := BuildHostRequest(TestInput{
		Method: "GET", Target: "http://api",
		Fields: map[string]any{"query": map[string]any{"q": "a b&c"}},
	})
	if !strings.Contains(req.URL, "q=a%20b%26c") {
		t.Errorf("special chars not escaped: %q", req.URL)
	}
}

func TestBuildHostRequest_MissingURLAndTarget(t *testing.T) {
	_, err := BuildHostRequest(TestInput{Method: "GET"})
	if err == nil {
		t.Error("expected errMissingURL")
	}
}

func TestBuildHostRequest_DefaultMethod(t *testing.T) {
	req, _ := BuildHostRequest(TestInput{Target: "http://api"})
	if req.Method != "GET" {
		t.Errorf("default method should be GET, got %q", req.Method)
	}
}

func TestBuildHostRequest_MethodUppercased(t *testing.T) {
	req, _ := BuildHostRequest(TestInput{Method: "post", Target: "http://api"})
	if req.Method != "POST" {
		t.Errorf("method should be uppercased, got %q", req.Method)
	}
}

func TestBuildHostRequest_BodyJSON(t *testing.T) {
	req, _ := BuildHostRequest(TestInput{
		Method: "POST", Target: "http://api",
		Fields: map[string]any{"body": map[string]any{"name": "alice"}},
	})
	if req.Headers["Content-Type"] != "application/json" {
		t.Errorf("default Content-Type should be JSON, got %q", req.Headers["Content-Type"])
	}
	var got map[string]any
	if err := json.Unmarshal(req.Body, &got); err != nil {
		t.Fatal(err)
	}
	if got["name"] != "alice" {
		t.Errorf("body wrong: %v", got)
	}
}

func TestBuildHostRequest_BodyString(t *testing.T) {
	req, _ := BuildHostRequest(TestInput{
		Method: "POST", Target: "http://api",
		Fields: map[string]any{"body": "raw text"},
	})
	if string(req.Body) != "raw text" {
		t.Errorf("body should pass through verbatim, got %q", req.Body)
	}
	if !strings.HasPrefix(req.Headers["Content-Type"], "text/plain") {
		t.Errorf("string body Content-Type should be text/plain, got %q", req.Headers["Content-Type"])
	}
}

func TestBuildHostRequest_BodyContentTypeOverride(t *testing.T) {
	req, _ := BuildHostRequest(TestInput{
		Method:  "POST", Target: "http://api",
		Headers: map[string]string{"Content-Type": "application/xml"},
		Fields:  map[string]any{"body": "<x/>"},
	})
	if req.Headers["Content-Type"] != "application/xml" {
		t.Errorf("user Content-Type should be preserved, got %q", req.Headers["Content-Type"])
	}
}

func TestBuildHostRequest_NoBody(t *testing.T) {
	req, _ := BuildHostRequest(TestInput{Method: "GET", Target: "http://api"})
	if req.Body != nil {
		t.Errorf("no body should yield nil, got %v", req.Body)
	}
}

func TestBuildHostRequest_HeadersForwarded(t *testing.T) {
	req, _ := BuildHostRequest(TestInput{
		Method:  "GET", Target: "http://api",
		Headers: map[string]string{"Authorization": "Bearer t", "X-Trace": "abc"},
	})
	if req.Headers["Authorization"] != "Bearer t" || req.Headers["X-Trace"] != "abc" {
		t.Errorf("headers not forwarded: %v", req.Headers)
	}
}

func TestBuildHostRequest_FollowRedirectsAndMax(t *testing.T) {
	follow := false
	req, _ := BuildHostRequest(TestInput{
		Method: "GET", Target: "http://api",
		Fields: map[string]any{"followRedirects": follow, "maxRedirects": 3},
	})
	if req.FollowRedirects == nil || *req.FollowRedirects != false {
		t.Errorf("followRedirects wrong: %v", req.FollowRedirects)
	}
	if req.MaxRedirects != 3 {
		t.Errorf("maxRedirects wrong: %d", req.MaxRedirects)
	}
}

func TestBuildHostRequest_TimeoutParsing(t *testing.T) {
	cases := []struct {
		in   string
		want int64
	}{
		{"100ms", 100},
		{"5s", 5000},
		{"2m", 120_000},
		{"1h", 3_600_000},
		{"500", 500}, // bare number assumed ms
		{"", 0},
		{"garbage", 0},
	}
	for _, c := range cases {
		t.Run(c.in, func(t *testing.T) {
			got := parseTimeoutMs(c.in)
			if got != c.want {
				t.Errorf("parseTimeoutMs(%q) = %d; want %d", c.in, got, c.want)
			}
		})
	}
}

func TestBuildOutput_JSONBody(t *testing.T) {
	resp := HostHTTPResponse{
		Status:  200,
		Headers: map[string]string{"Content-Type": "application/json"},
		Body:    []byte(`{"name":"alice","id":42}`),
	}
	out := BuildOutput(resp)
	body, ok := out.Body.(map[string]any)
	if !ok {
		t.Fatalf("body should decode as map, got %T", out.Body)
	}
	if body["name"] != "alice" {
		t.Errorf("decoded body wrong: %v", body)
	}
}

func TestBuildOutput_PlainTextBody(t *testing.T) {
	resp := HostHTTPResponse{
		Status:  200,
		Headers: map[string]string{"Content-Type": "text/plain"},
		Body:    []byte("hello"),
	}
	out := BuildOutput(resp)
	if out.Body != "hello" {
		t.Errorf("plain body should be string, got %v (%T)", out.Body, out.Body)
	}
}

func TestBuildOutput_ErrorPropagation(t *testing.T) {
	resp := HostHTTPResponse{Error: "connection refused", DurationMs: 1}
	out := BuildOutput(resp)
	if out.Error != "connection refused" {
		t.Errorf("error not propagated: %q", out.Error)
	}
}

func TestBuildOutput_EmptyBody(t *testing.T) {
	resp := HostHTTPResponse{Status: 204}
	out := BuildOutput(resp)
	if out.Body != nil {
		t.Errorf("empty body should be nil, got %v", out.Body)
	}
}

func TestLookupHeader_CaseInsensitive(t *testing.T) {
	headers := map[string]string{"content-type": "application/json"}
	if got := lookupHeader(headers, "Content-Type"); got != "application/json" {
		t.Errorf("case-insensitive lookup failed: %q", got)
	}
}

func TestSetHeader_CaseInsensitiveOverwrite(t *testing.T) {
	headers := map[string]string{"content-type": "old"}
	headers = setHeader(headers, "Content-Type", "new")
	// must overwrite the existing case-insensitive key, not add a duplicate.
	if len(headers) != 1 {
		t.Errorf("setHeader should not duplicate keys: %v", headers)
	}
	for _, v := range headers {
		if v != "new" {
			t.Errorf("value not updated: %q", v)
		}
	}
}

func TestPercentEscape(t *testing.T) {
	cases := map[string]string{
		"hello":   "hello",
		"a b":     "a%20b",
		"a&b":     "a%26b",
		"a/b":     "a%2Fb",
		"a-b_c.~": "a-b_c.~",
	}
	for in, want := range cases {
		if got := percentEscape(in); got != want {
			t.Errorf("percentEscape(%q) = %q; want %q", in, got, want)
		}
	}
}

func TestFlattenQuery_NilValueSkipped(t *testing.T) {
	got := flattenQuery(map[string]any{"a": 1, "b": nil, "c": "x"})
	if !reflect.DeepEqual(got, []string{"a=1", "c=x"}) {
		t.Errorf("nil value should be skipped: %v", got)
	}
}

func TestFlattenQuery_NotAMap(t *testing.T) {
	if got := flattenQuery("not a map"); len(got) != 0 {
		t.Errorf("non-map query should yield no pairs, got %v", got)
	}
}
