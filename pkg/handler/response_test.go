package handler_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"os"
	"testing"

	"github.com/mozey/httprouter-util/pkg/config"
	"github.com/mozey/httprouter-util/pkg/handler"
	"github.com/mozey/httprouter-util/pkg/share"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
)

func TestJSON(t *testing.T) {
	conf, err := config.LoadFile("dev")
	require.NoError(t, err)
	h := handler.NewHandler(conf)
	defer h.Cleanup()

	// Create request
	req, err := http.NewRequest("GET", "/", nil)
	require.NoError(t, err)

	// Logger
	ctx := req.Context()
	logger := log.With().
		Str("foo", "bar").
		Logger()
	ctx = logger.WithContext(ctx)
	req = req.WithContext(ctx)

	// JSONRaw
	rec := httptest.NewRecorder()
	resp := "{\"a\": 1}"
	h.JSON(http.StatusOK, rec, req, share.JSONRaw(resp))
	require.Equal(t, rec.Code, http.StatusOK, "invalid status code")
	require.Contains(t, rec.Body.String(), resp, "unexpected body")

	// String
	rec = httptest.NewRecorder()
	h.JSON(http.StatusOK, rec, req, "foo")

	if os.Getenv("APP_VERBOSE_TESTS") == "true" {
		dump, err := httputil.DumpResponse(rec.Result(), true)
		require.NoError(t, err)
		fmt.Println(string(dump))
	}

	require.Equal(t, rec.Code, http.StatusOK, "invalid status code")
	require.Contains(t, rec.Body.String(), "foo", "unexpected body")
	require.Equal(t, rec.Header().Get("Content-Type"),
		"application/json; charset=UTF-8")

	// Error
	rec = httptest.NewRecorder()
	err = errors.New("buz")
	h.JSON(http.StatusBadRequest, rec, req, errors.Wrap(err, "fiz"))
	require.Equal(t, rec.Code, http.StatusBadRequest, "invalid status code")
	require.Contains(t, rec.Body.String(), "fiz: buz", "unexpected body")
	require.Equal(t, rec.Header().Get("Content-Type"),
		"application/json; charset=UTF-8")

	// Struct
	type Custom struct {
		Message string `json:"msg"`
		Foo     string `json:"foo"`
	}
	rec = httptest.NewRecorder()
	h.JSON(http.StatusAccepted, rec, req,
		Custom{Message: "baz", Foo: "bar"})
	require.Equal(t, rec.Code, http.StatusAccepted, "invalid status code")
	require.Contains(t, rec.Body.String(),
		"\"msg\": \"baz\"", "unexpected body")
	require.Contains(t, rec.Body.String(),
		"\"foo\": \"bar\"", "unexpected body")
	require.Equal(t, rec.Header().Get("Content-Type"),
		"application/json; charset=UTF-8")
}
