package model

import (
	"sync"
	"testing"

	json "github.com/lookinlabs/status-page-middleware/pkg/json"
	"github.com/stretchr/testify/assert"
)

func TestModelMarshalling(test *testing.T) {
	var wg sync.WaitGroup
	wg.Add(3) // We have three tests to run in parallel

	go func() {
		defer wg.Done()
		test.Run("RequestMarshalling", func(requestCase *testing.T) {
			request := &Request{
				Method: "GET",
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
				Body: map[string]interface{}{
					"key": "value",
				},
			}

			data, err := json.Encode(request)
			assert.NoError(requestCase, err)

			var unmarshalledRequest Request
			if err = json.Decode(data, &unmarshalledRequest); err != nil {
				requestCase.Errorf("Expected no error, got %v", err)
			}
			assert.NoError(requestCase, err)
			assert.Equal(requestCase, request, &unmarshalledRequest)
		})
	}()

	go func() {
		defer wg.Done()
		test.Run("BasicAuthMarshalling", func(authCase *testing.T) {
			basicAuth := &BasicAuth{
				Username: "user",
				Password: "pass",
			}

			data, err := json.Encode(basicAuth)
			assert.NoError(authCase, err)

			var unmarshalledBasicAuth BasicAuth
			if err = json.Decode(data, &unmarshalledBasicAuth); err != nil {
				authCase.Errorf("Expected no error, got %v", err)
			}
			assert.NoError(authCase, err)
			assert.Equal(authCase, basicAuth, &unmarshalledBasicAuth)
		})
	}()

	go func() {
		defer wg.Done()
		test.Run("ServiceMarshalling", func(serviceCase *testing.T) {
			service := &Service{
				Name:   "Test Service",
				URL:    "http://example.com",
				Type:   "REST",
				Status: "active",
				Request: &Request{
					Method: "POST",
					Headers: map[string]string{
						"Authorization": "Bearer token",
					},
					Body: map[string]interface{}{
						"data": "test",
					},
				},
				BasicAuth: &BasicAuth{
					Username: "admin",
					Password: "adminpass",
				},
			}

			data, err := json.Encode(service)
			assert.NoError(serviceCase, err)

			var unmarshalledService Service
			if err = json.Decode(data, &unmarshalledService); err != nil {
				serviceCase.Errorf("Expected no error, got %v", err)
			}
			assert.NoError(serviceCase, err)
			assert.Equal(serviceCase, service, &unmarshalledService)
		})
	}()

	wg.Wait()
}
