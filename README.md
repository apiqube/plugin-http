# plugin-http

> HTTP executor plugin for [ApiQube](https://github.com/apiqube).

[![License](https://img.shields.io/badge/License-MIT-blue?style=flat-square)](LICENSE)
[![Status](https://img.shields.io/badge/Status-Scaffold-yellow?style=flat-square)]()

This is a built-in WASM plugin for the ApiQube testing engine. It handles any
test whose target starts with `http://` or `https://`.

## Protocols

`http`, `https`

## Manifest Fields

| Field              | Type   | Required | Description |
|--------------------|--------|----------|-------------|
| `method`           | string | yes      | HTTP method (GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS) |
| `endpoint`         | string | no       | Path relative to target base URL |
| `url`              | string | no       | Absolute URL override |
| `body`             | any    | no       | Request body (marshaled to JSON) |
| `query`            | map    | no       | URL query parameters |
| `follow_redirects` | bool   | no       | Follow HTTP redirects (default: true) |
| `max_redirects`    | number | no       | Max redirects (default: 10) |

## Example

```yaml
target: http://localhost:8081

tests:
  - name: Create user
    method: POST
    endpoint: /users
    body:
      name: "{{ fake.name }}"
      email: "{{ fake.email }}"
    expect:
      status: 201
```

## Build

```bash
tinygo build -o plugin-http.wasm -target=wasi ./
```

## Install

```bash
qube plugin install plugin-http.wasm
```

## License

[MIT](LICENSE)
