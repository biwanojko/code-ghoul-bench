//go:build !cgo

package main

// FallbackProcess is only compiled without CGO
func FallbackProcess(data []byte) int {
	return len(data)
}

// FallbackHelper is conditionally compiled
func FallbackHelper() string {
	return "fallback"
}
