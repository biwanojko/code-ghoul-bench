package server

import "net/http"

// Middleware is an HTTP middleware function
type Middleware func(http.Handler) http.Handler

// Chain chains multiple middlewares together
func Chain(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next
	}
}

// LoggingMiddleware logs each request
func LoggingMiddleware(l *Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			l.Log(LogInfo, r.URL.Path)
			next.ServeHTTP(w, r)
		})
	}
}

// AuthMiddleware checks for an auth token - dead code (never instantiated)
func AuthMiddleware(token string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Authorization") != token {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// RateLimitMiddleware applies rate limiting - dead code
func RateLimitMiddleware(rps int) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// TODO: implement rate limiting
			_ = rps
			next.ServeHTTP(w, r)
		})
	}
}
