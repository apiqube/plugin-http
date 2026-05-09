# plugin-http

> HTTP executor plugin for [ApiQube](https://github.com/apiqube).

[![License](https://img.shields.io/badge/License-MIT-blue?style=flat-square)](LICENSE)
[![Status](https://img.shields.io/badge/Status-Active%20Development-brightgreen?style=flat-square)]()

This is the first-party WASM plugin for the ApiQube testing engine. Engine routes any test whose target starts with `http://` or `https://` to this plugin via the `host` capability.

## Protocols

`http`, `https`

## Capabilities required

`http` (provided by the engine via `host_http_request`).

## Manifest fields

The HTTP method and path live on the core `TestInput.Method` / `TestInput.Resource` fields, not in `fields:` — the same as for any other ApiQube plugin. The fields below are HTTP-specific extensions:

| Field             | Type   | Description |
|-------------------|--------|-------------|
| `body`            | any    | Request body. Marshaled to JSON unless Content-Type overrides. |
| `query`           | map    | URL query parameters as key-value pairs. |
| `url`             | string | Absolute URL override; takes precedence over target+resource. |
| `followRedirects` | bool   | Whether to follow HTTP redirects (default true). |
| `maxRedirects`    | number | Maximum number of redirects to follow (default 10). |

## Example

```yaml
target: http://localhost:8081

tests:
  - name: Create user
    method: POST
    resource: /users
    body:
      name: "{{ fake.name }}"
      email: "{{ fake.email }}"
    expect:
      status: 201

  - name: Search active users
    method: GET
    resource: /users
    query:
      status: active
      limit: 50
    expect:
      status: 200
      body.length: "> 0"
```

## Build

```
tinygo build -o plugin-http.wasm -target=wasi ./
```

CI builds and uploads the artifact on every push.

## Test

```
go test ./...
```

The Go-side tests cover URL composition, header handling, body encoding, the contract types, and the WASM exports through their JSON wire interface. End-to-end tests run inside the engine module against a CI-built `plugin-http.wasm`.

## Install

```
qube plugin install plugin-http.wasm
```

## License

[MIT](LICENSE)
