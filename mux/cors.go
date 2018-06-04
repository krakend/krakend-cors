package mux

import (
	"net/http"

	krakendcors "github.com/devopsfaith/krakend-cors"
	"github.com/devopsfaith/krakend/config"
	"github.com/rs/cors"
)

// New returns a http.Handler with the defined CORS options in the ExtraConfig.
func New(e config.ExtraConfig) http.Handler {
	tmp := krakendcors.ConfigGetter(e)
	if tmp != nil {
		if cfg, ok := tmp.(krakendcors.Config); ok {
			return cors.New(cors.Options{
				AllowedOrigins:   cfg.AllowOrigins,
				AllowedMethods:   cfg.AllowMethods,
				AllowedHeaders:   cfg.AllowHeaders,
				ExposedHeaders:   cfg.ExposeHeaders,
				AllowCredentials: cfg.AllowCredentials,
				MaxAge:           int(cfg.MaxAge.Seconds()),
			}).Handler(http.NewServeMux())
		}
	}
	return nil
}
