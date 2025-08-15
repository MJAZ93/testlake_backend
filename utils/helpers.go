package utils

import (
	"crypto/rand"
	"encoding/hex"
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

// GenerateSecureToken generates a secure random token of the specified length
func GenerateSecureToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
