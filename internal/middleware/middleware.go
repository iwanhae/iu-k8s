package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

// Logger creates a new logger middleware
func Logger() func(next http.Handler) http.Handler {
	return middleware.Logger
}

// Recovery creates a panic recovery middleware
func Recovery() func(next http.Handler) http.Handler {
	return middleware.Recoverer
}

// Timeout creates a timeout middleware
func Timeout(timeout time.Duration) func(next http.Handler) http.Handler {
	return middleware.Timeout(timeout)
}

// RealIP creates a middleware to set real IP
func RealIP() func(next http.Handler) http.Handler {
	return middleware.RealIP
}

// RequestID creates a middleware to add request ID
func RequestID() func(next http.Handler) http.Handler {
	return middleware.RequestID
}
