package server

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/RedeployAB/container-apps-dapr/common/logger"
	"github.com/RedeployAB/container-apps-dapr/endpoint/report"
)

const (
	stopTimeout = time.Second * 10
)

// Defaults.
const (
	defaultPort         = 3000
	defaultReadTimeout  = time.Second * 15
	defaultWriteTimeout = time.Second * 15
	defaultIdleTimeout  = time.Second * 30
)

// log is the interface that wraps around methods Error and Info.
type log interface {
	Error(err error, msg string, keysAndValues ...any)
	Info(msg string, keysAndValues ...any)
}

// router is the interface that wraps around methods Handle and ServeHTTP.
type router interface {
	Handle(pattern string, handler http.Handler)
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

// server represents a server containing a *http.Server, a router (handler) and
// a logger.
type server struct {
	httpServer *http.Server
	router     router
	log        log
	reporter   report.Service
}

// Options for the server.
type Options struct {
	Logger       log
	Reporter     report.Service
	Host         string
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// New returns a new *server with the provided router and Options.
func New(router router, options Options) (*server, error) {
	if options.Reporter == nil {
		return nil, errors.New("reporter is required")
	}
	if options.Port == 0 {
		options.Port = defaultPort
	}
	if options.ReadTimeout == 0 {
		options.ReadTimeout = defaultReadTimeout
	}
	if options.WriteTimeout == 0 {
		options.WriteTimeout = defaultWriteTimeout
	}
	if options.IdleTimeout == 0 {
		options.IdleTimeout = defaultIdleTimeout
	}
	if options.Logger == nil {
		options.Logger = logger.New()
	}

	srv := &http.Server{
		Addr:         options.Host + ":" + strconv.Itoa(options.Port),
		Handler:      router,
		ReadTimeout:  options.ReadTimeout,
		WriteTimeout: options.WriteTimeout,
		IdleTimeout:  options.IdleTimeout,
	}

	return &server{
		router:     router,
		httpServer: srv,
		log:        options.Logger,
		reporter:   options.Reporter,
	}, nil
}

// Start the server.
func (s server) Start() {
	s.routes()
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.log.Error(err, "Server failed to start.")
			os.Exit(1)
		}
	}()
	s.log.Info("Server started.", "type", "server", "address", s.httpServer.Addr)
	sig, err := s.stop()
	if err != nil {
		s.log.Error(err, "Error stopping server.")
	}
	s.log.Info("Server stopped.", "type", "server", "reason", sig.String())
}

// stop server on SIGINT and SIGTERM.
func (s server) stop() (os.Signal, error) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	sig := <-stop

	ctx, cancel := context.WithTimeout(context.Background(), stopTimeout)
	defer cancel()

	s.httpServer.SetKeepAlivesEnabled(false)
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return nil, err
	}
	return sig, nil
}
