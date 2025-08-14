package utils

import (
	"fmt"
	"os"
)

func GetBaseURL() string {
	scheme := os.Getenv("SCHEME")
	if scheme == "" {
		scheme = "http"
	}

	ip := os.Getenv("IP")
	if ip == "" {
		ip = "localhost"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	return fmt.Sprintf("%s://%s:%s", scheme, ip, port)
}
