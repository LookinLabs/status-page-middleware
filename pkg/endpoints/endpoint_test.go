package endpoints

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lookinlabs/status-page-middleware/pkg/helpers"
	"github.com/stretchr/testify/assert"
)

func TestEndpoints(test *testing.T) {
	// Set up environment variables for testing
	envVars := map[string]string{
		"STATUS_PAGE_CONFIG_PATH":   "testdata/endpoints.json",
		"STATUS_PAGE_TEMPLATE_PATH": "testdata/status.html",
		"STATUS_PAGE_PATH":          "/teststatus",
	}
	defer helpers.SetEnv(envVars)()

	// Create a temporary JSON file for testing
	filePath := filepath.Join(os.TempDir(), "test_endpoints.json")
	Data := `[
        {"name": "service1", "url": "http://service1.com"},
        {"name": "service2", "url": "http://service2.com"}
    ]`
	err := os.WriteFile(filePath, []byte(Data), 0o644)
	assert.NoError(test, err)
	defer os.Remove(filePath)

	// WaitGroup to synchronize goroutines
	var wg sync.WaitGroup
	wg.Add(2)

	// Run testNewStatusPageController in a goroutine
	go func() {
		defer wg.Done()
		testNewStatusPageController(test, filePath)
	}()

	// Run testStatusPageMiddleware in a goroutine
	go func() {
		defer wg.Done()
		testStatusPageMiddleware(test, filePath)
	}()

	// Wait for both goroutines to finish
	wg.Wait()
}

func testNewStatusPageController(test *testing.T, filePath string) {
	// Initialize the StatusPageController
	handler, err := NewStatusPageController(filePath)
	if err != nil {
		log.Fatalf("Failed to initialize StatusPageController: %v", err)
	}

	assert.NoError(test, err)
	assert.NotNil(test, handler)
	assert.Len(test, handler.services, 2)
	assert.Equal(test, "service1", handler.services[0].Name)
	assert.Equal(test, "http://service1.com", handler.services[0].URL)
	assert.Equal(test, "service2", handler.services[1].Name)
	assert.Equal(test, "http://service2.com", handler.services[1].URL)
}

func testStatusPageMiddleware(test *testing.T, testFilePath string) {
	// Initialize the StatusPageController
	handler, err := NewStatusPageController(testFilePath)
	if err != nil {
		log.Fatalf("Failed to initialize StatusPageController: %v", err)
	}
	assert.NoError(test, err)
	assert.NotNil(test, handler)

	// Set up the Gin router and apply the middleware
	router := gin.Default()
	handler.StatusPageMiddleware(router)

	// Create a test HTTP server
	response := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/teststatus", nil)
	router.ServeHTTP(response, req)

	assert.Equal(test, http.StatusOK, response.Code)
	assert.Contains(test, response.Body.String(), "<title>Status Page</title>")
	assert.Contains(test, response.Body.String(), "service1")
	assert.Contains(test, response.Body.String(), "service2")
}
