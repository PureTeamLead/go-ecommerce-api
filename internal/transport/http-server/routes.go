package http_server

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type userHandler interface {
	UserLogin(e echo.Context) error
	UserRegister(e echo.Context) error
	UserDeleteAccount(e echo.Context) error
	GetAllUsers(e echo.Context) error
	UserUpdate(e echo.Context) error
	CheckJWT(next echo.HandlerFunc) echo.HandlerFunc
}

type productHandler interface {
	AddProduct(e echo.Context) error
	DeleteProduct(e echo.Context) error
	UpdateProductInfo(e echo.Context) error
	GetProduct(e echo.Context) error
	GetAllProducts(e echo.Context) error
}

type handlerAbs interface {
	userHandler
	productHandler
}

type AppConfig struct {
	Host        string        `yaml:"serv_host" env-required:"true"`
	Port        string        `yaml:"serv_port" env-required:"true"`
	Username    string        `yaml:"serv_username" env-required:"true"`
	Password    string        `env-required:"true" env:"HTTP_USER_PASSWORD"`
	RWTimeout   time.Duration `yaml:"rw_timeout" env-required:"true"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-required:"true"`
	SecretJWT   string        `env:"SECRET_JWT_KEY"`
}

type Router struct {
	Handler handlerAbs
	E       *echo.Echo
	config  AppConfig
	logger  *zap.Logger
	srv     *http.Server
}

func NewRouter(cfg AppConfig, handler handlerAbs, logger *zap.Logger) *Router {
	e := echo.New()

	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	srv := &http.Server{
		IdleTimeout:  cfg.IdleTimeout,
		ReadTimeout:  cfg.RWTimeout,
		WriteTimeout: cfg.RWTimeout,
		Addr:         addr,
	}

	ug := e.Group("/user")
	ug.POST("/login", handler.UserLogin)
	ug.POST("/register", handler.UserRegister)
	ug.DELETE("/delete", handler.UserDeleteAccount, handler.CheckJWT)
	ug.GET("/all", handler.GetAllUsers)
	ug.PUT("/update", handler.UserUpdate, handler.CheckJWT)

	e.GET("/getProduct", handler.GetProduct)

	pg := e.Group("/product", handler.CheckJWT)
	pg.POST("/add", handler.AddProduct)
	pg.DELETE("/delete", handler.DeleteProduct)
	pg.PUT("/update", handler.UpdateProductInfo)
	pg.GET("/all", handler.GetAllProducts)

	return &Router{
		E:       e,
		config:  cfg,
		logger:  logger,
		srv:     srv,
		Handler: handler,
	}
}

func (r *Router) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := r.E.Shutdown(ctx); err != nil {
		r.logger.Error("server shutdown failed", zap.Error(err))
	}
}

func (r *Router) Run() {
	err := r.E.StartServer(r.srv)
	if err != nil {
		r.logger.Error("Failed server serving", zap.Error(err))
		return
	}
}
