// Package config validates development settings without exposing secret values.
package config

import (
	"errors"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func DatabaseURL() (string, error) {
	value := os.Getenv("DATABASE_URL")
	if strings.TrimSpace(value) == "" {
		return "", errors.New("DATABASE_URL is required")
	}
	parsed, err := url.Parse(value)
	if err != nil || (parsed.Scheme != "postgres" && parsed.Scheme != "postgresql") || parsed.Hostname() == "" || strings.Trim(parsed.Path, "/") == "" {
		return "", errors.New("DATABASE_URL must be a Postgres URL with host and database")
	}
	return value, nil
}

func ListenAddress() (string, error) {
	value := os.Getenv("RELAY_ADDR")
	if value == "" {
		value = "127.0.0.1:8080"
	}
	host, port, err := net.SplitHostPort(value)
	number, portErr := strconv.Atoi(port)
	if err != nil || host == "" || portErr != nil || number < 1 || number > 65535 {
		return "", errors.New("RELAY_ADDR must contain a host and port from 1 to 65535")
	}
	return value, nil
}
