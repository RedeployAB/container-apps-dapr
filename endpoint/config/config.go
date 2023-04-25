package config

import (
	"fmt"
	"time"

	"github.com/RedeployAB/container-apps-dapr/endpoint/report"
	"github.com/caarlos0/env/v6"
)

const (
	defaultHost         = "0.0.0.0"
	defaultPort         = 3000
	defaultReadTimeout  = time.Second * 15
	defaultWriteTimeout = time.Second * 15
	defaultIdleTimeout  = time.Second * 30
)

const (
	reporterTypePubsub = "pubsub"
)

const (
	defaultReporterType          = reporterTypePubsub
	defaultReporterPubsubName    = "reports"
	defaultReporterPubsubTopic   = "create"
	defaultReporterPubsubTimeout = time.Second * 10
)

// Configuration contains the configuration for the application.
type Configuration struct {
	Server   Server
	Reporter Reporter
}

// Server contains the configuration for the server.
type Server struct {
	Host         string        `env:"ENDPOINT_HOST"`
	Port         int           `env:"ENDPOINT_PORT"`
	ReadTimeout  time.Duration `env:"ENDPOINT_READ_TIMEOUT"`
	WriteTimeout time.Duration `env:"ENDPOINT_WRITE_TIMEOUT"`
	IdleTimeout  time.Duration `env:"ENDPOINT_IDLE_TIMEOUT"`
}

// Reporter contains the configuration for the reporter service.
type Reporter struct {
	Type          string        `env:"ENDPOINT_REPORTER_TYPE"`
	PubsubName    string        `env:"ENDPOINT_REPORTER_PUBSUB_NAME"`
	PubsubTopic   string        `env:"ENDPOINT_REPORTER_PUBSUB_TOPIC"`
	PubsubTimeout time.Duration `env:"ENDPOINT_REPORTER_PUBSUB_TIMEOUT"`
}

// New creates a new *Configuration based on environment variables
// and default values.
func New() (*Configuration, error) {
	c := &Configuration{
		Server: Server{
			Host:         defaultHost,
			Port:         defaultPort,
			ReadTimeout:  defaultReadTimeout,
			WriteTimeout: defaultWriteTimeout,
			IdleTimeout:  defaultIdleTimeout,
		},
		Reporter: Reporter{
			Type:          defaultReporterType,
			PubsubName:    defaultReporterPubsubName,
			PubsubTopic:   defaultReporterPubsubTopic,
			PubsubTimeout: defaultReporterPubsubTimeout,
		},
	}

	if err := env.Parse(c); err != nil {
		return nil, err
	}

	return c, nil
}

// SetupReporter sets up a new report.Service based on the provided configuration.
func SetupReporter(c Reporter) (report.Service, error) {
	var r report.Reporter
	var err error
	if c.Type == reporterTypePubsub {
		r, err = report.NewPubsubReporter(report.PubsubReporterOptions{
			Name:    c.PubsubName,
			Topic:   c.PubsubTopic,
			Timeout: c.PubsubTimeout,
		})
		if err != nil {
			return nil, fmt.Errorf("setup service: %w", err)
		}
	} else {
		return nil, fmt.Errorf("setup service: unknown reporter type: %q", c.Type)
	}

	return report.NewService(r)
}
