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
        - Basic Auth [x]
        - Bearer Token [ ]
        - Custom Header [ ]
    - Health checks with Custom Headers [x]
    - Health checks via API Request against the endpoints [x]
- Add support for more protocols
    - HTTPS
    - FTP
    - SSH
    - WebSocket
- ChatOps - Send notifications when a service goes down or recovers
- Make more comprehensive UI for the status page
    - Add incident management support
    - Add metrics and monitoring support

## Quick Start

**1. Configure your endpoints in `config/endpoints.json` file**

You can get the examples from [Endpoints Config JSON File](./pkg/config/endpoints.json)

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

**2. Get the module**

```bash
go get -u github.com/lookinlabs/status-page-middleware
```

**3. Use the middleware**

```go
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/lookinlabs/status-page-middleware/pkg/endpoints"
	"github.com/lookinlabs/status-page-middleware/pkg/logger"
)

func main() {
	router := gin.Default()

	handler, err := endpoints.NewStatusPageController("path/to/config")
	if err != nil {
		log.Fatalf("Failed to initialize StatusPageController: %v", err)
	}
	handler.StatusPageMiddleware(router)

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

```

**4. Run your application**

```bash
make run
```

## Advanced

### Customizable Frontend

Status Page Middleware has default frontend written using HTML Template. You can customize the frontend by changing the HTML template file. Additionally you can copy the [Production Ready HTML Template](view/html/status.html) and edit according to your needs.

### Optional Environment Variables

Status Page Middleware can handle the following environment variables:
- `STATUS_PAGE_CONFIG_PATH` - Path to the JSON configuration file containing the endpoints to monitor. Default value is `config/endpoints.json`.
- `STATUS_PAGE_TEMPLATE_PATH` - Path to the HTML template file for the status page. Default value is `view/html/status.html`.
- `STATUS_PAGE_PATH` - URL Path to the status page. Default value is `/status`, which means that the status page will be available at `http://localhost:8080/status`.

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for detailed process for submitting pull requests to us.

## License

This project is licensed under the Apache 2.0 License - see the [LICENSE.md](LICENSE.md) file for details.
