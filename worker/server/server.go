package server

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/RedeployAB/container-apps-dapr/worker/report"
	"github.com/dapr/go-sdk/service/common"
	daprd "github.com/dapr/go-sdk/service/grpc"
)

const (
	defaultAddress = "0.0.0.0:3001"
)

// Type is the type of the server.
type Type string

const (
	TypeQueue  Type = "queue"
	TypePubsub Type = "pubsub"
)

const (
	defaultType  = TypeQueue
	defaultName  = "reports"
	defaultQueue = "create"
	defaultTopic = "create"
)

// log is the interface that wraps around methods Error and Info.
type log interface {
	Error(msg string, args ...any)
	Info(msg string, args ...any)
}

// service is the interface that wraps around methods Start, Stop and
// AddTopicEventHandler.
type service interface {
	Start() error
	Stop() error
	AddBindingInvocationHandler(name string, fn common.BindingInvocationHandler) error
	AddTopicEventHandler(sub *common.Subscription, fn common.TopicEventHandler) error
}

// server is the implementation of the Server interface. It contains
// the common.Service from the Dapr SDK.
type server struct {
	service  service
	reporter report.Service
	log      log
	address  string
	name     string
	queue    string
	topic    string
}

// Options for the server.
type Options struct {
	Reporter report.Service
	Logger   log
	Type     Type
	Address  string
	Name     string
	Queue    string
	Topic    string
}

// New creates and returns a server.
func New(options Options) (*server, error) {
	s, err := new(options)
	if err != nil {
		return nil, err
	}

	ds, err := daprd.NewService(options.Address)
	if err != nil {
		return nil, err
	}
	s.service = ds

	if options.Type == TypeQueue {
		if err := s.service.AddBindingInvocationHandler(s.name, s.queueReportHandler); err != nil {
			return nil, errors.New("adding binding handler: " + err.Error())
		}
	} else if options.Type == TypePubsub {
		subscription := &common.Subscription{
			PubsubName: s.name,
			Topic:      s.topic,
			Route:      "/" + s.name,
		}

		if err := s.service.AddTopicEventHandler(subscription, s.pubsubReportHandler); err != nil {
			return nil, errors.New("adding event handler: " + err.Error())
		}
	} else {
		return nil, fmt.Errorf("unsupported type: %v", options.Type)
	}

	return s, nil
}

// new creates and returns a server.
func new(options Options) (*server, error) {
	if options.Reporter == nil {
		return nil, errors.New("reporter is nil")
	}
	if len(options.Address) == 0 {
		options.Address = defaultAddress
	}
	if options.Logger == nil {
		options.Logger = slog.New(slog.NewJSONHandler(os.Stderr, nil))
	}
	if len(options.Name) == 0 {
		options.Name = defaultName
	}
	if len(options.Queue) == 0 {
		options.Queue = defaultQueue
	}
	if len(options.Topic) == 0 {
		options.Topic = defaultTopic
	}

	return &server{
		reporter: options.Reporter,
		log:      options.Logger,
		address:  options.Address,
		name:     options.Name,
		queue:    options.Queue,
		topic:    options.Topic,
	}, nil
}

// Start the server.
func (s server) Start() {
	go func() {
		if err := s.service.Start(); err != nil {
			s.log.Error("Server failed to start.", "error", err)
			os.Exit(1)
		}
	}()
	s.log.Info("Server started.", "type", "server", "address", s.address)
	sig, err := s.stop()
	if err != nil {
		s.log.Error("Error stopping server.", "error", err)
	}
	s.log.Info("Server stopped.", "type", "server", "reason", sig.String())
}

// stop the server.
func (s server) stop() (os.Signal, error) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	sig := <-stop

	if err := s.service.Stop(); err != nil {
		return nil, err
	}
	return sig, nil
}
