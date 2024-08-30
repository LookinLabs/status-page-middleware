package checks

import (
	"sync"
	"testing"

	"github.com/lookinlabs/status-page-middleware/pkg/model"
)

func TestChecks(testCase *testing.T) {
	var wg sync.WaitGroup

	runTest := func(name string, testFunc func(t *testing.T)) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			testCase.Run(name, testFunc)
		}()
	}

	runTest("HTTP Check", func(httpCase *testing.T) {
		method := "POST"
		headers := map[string]string{
			"Content-Type": "application/json",
		}
		body := `{"key": "value"}`
		basicAuth := &model.BasicAuth{
			Username: "your_username",
			Password: "your_password",
		}
		status, _ := HTTP("http://localhost:8080/ping", method, headers, body, basicAuth)
		if status != "up" && status != "down" {
			httpCase.Errorf("Expected status to be 'up' or 'down', got %s", status)
		}
	})

	runTest("DNS Check", func(dnsCase *testing.T) {
		status, _ := DNS("http://localhost:3000/ping")
		if status != "up" && status != "down" {
			dnsCase.Errorf("Expected status to be 'up' or 'down', got %s", status)
		}
	})

	runTest("TCP Check", func(tcpCase *testing.T) {
		status, _ := TCP("http://localhost:3000")
		if status != "up" && status != "down" {
			tcpCase.Errorf("Expected status to be 'up' or 'down', got %s", status)
		}
	})

	wg.Wait()
}
