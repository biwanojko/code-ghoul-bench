package server

import (
	"fmt"
	"net/http"
)

// Server is an HTTP server
type Server struct {
	host string
	port int
	mux  *http.ServeMux
}

// NewServer creates a new Server
func NewServer(host string, port int) *Server {
	return &Server{
		host: host,
		port: port,
		mux:  http.NewServeMux(),
	}
}

// Start starts the server
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	return http.ListenAndServe(addr, s.mux)
}

// Stop stops the server (not yet implemented - dead code)
func (s *Server) Stop() {
	// TODO: implement graceful shutdown
}

// RegisterHandler registers a handler for a path
func (s *Server) RegisterHandler(path string, h *Handler) {
	s.mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, h.Handle(r.URL.Path))
	})
}

// legacyStart is unused - dead code
func legacyStart(addr string) error {
	return http.ListenAndServe(addr, nil)
}
