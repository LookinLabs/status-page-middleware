package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/lookinlabs/status-page-middleware/config"
	"github.com/lookinlabs/status-page-middleware/controller"
)

func StatusPage(cfg *config.Environments) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.URL.Path == cfg.StatusPagePath {
			controller.ServiceStatuses(cfg, ctx)
			ctx.Abort()
		} else {
			ctx.Next()
		}
	}
}
