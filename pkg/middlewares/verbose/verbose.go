package verbose

import (
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

// Verbose is an HTTP middleware that logs request and response
// details. It records HTTP method, URL, client IP address,
// user-agent, referer, response status code, duration of request
// processing, and response size. This middleware is useful for
// detailed monitoring and debugging of HTTP requests and their
// handling behavior within the service.
func Verbose(logger *zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(rw http.ResponseWriter, req *http.Request) {
			start := time.Now()
			respData := responseData{
				status: 0,
				size:   0,
			}

			vrw := verboseResponseWriter{
				ResponseWriter: rw,
				responseData:   &respData,
			}

			next.ServeHTTP(&vrw, req)
			duration := time.Since(start)

			logger.Info().
				Str("method", req.Method).
				Str("url", req.URL.String()).
				Str("remote_addr", req.RemoteAddr).
				Str("user_agent", req.UserAgent()).
				Int("status", respData.status).
				Dur("duration", duration).
				Int("size", respData.size).
				Str("referer", req.Referer()).
				Msg("request handled")
		}

		return http.HandlerFunc(hfn)
	}
}
