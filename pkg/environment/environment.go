// Package env environment handle environment variables
package environment

import (
	os "os"
)

// GetEnv gets environment with fallback
func GetEnv(key string, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return fallback
}
