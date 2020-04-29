package response

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
	"net/url"
	"regexp"
)

type Response struct {
	Message string `json:"message"`
}

type ErrResponse struct {
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
}

type JSONRaw string

// HeaderXRequestID ...
const HeaderXRequestID = "X-Request-ID"

// JSON can be used by route handlers to respond to requests
func JSON(code int, w http.ResponseWriter, r *http.Request, resp interface{}) {
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

	// NOTE Don't use log.Panic().Err(err).Msg(""),
	// if ErrorFieldName is set to "message" this will override
	// err.Error() with an empty string.
	// Rather call panic(err) and let the PanicHandler do the logging

	// Marshal indented response JSON,
	// uses type switch to handle different resp types
	var b []byte
	var respStr interface{}
	var err error
	indent := "    "
	switch v := resp.(type) {
	case JSONRaw:
		respStr = resp

	case string:
		msg = v
		b, err = json.MarshalIndent(Response{Message: msg}, "", indent)
		if err != nil {
			panic(errors.WithStack(err))
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
		errResp := ErrResponse{Message: msg}
		requestID, ok := r.Context().Value(HeaderXRequestID).(string)
		if ok {
			errResp.RequestID = requestID
		}
		b, err = json.MarshalIndent(errResp, "", indent)
		if err != nil {
			panic(errors.WithStack(err))
		}
		respStr = string(b)

	default:
		b, err = json.MarshalIndent(resp, "", indent)
		if err != nil {
			panic(errors.WithStack(err))
		}
		// Rather than using reflection or type casting,
		// unmarshal the response to determine if a message property was set
		r := Response{}
		err := json.Unmarshal(b, &r)
		if err == nil {
			if r.Message != "" {
				msg = r.Message
			}
		}
		respStr = string(b)
	}

	// Write headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code) // Must be called after w.Header().Set?

	// Some of the properties below are also set in the logrequest middleware,
	// set them again in case this route does not call the middleware
	l := log.Ctx(ctx)
	l.UpdateContext(func(c zerolog.Context) zerolog.Context {
		// Set method and request_path on context for all logs?
		return c.
			Str("method", r.Method).
			Str("request_path", string(r.URL.Path))
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
		Str("request_path", string(r.URL.Path)).
		Str("request_query", query).
		Str("remote_addr", string(r.RemoteAddr)).
		Bool("log_to_es", false).
		Msg(msg)

	// Write response
	_, err = fmt.Fprint(w, respStr)
	if err != nil {
		panic(errors.Wrap(err, "write response"))
	}
}

// Write response bytes with specified code and content type headers
func Write(code int, contentType string, w http.ResponseWriter, r *http.Request, b []byte) {
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
		Str("request_uri", string(r.RequestURI)).
		Msg(http.StatusText(code))

	_, err := fmt.Fprint(w, string(b))
	if err != nil {
		panic(errors.Wrap(err, "write response"))
	}
}
