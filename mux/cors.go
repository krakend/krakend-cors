package mux

import (
	krakendcors "github.com/devopsfaith/krakend-cors"
	"github.com/devopsfaith/krakend/config"
	"github.com/devopsfaith/krakend/router/mux"
	"github.com/rs/cors"
)

// New returns a mux.HandlerMiddleware (wich implements the http.Handler interface)
// with the CORS configuration defined in the ExtraConfig.
func New(e config.ExtraConfig) mux.HandlerMiddleware {
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
			})
		}
	}
	return nil
}
