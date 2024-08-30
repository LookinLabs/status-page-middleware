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

func Services(_ *config.Environments, services []model.Service, ctx *gin.Context) {
	var wg sync.WaitGroup
	for each := range services {
		wg.Add(1)
		go func(service *model.Service) {
			defer wg.Done()
			checkService(service)
		}(&services[each])
	}

	wg.Wait()

	ctx.HTML(http.StatusOK, "status.html", gin.H{
		"services": services,
	})
}

func checkService(service *model.Service) {
	var wg sync.WaitGroup

	switch service.Type {
	case "http":
		wg.Add(1)
		go func() {
			checkHTTPService(service)
			defer wg.Done()
		}()
	case "dns":
		wg.Add(1)
		go func() {
			checkDNSService(service)
			defer wg.Done()
		}()
	case "tcp":
		wg.Add(1)
		go func() {
			checkTCPService(service)
			defer wg.Done()
		}()
	default:
		service.Status = "unknown"
		service.Error = "unknown service type"
	}

	wg.Wait()
}

func checkHTTPService(service *model.Service) {
	method, headers, body, basicAuth := prepareHTTPRequest(service)

	status, err := checks.HTTP(service.URL, method, headers, body, basicAuth)
	service.Status = status
	if err != nil {
		logger.Errorf("StatusMiddleware: Error checking HTTP status: %v", err)
		service.Error = err.Error()
	}
}

func prepareHTTPRequest(service *model.Service) (string, map[string]string, string, *model.BasicAuth) {
	method := "GET"
	headers := map[string]string{}
	body := ""

	if service.Request != nil {
		method = service.Request.Method
		headers = service.Request.Headers
		bodyBytes, err := json.Encode(service.Request.Body)
		if err != nil {
			logger.Errorf("StatusMiddleware: Error encoding JSON request body: %v", err)
			service.Status = "error"
			service.Error = err.Error()
			return method, headers, body, nil
		}
		body = string(bodyBytes)
	}

	var basicAuth *model.BasicAuth
	if service.BasicAuth != nil {
		logger.Infof("StatusMiddleware: Basic authentication provided for %s", service.URL)
		basicAuth = service.BasicAuth
	}

	return method, headers, body, basicAuth
}

func checkDNSService(service *model.Service) {
	status, err := checks.DNS(service.URL)
	service.Status = status
	if err != nil {
		logger.Errorf("StatusMiddleware: Error checking DNS status: %v", err)
		service.Error = err.Error()
	}
}

func checkTCPService(service *model.Service) {
	status, err := checks.TCP(service.URL)
	service.Status = status
	if err != nil {
		logger.Errorf("StatusMiddleware: Error checking TCP status: %v", err)
		service.Error = err.Error()
	}
}
