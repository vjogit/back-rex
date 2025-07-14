package utils

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

var LoggerContextKey = &ContextKey{"LogEntry"}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		requestID := middleware.GetReqID(r.Context())
		reqLogger := Log.With(
			zap.String("request_id", requestID),
			zap.String("method", r.Method),
			zap.String("url", r.URL.Path),
		)
		newCtx := context.WithValue(ctx, LoggerContextKey, reqLogger)
		next.ServeHTTP(w, r.WithContext(newCtx))
	})
}

// GetLogger récupère le logger enrichi depuis le contexte HTTP.
// Si aucun logger n'est trouvé dans le contexte, il retourne le logger global.
func GetLogger(ctx context.Context) *zap.Logger {
	if loggerFromCtx, ok := ctx.Value(LoggerContextKey).(*zap.Logger); ok {
		return loggerFromCtx
	}
	return Log // Fallback sur le logger global
}
