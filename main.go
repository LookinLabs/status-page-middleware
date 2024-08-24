package main

import (
	"log"

	"github.com/gin-gonic/gin"
	controller "github.com/lookinlabs/status-page-middleware/controller"
	"github.com/lookinlabs/status-page-middleware/pkg/endpoints"
	"github.com/lookinlabs/status-page-middleware/pkg/logger"
)

func main() {
	router := gin.Default()

	router.GET("/ping", controller.Ping)
	router.POST("/v2/ping", controller.PingV2)
	router.POST("/v3/ping", controller.PingV3)

	controller, err := endpoints.NewStatusPageController("path/to/config")
	if err != nil {
		log.Fatalf("Failed to initialize StatusPageController: %v", err)
	}
	controller.StatusPageMiddleware(router)

	if err := router.Run(":8080"); err != nil {
		logger.Fatalf("StatusMiddleware: Failed to run server: %v", err)
	}
}
