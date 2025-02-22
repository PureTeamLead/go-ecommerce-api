package handlers

import (
	"eshop/internal/domain/entities"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type productService interface {
	AddProduct(name, category string, price float64) (uuid.UUID, error)
	DeleteProduct(id uuid.UUID) (uuid.UUID, error)
	UpdateProduct(id uuid.UUID, name, category string, price float64) (*entities.Product, error)
	GetProduct(id uuid.UUID) (*entities.Product, error)
	GetAllProducts() ([]entities.Product, error)
}

type userService interface {
	Register(username, password string, email string, isAdmin bool) (uuid.UUID, error)
	Login(id uuid.UUID, username, password string) error
	DeleteAccount(id uuid.UUID, password string) (uuid.UUID, error)
	UpdateInfo(id uuid.UUID, username string, oldPassword string, newPassword string, email string, isAdmin bool) (*entities.User, error)
	GetAll() ([]entities.User, error)
}

type Handler struct {
	usrs       userService
	prrs       productService
	logger     *zap.Logger
	signingKey string
}

func NewHandler(usrs userService, prrs productService, logger *zap.Logger, jwtToken string) *Handler {
	return &Handler{usrs: usrs, logger: logger, prrs: prrs, signingKey: jwtToken}
}

func (h *Handler) GetSigningKey() string {
	return h.signingKey
}
