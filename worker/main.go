package main

import (
	"os"
	"strconv"

	"github.com/RedeployAB/container-apps-dapr/common/logger"
	"github.com/RedeployAB/container-apps-dapr/worker/config"
	"github.com/RedeployAB/container-apps-dapr/worker/server"
)

func main() {
	log := logger.New()

	cfg, err := config.New()
	if err != nil {
		log.Error(err, "Error loading configuration.")
		os.Exit(1)
	}

	reporter, err := config.SetupReporter(cfg.Storer)
	if err != nil {
		log.Error(err, "Error setting up reporter.")
		os.Exit(1)
	}

	srv, err := server.New(server.Options{
		Reporter: reporter,
		Logger:   log,
		Address:  cfg.Server.Host + ":" + strconv.Itoa(cfg.Server.Port),
		Type:     server.Type(cfg.Server.Type),
		Name:     cfg.Server.Name,
		Topic:    cfg.Server.Topic,
		Queue:    cfg.Server.Queue,
	})
	if err != nil {
		log.Error(err, "Error creating server.")
	}

	srv.Start()
}
