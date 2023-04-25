package config

import (
	"fmt"
	"time"

	"github.com/RedeployAB/container-apps-dapr/worker/report"
	"github.com/caarlos0/env/v6"
)

const (
	defaultHost = "0.0.0.0"
	defaultPort = 3001
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

const (
	storerTypeBlob = "blob"
)

const (
	defaultReporterStorerType    = storerTypeBlob
	defaultReporterStorerName    = "reports"
	defaultReporterStorerTimeout = time.Second * 10
)

// Configuration contains the configuration for the application.
type Configuration struct {
	Server   Server
	Reporter Reporter
}

// Server contains the configuration for the server.
type Server struct {
	Host string `env:"WORKER_HOST"`
	Port int    `env:"WORKER_PORT"`
}

// Reporter contains the configuration for the reporter service.
type Reporter struct {
	Type          string        `env:"WORKER_REPORTER_TYPE"`
	PubsubName    string        `env:"WORKER_REPORTER_PUBSUB_NAME"`
	PubsubTopic   string        `env:"WORKER_REPORTER_PUBSUB_TOPIC"`
	PubsubTimeout time.Duration `env:"WORKER_REPORTER_PUBSUB_TIMEOUT"`
	Storer        Storer
}

type Storer struct {
	Type string `env:"WORKER_REPORTER_STORER_TYPE"`
}

// New creates a new *Configuration based on environment variables
// and default values.
func New() (*Configuration, error) {
	c := &Configuration{
		Server: Server{
			Host: defaultHost,
			Port: defaultPort,
		},
		Reporter: Reporter{
			Type:          defaultReporterType,
			PubsubName:    defaultReporterPubsubName,
			PubsubTopic:   defaultReporterPubsubTopic,
			PubsubTimeout: defaultReporterPubsubTimeout,
			Storer: Storer{
				Type: defaultReporterStorerType,
			},
		},
	}

	if err := env.Parse(c); err != nil {
		return nil, err
	}

	return c, nil
}

// SetupReporter creates a new report.Service based on the provided configuration.
func SetupReporter(c Reporter) (report.Service, error) {
	var err error
	var storer report.Storer
	if c.Storer.Type == storerTypeBlob {
		storer, err = report.NewBlobStorer(report.BlobStorerOptions{
			Name: "reports",
		})
		if err != nil {
			return nil, fmt.Errorf("setup service: %w", err)
		}
	} else {
		return nil, fmt.Errorf("setup service: unknown storer type: %q", c.Storer.Type)
	}

	return report.NewService(storer)
}
