package status

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/lookinlabs/status-page-middleware/pkg/checks"
	"github.com/lookinlabs/status-page-middleware/pkg/config"
	json "github.com/lookinlabs/status-page-middleware/pkg/json"
	"github.com/lookinlabs/status-page-middleware/pkg/logger"
	"github.com/lookinlabs/status-page-middleware/pkg/model"
)

func Services(cfg *config.Environments, ctx *gin.Context) {
	services, err := config.LoadEndpoints(cfg.StatusPageConfigPath)
	if err != nil {
		logger.Errorf("StatusMiddleware: Error loading config: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error loading config",
		})

		return
	}

	var wg sync.WaitGroup
	for each := range services {
		wg.Add(1)
		go func(service *model.Service) {
			defer wg.Done()

			switch service.Type {
			case "http":
				method := "GET"
				headers := map[string]string{}
				body := ""

				var basicAuth *model.BasicAuth
				if service.Request != nil {
					method = service.Request.Method
					headers = service.Request.Headers
					bodyBytes, err := json.Encode(service.Request.Body)
					if err != nil {
						logger.Errorf("StatusMiddleware: Error encoding JSON request body: %v", err)
						service.Status = "error"
						service.Error = err.Error()

						return
					}
					body = string(bodyBytes)
				}

				if service.BasicAuth != nil {
					logger.Infof("StatusMiddleware: Basic authentication provided for %s", service.URL)
					basicAuth = service.BasicAuth
				}

				status, err := checks.HTTP(service.URL, method, headers, body, basicAuth)
				service.Status = status
				if err != nil {
					logger.Errorf("StatusMiddleware: Error checking HTTP status: %v", err)
					service.Error = err.Error()
				}

			case "dns":
				status, err := checks.DNS(service.URL)
				service.Status = status
				if err != nil {
					logger.Errorf("StatusMiddleware: Error checking DNS status: %v", err)
					service.Error = err.Error()
				}

			case "tcp":
				status, err := checks.TCP(service.URL)
				service.Status = status
				if err != nil {
					logger.Errorf("StatusMiddleware: Error checking TCP status: %v", err)
					service.Error = err.Error()
				}

			default:
				service.Status = "unknown"
				service.Error = "unknown service type"
			}
		}(&services[each])
	}

	wg.Wait()

	ctx.HTML(http.StatusOK, "status.html", gin.H{
		"services": services,
	})
}
