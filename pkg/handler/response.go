package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"

	"github.com/mozey/httprouter-util/pkg/share"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// JSON can be used by route handlers to respond to requests
func (h *Handler) JSON(code int, w http.ResponseWriter, r *http.Request, resp interface{}) {
	ctx := r.Context()

	// Default message
	msg := http.StatusText(code)

	// Log request here instead of in middleware,
	// otherwise status code can not be logged.
	// Set log level according to HTTP code
	logEvent := log.Ctx(ctx).Info()
	if code > 299 {
		logEvent = log.Ctx(ctx).Error()
	}

	// Marshal indented response JSON,
	// uses type switch to handle different resp types
	var b []byte
	var respStr interface{}
	var err error
	indent := "    "
	switch v := resp.(type) {
	case share.JSONRaw:
		respStr = resp

	case string:
		msg = v
		b, err = json.MarshalIndent(share.Response{Message: msg}, "", indent)
		if err != nil {
			log.Ctx(ctx).Error().Stack().Err(errors.WithStack(err)).Msg("")
			respStr = err.Error()
			break
		}
		respStr = string(b)

	case error:
		msg = v.Error()
		logEvent = log.Ctx(ctx).Error().Stack().Err(v)
		if code < 400 {
			// The caller should pass in an error code if resp is an error,
			// if not then override code
			code = http.StatusInternalServerError
		}
		errResp := share.ErrResponse{Message: msg}
		requestID, ok := r.Context().Value(share.HeaderXRequestID).(string)
		if ok {
			errResp.RequestID = requestID
		}
		b, err = json.MarshalIndent(errResp, "", indent)
		if err != nil {
			log.Ctx(ctx).Error().Stack().Err(errors.WithStack(err)).Msg("")
			respStr = err.Error()
			break
		}
		respStr = string(b)

	default:
		b, err = json.MarshalIndent(resp, "", indent)
		if err != nil {
			log.Ctx(ctx).Error().Stack().Err(errors.WithStack(err)).Msg("")
			respStr = err.Error()
			break
		}
		// Rather than using reflection or type casting,
		// unmarshal the response to determine if a message property was set
		r := share.Response{}
		err := json.Unmarshal(b, &r)
		if err == nil {
			if r.Message != "" {
				msg = r.Message
			}
		}
		respStr = string(b)
	}

	// Write headers
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(code) // Must be called after w.Header().Set?

	// Some of the properties below are also set in the logrequest middleware,
	// set them again in case this route does not call the middleware
	l := log.Ctx(ctx)
	l.UpdateContext(func(c zerolog.Context) zerolog.Context {
		// Set method and request_path on context for all logs?
		return c.
			Str("method", r.Method).
			Str("request_path", r.URL.Path)
	})

	// Get query and remove token
	query, err := url.QueryUnescape(r.URL.RawQuery)
	if err != nil {
		query = err.Error()
	}
	var re = regexp.MustCompile(`token=\w+`)
	query = re.ReplaceAllString(query, `token=xxx`)

	logEvent.Int("code", code).
		Str("method", r.Method).
		Str("request_path", r.URL.Path).
		Str("request_query", query).
		Str("remote_addr", r.RemoteAddr).
		Msg(msg)

	// Write response
	_, err = fmt.Fprint(w, respStr)
	if err != nil {
		log.Ctx(ctx).Error().Stack().Err(errors.WithStack(err)).Msg("")
	}
}

// Write response bytes with specified code and content type headers
func (h *Handler) Write(code int, contentType string, w http.ResponseWriter, r *http.Request, b []byte) {
	ctx := r.Context()

	if contentType == "" {
		// Default content type
		contentType = "text/html; charset=UTF-8"
	}

	// Write headers
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(code) // Must be called after w.Header().Set?

	// Log request here instead of in middleware,
	// otherwise status code is not logged.
	log.Ctx(ctx).Info().Int("code", code).
		Str("method", r.Method).
		Str("request_uri", r.RequestURI).
		Msg(http.StatusText(code))

	_, err := fmt.Fprint(w, string(b))
	if err != nil {
		log.Ctx(ctx).Error().Stack().Err(errors.WithStack(err)).Msg("")
	}
}
