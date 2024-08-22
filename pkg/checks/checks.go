package checks

import (
	"bytes"
	"encoding/base64"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/lookinlabs/status-page-middleware/pkg/logger"
	"github.com/lookinlabs/status-page-middleware/pkg/model"
)

func HTTP(urlString, method string, headers map[string]string, requestBody string, basicAuth *model.BasicAuth) (string, error) {
	parsedURL, err := url.ParseRequestURI(urlString)
	if err != nil {
		logger.Errorf("StatusMiddleware: Invalid URL: %v", err)
		return "down", err
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, parsedURL.String(), bytes.NewBuffer([]byte(requestBody)))
	if err != nil {
		logger.Errorf("StatusMiddleware: Error creating request: %v", err)
		return "down", err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Set basic authentication if provided
	if basicAuth != nil {
		auth := basicAuth.Username + ":" + basicAuth.Password
		encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))
		req.Header.Set("Authorization", "Basic "+encodedAuth)
	} else {
		logger.Infof("StatusMiddleware: No basic authentication provided for %s", urlString)
	}

	resp, err := client.Do(req)
	if err != nil {
		logger.Errorf("StatusMiddleware: Error sending request: %v", err)
		return "down", err
	}

	if resp.StatusCode != http.StatusOK {
		logger.Warnf("StatusMiddleware: Non-200 status code: %d", resp.StatusCode)
		return "down", err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			logger.Errorf("StatusMiddleware: Error closing response body: %v", err)
		}
	}()

	return "up", nil
}

func DNS(rawURL string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		logger.Errorf("StatusMiddleware: Invalid URL: %v", err)
		return "down", err
	}

	host := parsedURL.Hostname()
	_, err = net.LookupHost(host)
	if err != nil {
		logger.Errorf("StatusMiddleware: DNS lookup failed: %v", err)
		return "down", err
	}

	return "up", nil
}

func TCP(rawURL string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		logger.Errorf("StatusMiddleware: Invalid URL: %v", err)
		return "down", err
	}

	host := parsedURL.Hostname()
	port := parsedURL.Port()
	if port == "" {
		port = "80"
	}

	address := net.JoinHostPort(host, port)
	conn, err := net.DialTimeout("tcp", address, 5*time.Second)
	if err != nil {
		logger.Errorf("StatusMiddleware: TCP connection failed: %v", err)
		return "down", err
	}

	if err := conn.Close(); err != nil {
		logger.Errorf("StatusMiddleware: Error closing TCP connection: %v", err)
		return "up", err
	}

	return "up", nil
}
