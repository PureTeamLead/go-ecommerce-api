package main

import (
	"eshop/internal/infrastructure/config"
	"eshop/internal/infrastructure/shutdown"
	"eshop/internal/repositories"
	"eshop/internal/services"
	httpServer "eshop/internal/transport/http-server"
	"eshop/internal/transport/http-server/handlers"
	"eshop/migrations"
	logging "eshop/pkg/logger"
	"eshop/pkg/postgre"
	"flag"
	"go.uber.org/zap"
)

var (
	configPath = flag.String("config", "./config/config.yaml", "specifying path to config.yaml")
)

func main() {

	flag.Parse()
	cfg := config.LoadConfig(*configPath)

	logger := logging.NewLogger(cfg.Env)

	logger.Info("Logger is successfully set up",
		zap.String("env", cfg.Env))

	db, err := postgre.NewPostgres(cfg.DB)
	if err != nil {
		logger.Fatal("connecting db:", zap.Error(err))
	}

	logger.Info("Database is successfully connected")

	if err = migrations.PostgreMigrate(db, migrations.EmbedFS); err != nil {
		logger.Fatal("migrations fault", zap.Error(err))
	}

	userRepository := repositories.NewUserRepository(db)
	productRepository := repositories.NewProductRepository(db)

	userService := services.NewUserService(userRepository)
	productService := services.NewProductService(productRepository)

	handler := handlers.NewHandler(userService, productService, logger, cfg.App.SecretJWT)

	r := httpServer.NewRouter(cfg.App, handler, logger)

	go r.Run()

	shutdown.Stop(db, logger, r.Shutdown)
}
