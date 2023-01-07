package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/NYTimes/gziphandler"
	"github.com/alecthomas/units"
	"github.com/julienschmidt/httprouter"
	"github.com/mozey/httprouter-util/pkg/config"
	"github.com/mozey/httprouter-util/pkg/handler"
	"github.com/mozey/httprouter-util/pkg/middleware"
	"github.com/mozey/httprouter-util/pkg/share"
	"github.com/pkg/errors"
	"github.com/rs/cors"
	"github.com/rs/zerolog/log"
)

// Handler for this service
type Handler struct {
	*handler.Handler
}

// NewHandler creates a new top level handler
func NewHandler(conf *config.Config) (h *Handler) {
	h = &Handler{}
	h.Handler = handler.NewHandler(conf)
	return h
}

func CreateRouter(conf *config.Config) (h *Handler, cleanup func()) {
	h = NewHandler(conf)

	// Routes
	h.Routes()

	// Router setup
	h.Router.PanicHandler = middleware.PanicHandler(h.Handler)
	h.Router.NotFound = middleware.NotFound(h.Handler)

	// Middleware
	SetupMiddleware(h)

	return h, h.Cleanup
}

func (h *Handler) Routes() {
	// Routes

	// Index page requires special routes
	h.Router.HandlerFunc("GET", "/", h.Index)
	h.Router.HandlerFunc("GET", "/index.html", h.Index)
	h.Router.HandlerFunc("GET", "/favicon.ico", h.Favicon)

	// Misc
	h.Router.HandlerFunc("GET", "/api", h.API)
	h.Router.HandlerFunc("POST", "/api", h.API)
	h.Router.HandlerFunc("GET", "/panic", h.Panic)
	h.Router.HandlerFunc("GET", "/hello/:name", h.Hello)

	// Static content
	h.Router.ServeFiles("/www/*filepath", http.Dir(
		filepath.Join(h.Config.Dir(), "www")))

	// Client
	h.Router.HandlerFunc("GET", "/client/download", h.ClientDownload)
	h.Router.HandlerFunc("GET", "/client/version", h.ClientVersion)
}

// SetupMiddleware configures the middleware given a route handler
func SetupMiddleware(h *Handler) {
	// Middleware
	var httpHandler http.Handler = h.Router
	// WARNING Allows all origins
	httpHandler = cors.Default().Handler(h.Router)
	maxBytes, err := h.Config.FnMaxBytesKb().Int64()
	if err != nil {
		log.Error().Stack().Err(err).Msg("")
		os.Exit(1)
	}
	httpHandler = middleware.MaxBytes(httpHandler, &middleware.MaxBytesOptions{
		MaxBytes: maxBytes * int64(units.KiB),
	})
	httpHandler = middleware.LogRequest(httpHandler)
	httpHandler = middleware.Logger(httpHandler)
	httpHandler = middleware.Auth(httpHandler, &middleware.AuthOptions{
		Skipper: middleware.AuthSkipper,
	})
	httpHandler = gziphandler.GzipHandler(httpHandler)
	httpHandler = middleware.RequestID(httpHandler)

	h.HTTPHandler = httpHandler
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	f, err := os.Open(filepath.Join(h.Config.Dir(), "www", "index.html"))
	if err != nil {
		h.JSON(http.StatusInternalServerError, w, r,
			errors.WithStack(err))
		return
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		h.JSON(http.StatusInternalServerError, w, r,
			errors.WithStack(err))
		return
	}
	h.Write(http.StatusOK, "", w, r, b)
}

func (h *Handler) Favicon(w http.ResponseWriter, r *http.Request) {
	faviconPath := filepath.Join(h.Config.Dir(), "www", "favicon.ico")
	http.ServeFile(w, r, faviconPath)
}

func (h *Handler) API(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Read body
	b, err := h.GetBody(r)
	if err != nil {
		h.JSON(http.StatusInternalServerError, w, r, err)
		return
	}
	if len(b) > 0 {
		// Body must be valid JSON
		var data map[string]interface{}
		err = json.Unmarshal(b, &data)
		if err != nil {
			h.JSON(http.StatusInternalServerError, w, r, errors.WithStack(err))
			return
		}
		log.Ctx(ctx).Info().Interface("body", data).Msg("")
	}

	h.JSON(http.StatusOK, w, r, share.Response{
		Message: "Welcome",
	})
}

// ClientVersion prints the latest client version
func (h *Handler) ClientVersion(w http.ResponseWriter, r *http.Request) {
	f, err := os.Open(filepath.Join(h.Config.Dir(), "dist", "client.version"))
	if err != nil {
		h.JSON(http.StatusInternalServerError, w, r,
			errors.WithStack(err))
		return
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		h.JSON(http.StatusInternalServerError, w, r,
			errors.WithStack(err))
		return
	}
	h.Write(http.StatusOK, "", w, r, b)
}

// ClientDownload serves the latest client
func (h *Handler) ClientDownload(w http.ResponseWriter, r *http.Request) {
	clientPath := filepath.Join(h.Config.Dir(), "dist", "client")
	http.ServeFile(w, r, clientPath)
}

func (h *Handler) Panic(w http.ResponseWriter, r *http.Request) {
	panic("Oops!")
}

func (h *Handler) Hello(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	h.Write(http.StatusOK, "", w, r,
		[]byte(fmt.Sprintf("hello, %s!\n", params.ByName("name"))))
}

func (h *Handler) NotImplemented(w http.ResponseWriter, r *http.Request) {
	h.JSON(http.StatusNotImplemented, w, r,
		errors.Errorf(http.StatusText(http.StatusNotImplemented)))
}
