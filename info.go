package main

// info returns the PluginInfo for this plugin.
// Called from exported plugin_info() and defines all manifest fields
// users can write in their YAML when targeting http:// or https:// URLs.
func info() PluginInfo {
	return PluginInfo{
		Name:        "http",
		Version:     "0.1.0",
		Description: "HTTP executor for ApiQube — sends requests, collects responses.",
		Protocols:   []string{"http", "https"},
		Fields: map[string]FieldSpec{
			"method": {
				Type:        "string",
				Required:    true,
				Description: "HTTP method (GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS).",
			},
			"endpoint": {
				Type:        "string",
				Required:    false,
				Description: "Path relative to the target base URL (e.g. /users/1). Use this OR url.",
			},
			"url": {
				Type:        "string",
				Required:    false,
				Description: "Absolute URL override, ignoring target. Use this OR endpoint.",
			},
			"body": {
				Type:        "any",
				Required:    false,
				Description: "Request body. Marshaled to JSON unless Content-Type says otherwise.",
			},
			"query": {
				Type:        "map",
				Required:    false,
				Description: "URL query parameters as key-value pairs.",
			},
			"follow_redirects": {
				Type:        "bool",
				Required:    false,
				Description: "Whether to follow HTTP redirects (default: true).",
			},
			"max_redirects": {
				Type:        "number",
				Required:    false,
				Description: "Maximum number of redirects to follow (default: 10).",
			},
		},
	}
}
