package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func RequestLogger(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			next.ServeHTTP(ww, r)

			logger.Info("Request",
				zap.String("method", r.Method),
				zap.String("uri", r.RequestURI),
				zap.Int("status", ww.Status()),
				zap.Int("size", ww.BytesWritten()),
				zap.Duration("duration", time.Since(start)),
			)
		}

		return http.HandlerFunc(fn)
	}
}
