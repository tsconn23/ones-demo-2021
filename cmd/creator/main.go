package main

import (
	"flag"
	"github.com/project-alvarium/ones-demo-2021/internal/config"
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
}
