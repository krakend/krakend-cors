package mux

import (
	"io"
	"log"

	krakendcors "github.com/devopsfaith/krakend-cors"
	"github.com/luraproject/lura/config"
	"github.com/luraproject/lura/logging"
	"github.com/luraproject/lura/router/mux"
	"github.com/rs/cors"
)

// New returns a mux.HandlerMiddleware (which implements the http.Handler interface)
// with the CORS configuration defined in the ExtraConfig.
func New(e config.ExtraConfig) mux.HandlerMiddleware {
	return NewWithLogger(e, nil)
}

func NewWithLogger(e config.ExtraConfig, l logging.Logger) mux.HandlerMiddleware {
	tmp := krakendcors.ConfigGetter(e)
	if tmp == nil {
		return nil
	}
	cfg, ok := tmp.(krakendcors.Config)
	if !ok {
		return nil
	}
	c := cors.New(cors.Options{
		AllowedOrigins:   cfg.AllowOrigins,
		AllowedMethods:   cfg.AllowMethods,
		AllowedHeaders:   cfg.AllowHeaders,
		ExposedHeaders:   cfg.ExposeHeaders,
		AllowCredentials: cfg.AllowCredentials,
		MaxAge:           int(cfg.MaxAge.Seconds()),
	})
	if l == nil || !cfg.Debug {
		return c
	}
	r, w := io.Pipe()
	c.Log = log.New(w, "", log.LstdFlags)
	go func() {
		msg := make([]byte, 1024)
		for {
			r.Read(msg)
			l.Debug(string(msg))
		}
	}()
	return c
}
