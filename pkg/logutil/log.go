package logutil

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"time"
)

// SetupLogger sets up logging using zerolog.
//
// Wrap new errors with WithStack
//		errors.WithStack(fmt.Errorf("foo"))
//
// Route handlers should not override the stack,
// just return the error.
//
// Errors returned from build-in or third party packages
// can be wrapped using `errors.Wrap` or `.WithStack`.
// Try to wrap errors on the boundaries of this project
//
func SetupLogger(consoleWriter bool) {
	// Prod
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.TimestampFieldName = "created"
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	zerolog.ErrorFieldName = "message"
	zerolog.ErrorStackMarshaler = MarshalStack
	log.Logger = log.With().Caller().Logger()

	if consoleWriter {
		// Dev
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Logger = log.Output(ConsoleWriter{
			Out:           os.Stderr,
			NoColor:       false,
			TimeFormat:    "2006-01-02 15:04:05",
			MarshalIndent: true,
		})
	}
}
