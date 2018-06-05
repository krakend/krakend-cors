package cors

import (
	"time"

	"github.com/devopsfaith/krakend/config"
)

// Namespace is the key to look for extra configuration details
const Namespace = "github_com/devopsfaith/krakend-cors"

// Config holds the configuration of CORS
type Config struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	ExposeHeaders    []string
	AllowCredentials bool
	MaxAge           time.Duration
}

// ConfigGetter implements the config.ConfigGetter interface. It parses the extra config an allowed
// origin must be defined, the rest of the options will use a default if not defined.
func ConfigGetter(e config.ExtraConfig) interface{} {
	v, ok := e[Namespace]
	if !ok {
		return nil
	}

	tmp, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}

	cfg := Config{}

	if allowOrigins, ok := tmp["allow_origins"]; ok {
		if v, ok := allowOrigins.([]string); ok {
			cfg.AllowOrigins = v
		}
	} else {
		return nil
	}

	if allowMethods, ok := tmp["allow_methods"]; ok {
		if v, ok := allowMethods.([]string); ok {
			cfg.AllowMethods = v
		}
	}

	if allowHeaders, ok := tmp["allow_headers"]; ok {
		if v, ok := allowHeaders.([]string); ok {
			cfg.AllowHeaders = v
		}
	}

	if exposeHeaders, ok := tmp["expose_headers"]; ok {
		if v, ok := exposeHeaders.([]string); ok {
			cfg.ExposeHeaders = v
		}
	}

	//cfg.AllowCredentials = true
	if allowCredentials, ok := tmp["allow_credentials"]; ok {
		if v, ok := allowCredentials.(bool); ok {
			cfg.AllowCredentials = v
		}
	}

	//cfg.MaxAge = 12 * time.Hour
	if maxAge, ok := tmp["max_age"]; ok {
		if d, err := time.ParseDuration(maxAge.(string)); err == nil {
			cfg.MaxAge = d
		}
	}
	return cfg
}
