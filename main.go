package main

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"runtime/debug"
	"time"
)

var dev bool

type Response struct {
	Message string `json:"message"`
}

// RespondJSON can be used by route handler to respond to requests
func RespondJSON(w http.ResponseWriter, r *http.Request, i interface{}) {
	j, err := json.MarshalIndent(i, "", "    ")
	if err != nil {
		log.Panic().Err(err)
	}
	fmt.Fprint(w, string(j))
}

func Index(w http.ResponseWriter, r *http.Request) {
	RespondJSON(w, r, map[string]interface{}{
		"message": "Welcome",
	})
}

func Panic(w http.ResponseWriter, r *http.Request) {
	panic("Oops!")
}

func Hello(w http.ResponseWriter, r *http.Request) {
	ps := httprouter.ParamsFromContext(r.Context())
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}

func NotFound(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//fmt.Fp
}

func PanicHandler(w http.ResponseWriter, r *http.Request, rcv interface{}) {
	w.WriteHeader(http.StatusInternalServerError)
	RespondJSON(w, r, Response{
		Message: fmt.Sprintf("%s", rcv),
	})
	log.Error().Msg(rcv.(string))
	if dev {
		debug.PrintStack()
	}
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log the request
		log.Info().
			Str("method", r.Method).
			Str("request_uri", string(r.RequestURI))

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Authenticate
		token := r.URL.Query().Get("token")
		if token == "123" {
			// Call the next handler
			next.ServeHTTP(w, r)
			return
		}

		w.WriteHeader(http.StatusUnauthorized)
		msg := "Invalid or missing token"
		RespondJSON(w, r, Response{
			Message: msg,
		})
		log.Error().Msg(msg)
	})
}

func main() {
	dev = true

	// Logger
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if dev {
		log.Logger = log.Output(zerolog.ConsoleWriter{
			Out: os.Stderr,
			NoColor: false,
			TimeFormat: time.RFC3339,
		})
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	// Router
	router := httprouter.New()
	router.PanicHandler = PanicHandler
	router.HandlerFunc("GET", "/", Index)
	router.HandlerFunc("GET", "/panic", Panic)
	router.HandlerFunc("GET", "/hello/:name", Hello)

	// Middleware
	handler := cors.Default().Handler(router)
	handler = LoggingMiddleware(handler)
	handler = AuthMiddleware(handler)

	listen := ":8080"
	if dev {
		// Header to make reloads more visible
		fmt.Println(".")
		fmt.Println(".")
		fmt.Println(".")
		fmt.Println(".")
		fmt.Println(".")
	}
	log.Info().Msgf("listening on %s", listen)
	log.Fatal().Err(http.ListenAndServe(listen, handler))
}

