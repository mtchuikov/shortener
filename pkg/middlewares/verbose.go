package middlewares

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
)

func ChiVerbose(logger zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(rw http.ResponseWriter, req *http.Request) {
			start := time.Now()
			wrw := middleware.NewWrapResponseWriter(rw, req.ProtoMajor)
			next.ServeHTTP(wrw, req)
			duration := time.Since(start)

			event := logger.Info()
			id := req.Context().Value(middleware.RequestIDKey)
			if id != nil {
				event.Str("id", id.(string))
			}

			event.
				Str("from", req.RemoteAddr).Str("method", req.Method).
				Int("status", wrw.Status()).Str("remote", req.RemoteAddr).
				Str("url", req.RequestURI).Int("response bytes", wrw.BytesWritten()).
				Dur("duration", duration).Msg("request completed")
		}

		return http.HandlerFunc(fn)
	}
}
