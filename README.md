[![Go Report Card](https://goreportcard.com/badge/github.com/devopsfaith/krakend-cors)](https://goreportcard.com/report/github.com/devopsfaith/krakend-cors)  [![GoDoc](https://godoc.org/github.com/devopsfaith/krakend-cors?status.svg)](https://godoc.org/github.com/devopsfaith/krakend-cors)

KrakenD CORS
====

A set of building blocks for instrumenting [KrakenD](http://www.krakend.io) gateways

## Available flavours

1. [mux](github.com/devopsfaith/krakend-cors/blob/master/mux) Mux based handlers
2. [gin](github.com/devopsfaith/krakend-cors/blob/master/gin) Gin based handlers

Check the tests and the documentation for more details

## Configuration

You need to add an ExtraConfig section to the configuration to enable the CORS middleware.
At least one option should be defined.

- `allow_origins` list of strings (you can also use a wildcard, leaving it empty allows all origins too)
- `allow_headers` list of strings
- `allow_methods` list of strings
- `expose_headers` list of strings
- `allow_credentials` bool
- `max_age` duration (Ex: "12h", "5m", "3600s", ...)

### Configuration Example

```
  "extra_config": {
    "github_com/devopsfaith/krakend-cors": {
      "allow_origins": [ "http://foobar.com" ],
      "allow_methods": [ "POST", "GET"],
      "max_age": "12h"
    }
  }
  ```

  or leave the defaults (the defaults allows all origins):
  ```
  "extra_config": {
    github_com/devopsfaith/krakend-metrics": {
      "allow_origins": []
    }
  }
  ```
