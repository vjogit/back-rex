package utils

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

// ZapLogger est un middleware qui log les détails de la requête HTTP en utilisant zap.
func ZapLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		req := *r // Copie la requête pour ne pas modifier l'original

		ww := &responseWriter{ResponseWriter: w}
		next.ServeHTTP(ww, &req)

		latency := time.Since(start)

		Log.Info("Requête traitée",
			zap.String("method", r.Method),
			zap.String("uri", r.RequestURI),
			zap.String("protocol", r.Proto),
			zap.Int("status", ww.Status()),
			zap.Duration("latency", latency),
			zap.String("user_agent", r.UserAgent()),
			zap.String("remote_addr", r.RemoteAddr),
			zap.String("request_id", middleware.GetReqID(r.Context())),
		)
	})
}

// responseWriter capture le statut de la réponse pour le logging.
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (ww *responseWriter) WriteHeader(status int) {
	ww.status = status
	ww.ResponseWriter.WriteHeader(status)
}

func (ww *responseWriter) Status() int {
	return ww.status
}
