package middleware

import (
	"net/http"

	"github.com/mozey/httprouter-util/pkg/share"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Pass a sub-logger by context
		// https://github.com/rs/zerolog#pass-a-sub-logger-by-context
		var logger zerolog.Logger

		// Set request_id on logger context
		requestID, ok :=
			r.Context().Value(share.HeaderXRequestID).(string)
		if ok {
			logger = log.With().
				Str("request_id", requestID).
				Logger()
		} else {
			logger = log.With().Logger()
		}

		// Logger must be set on context for all requests
		// otherwise level is set to "disabled"
		// when calling log.Ctx
		ctx = logger.WithContext(ctx)

		// Call the next handler
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
