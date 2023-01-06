package middleware

import (
	"net/http"

	"github.com/rs/zerolog/log"
)

// LogRequest middleware logs details about the request before handling it.
// Useful when debugging invocations that never get a response, e.g. timeouts
func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Set method and request_path on context.
		// Use With instead of UpdateContext
		logger := log.Ctx(ctx).With().
			Str("method", r.Method).
			Str("request_path", r.URL.Path).
			Logger()
		ctx = logger.WithContext(ctx)

		log.Ctx(ctx).Info().
			Str("remote_addr", r.RemoteAddr).
			Msg("")

		// Call the next handler
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
