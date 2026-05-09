package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
)

// HostHTTPRequest is the wire shape sent through host_http_request. It
// mirrors engine's internal/plugin/capabilities.HTTPRequest exactly — a
// drift here breaks the host call.
type HostHTTPRequest struct {
	Method          string            `json:"method"`
	URL             string            `json:"url"`
	Headers         map[string]string `json:"headers,omitempty"`
	Body            []byte            `json:"body,omitempty"`
	TimeoutMs       int64             `json:"timeout_ms,omitempty"`
	FollowRedirects *bool             `json:"follow_redirects,omitempty"`
	MaxRedirects    int               `json:"max_redirects,omitempty"`
}

// HostHTTPResponse is the wire shape returned by host_http_request.
type HostHTTPResponse struct {
	Status     int               `json:"status"`
	Headers    map[string]string `json:"headers,omitempty"`
	Body       []byte            `json:"body,omitempty"`
	DurationMs int64             `json:"duration_ms"`
	Error      string            `json:"error,omitempty"`
}

// errMissingURL signals the input had neither a Target+Resource pair nor a
// fields.url override.
var errMissingURL = errors.New("plugin-http: cannot derive URL — set target+resource or fields.url")

// BuildHostRequest composes a HostHTTPRequest from a TestInput.
//
// Precedence for the URL:
//  1. fields.url (absolute) — used as-is.
//  2. Target + Resource concatenated with a slash boundary normalized.
//
// Query parameters are sourced from fields.query (map[string]any, value
// stringified). Body is sourced from fields.body and JSON-encoded unless a
// Content-Type header overrides — bytes/strings then pass through verbatim.
func BuildHostRequest(input TestInput) (HostHTTPRequest, error) {
	url, err := composeURL(input)
	if err != nil {
		return HostHTTPRequest{}, err
	}

	body, headers, err := composeBody(input)
	if err != nil {
		return HostHTTPRequest{}, err
	}

	return HostHTTPRequest{
		Method:          methodOrDefault(input.Method),
		URL:             url,
		Headers:         headers,
		Body:            body,
		TimeoutMs:       parseTimeoutMs(input.Timeout),
		FollowRedirects: boolField(input.Fields, "followRedirects"),
		MaxRedirects:    intField(input.Fields, "maxRedirects"),
	}, nil
}

// BuildOutput shapes a HostHTTPResponse into a TestOutput. Errors at the host
// layer (e.g. connection refused) come back via resp.Error and surface as
// TestOutput.Error.
func BuildOutput(resp HostHTTPResponse) TestOutput {
	out := TestOutput{
		Status:     resp.Status,
		Headers:    resp.Headers,
		DurationMs: resp.DurationMs,
		Error:      resp.Error,
	}
	if len(resp.Body) > 0 {
		out.Body = decodeBody(resp.Headers, resp.Body)
	}
	return out
}

func methodOrDefault(method string) string {
	if method == "" {
		return "GET"
	}
	return strings.ToUpper(method)
}

// composeURL handles the three URL sources in precedence order.
func composeURL(input TestInput) (string, error) {
	if abs, ok := stringField(input.Fields, "url"); ok && abs != "" {
		return appendQuery(abs, input.Fields), nil
	}

	target := input.Target
	resource := input.Resource
	if target == "" && resource == "" {
		return "", errMissingURL
	}

	url := joinTargetResource(target, resource)
	return appendQuery(url, input.Fields), nil
}

// joinTargetResource concatenates target and resource with exactly one
// boundary slash. Either side may already include trailing/leading slashes.
func joinTargetResource(target, resource string) string {
	if resource == "" {
		return target
	}
	if target == "" {
		return resource
	}
	t := strings.TrimRight(target, "/")
	r := strings.TrimLeft(resource, "/")
	return t + "/" + r
}

func appendQuery(url string, fields map[string]any) string {
	q, ok := fields["query"]
	if !ok {
		return url
	}
	pairs := flattenQuery(q)
	if len(pairs) == 0 {
		return url
	}
	sep := "?"
	if strings.Contains(url, "?") {
		sep = "&"
	}
	return url + sep + strings.Join(pairs, "&")
}

// flattenQuery converts a map[string]any into a stable, encoded "k=v" list.
// Keys are sorted for determinism. Values are stringified; nil entries skip.
func flattenQuery(q any) []string {
	m, ok := q.(map[string]any)
	if !ok {
		return nil
	}
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	out := make([]string, 0, len(keys))
	for _, k := range keys {
		v := m[k]
		if v == nil {
			continue
		}
		out = append(out, escapeQueryKV(k, fmt.Sprint(v)))
	}
	return out
}

