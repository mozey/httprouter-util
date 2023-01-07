package middleware_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/mozey/httprouter-util/internal/app"
	"github.com/mozey/httprouter-util/pkg/config"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
)

// TestGzipHandler, see SetupMiddleware in internal/app/handler.go
func TestGzipHandler(t *testing.T) {
	conf, err := config.LoadFile("dev")
	require.NoError(t, err)
	h := app.NewHandler(conf)
	h.Routes()
	app.SetupMiddleware(h)
	defer h.Cleanup()

	// Download image file.
	// Note that "www/index.html" is not compressed, possibly because the
	// GzipHandler is configured to use gzip.DefaultCompression?
	imagePath := filepath.Join(h.Config.Dir(), "www", "temple.jpg")
	if _, err := os.Open(imagePath); os.IsNotExist(err) {
		out, err := os.Create(imagePath)
		require.NoError(t, err)
		defer out.Close()

		resp, err := http.Get(
			"https://upload.wikimedia.org/wikipedia/commons/4/44/Ankor_Wat_temple.jpg")
		require.NoError(t, err)
		defer resp.Body.Close()

		// Check server response
		require.Equal(t, http.StatusOK, resp.StatusCode)

		// Writer the body to file
		_, err = io.Copy(out, resp.Body)
		require.NoError(t, err)
	}

	// Create request
	u := url.URL{
		Path: "/www/temple.jpg",
	}
	req, err := http.NewRequest("GET", u.String(), nil)
	require.NoError(t, err)
	req.Header.Add("Accept-Encoding", "gzip")

	if os.Getenv("APP_VERBOSE_TESTS") == "true" {
		dump, err := httputil.DumpRequest(req, true)
		require.NoError(t, err)
		fmt.Println(string(dump))
	}

	// Record handler response
	rec := httptest.NewRecorder()
	h.HTTPHandler.ServeHTTP(rec, req)

	// Dump request and response
	if os.Getenv("APP_VERBOSE_TESTS") == "true" {
		log.Info().Str("req.URL.Path", req.URL.Path).
			Msg("Note the path rewrite")
		for k, v := range rec.Result().Header {
			log.Info().Str("key", k).Strs("value", v).Msg("Response header")
		}
		fmt.Println("Content-Encoding", rec.Header().Get("Content-Encoding"))
	}

	// Verify response
	require.Equal(t, rec.Code, http.StatusOK)
	require.Equal(t, "gzip", rec.Header().Get("Content-Encoding"))
}
