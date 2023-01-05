package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/NYTimes/gziphandler"
	"github.com/alecthomas/units"
	"github.com/julienschmidt/httprouter"
	"github.com/mozey/httprouter-example/pkg/config"
	"github.com/mozey/httprouter-example/pkg/response"
	"github.com/mozey/logutil"
	"github.com/pkg/errors"
	"github.com/rs/cors"
	"github.com/rs/zerolog/log"
)

type Handler struct {
	Config    *config.Config
	Router    *httprouter.Router
	FlushLogs func()
}

func NewHandler(conf *config.Config) (h *Handler) {
	h = &Handler{}
	h.Config = conf
	h.Router = httprouter.New()

	flushLogs, err := SetupLogger(conf)
	if err != nil {
		// Continue if setup of persistance log-writer fails
		log.Error().Stack().Err(err).Msg("")
		h.FlushLogs = func() {}
	} else {
		h.FlushLogs = flushLogs
	}

	return h
}

// SetupLogger sets up the logger to write to console and ...
func SetupLogger(conf *config.Config) (flushLogs func(), err error) {
	// Main logger must always use console writer,
	// for readability inside tmux on both dev and prod
	logutil.SetupLogger(true)

	// TODO Second writer to persist logs?
	return func() {}, nil
}

// Cleanup function must be called before the application exits
func (h *Handler) Cleanup() {
	h.FlushLogs()
}

func CreateRouter(conf *config.Config) (h *Handler, cleanup func()) {
	h = NewHandler(conf)

	// Routes
	h.Routes()

	// Router setup
	h.Router.PanicHandler = PanicHandler()
	h.Router.NotFound = http.HandlerFunc(h.NotFound)

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
		panic(err)
	}
	handler = MaxBytesHandler(handler, &MaxBytesHandlerOptions{
		MaxBytes: maxBytes,
	})
	handler = LogRequestMiddleware(handler)
	handler = LoggerMiddleware(handler)
	handler = AuthMiddleware(handler, &AuthOptions{
		Skipper: AuthSkipper,
	})
	handler = gziphandler.GzipHandler(handler)
	handler = RequestIDMiddleware(handler)
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

func (h *Handler) NotFound(w http.ResponseWriter, r *http.Request) {
	response.JSON(http.StatusNotFound, w, r, response.Response{
		Message: fmt.Sprintf("invalid route %v", r.URL.Path),
	})
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

	// Settings to protect against malicious clients
	srv.ReadTimeout = 5 * time.Second
	srv.WriteTimeout = 2 * srv.ReadTimeout
	srv.MaxHeaderBytes = int(1 * units.KiB)

	log.Info().Msgf("listening on %s", h.Config.Addr())
	err := errors.WithStack(srv.ListenAndServe())
	log.Fatal().Stack().Err(err).Msg("") // Don't override err, use empty msg
}
