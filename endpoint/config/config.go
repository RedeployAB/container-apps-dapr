package config

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/RedeployAB/container-apps-dapr/endpoint/report"
	"github.com/caarlos0/env/v8"
)

const (
	defaultHost         = "0.0.0.0"
	defaultPort         = 3000
	defaultReadTimeout  = time.Second * 15
	defaultWriteTimeout = time.Second * 15
	defaultIdleTimeout  = time.Second * 30
)

const (
	reporterTypeQueue  = "queue"
	reporterTypePubsub = "pubsub"
)

const (
	defaultReporterType    = reporterTypeQueue
	defaultReporterName    = "reports"
	defaultReporterTimeout = time.Second * 10
	defaultReporterQueue   = "create"
	defaultReporterTopic   = "create"
)

// Configuration contains the configuration for the application.
type Configuration struct {
	Server   Server
	Reporter Reporter
}

// Server contains the configuration for the server.
type Server struct {
	Security     Security
	Host         string        `env:"ENDPOINT_HOST"`
	Port         int           `env:"ENDPOINT_PORT"`
	ReadTimeout  time.Duration `env:"ENDPOINT_READ_TIMEOUT"`
	WriteTimeout time.Duration `env:"ENDPOINT_WRITE_TIMEOUT"`
	IdleTimeout  time.Duration `env:"ENDPOINT_IDLE_TIMEOUT"`
}

// Security contains the configuration for server security.
type Security struct {
	Keys map[string]struct{} `env:"ENDPOINT_SECURITY_KEYS"`
}

// Reporter contains the configuration for the reporter service.
type Reporter struct {
	Type    string        `env:"ENDPOINT_REPORTER_TYPE"`
	Name    string        `env:"ENDPOINT_REPORTER_NAME"`
	Timeout time.Duration `env:"ENDPOINT_REPORTER_TIMEOUT"`
	Queue   string        `env:"ENDPOINT_REPORTER_QUEUE"`
	Topic   string        `env:"ENDPOINT_REPORTER_TOPIC"`
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
			Type:    defaultReporterType,
			Name:    defaultReporterName,
			Timeout: defaultReporterTimeout,
			Queue:   defaultReporterQueue,
			Topic:   defaultReporterTopic,
		},
	}

	if err := parseEnv(c); err != nil {
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
			Name:    c.Name,
			Topic:   c.Topic,
			Timeout: c.Timeout,
		})
		if err != nil {
			return nil, fmt.Errorf("setup service: %w", err)
		}
	} else if c.Type == reporterTypeQueue {
		r, err = report.NewQueueReporter(report.QueueReporterOptions{
			Name:    c.Name,
			Queue:   c.Queue,
			Timeout: c.Timeout,
		})
		if err != nil {
			return nil, fmt.Errorf("setup service: %w", err)
		}
	} else {
		return nil, fmt.Errorf("setup service: unknown reporter type: %q", c.Type)
	}

	return report.NewService(r)
}

// parseEnv parses the provided value using the env package.
func parseEnv(v any) error {
	return env.ParseWithOptions(v, env.Options{
		FuncMap: map[reflect.Type]env.ParserFunc{
			reflect.TypeOf(map[string]struct{}{}): parseStructMap,
		},
	})
}

// parseStructMap parses the provided comma separated string into a map[string]struct{}.
func parseStructMap(v string) (any, error) {
	parts := strings.Split(strings.ReplaceAll(v, " ", ""), ",")
	m := make(map[string]struct{}, len(parts))
	for _, v := range parts {
		m[v] = struct{}{}
	}
	return m, nil
}
