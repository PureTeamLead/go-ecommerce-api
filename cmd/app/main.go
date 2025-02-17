package main

import (
	"eshop/internal/infrastructure/config"
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
	// TODO: setup Singleton
	// TODO: implement validating passwords/username/emails
	// TODO: write some tests
	// TODO: replace closing db to graceful shutdown

	// TODO: implement seller in database
	// TODO: create aggregates

	// TODO: middlewares
	// TODO: add JWT token

	flag.Parse()
	cfg := config.LoadConfig(*configPath)

	logger := logging.NewLogger(cfg.Env)
	defer logger.Sync()

	logger.Info("Logger is successfully set up",
		zap.String("env", cfg.Env))

	db, err := postgre.NewPostgres(cfg.DB)
	if err != nil {
		logger.Fatal("connecting db:", zap.Error(err))
	}
	defer db.Close()

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
	r.Run()
}
