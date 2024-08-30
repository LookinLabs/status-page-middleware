package config

import (
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/lookinlabs/status-page-middleware/pkg/helpers"
	"github.com/stretchr/testify/assert"
)

func TestConfig(test *testing.T) {
	var wg sync.WaitGroup
	wg.Add(2) // We have two tests to run in parallel

	go func() {
		defer wg.Done()
		testLoadStatusPage(test)
	}()

	go func() {
		defer wg.Done()
		testLoadEndpoints(test)
	}()

	wg.Wait()
}

func testLoadStatusPage(test *testing.T) {
	// Set up environment variables for testing
	envVars := map[string]string{
		"STATUS_PAGE_CONFIG_PATH":   "testdata/endpoints.json",
		"STATUS_PAGE_TEMPLATE_PATH": "testdata/status.html",
		"STATUS_PAGE_PATH":          "/teststatus",
	}
	defer helpers.SetEnv(envVars)()

	env, err := LoadStatusPage()
	assert.NoError(test, err)
	assert.Equal(test, "testdata/endpoints.json", env.StatusPageConfigPath)
	assert.Equal(test, "testdata/status.html", env.StatusPageTemplatePath)
	assert.Equal(test, "/teststatus", env.StatusPagePath)
}

func testLoadEndpoints(test *testing.T) {
	// Create a temporary JSON file for testing
	filePath := filepath.Join(os.TempDir(), "test_endpoints.json")
	data := `[
        {"name": "service1", "url": "http://service1.com"},
        {"name": "service2", "url": "http://service2.com"}
    ]`
	err := os.WriteFile(filePath, []byte(data), 0o644)
	assert.NoError(test, err)
	defer os.Remove(filePath)

	services, err := LoadEndpoints(filePath)
	assert.NoError(test, err)
	assert.Len(test, services, 2)
	assert.Equal(test, "service1", services[0].Name)
	assert.Equal(test, "http://service1.com", services[0].URL)
	assert.Equal(test, "service2", services[1].Name)
	assert.Equal(test, "http://service2.com", services[1].URL)
}
