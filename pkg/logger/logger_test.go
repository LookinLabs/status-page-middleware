package logger

import (
	"bytes"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLogger(test *testing.T) {
	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()
		testLogLevel(test)
	}()

	go func() {
		defer wg.Done()
		testLogFileCreationAndContent(test)
	}()

	time.Sleep(100 * time.Millisecond)

	testLoggingMessages(test)
	wg.Wait()
}

func testLogLevel(test *testing.T) {
	logLevels := []string{"debug", "info", "warn", "error", "dpanic", "panic", "fatal"}

	for _, level := range logLevels {
		test.Run("LogLevel_"+level, func(test *testing.T) {
			logPath := "./test_logs/test_log_level_" + level + ".log"
			test.Setenv("LOG_LEVEL", level)
			defer os.Unsetenv("LOG_LEVEL")

			test.Setenv("LOG_PATH", logPath)
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
		// Create a buffer to capture the logs
		var buffer bytes.Buffer

		// Redirect stdout to the buffer
		oldStdout := os.Stdout
		readPipe, writePipe, _ := os.Pipe()
		os.Stdout = writePipe

		// Initialize the logger to write to stdout
		initializeLogger()

		// Log messages
		Infof("This is an info message")
		Warnf("This is a warning message")
		Errorf("This is an error message")

		// Close the writer and restore stdout
		writePipe.Close()
		os.Stdout = oldStdout

		// Read the captured logs
		_, err := buffer.ReadFrom(readPipe)
		if err != nil {
			test.Errorf("Failed to read from pipe: %v", err)
		}

		// Check if the buffer contains the log messages
		logOutput := buffer.String()
		assert.Contains(test, logOutput, "This is an info message")
		assert.Contains(test, logOutput, "This is a warning message")
		assert.Contains(test, logOutput, "This is an error message")
	})
}

func testLogFileCreationAndContent(test *testing.T) {
	test.Run("LogFileCreationAndContent", func(test *testing.T) {
		logPath := "./test_logs/test_log_file_creation.log"
		test.Setenv("LOG_PATH", logPath)
		defer os.Unsetenv("LOG_PATH")

		initializeLogger() // Initialize the logger

		Infof("Testing log file creation")

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
