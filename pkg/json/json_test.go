package json

import (
	"sync"
	"testing"
)

func TestJSONOperations(testCase *testing.T) {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		testCase.Run("Encode JSON", func(encodeCase *testing.T) {
			input := map[string]string{"key": "value"}
			_, err := Encode(input)
			if err != nil {
				encodeCase.Errorf("Expected no error, got %v", err)
			}
		})
	}()

	go func() {
		defer wg.Done()
		testCase.Run("Decode JSON", func(decodeCase *testing.T) {
			input := []byte(`{"key": "value"}`)
			var output map[string]string
			if err := Decode(input, &output); err != nil {
				decodeCase.Errorf("Expected no error, got %v", err)
			}

			if output["key"] != "value" {
				decodeCase.Errorf("Expected 'value', got %s", output["key"])
			}
		})
	}()

	wg.Wait()
}
