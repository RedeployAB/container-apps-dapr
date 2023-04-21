package logger

import (
	"os"

	"github.com/go-logr/logr"
	"github.com/go-logr/zerologr"
	"github.com/rs/zerolog"
)

// Options contains options for the logr.Logger.
type Options struct {
	// Caller sets if calling line of code should be included with the log.
	Caller bool
}

// New returns a new logr.Logger.
func New(options ...Options) logr.Logger {
	return newZerologr(options...)
}

// newZerologr returns a new logr.Logger based on zerolog together with
// zerologr.
func newZerologr(options ...Options) logr.Logger {
	opts := Options{}
	for _, o := range options {
		if o.Caller {
			opts.Caller = o.Caller
		}
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	zerologr.NameFieldName = "logger"
	zerologr.NameSeparator = "/"

	zl := zerolog.New(os.Stderr)
	zl = zl.With().Timestamp().Logger()
	if opts.Caller {
		zl = zl.With().Caller().Logger()
	}
	return zerologr.New(&zl)
}
