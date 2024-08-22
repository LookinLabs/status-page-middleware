package json

import (
	json "github.com/goccy/go-json"
	"github.com/lookinlabs/status-page-middleware/pkg/logger"
)

func Encode(input interface{}) ([]byte, error) {
	bodyBytes, err := json.Marshal(input)
	if err != nil {
		logger.Errorf("unable to marshal object: %v", err)
		return nil, err
	}

	return bodyBytes, nil
}

func Decode(input []byte, output interface{}) error {
	if err := json.Unmarshal(input, output); err != nil {
		logger.Errorf("unable to unmarshal object: %v", err)
		return err
	}

	return nil
}
