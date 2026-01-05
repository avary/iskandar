package config

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

func ParseDestination(destination string) (string, error) {
	if portNum, err := strconv.Atoi(destination); err == nil {
		if portNum < 1 || portNum > 65535 {
			return "", fmt.Errorf("port must be between 1 and 65535")
		}
		return fmt.Sprintf("http://localhost:%d", portNum), nil
	}

	if !strings.Contains(destination, "://") {
		destination = "http://" + destination
	}

	destURL, err := url.Parse(destination)
	if err != nil {
		return "", fmt.Errorf("invalid destination format: %w", err)
	}

	port := destURL.Port()
	if port == "" {
		return "", fmt.Errorf("port is required in destination")
	}

	portNum, err := strconv.Atoi(port)
	if err != nil {
		return "", fmt.Errorf("port must be a number: %w", err)
	}
	if portNum < 1 || portNum > 65535 {
		return "", fmt.Errorf("port must be between 1 and 65535")
	}

	return fmt.Sprintf("http://%s:%s", destURL.Hostname(), port), nil
}

func ParseServerURL(serverURL string) (string, error) {
	scheme := "ws://"
	serverAddr := serverURL

	if strings.Contains(serverURL, "://") {
		if strings.HasPrefix(serverURL, "https://") {
			scheme = "wss://"
			serverAddr = strings.TrimPrefix(serverURL, "https://")
		} else if strings.HasPrefix(serverURL, "http://") {
			scheme = "ws://"
			serverAddr = strings.TrimPrefix(serverURL, "http://")
		} else {
			return "", fmt.Errorf("invalid scheme: only http:// and https:// are allowed")
		}
	}

	if !strings.Contains(serverAddr, "://") {
		serverAddr = "http://" + serverAddr
	}

	u, err := url.Parse(serverAddr)
	if err != nil {
		return "", fmt.Errorf("invalid server URL: %w", err)
	}

	return scheme + u.Host + "/tunnel/connect", nil
}
