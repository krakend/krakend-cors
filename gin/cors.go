package gin

import (
	krakendcors "github.com/devopsfaith/krakend-cors"
	"github.com/devopsfaith/krakend/config"
	"github.com/gin-gonic/gin"
	"gopkg.in/gin-contrib/cors.v1"
)

// New returns a gin.HandlerFunc with the CORS configuration provided in the ExtraConfig
func New(e config.ExtraConfig) gin.HandlerFunc {
	tmp := krakendcors.ConfigGetter(e)
	if tmp != nil {
		if cfg, ok := tmp.(cors.Config); ok {
			return cors.New(cfg)
		}
	}
	return nil
}
