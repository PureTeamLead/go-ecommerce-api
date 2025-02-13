package main

import (
	"eshop/internal/infrastructure/config"
	logging "eshop/internal/infrastructure/logger"
	"eshop/internal/repositories"
	"eshop/internal/services"
	httpServer "eshop/internal/transport/http-server"
	"eshop/migrations"
	"eshop/pkg/postgre"
	"flag"
	"go.uber.org/zap"
)

var (
	configPath = flag.String("config", "./config/config.yaml", "specifying path to config.yaml")
)

func main() {
	// TODO: add some getters and setters(if needed)
	// TODO: setup Singleton
	// TODO: implement validating passwords/username/emails
	// TODO: write some tests
	// TODO: use logger only on high level -> delete from repositories, services
	// TODO: write graceful shutdown

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
	// TODO: replace closing db to graceful shutdown
	defer db.Close()

	logger.Info("Database is successfully connected")

	if err = migrations.PostgreMigrate(db, migrations.EmbedFS); err != nil {
		logger.Fatal("migrations fault", zap.Error(err))
	}

	// TODO: implement seller in database
	userRepository := repositories.NewUserRepository(db)
	//productRepository := repositories.NewProductRepository(db)

	userService := services.NewUserService(userRepository, logger)
	//productService := services.NewProductService(productRepository, logger)

	userHandler := httpServer.NewUserHandler(userService, logger)
	// TODO: construct router, middlewares, routes
	// TODO: add JWT token
	r := httpServer.NewRouter(cfg.App, userHandler, logger)
	r.Run()
}
