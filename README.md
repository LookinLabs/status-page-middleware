# Status Page Middleware

Status Page middleware is used to generate a simple status page for your application based on a JSON configuration file and environment variables. It allows you to define multiple endpoints to monitor and display their status on a single page. The middleware periodically checks the status of each endpoint and updates the status page accordingly. It also provides options to customize the look and feel of the status page and integrate with popular dashboard tools for visualizing the status and metrics.

The unique functionality of this status page is that it can be used for any application which uses Gin framework. It's also highly flexible and can be customized by the development team themselves.

## Table of Contents

- [Roadmap](#roadmap)
- [Quick Start](#quick-start)
- [Using Middleware with Viper Configuration](#using-middleware-with-viper-configuration)
- [Contributing](#contributing)
- [License](#license)

## Roadmap

- Add embed support for the status page template
- Add more advanced health checks
    - Health checks with Authentication
    - Health checks with Custom Headers
    - Health checks via API Request against the endpoints
- Add support for more protocols
    - HTTPS
    - FTP
    - SSH
    - WebSocket
- ChatOps - Send notifications when a service goes down or recovers
- Make more comprehensive UI for the status page

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

**2. Configure your environment variables (optional)**

By default STATUS_PAGE_CONFIG_PATH, STATUS_PAGE_TEMPLATE_PATH, and STATUS_PAGE_PATH are set to the following values:

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
	"github.com/lookinlabs/status-page-middleware/pkg/config"
	"github.com/lookinlabs/status-page-middleware/pkg/endpoints"
)

func main() {
	router := gin.Default()

	router.GET("/ping", controller.Ping)
	config, err := config.LoadStatusPage()
	if err != nil {
		return nil, fmt.Errorf("failed to load status page environment variables: %v", err)
	}

	router.LoadHTMLGlob(config.StatusPageTemplatePath)
	router.Use(endpoints.StatusPage(config))
	router.Run(":8080")
}
```

**5. Write the status page**

Write the status page via HTML template that you've specified in `STATUS_PAGE_TEMPLATE_PATH` environment variable.

By default you can use the [Production Ready HTML Template](view/html/status.html) for your status page.

**6. Run your application**

```bash
make run
```

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for detailed process for submitting pull requests to us.

## License

This project is licensed under the Apache 2.0 License - see the [LICENSE.md](LICENSE.md) file for details.
