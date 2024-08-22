package endpoints

import (
	"log"
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

func PingV3(ctx *gin.Context) {
	// Check for Basic Authentication
	username, password, hasAuth := ctx.Request.BasicAuth()
	if !hasAuth || username != "admin" || password != "" {
		var errorMessage string
		if !hasAuth {
			errorMessage = "No authentication provided"
			log.Println("Ping v3 endpoint: No authentication provided")
		} else if username != "admin" {
			errorMessage = "Invalid username"
			log.Println("Ping v3 endpoint: Invalid username")
		} else if password != "" {
			errorMessage = "Password should be empty"
			log.Println("Ping v3 endpoint: Password should be empty")
		}

		ctx.JSON(http.StatusForbidden, gin.H{
			"error":   "Forbidden",
			"message": errorMessage,
		})
		return
	}

	// Check Headers
	messageHeader := ctx.GetHeader("Message")
	pingPongHeader := ctx.GetHeader("PingPong")
	if messageHeader != "ping" || pingPongHeader != "pong" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid headers"})
		return
	}

	// Check Request Body
	var requestBody map[string]string
	if err := ctx.BindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	if requestBody["message"] != "ping" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid body"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "pong"})
}
