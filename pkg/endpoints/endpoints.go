package endpoints

import (
	"html/template" // Use html/template instead of text/template
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/lookinlabs/status-page-middleware/pkg/config"
	"github.com/lookinlabs/status-page-middleware/pkg/logger"
	"github.com/lookinlabs/status-page-middleware/pkg/status"
	"github.com/lookinlabs/status-page-middleware/view" // Import the view package
)

func StatusPageMiddleware(router *gin.Engine) {
	env, err := config.LoadStatusPage()
	if err != nil {
		logger.Fatalf("StatusMiddleware: Failed to load status page configuration: %v", err)
		return
	}

	// Check if custom template path is provided
	if _, err := os.Stat(env.StatusPageTemplatePath); os.IsNotExist(err) {
		// Use embedded template if custom template is not provided
		tmpl, err := template.ParseFS(view.StatusPageHTML, "html/status.html")
		if err != nil {
			logger.Fatalf("StatusMiddleware: Failed to parse embedded HTML template: %v", err)
			return
		}
		router.SetHTMLTemplate(tmpl)
	} else {
		// Use custom template if provided
		router.LoadHTMLGlob(env.StatusPageTemplatePath)
		if err != nil {
			logger.Fatalf("StatusMiddleware: Failed to load HTML templates: %v", err)
			return
		}
	}

	router.Use(getStatusPage(env))
}

func getStatusPage(cfg *config.Environments) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.URL.Path == cfg.StatusPagePath {
			status.Services(cfg, ctx)
			ctx.HTML(http.StatusOK, "status.html", gin.H{})
			ctx.Abort()
		} else {
			ctx.Next()
		}
	}
}
