package gin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestInvalidCfg(t *testing.T) {
	sampleCfg := map[string]interface{}{}
	corsMw := New(sampleCfg)
	if corsMw != nil {
		t.Error("The corsMw should be nil.\n")
	}
}

func TestNew(t *testing.T) {
	sampleCfg := map[string]interface{}{}
	serialized := []byte(`{ 
        "github_com/devopsfaith/krakend-cors": {
            "allow_origins": ["http://foobar.com", "http://example.com"],
            "allow_headers": ["origin"],
			"allow_methods": ["GET"],
			"max_age": "2h"
			}
		}`)
	err := json.Unmarshal(serialized, &sampleCfg)
	if err != nil {
		t.Errorf("cannot unmarshal sampleCfg: %s", err.Error())
		return
	}
	e := gin.Default()
	corsMw := New(sampleCfg)
	if corsMw == nil {
		t.Error("The cors middleware should not be nil.\n")
		return
	}
	e.Use(corsMw)
	e.GET("/foo", func(c *gin.Context) { c.String(200, "Yeah") })

	res := httptest.NewRecorder()

	req, _ := http.NewRequest("OPTIONS", "http://example.com/foo", nil)
	req.Header.Add("Origin", "http://foobar.com")
	req.Header.Add("Access-Control-Request-Method", "GET")
	req.Header.Add("Access-Control-Request-Headers", "origin")
	e.ServeHTTP(res, req)
	fmt.Printf("METHOD -> %s\n", req.Method)
	if res.Code != 200 && res.Code != 204 {
		t.Errorf("Invalid status code: %d should be 200 or 204", res.Code)
		return
	}

	assertHeaders(t, res.Header(), map[string]string{
		"Vary":                         "Origin, Access-Control-Request-Method, Access-Control-Request-Headers",
		"Access-Control-Allow-Origin":  "http://foobar.com",
		"Access-Control-Allow-Methods": "GET",
		"Access-Control-Allow-Headers": "origin",
		"Access-Control-Max-Age":       "7200",
	})
}

func TestAllowOriginWildcard(t *testing.T) {
	sampleCfg := map[string]interface{}{}
	// WARNING: even if we allow all origins, we still have to specify
	// the allow_headers config
	serialized := []byte(`{ "github_com/devopsfaith/krakend-cors": {
            "allow_origins": [ "*" ],
            "allow_headers": ["origin"]
			}
		}`)
	json.Unmarshal(serialized, &sampleCfg)
	e := gin.Default()
	corsMw := New(sampleCfg)
	if corsMw == nil {
		t.Error("The cors middleware should not be nil.\n")
	}
	e.Use(corsMw)
	e.GET("/wildcard", func(c *gin.Context) { c.String(200, "Yeah") })
	res := httptest.NewRecorder()
	req, _ := http.NewRequest("OPTIONS", "http://example.com/wildcard", nil)
	req.Header.Add("Origin", "http://foobar.com")
	req.Header.Add("Access-Control-Request-Method", "GET")
	req.Header.Add("Access-Control-Request-Headers", "origin")
	e.ServeHTTP(res, req)
	if res.Code != 200 && res.Code != 204 {
		t.Errorf("Invalid status code: %d should be 200", res.Code)
	}

	assertHeaders(t, res.Header(), map[string]string{
		"Vary":                         "Origin, Access-Control-Request-Method, Access-Control-Request-Headers",
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "GET",
		"Access-Control-Allow-Headers": "origin",
	})
}

func TestAllowOriginEmpty(t *testing.T) {
	// WARNING: with an empty config, the library now falls back
	// to "secure" defaults, not allowing the request
	// (in the mux/cors_test.go, we did the reverse, we specified
	// the test to allow everything).
	sampleCfg := map[string]interface{}{}
	serialized := []byte(`{ "github_com/devopsfaith/krakend-cors": {
			}
		}`)
	json.Unmarshal(serialized, &sampleCfg)
	e := gin.Default()
	corsMw := New(sampleCfg)
	if corsMw == nil {
		t.Error("The cors middleware should not be nil.\n")
	}
	e.Use(corsMw)
	e.GET("/foo", func(c *gin.Context) { c.String(200, "Yeah") })
	res := httptest.NewRecorder()
	req, _ := http.NewRequest("OPTIONS", "http://example.com/foo", nil)
	req.Header.Add("Origin", "http://foobar.com")
	req.Header.Add("Access-Control-Request-Method", "GET")
	req.Header.Add("Access-Control-Request-Headers", "origin")
	e.ServeHTTP(res, req)
	if res.Code != 200 && res.Code != 204 {
		t.Errorf("Invalid status code: %d should be 200", res.Code)
	}

	assertHeaders(t, res.Header(), map[string]string{
		"Vary":                         "Origin, Access-Control-Request-Method, Access-Control-Request-Headers",
		"Access-Control-Allow-Origin":  "",
		"Access-Control-Allow-Methods": "",
		"Access-Control-Allow-Headers": "",
	})
}

var allHeaders = []string{
	"Vary",
	"Access-Control-Allow-Origin",
	"Access-Control-Allow-Methods",
	"Access-Control-Allow-Headers",
	"Access-Control-Allow-Credentials",
	"Access-Control-Max-Age",
	"Access-Control-Expose-Headers",
}

func assertHeaders(t *testing.T, resHeaders http.Header, expHeaders map[string]string) {
	for _, name := range allHeaders {
		got := strings.Join(resHeaders[name], ", ")
		want := expHeaders[name]
		if got != want {
			t.Errorf("Response header %q = %q, want %q", name, got, want)
		}
	}
}
