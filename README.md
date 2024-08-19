# Status Page Middleware

Status Page middleware is based on Gin-Gonic framework and used to generate simple status page based on JSON configuration file and environment variables.

## Quick Start

**1. Configure your endpoints in `config.json` file**

```json
[
    {
        "name": "MyFancy HTTP API",
        "url": "http://localhost:8080/ping",
        "type": "http"
    },
    {
        "name": "MyFancy DNS Check",
        "url": "http://localhost/",
        "type": "dns"
    },
    {
        "name": "MyFancy TCP Check for API",
        "url": "http://localhost:8080/health",
        "type": "tcp"
    }
]
```

**2. Configure your environment variables**

```bash
STATUS_PAGE_CONFIG_PATH="config/endpoints.json"
STATUS_PAGE_TEMPLATE_PATH="view/html/status.html"
STATUS_PAGE_PATH="/status"
```

**3. Get the module**

```bash
go get -u github.com/lookinlabs/status-page-middleware
```

**4. Use the middleware**

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/lookinlabs/status-page-middleware"
)

func main() {
	router := gin.Default()

	router.GET("/ping", controller.Ping)
	cfg, err := config.LoadEnvironmentals()
	if err != nil {
		log.Fatalf("Error loading environmentals: %v", err)
	}

	router.LoadHTMLGlob(cfg.StatusPageTemplatePath)
	router.Use(middleware.StatusPage(cfg))

	router.Run(":8080")
}
```

**5. Run your application**

```bash
make run
```

## Using Middleware with Viper Configuration

If you have already Viper configuration done, than it's recommended to use it for loading configuration.

**1. Add Status Page environmentals to Viper Config**
```go
// config/config.go
package config

import (
	"encoding/json"
	"os"

	"github.com/joho/godotenv"
	"github.com/lookinlabs/status-page-middleware/model"
	"github.com/spf13/viper"
)

type Environmentals struct {
    MyAppPort              string
    MyAppHost              string
    MyAppEnv               string
	StatusPageConfigPath   string
	StatusPageTemplatePath string
    StatusPagePath         string
}

func LoadEnvironmentals() (*Environmentals, error) {
	viper.SetConfigFile(".env")

	if err := viper.ReadInConfig(); err == nil {
		if err := godotenv.Load(); err != nil {
			return nil, err
		}
	}

	viper.AutomaticEnv()

	return &Environmentals{
        MyAppPort:              viper.GetString("MY_APP_PORT"),
        MyAppHost:              viper.GetString("MY_APP_HOST"),
        MyAppEnv:               viper.GetString("MY_APP_ENV"),
		StatusPageConfigPath:   viper.GetString("STATUS_PAGE_CONFIG_PATH"),
		StatusPageTemplatePath: viper.GetString("STATUS_PAGE_TEMPLATE_PATH"),
        StatusPagePath:         viper.GetString("STATUS_PAGE_PATH"),
	}, nil
}
```

**2. Use the middleware with Viper configuration**
```go
// middleware/router.go
    router := gin.Default()

    // Add the status page middleware
    router.LoadHTMLGlob(cfg.StatusPageTemplatePath)
    router.Use(StatusPage(cfg))
    
    return router
```

**3. Import the NewRouter to main**
```go
// main.go
package main

import (
    "log"

    "github.com/gin-gonic/gin"
    "github.com/lookinlabs/status-page-middleware/middleware"
    "github.com/lookinlabs/status-page-middleware/model"
    "github.com/lookinlabs/status-page-middleware/config"
)

func main() {
    cfg, err := config.LoadEnvironmentals()
    if err != nil {
        log.Fatalf("Error loading environmentals: %v", err)
    }

    router := middleware.NewRouter(cfg)
    router.Run(cfg.MyAppPort)
}
```

**4. Run your application**
```bash
go run main.go
```