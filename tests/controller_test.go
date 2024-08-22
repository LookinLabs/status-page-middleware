package tests

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/gin-gonic/gin"
	controller "github.com/lookinlabs/status-page-middleware/controller"
	"github.com/lookinlabs/status-page-middleware/pkg/config"
	"github.com/lookinlabs/status-page-middleware/pkg/status"
)

func TestControllers(testCase *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	var wg sync.WaitGroup

	runTest := func(name string, testFunc func(t *testing.T)) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			testCase.Run(name, testFunc)
		}()
	}

	runTest("Ping", func(testPing *testing.T) {
		router := gin.Default()
		router.GET("/ping", controller.Ping)

		response := httptest.NewRecorder()
		request, _ := http.NewRequest(http.MethodGet, "/ping", nil)
		router.ServeHTTP(response, request)

		if response.Code != http.StatusOK {
			testPing.Errorf("Expected status code 200, got %d", response.Code)
		}
		if response.Body.String() != `{"message":"pong"}` {
			testPing.Errorf("Expected body '{\"message\":\"pong\"}', got %s", response.Body.String())
		}
	})

	runTest("Service Statuses", func(testServices *testing.T) {
		router := gin.Default()
		router.LoadHTMLGlob("../view/html/status.html")

		cfg := &config.Environments{
			StatusPageConfigPath: "../pkg/config/endpoints.json",
		}
		router.GET("/status", func(ctx *gin.Context) {
			status.Services(cfg, ctx)
		})

		response := httptest.NewRecorder()
		request, _ := http.NewRequest(http.MethodGet, "/status", nil)
		router.ServeHTTP(response, request)

		if response.Code != http.StatusOK {
			testServices.Errorf("Expected status code 200, got %d", response.Code)
		}
	})

	wg.Wait()
}
