package server

import (
	"crypto/rand"
	"encoding/hex"
	"strings"
)

// GenerateID generates a random hex ID
func GenerateID(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// SanitizePath removes dangerous path components
func SanitizePath(path string) string {
	parts := strings.Split(path, "/")
	var safe []string
	for _, p := range parts {
		if p != ".." && p != "." && p != "" {
			safe = append(safe, p)
		}
	}
	return "/" + strings.Join(safe, "/")
}

// TruncateString truncates a string to max length - dead code
func TruncateString(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}

// PadLeft pads a string on the left - dead code
func PadLeft(s string, width int, pad byte) string {
	if len(s) >= width {
		return s
	}
	padding := strings.Repeat(string(pad), width-len(s))
	return padding + s
}
