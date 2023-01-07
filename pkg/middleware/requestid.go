package middleware

import (
	"context"
	"net/http"

	"github.com/mozey/httprouter-util/pkg/share"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/ksuid"
)

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Use existing header if available.
		requestID := r.Header.Get(share.HeaderXRequestID)

		if requestID == "" {
			// Generate new id
			id, err := ksuid.NewRandom()
			if err != nil {
				requestID = err.Error()
			} else {
				requestID = id.String()
			}
		}

		// Set header
		w.Header().Set(share.HeaderXRequestID, requestID)

		// Also set on context
		ctx = context.WithValue(ctx, share.HeaderXRequestID, requestID)

		// Logger must be set on context for all requests
		// otherwise level is set to "disabled"
		// when calling log.Ctx.
		// With auth failures the Logger middleware won't run,
		// initialise logger on the ctx here in case that happens
		logger := log.With().
			Str("request_id", w.Header().Get(share.HeaderXRequestID)).
			Logger()
		ctx = logger.WithContext(ctx)

		// Call the next handler
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
