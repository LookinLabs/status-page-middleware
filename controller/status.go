package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lookinlabs/status-page-middleware/config"
	"github.com/lookinlabs/status-page-middleware/pkg/checks"
)

func ServiceStatuses(cfg *config.Environments, ctx *gin.Context) {
	services, err := config.LoadConfig(cfg.StatusPageConfigPath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error loading config",
		})

		return
	}

	for i := range services {
		switch services[i].Type {
		case "http":
			services[i].Status = checks.HTTP(services[i].URL)
		case "dns":
			services[i].Status = checks.DNS(services[i].URL)
		case "tcp":
			services[i].Status = checks.TCP(services[i].URL)
		default:
			services[i].Status = "unknown"
		}
	}

	ctx.HTML(http.StatusOK, "status.html", gin.H{
		"services": services,
	})
}

func Ping(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
