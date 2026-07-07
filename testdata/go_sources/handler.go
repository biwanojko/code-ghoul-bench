package server

import "fmt"

// Handler processes incoming requests
type Handler struct {
	name string
}

// NewHandler creates a new Handler
func NewHandler(name string) *Handler {
	return &Handler{name: name}
}

// Handle processes a request and returns a response
func (h *Handler) Handle(req string) string {
	return fmt.Sprintf("%s: handled %s", h.name, req)
}

// internalProcess is an unexported helper - not called anywhere
func internalProcess(data []byte) []byte {
	return data
}
