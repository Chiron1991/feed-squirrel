package utils

import "os"

// Getenv is like os.Getenv, but returns the given fallback for failed lookups
func Getenv(key string, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
