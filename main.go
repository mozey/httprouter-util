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
	Config *config.Config
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
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
	// Handler
	h := Handler{}
	h.Config = config.New()

	// Logger
	logutil.SetupLogger(h.Config.Dev() == "true")

	// Router
	router := httprouter.New()
	router.PanicHandler = PanicHandler()
	router.NotFound = http.HandlerFunc(h.NotFound)

	// Routes
	// Index page requires special routes
	router.HandlerFunc("GET", "/", h.Index)
	router.HandlerFunc("GET", "/index.html", h.Index)
	router.HandlerFunc("GET", "/favicon.ico", h.Favicon)
	// Misc
	router.HandlerFunc("GET", "/api", h.API)
	router.HandlerFunc("POST", "/api", h.API)
	router.HandlerFunc("GET", "/panic", h.Panic)
	router.HandlerFunc("GET", "/hello/:name", h.Hello)
	// TODO Example endpoint to proxy external service,
	// probably better to use Caddy for this?
	// https://github.com/mozey/httprouter-example/issues/6
	router.HandlerFunc("GET", "/proxy", h.NotImplemented)
	router.HandlerFunc("GET", "/proxy/*filepath", h.NotImplemented)
	// TODO Example endpoint for basic auth
	// ...
	// TODO Example endpoint for identity management (e.g. AWS IAM)
	// ...
	// Static content
	router.ServeFiles("/www/*filepath", http.Dir("www"))
	// Client
	router.HandlerFunc("GET", "/client/download", h.ClientDownload)
	router.HandlerFunc("GET", "/client/version", h.ClientVersion)

	// TODO Move max bytes and timeout settings below to config file
	// Middleware
	var handler http.Handler = router
	// WARNING Allows all origins
	handler = cors.Default().Handler(router)
	handler = MaxBytesHandler(handler, &MaxBytesHandlerOptions{
		MaxBytes: int64(1 * units.KiB),
	})
	handler = LogRequestMiddleware(handler)
	handler = LoggerMiddleware(handler)
	handler = AuthMiddleware(handler, &AuthOptions{
		Skipper: AuthSkipper,
	})
	handler = gziphandler.GzipHandler(handler)
	handler = RequestIDMiddleware(handler)

	var srv http.Server
	srv.Handler = handler
	srv.Addr = h.Config.Addr()

	// Settings to protect against malicious clients
	srv.ReadTimeout = 5 * time.Second
	srv.WriteTimeout = 2 * srv.ReadTimeout
	srv.MaxHeaderBytes = int(1 * units.KiB)

	if h.Config.Dev() == "true" {
		// Header to make app reloads more visible on dev
		fmt.Println(".")
		fmt.Println("..")
		fmt.Println("...")
		fmt.Println("....")
		fmt.Println(".....")
	}

	log.Info().Msgf("listening on %s", h.Config.Addr())
	err := errors.WithStack(srv.ListenAndServe())
	log.Fatal().Stack().Err(err).Msg("") // Don't override err, use empty msg
}
