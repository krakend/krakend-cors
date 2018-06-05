package cors

import (
	"testing"
	"time"
)

func TestConfigGetter(t *testing.T) {
	sampleCfg := map[string]interface{}{
		Namespace: map[string]interface{}{
			"allow_origins":     []string{"http://localhost", "http://www.example.com"},
			"allow_headers":     []string{"X-Test", "Content-Type"},
			"allow_methods":     []string{"POST", "GET"},
			"expose_headers":    []string{"Content-Type"},
			"allow_credentials": false,
			"max_age":           "24h",
		},
	}
	testCfg := ConfigGetter(sampleCfg).(Config)
	if len(testCfg.AllowOrigins) != 2 {
		t.Error("Should have exactly 2 allowed origins.\n")
	}
	for i, v := range []string{"http://localhost", "http://www.example.com"} {
		if testCfg.AllowOrigins[i] != v {
			t.Errorf("Invalid value %s should be %s\n", testCfg.AllowOrigins[i], v)
		}
	}
	if len(testCfg.AllowHeaders) != 2 {
		t.Error("Should have exactly 2 allowed headers.\n")
	}
	for i, v := range []string{"X-Test", "Content-Type"} {
		if testCfg.AllowHeaders[i] != v {
			t.Errorf("Invalid value %s should be %s\n", testCfg.AllowHeaders[i], v)
		}
	}
	if len(testCfg.AllowMethods) != 2 {
		t.Error("Should have exactly 2 allowed headers.\n")
	}
	for i, v := range []string{"POST", "GET"} {
		if testCfg.AllowMethods[i] != v {
			t.Errorf("Invalid value %s should be %s", testCfg.AllowMethods[i], v)
		}
	}
	if len(testCfg.ExposeHeaders) != 1 {
		t.Error("Should have exactly 2 allowed headers.\n")
	}
	for i, v := range []string{"Content-Type"} {
		if testCfg.ExposeHeaders[i] != v {
			t.Errorf("Invalid value %s should be %s", testCfg.ExposeHeaders[i], v)
		}
	}
	if testCfg.AllowCredentials {
		t.Error("Allow Credentials should be disabled.\n")
	}

	if testCfg.MaxAge != 24*time.Hour {
		t.Errorf("Unexpected collection time: %v\n", testCfg.MaxAge)
	}
}

func TestDefaultConfiguration(t *testing.T) {
	sampleCfg := map[string]interface{}{
		Namespace: map[string]interface{}{
			"allow_origins": []string{"http://www.example.com"},
		},
	}
	defaultCfg := ConfigGetter(sampleCfg).(Config)
	if defaultCfg.AllowOrigins[0] != "http://www.example.com" {
		t.Error("Wrong AllowOrigin.\n")
	}
}

func TestWrongOrEmptyConfiguration(t *testing.T) {
	sampleCfg := map[string]interface{}{}
	if _, ok := ConfigGetter(sampleCfg).(Config); ok {
		t.Error("The config should be nil\n")
	}
	badCfg := map[string]interface{}{Namespace: "test"}
	if _, ok := ConfigGetter(badCfg).(Config); ok {
		t.Error("The config should be nil\n")
	}
	noOriginCfg := map[string]interface{}{
		Namespace: map[string]interface{}{
			"allow_origin":  "",
			"allow_headers": []string{"Content-Type"},
		},
	}
	if v, ok := ConfigGetter(noOriginCfg).(Config); ok {
		t.Errorf("The configuration should be nil, the Origin must not be empty: %v", v)
	}
}
