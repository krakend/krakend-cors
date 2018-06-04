package mux

import (
	"net/http"

	krakendcors "github.com/devopsfaith/krakend-cors"
	"github.com/devopsfaith/krakend/config"
	"github.com/rs/cors"
)

// New returns a http.Handler with the defined CORS options in the ExtraConfig.
func New(e config.ExtraConfig) http.Handler {
	var cfg *krakendcors.Config
	if tmp, ok := krakendcors.ConfigGetter(e).(*krakendcors.Config); ok {
		cfg = tmp
	}

	if cfg == nil {
		return nil
	}
	return cors.New(cors.Options{
		AllowedOrigins:   cfg.AllowOrigins,
		AllowedMethods:   cfg.AllowMethods,
		AllowedHeaders:   cfg.AllowHeaders,
		ExposedHeaders:   cfg.ExposeHeaders,
		AllowCredentials: cfg.AllowCredentials,
		MaxAge:           int(cfg.MaxAge.Seconds()),
	}).Handler(http.NewServeMux())
}
