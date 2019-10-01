package main

import (
	"context"
	"errors"
	"gitlab.com/verygoodsoftwarenotvirus/logging/v1/zerolog"
	v1 "gitlab.com/verygoodsoftwarenotvirus/todo/internal/config/v1"
	trace "go.opencensus.io/trace"
	"log"
	"os"
)

func main() {
	// initialize our logger of choice
	logger := zerolog.NewZeroLogger()

	// find and validate our configuration filepath
	configFilepath := os.Getenv("CONFIGURATION_FILEPATH")
	if configFilepath == "" {
		logger.Fatal(errors.New("no configuration file provided"))
	}

	// only allow initialization to take so long
	cfg, err := v1.ParseConfigFile(configFilepath)
	if err != nil || cfg == nil {
		logger.Fatal(err)
	}

	// only allow initialization to take so long
	tctx, cancel := context.WithTimeout(context.Background(), cfg.Meta.StartupDeadline)
	ctx, span := trace.StartSpan(tctx, "initialization")

	// connect to our database
	db, err := cfg.ProvideDatabase(ctx, logger)
	if err != nil {
		logger.Fatal(err)
	}

	// build our server struct
	server, err := BuildServer(ctx, cfg, logger, db)
	span.End()
	cancel()

	if err != nil {
		log.Fatal(err)
	}

	// I slept and dreamt that life was joy.
	//    I awoke and saw that life was service.
	//    	I acted and behold, service deployed.
	server.Serve()
}
