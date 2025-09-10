package main

import (
	"log"

	"github.com/Himany/gofermart/internal/logger"
	"github.com/Himany/gofermart/internal/server"

	"go.uber.org/zap"
)

func main() {
	cfg, err := parseFlags()
	if err != nil {
		log.Fatal("failed to initialize flags: " + err.Error())
	}

	if err := logger.Initialize(cfg.LogLevel); err != nil {
		log.Fatal("failed to initialize logger: " + err.Error())
	}

	logger.Log.Info("flags", zap.Object("config", cfg))

	if err := server.Run(cfg); err != nil {
		logger.Log.Fatal("main", zap.Error(err))
	}
}
