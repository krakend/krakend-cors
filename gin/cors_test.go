package gin

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/devopsfaith/krakend-cors"
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
	gin.SetMode(gin.ReleaseMode)
	sampleCfg := map[string]interface{}{
		cors.Namespace: map[string]interface{}{
			"allow_origins": []string{"http://foobar.com"},
			"allow_methods": []string{"POST", "GET"},
			"max_age":       "2h",
		},
	}
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
	e.ServeHTTP(res, req)
	if res.Code != 200 {
		t.Errorf("Invalid status code: %d should be 200", res.Code)
	}

	assertHeaders(t, res.Header(), map[string]string{
		"Access-Control-Allow-Methods": "POST,GET",
		"Access-Control-Max-Age":       "7200",
		"Vary": "Origin, Access-Control-Request-Method, Access-Control-Request-Headers",
		"Access-Control-Allow-Origin": "http://foobar.com",
	},
	)

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
