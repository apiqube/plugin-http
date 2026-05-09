package main

// Info returns the PluginInfo for plugin-http.
//
// This is the metadata the host reads via plugin_info. It declares:
//   - the protocols this plugin handles (http, https)
//   - the host capabilities required (http)
//   - the manifest fields user tests can set
//   - the events this plugin can emit during execute
//
// method and resource are NOT declared here — they are core fields on every
// TestInput. Plugin-http reads them via input.Method and input.Resource.
func Info() PluginInfo {
	return PluginInfo{
		Name:         "http",
		Version:      "0.1.0",
		Description:  "HTTP executor for ApiQube — sends requests, collects responses.",
		Protocols:    []string{"http", "https"},
		Capabilities: []string{"http"},
		Fields: map[string]FieldSpec{
			"body": {
				Type:        "any",
				Description: "Request body. Marshaled to JSON unless Content-Type overrides.",
			},
			"query": {
				Type:        "map",
				Description: "URL query parameters as key-value pairs.",
			},
			"url": {
				Type:        "string",
				Description: "Absolute URL override; takes precedence over target+resource.",
			},
			"followRedirects": {
				Type:        "bool",
				Description: "Whether to follow HTTP redirects (default true).",
			},
			"maxRedirects": {
				Type:        "number",
				Description: "Maximum number of redirects to follow (default 10).",
			},
		},
		Events: map[string]EventSpec{
			"Sent": {
				Description: "Request was sent to the target.",
				Fields: map[string]FieldSpec{
					"method": {Type: "string"},
					"url":    {Type: "string"},
				},
			},
			"Received": {
				Description: "Response was received from the target.",
				Fields: map[string]FieldSpec{
					"status":      {Type: "number"},
					"duration_ms": {Type: "number"},
				},
			},
		},
	}
}
