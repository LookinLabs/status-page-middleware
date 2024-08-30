package endpoints

import (
	"html/template"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/lookinlabs/status-page-middleware/pkg/config"
	"github.com/lookinlabs/status-page-middleware/pkg/logger"
	"github.com/lookinlabs/status-page-middleware/pkg/model"
	"github.com/lookinlabs/status-page-middleware/pkg/status"
	"github.com/lookinlabs/status-page-middleware/view"
)

type StatusPage struct {
	env      *config.Environments
	services []model.Service
}

func NewStatusPageController(configPath string) (*StatusPage, error) {
	env, err := config.LoadStatusPage()
	if err != nil {
		logger.Errorf("StatusMiddleware: Failed to load status page configuration: %v", err)
		return nil, err
	}

	services, err := config.LoadEndpoints(configPath)
	if err != nil {
		logger.Errorf("StatusMiddleware: Failed to load endpoints: %v", err)
		return nil, err
	}

	return &StatusPage{env: env, services: services}, nil
}

func (controller *StatusPage) StatusPageMiddleware(router *gin.Engine) {
	// Check if custom template path is provided
	if _, err := os.Stat(controller.env.StatusPageTemplatePath); os.IsNotExist(err) {
		// Use embedded template if custom template is not provided
		tmpl, err := template.ParseFS(view.StatusPageHTML, "html/status.html")
		if err != nil {
			logger.Fatalf("StatusMiddleware: Failed to parse embedded HTML template: %v", err)
			return
		}
		router.SetHTMLTemplate(tmpl)
	} else {
		// Use custom template if provided
		router.LoadHTMLGlob(controller.env.StatusPageTemplatePath)
		if err != nil {
			logger.Fatalf("StatusMiddleware: Failed to load HTML templates: %v", err)
			return
		}
	}

	router.Use(controller.getStatusPage())
}

func (controller *StatusPage) getStatusPage() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.URL.Path == controller.env.StatusPagePath {
			status.Services(controller.env, controller.services, ctx)
			ctx.Abort()
		} else {
			ctx.Next()
		}
	}
}
