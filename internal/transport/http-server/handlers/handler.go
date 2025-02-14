package handlers

import (
	"eshop/internal/services"
	"go.uber.org/zap"
)

// TODO: use err creating user and err updating user

// TODO: add other methods to handler

type Handler interface {
	UserHandler
	ProductHandler
}

type HandlerStruct struct {
	usrs   services.UserService
	prrs   services.ProductService
	logger *zap.Logger
}

func NewHandler(usrs services.UserService, prrs services.ProductService, logger *zap.Logger) Handler {
	return &HandlerStruct{usrs: usrs, logger: logger, prrs: prrs}
}
