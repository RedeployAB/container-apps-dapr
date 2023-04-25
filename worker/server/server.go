package server

import (
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/RedeployAB/container-apps-dapr/common/logger"
	"github.com/RedeployAB/container-apps-dapr/worker/report"
	"github.com/dapr/go-sdk/service/common"
	daprd "github.com/dapr/go-sdk/service/grpc"
)

const (
	defaultAddress = "0.0.0.0:3001"
)

const (
	defaultPubsubName  = "reports"
	defaultPubsubTopic = "create"
)

// log is the interface that wraps around methods Error and Info.
type log interface {
	Error(err error, msg string, keysAndValues ...any)
	Info(msg string, keysAndValues ...any)
}

// service is the interface that wraps around methods Start, Stop and
// AddTopicEventHandler.
type service interface {
	Start() error
	Stop() error
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
	topic    string
}

// Options for the server.
type Options struct {
	Reporter report.Service
	Logger   log
	Address  string
	Name     string
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

	subscription := &common.Subscription{
		PubsubName: s.name,
		Topic:      s.topic,
		Route:      "/" + s.name,
	}

	if err := s.service.AddTopicEventHandler(subscription, s.reportHandler); err != nil {
		return nil, errors.New("adding event handler: " + err.Error())
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
		options.Logger = logger.New()
	}

	if len(options.Name) == 0 {
		options.Name = defaultPubsubName
	}
	if len(options.Topic) == 0 {
		options.Topic = defaultPubsubTopic
	}

	return &server{
		reporter: options.Reporter,
		log:      options.Logger,
		address:  options.Address,
		name:     options.Name,
		topic:    options.Topic,
	}, nil
}

// Start the server.
func (s server) Start() {
	go func() {
		if err := s.service.Start(); err != nil {
			s.log.Error(err, "Server failed to start.")
			os.Exit(1)
		}
	}()
	s.log.Info("Server started.", "type", "server", "address", s.address)
	sig, err := s.stop()
	if err != nil {
		s.log.Error(err, "Error stopping server.")
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
