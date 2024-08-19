package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/lookinlabs/status-page-middleware/config"
	"github.com/lookinlabs/status-page-middleware/controller"
	"github.com/lookinlabs/status-page-middleware/middleware"
)

func main() {
	router := gin.Default()

	router.GET("/ping", controller.Ping)
	cfg, err := config.LoadEnvironmentVariables()
	if err != nil {
		log.Fatalf("Error loading environments: %v", err)
	}

	router.LoadHTMLGlob(cfg.StatusPageTemplatePath)
	router.Use(middleware.StatusPage(cfg))

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
