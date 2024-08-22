package checks

import (
	"bytes"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"
)

func HTTP(urlString, method string, headers map[string]string, requestBody string) string {
	parsedURL, err := url.ParseRequestURI(urlString)
	if err != nil {
		log.Printf("Invalid URL: %v", err)
		return "down"
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, parsedURL.String(), bytes.NewBuffer([]byte(requestBody)))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return "down"
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error performing request: %v", err)
		return "down"
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("Unexpected status code: %d", resp.StatusCode)
		return "down"
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("Error closing response body: %v", err)
		}
	}()

	return "up"
}

func DNS(rawURL string) string {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		log.Printf("Invalid URL: %v", err)
		return "down"
	}
	host := parsedURL.Hostname()
	_, err = net.LookupHost(host)
	if err != nil {
		log.Printf("DNS lookup failed for host: %v", err)
		return "down"
	}

	return "up"
}

func TCP(rawURL string) string {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		log.Printf("Invalid URL: %v", err)
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
		log.Printf("TCP connection failed: %v", err)
		return "down"
	}

	if err := conn.Close(); err != nil {
		log.Printf("Error closing connection: %v", err)
	}

	return "up"
}
