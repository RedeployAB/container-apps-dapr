package config

import (
	"fmt"
	"time"

	"github.com/RedeployAB/container-apps-dapr/worker/report"
	"github.com/caarlos0/env/v10"
)

const (
	defaultHost = "0.0.0.0"
	defaultPort = 3001
)

const (
	typeServerQueue  = "queue"
	typeServerPubsub = "pubsub"
)

const (
	defaultType  = typeServerQueue
	defaultName  = "reports"
	defaultQueue = "create"
	defaultTopic = "create"
)

const (
	storerTypeBlob = "blob"
)

const (
	defaultStorerType    = storerTypeBlob
	defaultStorerName    = "reports-output"
	defaultStorerTimeout = time.Second * 10
)

// Configuration contains the configuration for the application.
type Configuration struct {
	Server Server
	Storer Storer
}

// Server contains the configuration for the server.
type Server struct {
	Host  string `env:"WORKER_HOST"`
	Port  int    `env:"WORKER_PORT"`
	Type  string `env:"WORKER_TYPE"`
	Name  string `env:"WORKER_NAME"`
	Queue string `env:"WORKER_QUEUE"`
	Topic string `env:"WORKER_TOPIC"`
}

// Storer contains the configuration for the storer.
type Storer struct {
	Type    string        `env:"WORKER_STORER_TYPE"`
	Name    string        `env:"WORKER_STORER_NAME"`
	Timeout time.Duration `env:"WORKER_STORER_TIMEOUT"`
}

// New creates a new *Configuration based on environment variables
// and default values.
func New() (*Configuration, error) {
	c := &Configuration{
		Server: Server{
			Host:  defaultHost,
			Port:  defaultPort,
			Type:  defaultType,
			Name:  defaultName,
			Queue: defaultQueue,
			Topic: defaultTopic,
		},
		Storer: Storer{
			Type:    defaultStorerType,
			Name:    defaultStorerName,
			Timeout: defaultStorerTimeout,
		},
	}

	if err := env.Parse(c); err != nil {
		return nil, err
	}

	return c, nil
}

// SetupReporter creates a new report.Service based on the provided configuration.
func SetupReporter(c Storer) (report.Service, error) {
	var err error
	var storer report.Storer
	if c.Type == storerTypeBlob {
		storer, err = report.NewBlobStorer(func(o *report.BlobStorerOptions) {
			o.Name = c.Name
			o.Timeout = c.Timeout
		})
		if err != nil {
			return nil, fmt.Errorf("setup service: %w", err)
		}
	} else {
		return nil, fmt.Errorf("setup service: unknown storer type: %q", c.Type)
	}

	return report.NewService(storer)
}
