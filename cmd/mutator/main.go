package main

import (
	"context"
	"flag"
	"github.com/gorilla/mux"
	"github.com/project-alvarium/alvarium-sdk-go/pkg"
	"github.com/project-alvarium/alvarium-sdk-go/pkg/factories"
	"github.com/project-alvarium/alvarium-sdk-go/pkg/interfaces"
	"github.com/project-alvarium/ones-demo-2021/internal/bootstrap"
	"github.com/project-alvarium/ones-demo-2021/internal/config"
	"github.com/project-alvarium/ones-demo-2021/internal/db"
	"github.com/project-alvarium/ones-demo-2021/internal/mutator"
	logConfig "github.com/project-alvarium/provider-logging/pkg/config"
	logFactory "github.com/project-alvarium/provider-logging/pkg/factories"
	"github.com/project-alvarium/provider-logging/pkg/logging"
	"os"
)

func main() {
	// Load config
	var configPath string
	flag.StringVar(&configPath,
		"cfg",
		"./res/config.json",
		"Path to JSON configuration file.")
	flag.Parse()

	fileFormat := config.GetFileExtension(configPath)
	reader, err := config.NewReader(fileFormat)
	if err != nil {
		tmpLog := logFactory.NewLogger(logConfig.LoggingInfo{MinLogLevel: logging.ErrorLevel})
		tmpLog.Error(err.Error())
		os.Exit(1)
	}

	cfg := config.ApplicationConfig{}
	err = reader.Read(configPath, &cfg)
	if err != nil {
		tmpLog := logFactory.NewLogger(logConfig.LoggingInfo{MinLogLevel: logging.ErrorLevel})
		tmpLog.Error(err.Error())
		os.Exit(1)
	}

	logger := logFactory.NewLogger(cfg.Logging)
	logger.Write(logging.DebugLevel, "config loaded successfully")
	logger.Write(logging.DebugLevel, cfg.AsString())

	// List of annotators driven from config, eventually support dist. policy.
	var annotators []interfaces.Annotator
	for _, t := range cfg.Sdk.Annotators {
		instance, err := factories.NewAnnotator(t, cfg.Sdk)
		if err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}
		annotators = append(annotators, instance)
	}
	sdk := pkg.NewSdk(annotators, cfg.Sdk, logger)

	// Connect to database
	database, err := db.NewMongoProvider(cfg.Mongo, logger)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	r := mux.NewRouter()
	chMutate := make(chan []byte)
	mutator.LoadRestRoutes(r, chMutate, logger)
	mutate := mutator.NewMutateWorker(sdk, chMutate, cfg.Sdk, cfg.NextHop, database, logger)
	ctx, cancel := context.WithCancel(context.Background())
	bootstrap.Run(
		ctx,
		cancel,
		cfg,
		[]bootstrap.BootstrapHandler{
			mutator.NewHttpServer(r, chMutate, cfg.Endpoint, database, logger).BootstrapHandler,
			sdk.BootstrapHandler,
			database.BootstrapHandler,
			mutate.BootstrapHandler,
		})
}
