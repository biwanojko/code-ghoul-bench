package server

import "fmt"

// OldHandler is a deprecated handler - dead code (entire file is dead)
type OldHandler struct{}

// Process is a deprecated process method - dead code
func (h *OldHandler) Process(data string) string {
	return fmt.Sprintf("old: %s", data)
}

// legacyInit was used in the old startup sequence - dead code
func legacyInit() error {
	return nil
}

// OldRouter is the old routing implementation - dead code
type OldRouter struct {
	routes map[string]string
}

// NewOldRouter creates a legacy router - dead code
func NewOldRouter() *OldRouter {
	return &OldRouter{routes: make(map[string]string)}
}
