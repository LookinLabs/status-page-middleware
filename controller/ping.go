package endpoints

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Ping(ctx *gin.Context) {
	ctx.String(http.StatusOK, "pong")
}

func PingV2(ctx *gin.Context) {
	messageHeader := ctx.GetHeader("Message")
	pingPongHeader := ctx.GetHeader("PingPong")

	var requestBody map[string]string
	if err := ctx.BindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if messageHeader != "ping" || pingPongHeader != "pong" || requestBody["message"] != "ping" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid headers or body"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "pong"})
}