// escapeQueryKV percent-encodes the few characters most likely to hurt — a
// minimal pure-Go encoder so we don't pull net/url under TinyGo.
func escapeQueryKV(key, value string) string {
	return percentEscape(key) + "=" + percentEscape(value)
}

func percentEscape(s string) string {
	const upper = "0123456789ABCDEF"
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case 'a' <= c && c <= 'z',
			'A' <= c && c <= 'Z',
			'0' <= c && c <= '9',
			c == '-', c == '_', c == '.', c == '~':
			b.WriteByte(c)
		default:
			b.WriteByte('%')
			b.WriteByte(upper[c>>4])
			b.WriteByte(upper[c&0xF])
		}
	}
	return b.String()
}

// composeBody encodes the body and finalizes Content-Type. Returns the
// effective headers (copy + Content-Type insertion) and body bytes.
func composeBody(input TestInput) ([]byte, map[string]string, error) {
	headers := copyHeaders(input.Headers)
	raw, ok := input.Fields["body"]
	if !ok || raw == nil {
		return nil, headers, nil
	}

	contentType := lookupHeader(headers, "Content-Type")

	switch v := raw.(type) {
	case string:
		if contentType == "" {
			headers = setHeader(headers, "Content-Type", "text/plain; charset=utf-8")
		}
		return []byte(v), headers, nil
	case []byte:
		return v, headers, nil
	}

	// Default: JSON-encode.
	data, err := json.Marshal(raw)
	if err != nil {
		return nil, nil, fmt.Errorf("encode body: %w", err)
	}
	if contentType == "" {
		headers = setHeader(headers, "Content-Type", "application/json")
	}
	return data, headers, nil
}

// decodeBody attempts to JSON-parse the response body when Content-Type says
// JSON; otherwise returns the body as a string. The engine receives the
// resulting any directly into its assertion machinery.
func decodeBody(headers map[string]string, body []byte) any {
	if len(body) == 0 {
		return nil
	}
	if isJSON(headers) {
		var v any
		if err := json.Unmarshal(body, &v); err == nil {
			return v
		}
	}
	return string(body)
}

func isJSON(headers map[string]string) bool {
	ct := lookupHeader(headers, "Content-Type")
	return strings.Contains(strings.ToLower(ct), "json")
}

func copyHeaders(in map[string]string) map[string]string {
	if len(in) == 0 {
		return map[string]string{}
	}
	out := make(map[string]string, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}

// lookupHeader returns the value for a header name, case-insensitively.
func lookupHeader(headers map[string]string, name string) string {
	for k, v := range headers {
		if strings.EqualFold(k, name) {
			return v
		}
	}
	return ""
}

func setHeader(headers map[string]string, name, value string) map[string]string {
	for k := range headers {
		if strings.EqualFold(k, name) {
			headers[k] = value
			return headers
		}
	}
	headers[name] = value
	return headers
}

// stringField, intField, boolField extract typed values from Fields.
func stringField(fields map[string]any, key string) (string, bool) {
	if v, ok := fields[key]; ok {
		if s, ok := v.(string); ok {
			return s, true
		}
	}
	return "", false
}

func boolField(fields map[string]any, key string) *bool {
	v, ok := fields[key]
	if !ok {
		return nil
	}
	if b, ok := v.(bool); ok {
		return &b
	}
	return nil
}

func intField(fields map[string]any, key string) int {
	v, ok := fields[key]
	if !ok {
		return 0
	}
	switch x := v.(type) {
	case int:
		return x
	case int64:
		return int(x)
	case float64:
		return int(x)
	}
	return 0
}

// parseTimeoutMs converts duration strings like "30s" / "500ms" into ms.
// Returns 0 if unset or unparseable. Minimal pure-Go parser to avoid pulling
// time.ParseDuration under TinyGo.
func parseTimeoutMs(s string) int64 {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}

	for _, suffix := range []struct {
		s string
		m int64
	}{
		{"ms", 1},
		{"s", 1000},
		{"m", 60_000},
		{"h", 3_600_000},
	} {
		if strings.HasSuffix(s, suffix.s) {
			num := strings.TrimSuffix(s, suffix.s)
			n, err := parseInt64(num)
			if err != nil {
				return 0
			}
			return n * suffix.m
		}
	}
	n, err := parseInt64(s)
	if err != nil {
		return 0
	}
	return n
}

func parseInt64(s string) (int64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, errors.New("empty number")
	}
	var n int64
	var sign int64 = 1
	for i, c := range s {
		if i == 0 && (c == '-' || c == '+') {
			if c == '-' {
				sign = -1
			}
			continue
		}
		if c < '0' || c > '9' {
			return 0, fmt.Errorf("invalid digit %q", c)
		}
		n = n*10 + int64(c-'0')
	}
	return n * sign, nil
}
