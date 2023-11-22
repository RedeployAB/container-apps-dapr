package main

import (
	"log/slog"
	"os"
	"strconv"

	"github.com/RedeployAB/container-apps-dapr/worker/config"
	"github.com/RedeployAB/container-apps-dapr/worker/server"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stderr, nil))

	cfg, err := config.New()
	if err != nil {
		log.Error("Error loading configuration.", "error", err)
		os.Exit(1)
	}

	reporter, err := config.SetupReporter(cfg.Storer)
	if err != nil {
		log.Error("Error setting up reporter.", "error", err)
		os.Exit(1)
	}

	srv, err := server.New(server.Options{
		Reporter: reporter,
		Logger:   log,
		Address:  cfg.Server.Host + ":" + strconv.Itoa(cfg.Server.Port),
		Type:     server.Type(cfg.Server.Type),
		Name:     cfg.Server.Name,
		Queue:    cfg.Server.Queue,
		Topic:    cfg.Server.Topic,
	})
	if err != nil {
		log.Error("Error creating server.", "error", err)
	}

	srv.Start()
}
