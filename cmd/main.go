package main

import (
	"eshop/internal/config"
	logPack "eshop/internal/logger"
	"eshop/internal/storage"
	"flag"
	"go.uber.org/zap"
)

var (
	configPath = flag.String("config", "./config/config.yaml", "specifying path to config.yaml")
)

func main() {
	flag.Parse()
	// If error occurs, crushes the program
	cfg := config.LoadConfig(*configPath)

	logger := logPack.NewLogger(cfg.Env)
	defer logger.Sync()

	logger.Info("Logger is successfully set up",
		zap.String("env", cfg.Env))

	// TODO: make migrations func, methods with db
	db, err := storage.InitDB(cfg.MakeUrlDB())
	if err != nil {
		logger.Fatal("connecting db:", zap.Error(err))
	}

	_ = db

	// TODO: construct router, middlewares, routes
	// TODO: start server

	//httpLogger := logger.With(zap.Int("request_id", generateID()))
}
