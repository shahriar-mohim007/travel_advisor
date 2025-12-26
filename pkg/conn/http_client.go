package conn

import (
	"net/http"
	"time"
)

var client *http.Client

// GetHTTClient return a http client
func GetHTTClient() *http.Client {
	return client
}

// InitClient init the http client
func InitClient() {
	client = &http.Client{
		Timeout: 60 * time.Second,
	}
}
