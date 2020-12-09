package main

import (
	"context"
	"fmt"
	"github.com/mozey/httprouter-example/pkg/response"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/ksuid"
	"net/http"
	"strings"
)

type PanicHandlerFunc func(http.ResponseWriter, *http.Request, interface{})

func PanicHandler() PanicHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, rcv interface{}) {
		err := errors.WithStack(fmt.Errorf("%s", rcv))
		// response.JSON must print stack if resp type is error
		response.JSON(http.StatusInternalServerError, w, r, err)
	}
}

// LogRequest middleware logs details about the request before handling it.
// This is useful when debugging invocations that never get a response,
// e.g. lambda function timeouts
func LogRequestMiddleware(next http.Handler) http.Handler {
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

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Use existing header if available.
		requestID := r.Header.Get(response.HeaderXRequestID)

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
		w.Header().Set(response.HeaderXRequestID, requestID)

		// Also set on context
		ctx = context.WithValue(ctx, response.HeaderXRequestID, requestID)

		// Logger must be set on context for all requests
		// otherwise level is set to "disabled"
		// when calling log.Ctx.
		// With auth failures the LoggerMiddleware won't run,
		// initialise logger on the ctx here in case that happens
		logger := log.With().
			Str("request_id", w.Header().Get(response.HeaderXRequestID)).
			Logger()
		ctx = logger.WithContext(ctx)

		// Call the next handler
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Pass a sub-logger by context
		// https://github.com/rs/zerolog#pass-a-sub-logger-by-context
		var logger zerolog.Logger

		// Set request_id on logger context
		requestID, ok :=
			r.Context().Value(response.HeaderXRequestID).(string)
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

// AuthSkipper lists endpoints that do not require validation,
// return true if auth should be skipped
func AuthSkipper(r *http.Request) bool {
	path := r.URL.Path

	// Rewrite path to skip routes starting with a specified string
	if strings.Index(path, "/proxy") > -1 {
		// Let the proxy endpoint will do it's own auth
		path = "/proxy"
	}
	if strings.Index(path, "/www") > -1 {
		// Static files are public
		path = "/www"
	}

	switch path {
	case
		"/",
		"/index.html",
		"/favicon.ico",
		"/panic",
		"/www",
		"/proxy":
		return true
	}
	return false
}

type AuthOptions struct {
	Skipper func(r *http.Request) bool
}

func AuthMiddleware(next http.Handler, o *AuthOptions) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		if o.Skipper != nil {
			if o.Skipper(r) {
				// Skip auth for this request.
				// Call the next handler
				r = r.WithContext(ctx)
				next.ServeHTTP(w, r)
				return
			}
		}

		// Authenticate
		token := r.URL.Query().Get("token")
		if token != "123" {
			resp := response.ErrResponse{
				Message: "invalid token",
			}
			requestID, ok :=
				r.Context().Value(response.HeaderXRequestID).(string)
			if ok {
				// Set request_id from context
				resp.RequestID = requestID
			}
			response.JSON(http.StatusBadRequest, w, r, resp)
			return
		}

		// Call the next handler
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
