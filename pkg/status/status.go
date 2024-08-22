package status

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lookinlabs/status-page-middleware/pkg/checks"
	"github.com/lookinlabs/status-page-middleware/pkg/config"
	"github.com/lookinlabs/status-page-middleware/pkg/model"
)

func Services(cfg *config.Environments, ctx *gin.Context) {
	services, err := config.LoadEndpoints(cfg.StatusPageConfigPath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error loading config",
		})
		return
	}

	for i := range services {
		switch services[i].Type {
		case "http":
			method := "GET"
			headers := map[string]string{}
			body := ""
			var basicAuth *model.BasicAuth

			if services[i].Request != nil {
				method = services[i].Request.Method
				headers = services[i].Request.Headers
				bodyBytes, err := json.Marshal(services[i].Request.Body)
				if err != nil {
					services[i].Status = "error"
					continue
				}
				body = string(bodyBytes)
			}

			if services[i].BasicAuth != nil {
				basicAuth = services[i].BasicAuth
			}

			services[i].Status = checks.HTTP(services[i].URL, method, headers, body, basicAuth)
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
