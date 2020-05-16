package main

import (
	"compress/gzip"
	"fmt"
	gh "github.com/gorilla/handlers"
	"github.com/julienschmidt/httprouter"
	"github.com/mozey/httprouter-example/pkg/config"
	"github.com/mozey/httprouter-example/pkg/logutil"
	"github.com/mozey/httprouter-example/pkg/response"
	"github.com/rs/cors"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"net/http"
	"os"
)

type Handler struct {
	Config *config.Config
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	f, err := os.Open("www/index.html")
	if err != nil {
		response.JSON(http.StatusInternalServerError, w, r, response.Response{
			Message: "Index page not found",
		})
		return
	}
	w.WriteHeader(http.StatusOK)
	b, err := ioutil.ReadAll(f)
	if err != nil {
		response.JSON(http.StatusInternalServerError, w, r, response.Response{
			Message: "Error reading index page",
		})
		return
	}
	_, _ = fmt.Fprintf(w, string(b)) // Write string
}

func (h *Handler) Favicon(w http.ResponseWriter, r *http.Request) {
	f, err := os.Open("www/favicon.ico")
	if err != nil {
		response.JSON(http.StatusInternalServerError, w, r, response.Response{
			Message: "favicon not found",
		})
		return
	}
	w.WriteHeader(http.StatusOK)
	b, err := ioutil.ReadAll(f)
	if err != nil {
		response.JSON(http.StatusInternalServerError, w, r, response.Response{
			Message: "Error reading favicon",
		})
		return
	}
	_, _ = w.Write(b) // Write bytes
}

func (h *Handler) API(w http.ResponseWriter, r *http.Request) {
	response.JSON(http.StatusOK, w, r, response.Response{
		Message: "Welcome",
	})
}

func (h *Handler) Panic(w http.ResponseWriter, r *http.Request) {
	panic("Oops!")
}

func (h *Handler) Hello(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	// DANGER Don't use the response helper.
	// Write headers and response right here in the handler
	w.WriteHeader(http.StatusAccepted)
	_, _ = fmt.Fprintf(w, "hello, %s!\n", params.ByName("name"))
}

func (h *Handler) NotFound(w http.ResponseWriter, r *http.Request) {
	response.JSON(http.StatusNotFound, w, r, response.Response{
		Message: fmt.Sprintf("%v not found", r.URL.Path),
	})
}

func (h *Handler) Proxy(w http.ResponseWriter, r *http.Request) {
	response.JSON(http.StatusNotImplemented, w, r, map[string]string{
		"Message": "Not implemented",
		"Proxy": h.Config.Proxy(),
	})
}

func main() {
	// Handler
	h := Handler{}
	h.Config = config.New()

	// Logger
	logutil.SetupLogger(h.Config.Dev() == "true")

	// Router
	router := httprouter.New()
	router.PanicHandler = PanicHandler(&PanicHandlerOptions{
		PrintStack: h.Config.Dev() == "true",
	})
	router.NotFound = http.HandlerFunc(h.NotFound)

	// Routes
	// Index page requires special routes
	router.HandlerFunc("GET", "/", h.Index)
	router.HandlerFunc("GET", "/index.html", h.Index)
	router.HandlerFunc("GET", "/favicon.ico", h.Favicon)
	// Misc
	router.HandlerFunc("GET", "/api", h.API)
	router.HandlerFunc("GET", "/panic", h.Panic)
	router.HandlerFunc("GET", "/hello/:name", h.Hello)
	// TODO Proxy
	router.HandlerFunc("GET", "/proxy", h.Proxy)
	router.HandlerFunc("GET", "/proxy/*filepath", h.Proxy)
	// Static content
	router.ServeFiles("/www/*filepath", http.Dir("www"))

	// Middleware
	var handler http.Handler = router
	// WARNING Allows all origins
	handler = cors.Default().Handler(router)
	handler = LogRequestMiddleware(handler)
	handler = LoggerMiddleware(handler)
	handler = AuthMiddleware(handler, &AuthOptions{
		Skipper: AuthSkipper,
	})
	handler = gh.CompressHandlerLevel(handler, gzip.BestSpeed)
	handler = RequestIDMiddleware(handler)

	if h.Config.Dev() == "true" {
		// Header to make apps more visible on dev
		fmt.Println(".")
		fmt.Println("..")
		fmt.Println("...")
		fmt.Println("....")
		fmt.Println(".....")
	}
	log.Info().Msgf("listening on %s", h.Config.Addr())
	log.Fatal().Err(http.ListenAndServe(h.Config.Addr(), handler))
}
