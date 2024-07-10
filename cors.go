package cors

import (
	"time"

	"github.com/luraproject/lura/v2/config"
)

// Namespace is the key to look for extra configuration details
const Namespace = "github_com/devopsfaith/krakend-cors"

// Config holds the configuration of CORS
type Config struct {
	AllowOrigins         []string
	AllowMethods         []string
	AllowHeaders         []string
	ExposeHeaders        []string
	AllowCredentials     bool
	AllowPrivateNetwork  bool
	OptionsPassthrough   bool
	OptionsSuccessStatus int
	MaxAge               time.Duration
	Debug                bool
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
	cfg.AllowOrigins = getList(tmp, "allow_origins")
	cfg.AllowMethods = getList(tmp, "allow_methods")
	cfg.AllowHeaders = getList(tmp, "allow_headers")
	cfg.ExposeHeaders = getList(tmp, "expose_headers")

	if allowCredentials, ok := tmp["allow_credentials"]; ok {
		if v, ok := allowCredentials.(bool); ok {
			cfg.AllowCredentials = v
		}
	}

	if debug, ok := tmp["debug"]; ok {
		v, ok := debug.(bool)
		cfg.Debug = ok && v
	}

	if allowPrivateNetwork, ok := tmp["allow_private_network"]; ok {
		v, ok := allowPrivateNetwork.(bool)
		cfg.AllowPrivateNetwork = ok && v
	}

	if optionsPassthrough, ok := tmp["options_passthrough"]; ok {
		v, ok := optionsPassthrough.(bool)
		cfg.OptionsPassthrough = ok && v
	}

	if optionsSuccessStatus, ok := tmp["options_success_status"]; ok {
		if v, ok := optionsSuccessStatus.(float64); ok {
			cfg.OptionsSuccessStatus = int(v)
		}
	}

	if maxAge, ok := tmp["max_age"]; ok {
		if d, err := time.ParseDuration(maxAge.(string)); err == nil {
			cfg.MaxAge = d
		}
	}
	return cfg
}

func getList(data map[string]interface{}, name string) []string {
	var out []string
	if vs, ok := data[name]; ok {
		if v, ok := vs.([]interface{}); ok {
			for _, s := range v {
				if j, ok := s.(string); ok {
					out = append(out, j)
				}
			}
		}
	}
	return out
}
