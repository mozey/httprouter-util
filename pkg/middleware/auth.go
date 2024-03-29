package middleware

import (
	"net/http"
	"strings"

	"github.com/mozey/httprouter-util/pkg/handler"
	"github.com/mozey/httprouter-util/pkg/share"
)

// AuthSkipper lists endpoints that do not require validation,
// return true if auth should be skipped
func AuthSkipper(r *http.Request) bool {
	path := r.URL.Path

	// Rewrite path to skip routes starting with the listed prefixes
	if strings.HasPrefix(path, "/www") {
		// Static files are public
		path = "/www"
	}

	switch path {
	case
		"/",
		"/index.html",
		"/favicon.ico",
		"/panic",
		"/www":
		return true
	}
	return false
}

type AuthOptions struct {
	H       *handler.Handler
	Skipper func(r *http.Request) bool
}

func Auth(next http.Handler, o *AuthOptions) http.Handler {
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
			resp := share.ErrResponse{
				Message: "invalid token",
			}
			requestID, ok :=
				r.Context().Value(share.HeaderXRequestID).(string)
			if ok {
				// Set request_id from context
				resp.RequestID = requestID
			}
			o.H.JSON(http.StatusBadRequest, w, r, resp)
			return
		}

		// Call the next handler
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
