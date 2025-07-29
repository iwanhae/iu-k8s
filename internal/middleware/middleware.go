package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"iu-k8s.linecorp.com/server/internal/log"
)

var (
	RealIP    = middleware.RealIP
	RequestID = middleware.RequestID
	Recovery  = middleware.Recoverer
	GetReqID  = middleware.GetReqID
)

// Logger creates a new logger middleware
func Logger(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		reqID := GetReqID(r.Context())
		logger := slog.With("req_id", reqID, "component", "http")
		attrs := []any{
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("remote_addr", r.RemoteAddr),
			slog.String("user_agent", r.UserAgent()),
		}

		logger.Debug("request received", attrs...)
		rCtx := r.WithContext(log.With(r.Context(), logger))

		rw := &responseWriter{ResponseWriter: w}
		now := time.Now()
		next.ServeHTTP(rw, rCtx)
		attrs = append(attrs,
			slog.Duration("duration", time.Since(now)),
			slog.Int("status", rw.statusCode),
		)
		logger.Info("request", attrs...)
	}
	return http.HandlerFunc(fn)
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
