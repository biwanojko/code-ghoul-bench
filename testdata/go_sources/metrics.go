package server

import "sync/atomic"

// Metrics tracks server metrics
type Metrics struct {
	requests int64
	errors   int64
}

// globalMetrics is a package-level metrics instance
var globalMetrics = &Metrics{}

// IncrRequests increments the request counter
func IncrRequests() {
	atomic.AddInt64(&globalMetrics.requests, 1)
}

// IncrErrors increments the error counter
func IncrErrors() {
	atomic.AddInt64(&globalMetrics.errors, 1)
}

// GetStats returns current metrics snapshot - never called (dead code)
func GetStats() (int64, int64) {
	return atomic.LoadInt64(&globalMetrics.requests),
		atomic.LoadInt64(&globalMetrics.errors)
}

// resetMetrics resets all counters - dead code
func resetMetrics() {
	atomic.StoreInt64(&globalMetrics.requests, 0)
	atomic.StoreInt64(&globalMetrics.errors, 0)
}
