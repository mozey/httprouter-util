package response_test

import (
	"github.com/mozey/httprouter-example/pkg/response"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestJSON(t *testing.T) {
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

	// String
	rec := httptest.NewRecorder()
	response.JSON(http.StatusOK, rec, req, "foo")
	//dump, err := httputil.DumpResponse(rec.Result(), true)
	//require.NoError(t, err)
	//fmt.Println(string(dump))
	require.Equal(t, rec.Code, http.StatusOK, "invalid status code")
	require.Contains(t, rec.Body.String(), "foo", "unexpected body")
	require.Equal(t, rec.Header().Get("Content-Type"),
		"application/json; charset=UTF-8")

	// Error
	rec = httptest.NewRecorder()
	err = errors.New("buz")
	response.JSON(http.StatusBadRequest, rec, req, errors.Wrap(err, "fiz"))
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
	response.JSON(http.StatusAccepted, rec, req,
		Custom{Message: "baz", Foo: "bar"})
	require.Equal(t, rec.Code, http.StatusAccepted, "invalid status code")
	require.Contains(t, rec.Body.String(),
		"\"msg\": \"baz\"", "unexpected body")
	require.Contains(t, rec.Body.String(),
		"\"foo\": \"bar\"", "unexpected body")
	require.Equal(t, rec.Header().Get("Content-Type"),
		"application/json; charset=UTF-8")
}
