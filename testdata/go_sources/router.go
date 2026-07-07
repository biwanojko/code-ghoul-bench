package server

import "net/http"

// Route represents an HTTP route
type Route struct {
	Method  string
	Path    string
	Handler http.HandlerFunc
}

// Router is a simple HTTP router
type Router struct {
	routes []Route
}

// NewRouter creates a new Router
func NewRouter() *Router {
	return &Router{}
}

// Add adds a route to the router
func (r *Router) Add(method, path string, handler http.HandlerFunc) {
	r.routes = append(r.routes, Route{
		Method:  method,
		Path:    path,
		Handler: handler,
	})
}

// ServeHTTP implements http.Handler
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for _, route := range r.routes {
		if route.Method == req.Method && route.Path == req.URL.Path {
			route.Handler(w, req)
			return
		}
	}
	http.NotFound(w, req)
}

// PrintRoutes prints all registered routes - dead code
func (r *Router) PrintRoutes() {
	for _, route := range r.routes {
		println(route.Method, route.Path)
	}
}
