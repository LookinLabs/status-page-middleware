package config

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	json "github.com/lookinlabs/status-page-middleware/pkg/json"
	"github.com/lookinlabs/status-page-middleware/pkg/logger"
	"github.com/lookinlabs/status-page-middleware/pkg/model"
	"github.com/spf13/viper"
)

type Environments struct {
	StatusPageConfigPath   string
	StatusPageTemplatePath string
	StatusPagePath         string
}

func LoadStatusPage() (*Environments, error) {
	viper.SetConfigFile(".env")

	if err := viper.ReadInConfig(); err == nil {
		if err := godotenv.Load(); err != nil {
			logger.Errorf("StatusMiddleware: Error loading .env file: %v", err)
			return nil, err
		}
	}

	viper.AutomaticEnv()

	return &Environments{
		StatusPageConfigPath:   viper.GetString("STATUS_PAGE_CONFIG_PATH"),
		StatusPageTemplatePath: viper.GetString("STATUS_PAGE_TEMPLATE_PATH"),
		StatusPagePath:         viper.GetString("STATUS_PAGE_PATH"),
	}, nil
}

func LoadEndpoints(filename string) ([]model.Service, error) {
	file, err := os.Open(filepath.Clean(filename))
	if err != nil {
		logger.Errorf("StatusMiddleware: Error opening file: %v", err)
		return nil, err
	}
	defer file.Close()

	data, err := os.ReadFile(filename)
	if err != nil {
		logger.Errorf("StatusMiddleware: Error reading file: %v", err)
		return nil, err
	}

	var services []model.Service
	err = json.Decode(data, &services)
	if err != nil {
		logger.Errorf("StatusMiddleware: Error decoding JSON: %v", err)
		return nil, err
	}

	return services, nil
}
