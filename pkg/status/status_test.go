package status

import (
	"sync"
	"testing"

	"github.com/lookinlabs/status-page-middleware/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestAllServicesConcurrently(test *testing.T) {
	var wg sync.WaitGroup

	wg.Add(5)
	go func() {
		defer wg.Done()
		service := &model.Service{Type: "http", URL: "http://example.com"}
		checkService(service)
		assert.NotEqual(test, "unknown", service.Status)
	}()

	go func() {
		defer wg.Done()
		service := &model.Service{Type: "http", URL: "http://example.com"}
		checkHTTPService(service)
		assert.NotEqual(test, "unknown", service.Status)
	}()

	go func() {
		defer wg.Done()
		service := &model.Service{
			Type: "http",
			URL:  "http://example.com",
			Request: &model.Request{
				Method:  "POST",
				Headers: map[string]string{"Content-Type": "application/json"},
				Body:    map[string]interface{}{"key": "value"},
			},
		}

		method, headers, body, basicAuth := prepareHTTPRequest(service)
		assert.Equal(test, "POST", method)
		assert.Equal(test, "application/json", headers["Content-Type"])
		assert.Contains(test, body, `"key":"value"`)
		assert.Nil(test, basicAuth)
	}()

	go func() {
		defer wg.Done()
		service := &model.Service{Type: "dns", URL: "example.com"}
		checkDNSService(service)
		assert.NotEqual(test, "unknown", service.Status)
	}()

	go func() {
		defer wg.Done()
		service := &model.Service{Type: "tcp", URL: "example.com:80"}
		checkTCPService(service)
		assert.NotEqual(test, "unknown", service.Status)
	}()

	wg.Wait()
}
