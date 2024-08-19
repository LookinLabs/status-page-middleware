package checks

import (
	"log"
	"net"
	"net/http"
	"net/url"
	"time"
)

func HTTP(urlString string) string {
	parsedURL, err := url.ParseRequestURI(urlString)
	if err != nil {
		log.Printf("Invalid URL: %v", err)
		return "down"
	}

	resp, err := http.Get(parsedURL.String())
	if err != nil || resp.StatusCode != http.StatusOK {
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
		return "down"
	}
	host := parsedURL.Hostname()
	_, err = net.LookupHost(host)
	if err != nil {
		return "down"
	}

	return "up"
}

func TCP(rawURL string) string {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
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
		return "down"
	}

	if err := conn.Close(); err != nil {
		log.Printf("Error closing connection: %v", err)
	}

	return "up"
}
