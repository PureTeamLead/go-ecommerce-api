package http_server

import (
	"eshop/internal/transport/http-server/handlers"
	"fmt"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type AppConfig struct {
	Host        string        `yaml:"serv_host" env-required:"true"`
	Port        string        `yaml:"serv_port" env-required:"true"`
	Username    string        `yaml:"serv_username" env-required:"true"`
	Password    string        `yaml:"serv_password" env-required:"true" env:"HTTP_USER_PASSWORD"`
	RWTimeout   time.Duration `yaml:"rw_timeout" env-required:"true"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-required:"true"`
}

type Router struct {
	Handler handlers.Handler
	E       *echo.Echo
	config  AppConfig
	logger  *zap.Logger
	srv     *http.Server
}

func NewRouter(cfg AppConfig, handler handlers.Handler, logger *zap.Logger) *Router {
	e := echo.New()

	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	srv := &http.Server{
		IdleTimeout:  cfg.IdleTimeout,
		ReadTimeout:  cfg.RWTimeout,
		WriteTimeout: cfg.RWTimeout,
		Addr:         addr,
	}

	// TODO: make routes for id, not methods names

	ug := e.Group("/users")
	ug.POST("/login", handler.UserLogin)
	ug.POST("/register", handler.UserRegister)
	ug.DELETE("/delete", handler.UserDeleteAccount)

	//pg := e.Group("/products")
	//pg.POST("/add", AddProductHandler)
	//pg.DELETE("/delete", DeleteProductHandler)
	//pg.PUT("/update", UpdateProductHandler)
	//pg.GET("/{id}", GetProductHandler)

	return &Router{
		E:       e,
		config:  cfg,
		logger:  logger,
		srv:     srv,
		Handler: handler,
	}
}

func (r *Router) Run() {
	r.logger.Fatal("Shutting down the server", zap.Error(r.E.StartServer(r.srv)))
}
