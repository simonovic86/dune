package main

import (
	"strings"
)

// NormalizeQuery simple query normalization
func NormalizeQuery(query string) string {
	return strings.ToLower(strings.TrimSpace(query))
}
