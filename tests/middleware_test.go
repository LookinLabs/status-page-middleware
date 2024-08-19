package tests

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lookinlabs/status-page-middleware/config"
	"github.com/lookinlabs/status-page-middleware/controller"
	"github.com/lookinlabs/status-page-middleware/middleware"
)

func TestStatusPageMiddleware(testCase *testing.T) {
	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	// Initialize router
	router := gin.Default()

	cfg := &config.Environments{
		StatusPageConfigPath: "../config/endpoints.json",
	}

	// Load HTML templates
	router.LoadHTMLGlob("../view/html/status.html")

	// Use middleware
	router.Use(middleware.StatusPage(cfg))

	// Define a test endpoint
	router.GET("/ping", controller.Ping)

	// Test the /ping endpoint
	testCase.Run("Ping Endpoint", func(pingCase *testing.T) {
		response := httptest.NewRecorder()
		request, _ := http.NewRequest(http.MethodGet, "/ping", nil)
		router.ServeHTTP(response, request)

		if response.Code != http.StatusOK {
			pingCase.Errorf("Expected status code 200, got %d", response.Code)
		}
		expectedBody := `{"message":"pong"}`
		if response.Body.String() != expectedBody {
			pingCase.Errorf("Expected body '%s', got %s", expectedBody, response.Body.String())
		}
	})

	// Test the status page endpoint
	testCase.Run("Status Page Endpoint", func(t *testing.T) {
		response := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, cfg.StatusPagePath, nil)
		router.ServeHTTP(response, req)

		if response.Code != http.StatusOK {
			t.Errorf("Expected status code 200, got %d", response.Code)
		}
		// Check if the body contains specific HTML content
		expectedContent := "<html>" // Adjust this to match the actual content of your HTML template
		if !strings.Contains(response.Body.String(), expectedContent) {
			t.Errorf("Expected body to contain '%s', got %s", expectedContent, response.Body.String())
		}
	})
}
