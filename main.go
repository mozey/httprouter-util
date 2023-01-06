package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/NYTimes/gziphandler"
	"github.com/alecthomas/units"
	"github.com/julienschmidt/httprouter"
	"github.com/mozey/httprouter-util/internal/handler"
	"github.com/mozey/httprouter-util/pkg/config"
	"github.com/mozey/httprouter-util/pkg/middleware"
	"github.com/mozey/httprouter-util/pkg/response"
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
	h.Router.PanicHandler = middleware.PanicHandler()
	h.Router.NotFound = middleware.NotFound()

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
	h.Router.ServeFiles("/www/*filepath", http.Dir("www"))

	// Client
	h.Router.HandlerFunc("GET", "/client/download", h.ClientDownload)
	h.Router.HandlerFunc("GET", "/client/version", h.ClientVersion)
}

// SetupMiddleware configures the middleware given a route handler
func SetupMiddleware(h *Handler) {
	// Middleware
	var handler http.Handler = h.Router
	// WARNING Allows all origins
	handler = cors.Default().Handler(h.Router)
	maxBytes, err := h.Config.FnMaxBytesKb().Int64()
	if err != nil {
		log.Error().Stack().Err(err).Msg("")
		os.Exit(1)
	}
	handler = middleware.MaxBytes(handler, &middleware.MaxBytesOptions{
		MaxBytes: maxBytes * int64(units.KiB),
	})
	handler = middleware.LogRequest(handler)
	handler = middleware.Logger(handler)
	handler = middleware.Auth(handler, &middleware.AuthOptions{
		Skipper: middleware.AuthSkipper,
	})
	handler = gziphandler.GzipHandler(handler)
	handler = middleware.RequestID(handler)
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// TODO
	log.Info().Msg("Log without context")
	log.Ctx(ctx).Info().Msg("Why is this log not printed to stdout?")
	f, err := os.Open(filepath.Join(h.Config.Dir(), "www", "index.html"))
	if err != nil {
		log.Error().Stack().Err(err).Msg("")
		response.JSON(http.StatusInternalServerError, w, r,
			errors.WithStack(fmt.Errorf("index not found")))
		return
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Error().Stack().Err(err).Msg("")
		response.JSON(http.StatusInternalServerError, w, r,
			errors.WithStack(fmt.Errorf("error reading index")))
		return
	}
	response.Write(http.StatusOK, "", w, r, b)
}

func (h *Handler) Favicon(w http.ResponseWriter, r *http.Request) {
	faviconPath := filepath.Join(h.Config.Dir(), "www", "favicon.ico")
	http.ServeFile(w, r, faviconPath)
}

func (h *Handler) API(w http.ResponseWriter, r *http.Request) {
	// Read body
	_, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response.JSON(
			http.StatusInternalServerError, w, r, errors.WithStack(err))
		return
	}

	response.JSON(http.StatusOK, w, r, response.Response{
		Message: "Welcome",
	})
}

// ClientVersion prints the latest client version
func (h *Handler) ClientVersion(w http.ResponseWriter, r *http.Request) {
	f, err := os.Open(filepath.Join(h.Config.Dir(), "dist", "client.version"))
	if err != nil {
		log.Error().Stack().Err(err).Msg("")
		response.JSON(http.StatusInternalServerError, w, r,
			errors.WithStack(fmt.Errorf("client.version not found")))
		return
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Error().Stack().Err(err).Msg("")
		response.JSON(http.StatusInternalServerError, w, r,
			errors.WithStack(fmt.Errorf("error reading client.version")))
		return
	}
	response.Write(http.StatusOK, "", w, r, b)
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
	response.Write(http.StatusOK, "", w, r,
		[]byte(fmt.Sprintf("hello, %s!\n", params.ByName("name"))))
}

func (h *Handler) NotImplemented(w http.ResponseWriter, r *http.Request) {
	response.JSON(http.StatusNotImplemented, w, r,
		errors.Errorf(http.StatusText(http.StatusNotImplemented)))
}

func main() {
	conf := config.New()

	h, cleanup := CreateRouter(conf)
	defer cleanup()

	// Header to make app reloads more visible
	fmt.Println(".")
	fmt.Println("..")
	fmt.Println("...")
	fmt.Println("....")
	fmt.Println(".....")

	var srv http.Server
	srv.Handler = http.Handler(h.Router)
	srv.Addr = h.Config.Addr()

	// Harden server for Internet exposure
	// https://github.com/mozey/httprouter-util/issues/2
	srv.ReadTimeout = 5 * time.Second
	srv.WriteTimeout = 2 * srv.ReadTimeout
	maxBytes, err := h.Config.FnMaxBytesKb().Int64()
	if err != nil {
		log.Error().Stack().Err(err).Msg("")
		os.Exit(1)
	}
	srv.MaxHeaderBytes = int(maxBytes * int64(units.KiB))

	shutdown := make(chan struct{})
	go func() {
		// "Shutdown gracefully shuts down the server without
		//	interrupting any active connections...
		// 	does not attempt to close nor wait for... WebSockets"
		// https://golang.org/pkg/net/http/#Server.Shutdown
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint
		log.Info().Msg("ctrl+c interrupt, shutting down...")

		// Interrupt signal received, shut down.
		err := srv.Shutdown(context.Background())
		if err != nil {
			// Error from closing listeners, or context timeout
			log.Error().Stack().Err(err).Msg("")
			os.Exit(1)
		}

		close(shutdown)
	}()

	log.Info().Msgf("listening on %s", h.Config.Addr())
	err = errors.WithStack(srv.ListenAndServe())
	if err.Error() != http.ErrServerClosed.Error() {
		log.Error().Stack().Err(err).Msg("")
		os.Exit(1)
	}
	<-shutdown
	log.Info().Msg("bye!")
	os.Exit(0)
}
