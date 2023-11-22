package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/RedeployAB/container-apps-dapr/endpoint/config"
	"github.com/RedeployAB/container-apps-dapr/endpoint/server"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stderr, nil))

	cfg, err := config.New()
	if err != nil {
		log.Error("Error loading configuration.", "error", err)
		os.Exit(1)
	}

	reporter, err := config.SetupReporter(cfg.Reporter)
	if err != nil {
		log.Error("Error setting up reporter.", "error", err)
		os.Exit(1)
	}

	srv, err := server.New(http.NewServeMux(), server.Options{
		Reporter:     reporter,
		Logger:       log,
		Host:         cfg.Server.Host,
		Port:         cfg.Server.Port,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
		Security: server.Security{
			Keys: cfg.Server.Security.Keys,
		},
	})
	if err != nil {
		log.Error("Error creating server.", "error", err)
	}

	srv.Start()
}
