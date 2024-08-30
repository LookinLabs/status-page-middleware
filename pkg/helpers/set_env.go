package helpers

import "os"

// SetEnv sets environment variables for testing and returns a function to unset them
func SetEnv(envVars map[string]string) func() {
	for key, value := range envVars {
		os.Setenv(key, value)
	}

	return func() {
		for key := range envVars {
			os.Unsetenv(key)
		}
	}
}
