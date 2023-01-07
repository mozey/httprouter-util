package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/mozey/httprouter-util/internal/app"
	"github.com/mozey/httprouter-util/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestAuth(t *testing.T) {
	conf, err := config.LoadFile("dev")
	require.NoError(t, err)
	h := app.NewHandler(conf)
	h.Routes()
	app.SetupMiddleware(h)
	defer h.Cleanup()

	// Valid token
	u := url.URL{
		Path: "/api",
	}
	q := u.Query()
	q.Set("token", "123")
	u.RawQuery = q.Encode()
	req, err := http.NewRequest("GET", u.String(), nil)
	require.NoError(t, err)

	rec := httptest.NewRecorder()
	h.HTTPHandler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	// Invalid token
	q = url.Values{}
	q.Set("token", "abc")
	u.RawQuery = q.Encode()
	req, err = http.NewRequest("GET", u.String(), nil)
	require.NoError(t, err)

	rec = httptest.NewRecorder()
	h.HTTPHandler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)

	// No token
	q = url.Values{}
	u.RawQuery = q.Encode()
	req, err = http.NewRequest("GET", u.String(), nil)
	require.NoError(t, err)

	rec = httptest.NewRecorder()
	h.HTTPHandler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)

	// Skip auth
	u.Path = "/index.html"
	req, err = http.NewRequest("GET", u.String(), nil)
	require.NoError(t, err)

	rec = httptest.NewRecorder()
	h.HTTPHandler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	u.Path = "/www/index.html"
	req, err = http.NewRequest("GET", u.String(), nil)
	require.NoError(t, err)

	rec = httptest.NewRecorder()
	h.HTTPHandler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusMovedPermanently, rec.Code)
}
