package logger

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLogger(test *testing.T) {
	testLogLevel(test)
	testLoggingMessages(test)
	testLogFileCreationAndContent(test)
}

func testLogLevel(test *testing.T) {
	logLevels := []string{"debug", "info", "warn", "error", "dpanic", "panic", "fatal"}

	for _, level := range logLevels {
		test.Run("LogLevel_"+level, func(test *testing.T) {
			logPath := "./test_logs/test_log_level_" + level + ".log"
			os.Setenv("LOG_LEVEL", level)
			defer os.Unsetenv("LOG_LEVEL")

			os.Setenv("LOG_PATH", logPath)
			defer os.Unsetenv("LOG_PATH")

			initializeLogger() // Initialize the logger with the new log level

			expectedLevel := getLogLevelFromEnv()
			assert.True(test, log.Desugar().Core().Enabled(expectedLevel), "Expected log level %v to be enabled", expectedLevel)

			// Cleanup log file
			defer os.Remove(logPath)
		})
	}
}

func testLoggingMessages(test *testing.T) {
	test.Run("LoggingMessages", func(test *testing.T) {
		logPath := "./test_logs/test_logging_messages.log"
		os.Setenv("LOG_PATH", logPath)
		defer os.Unsetenv("LOG_PATH")

		initializeLogger() // Initialize the logger

		Infof("This is an info message")
		Warnf("This is a warning message")
		Errorf("This is an error message")

		// Ensure the logger is flushed
		log.Sync()

		// Add a small delay to ensure log messages are written to the file
		time.Sleep(100 * time.Millisecond)

		// Check if the log file contains the messages
		content, err := os.ReadFile(logPath)
		assert.NoError(test, err)

		assert.Contains(test, string(content), "This is an info message")
		assert.Contains(test, string(content), "This is a warning message")
		assert.Contains(test, string(content), "This is an error message")

		// Cleanup log file
		defer os.Remove(logPath)
	})
}

func testLogFileCreationAndContent(test *testing.T) {
	test.Run("LogFileCreationAndContent", func(test *testing.T) {
		logPath := "./test_logs/test_log_file_creation.log"
		os.Setenv("LOG_PATH", logPath)
		defer os.Unsetenv("LOG_PATH")

		initializeLogger() // Initialize the logger

		Infof("Testing log file creation")

		// Ensure the logger is flushed
		log.Sync()

		// Add a small delay to ensure log messages are written to the file
		time.Sleep(100 * time.Millisecond)

		// Check if the log file is created
		_, err := os.Stat(logPath)
		assert.NoError(test, err)

		// Check if the log file contains the message
		content, err := os.ReadFile(logPath)
		assert.NoError(test, err)
		assert.Contains(test, string(content), "Testing log file creation")

		// Cleanup log file
		defer os.Remove(logPath)
	})
}
