package checks

import (
	"bytes"
	"encoding/base64"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/lookinlabs/status-page-middleware/pkg/model"
)

func HTTP(urlString, method string, headers map[string]string, requestBody string, basicAuth *model.BasicAuth) string {
	parsedURL, err := url.ParseRequestURI(urlString)
	if err != nil {
		log.Printf("StatusMiddleware: Invalid URL: %v", err)
		return "down"
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, parsedURL.String(), bytes.NewBuffer([]byte(requestBody)))
	if err != nil {
		log.Printf("StatusMiddleware: Error creating request: %v", err)
		return "down"
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Set basic authentication if provided
	if basicAuth != nil {
		log.Printf("StatusMiddleware: Setting basic auth: %s:%s", basicAuth.Username, basicAuth.Password)
		auth := basicAuth.Username + ":" + basicAuth.Password
		encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))
		req.Header.Set("Authorization", "Basic "+encodedAuth)
		log.Printf("StatusMiddleware: Authorization header set: Basic %s", encodedAuth)
	} else {
		log.Printf("StatusMiddleware: No basic auth provided")
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("StatusMiddleware: Error performing request: %v", err)
		return "down"
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("StatusMiddleware: Unexpected status code: %d", resp.StatusCode)
		return "down"
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("StatusMiddleware: Error closing response body: %v", err)
		}
	}()

	return "up"
}

func DNS(rawURL string) string {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		log.Printf("StatusMiddleware: Invalid URL: %v", err)
		return "down"
	}
	host := parsedURL.Hostname()
	_, err = net.LookupHost(host)
	if err != nil {
		log.Printf("StatusMiddleware: DNS lookup failed for host: %v", err)
		return "down"
	}

	return "up"
}

func TCP(rawURL string) string {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		log.Printf("StatusMiddleware: Invalid URL: %v", err)
		return "down"
	}

	host := parsedURL.Hostname()
	port := parsedURL.Port()
	if port == "" {
		port = "80"
	}

	address := net.JoinHostPort(host, port)
	conn, err := net.DialTimeout("tcp", address, 5*time.Second)
	if err != nil {
		log.Printf("StatusMiddleware: TCP connection failed: %v", err)
		return "down"
	}

	if err := conn.Close(); err != nil {
		log.Printf("StatusMiddleware: Error closing connection: %v", err)
	}

	return "up"
}
