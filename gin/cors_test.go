package gin

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/luraproject/lura/v2/config"
	"github.com/luraproject/lura/v2/logging"
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
	serialized := []byte(`{ "github_com/devopsfaith/krakend-cors": {
			"allow_origins": [ "http://foobar.com" ],
			"allow_methods": [ "GET" ],
			"max_age": "2h"
			}
		}`)
	json.Unmarshal(serialized, &sampleCfg)
	gin.SetMode(gin.TestMode)
	e := gin.New()
	corsMw := New(sampleCfg)
	if corsMw == nil {
		t.Error("The cors middleware should not be nil.\n")
	}
	e.Use(corsMw)
	e.GET("/foo", func(c *gin.Context) { c.String(200, "Yeah") })
	res := httptest.NewRecorder()
	req, _ := http.NewRequest("OPTIONS", "https://example.com/foo", http.NoBody)
	req.Header.Add("Origin", "http://foobar.com")
	req.Header.Add("Access-Control-Request-Method", "GET")
	req.Header.Add("Access-Control-Request-Headers", "origin")
	e.ServeHTTP(res, req)
	if res.Code != 200 {
		t.Errorf("Invalid status code: %d should be 200", res.Code)
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
	serialized := []byte(`{ "github_com/devopsfaith/krakend-cors": {
			"allow_origins": [ "*" ]
			}
		}`)
	json.Unmarshal(serialized, &sampleCfg)
	gin.SetMode(gin.TestMode)
	e := gin.New()
	corsMw := New(sampleCfg)
	if corsMw == nil {
		t.Error("The cors middleware should not be nil.\n")
	}
	e.Use(corsMw)
	e.GET("/wildcard", func(c *gin.Context) { c.String(200, "Yeah") })
	res := httptest.NewRecorder()
	req, _ := http.NewRequest("OPTIONS", "https://example.com/wildcard", http.NoBody)
	req.Header.Add("Origin", "http://foobar.com")
	req.Header.Add("Access-Control-Request-Method", "GET")
	req.Header.Add("Access-Control-Request-Headers", "origin")
	e.ServeHTTP(res, req)
	if res.Code != 200 {
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
	sampleCfg := map[string]interface{}{}
	serialized := []byte(`{ "github_com/devopsfaith/krakend-cors": {
			}
		}`)
	json.Unmarshal(serialized, &sampleCfg)
	gin.SetMode(gin.TestMode)
	e := gin.New()
	corsMw := New(sampleCfg)
	if corsMw == nil {
		t.Error("The cors middleware should not be nil.\n")
	}
	e.Use(corsMw)
	e.GET("/foo", func(c *gin.Context) { c.String(200, "Yeah") })
	res := httptest.NewRecorder()
	req, _ := http.NewRequest("OPTIONS", "https://example.com/foo", http.NoBody)
	req.Header.Add("Origin", "http://foobar.com")
	req.Header.Add("Access-Control-Request-Method", "GET")
	req.Header.Add("Access-Control-Request-Headers", "origin")
	e.ServeHTTP(res, req)
	if res.Code != 200 {
		t.Errorf("Invalid status code: %d should be 200", res.Code)
	}

	assertHeaders(t, res.Header(), map[string]string{
		"Vary":                         "Origin, Access-Control-Request-Method, Access-Control-Request-Headers",
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "GET",
		"Access-Control-Allow-Headers": "origin",
	})
}

func ExampleNewRunServerWithLogger() {
	var localHandler http.Handler
	next := func(_ context.Context, _ config.ServiceConfig, handler http.Handler) error {
		localHandler = handler
		return nil
	}

	buf := bytes.NewBuffer(nil)
	l, _ := logging.NewLogger("DEBUG", buf, "")
	corsRunServer := NewRunServerWithLogger(next, l)

	sampleCfg := map[string]interface{}{}
	serialized := []byte(`{ "github_com/devopsfaith/krakend-cors": {
			"allow_origins": [ "http://foobar.com" ],
			"allow_methods": [ "GET" ],
			"max_age": "2h",
			"debug": true
			}
		}`)
	json.Unmarshal(serialized, &sampleCfg)
	cfg := config.ServiceConfig{ExtraConfig: sampleCfg}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("Yeah"))
	})

	if err := corsRunServer(context.Background(), cfg, mux); err != nil {
		fmt.Println(err)
		return
	}

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("OPTIONS", "http://example.com/", nil)
	req.Header.Add("Origin", "http://foobar.com")
	req.Header.Add("Access-Control-Request-Method", "GET")
	req.Header.Add("Access-Control-Request-Headers", "origin")
	localHandler.ServeHTTP(res, req)
	if res.Code != 200 {
		fmt.Printf("Invalid status code: %d should be 200", res.Code)
		return
	}
	fmt.Println(res.Code)

	b, _ := json.MarshalIndent(res.Header(), "", "\t")
	fmt.Println(string(b))

	fmt.Println("'" + res.Body.String() + "'")

	res = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "http://example.com/", nil)
	req.Header.Add("Origin", "http://foobar.com")
	req.Header.Add("Access-Control-Request-Method", "GET")
	req.Header.Add("Access-Control-Request-Headers", "origin")
	localHandler.ServeHTTP(res, req)
	if res.Code != 200 {
		fmt.Printf("Invalid status code: %d should be 200", res.Code)
		return
	}
	fmt.Println(res.Code)

	b, _ = json.MarshalIndent(res.Header(), "", "\t")
	fmt.Println(string(b))

	fmt.Println("'" + res.Body.String() + "'")

	re := regexp.MustCompile(`(\d\d\d\d\/\d\d\/\d\d \d\d:\d\d:\d\d\s+)`)
	fmt.Println(re.ReplaceAllString(buf.String(), ""))

	// output:
	// 200
	// {
	// 	"Access-Control-Allow-Headers": [
	// 		"Origin"
	// 	],
	// 	"Access-Control-Allow-Methods": [
	// 		"GET"
	// 	],
	// 	"Access-Control-Allow-Origin": [
	// 		"http://foobar.com"
	// 	],
	// 	"Access-Control-Max-Age": [
	// 		"7200"
	// 	],
	// 	"Vary": [
	// 		"Origin",
	// 		"Access-Control-Request-Method",
	// 		"Access-Control-Request-Headers"
	// 	]
	// }
	// ''
	// 200
	// {
	// 	"Access-Control-Allow-Origin": [
	// 		"http://foobar.com"
	// 	],
	// 	"Content-Type": [
	// 		"text/plain; charset=utf-8"
	// 	],
	// 	"Vary": [
	// 		"Origin"
	// 	]
	// }
	// 'Yeah'
	// DEBUG: [SERVICE: Gin][CORS] Enabled CORS for all requests
	// DEBUG: [CORS] Handler: Preflight request
	// DEBUG: [CORS] Preflight response headers: map[Access-Control-Allow-Headers:[Origin] Access-Control-Allow-Methods:[GET] Access-Control-Allow-Origin:[http://foobar.com] Access-Control-Max-Age:[7200] Vary:[Origin Access-Control-Request-Method Access-Control-Request-Headers]]
	// DEBUG: [CORS] Handler: Actual request
	// DEBUG: [CORS] Actual response added headers: map[Access-Control-Allow-Origin:[http://foobar.com] Vary:[Origin]]

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
