package json

import (
	"log"

	json "github.com/goccy/go-json"
)

func Encode(input interface{}) ([]byte, error) {
	bodyBytes, err := json.Marshal(input)
	if err != nil {
		log.Printf("unable to marshal object: %v", err)
		return nil, err
	}

	return bodyBytes, nil
}

func Decode(input []byte, output interface{}) error {
	if err := json.Unmarshal(input, output); err != nil {
		log.Printf("unable to unmarshal data: %v", err)
		return err
	}

	return nil
}
