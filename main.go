package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/alecthomas/units"
	"github.com/mozey/httprouter-util/internal/app"
	"github.com/mozey/httprouter-util/pkg/config"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func main() {
	conf := config.New()

	h, cleanup := app.CreateRouter(conf)
	defer cleanup()

	// Header to make app reloads more visible
	fmt.Println(".")
	fmt.Println("..")
	fmt.Println("...")
	fmt.Println("....")
	fmt.Println(".....")

	var srv http.Server
	srv.Handler = h.HTTPHandler
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
