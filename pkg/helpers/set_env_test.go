package helpers

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetEnv(test *testing.T) {
	envVars := map[string]string{
		"TEST_KEY_1": "TEST_VALUE_1",
		"TEST_KEY_2": "TEST_VALUE_2",
	}

	// Set environment variables
	unsetEnv := SetEnv(envVars)

	// Verify that the environment variables are set
	for key, expectedValue := range envVars {
		actualValue, exists := os.LookupEnv(key)
		assert.True(test, exists, "Expected environment variable %s to be set", key)
		assert.Equal(test, expectedValue, actualValue, "Expected environment variable %s to have value %s, got %s", key, expectedValue, actualValue)
	}

	// Unset environment variables
	unsetEnv()

	// Verify that the environment variables are unset
	for key := range envVars {
		_, exists := os.LookupEnv(key)
		assert.False(test, exists, "Expected environment variable %s to be unset", key)
	}
}
