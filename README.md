KrakenD CORS
====

A set of building blocks for instrumenting [KrakenD](http://www.krakend.io) gateways

## Available flavours

1. [mux](github.com/devopsfaith/krakend-cors/blob/master/mux) Mux based handlers
2. [gin](github.com/devopsfaith/krakend-cors/blob/master/gin) Gin based handlers

Check the tests and the documentation for more details

## Configuration

You need to add an ExtraConfig section to the configuration to enable the CORS middleware.

- `allow_origins` list of strings (must be defined or the middleware will be disabled, you can also use a wildcard)
- `allow_headers` list of strings
- `allow_methods` list of strings
- `expose_headers` list of strings
- `allow_credentials` bool
- `max_age` duration (Ex: "12h", "5m", "3600s", ...)

### Configuration Example

This configuration will set the _collection time_ to 2 minutes and will disable the proxy metrics collector (backend and router metrics will be enabled since the default for all layers is to be enabled).
```
  "extra_config": {
    "github_com/devopsfaith/krakend-cors": {
      "allow_origins": [ "http://foobar.com" ],
      "allow_methods": [ "POST", "GET"],
      "max_age": "12h"
    }
  }
  ```

  or leave the defaults:
  ```
  "extra_config": {
    github_com/devopsfaith/krakend-metrics": {}
  }
  ```
