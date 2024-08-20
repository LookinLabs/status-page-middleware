package endpoints

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lookinlabs/status-page-middleware/pkg/config"
	"github.com/lookinlabs/status-page-middleware/pkg/status"
)

func StatusPage(cfg *config.Environments) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.URL.Path == cfg.StatusPagePath {
			status.Services(cfg, ctx)
			ctx.Abort()
		} else {
			ctx.Next()
		}
	}
}

func Ping(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
