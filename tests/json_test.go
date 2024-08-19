package tests

import (
	"encoding/json"
	"testing"
)

func TestJSONOperations(testCase *testing.T) {
	done := make(chan bool)

	go func() {
		testCase.Run("Encode JSON", func(encodeCase *testing.T) {
			input := map[string]string{"key": "value"}
			_, err := json.Marshal(input)
			if err != nil {
				encodeCase.Errorf("Expected no error, got %v", err)
			}
		})
		done <- true
	}()

	go func() {
		testCase.Run("Decode JSON", func(decodeCase *testing.T) {
			input := []byte(`{"key": "value"}`)
			var output map[string]string
			err := json.Unmarshal(input, &output)
			if err != nil {
				decodeCase.Errorf("Expected no error, got %v", err)
			}
			if output["key"] != "value" {
				decodeCase.Errorf("Expected 'value', got %s", output["key"])
			}
		})
		done <- true
	}()

	<-done
	<-done
}
