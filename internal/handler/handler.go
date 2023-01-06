package handler

import (
	"github.com/julienschmidt/httprouter"
	"github.com/mozey/httprouter-util/pkg/config"
	"github.com/mozey/logutil"
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
