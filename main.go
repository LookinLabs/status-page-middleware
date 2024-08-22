package main

import (
	"log"

	"github.com/gin-gonic/gin"
	controller "github.com/lookinlabs/status-page-middleware/controller"
	"github.com/lookinlabs/status-page-middleware/pkg/config"
	"github.com/lookinlabs/status-page-middleware/pkg/endpoints"
)

func main() {
	router := gin.Default()

	router.GET("/ping", controller.Ping)
	router.POST("/v2/ping", controller.PingV2)
	router.POST("/v3/ping", controller.PingV3)
	config, err := config.LoadStatusPage()
	if err != nil {
		log.Fatalf("failed to load status page environment variables: %v", err)
	}

	router.LoadHTMLGlob(config.StatusPageTemplatePath)
	router.Use(endpoints.StatusPage(config))
	router.Run(":8080")
}
