package helpers

import (
	"os"

	"github.com/lookinlabs/status-page-middleware/pkg/logger"
)

// SetEnv sets environment variables for testing and returns a function to unset them
func SetEnv(envVars map[string]string) func() {
	for key, value := range envVars {
		if err := os.Setenv(key, value); err != nil {
			logger.Errorf("StatusMiddleware: Failed to set environment variable %s: %v", key, err)
		}
	}

	return func() {
		for key := range envVars {
			if err := os.Unsetenv(key); err != nil {
				logger.Errorf("StatusMiddleware: Failed to unset environment variable %s: %v", key, err)
			}
		}
	}
}
