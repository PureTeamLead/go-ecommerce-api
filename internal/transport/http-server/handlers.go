package http_server

import (
	"eshop/internal/infrastructure/errs"
	"eshop/internal/services"
	"eshop/internal/transport/http-server/dto"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
)

// TODO: use err creating user and err updating user

// TODO: add other methods to handler

type Handler interface {
	UserLogin(e echo.Context) error
	UserRegister(e echo.Context) error
	UserDeleteAccount(e echo.Context) error
}

type userHandler struct {
	usrs   services.UserService
	logger *zap.Logger
}

func NewUserHandler(usrs services.UserService, logger *zap.Logger) Handler {
	return &userHandler{usrs: usrs, logger: logger}
}

func (uh *userHandler) UserLogin(e echo.Context) error {
	loginLog := uh.logger.With(zap.String("Handler", "Login"))
	var r dto.LoginRequest
	if err := e.Bind(&r); err != nil {
		loginLog.Error("fail on binding request", zap.Error(err))
		resp := dto.NewErrorResponse(errs.ErrBadRequest, "Failed accepting request")
		return e.JSON(http.StatusBadRequest, resp)
	}

	if err := uh.usrs.Login(&r); err != nil {
		loginLog.Error("failed on service operation", zap.Error(err))
		resp := dto.NewErrorResponse(err, "User login failed")
		return e.JSON(http.StatusUnauthorized, resp)
	}

	loginLog.Info("Logged in user", zap.Any("id", r.ID))
	return e.JSON(http.StatusOK, dto.NewOkReponse("Successfully logged in"))
}

func (uh *userHandler) UserRegister(e echo.Context) error {
	registerLog := uh.logger.With(zap.String("Handler", "Register"))

	var r dto.RegisterRequest
	if err := e.Bind(&r); err != nil {
		registerLog.Error("fail on binding request", zap.Error(err))
		resp := dto.NewErrorResponse(errs.ErrBadRequest, "Failed accepting request")
		return e.JSON(http.StatusBadRequest, resp)
	}

	id, err := uh.usrs.Register(&r)
	if err != nil {
		registerLog.Error("failed on service operation", zap.Error(err))
		resp := dto.NewErrorResponse(errs.ErrCreatingUser, "User register process failed")
		return e.JSON(http.StatusUnauthorized, resp)
	}

	registerLog.Info("Registered user", zap.Any("id", id))
	return e.JSON(http.StatusOK, dto.NewOkReponse("Successfully registered"))
}

func (uh *userHandler) UserDeleteAccount(e echo.Context) error {
	deletionAccountLog := uh.logger.With(zap.String("Handler", "DeleteAccount"))

	var r dto.DeleteRequest
	if err := e.Bind(&r); err != nil {
		deletionAccountLog.Error("fail on binding request", zap.Error(err))
		resp := dto.NewErrorResponse(errs.ErrBadRequest, "Failed accepting request")
		return e.JSON(http.StatusBadRequest, resp)
	}

	id, err := uh.usrs.DeleteAccount(&r)
	if err != nil {
		deletionAccountLog.Error("failed on service operation", zap.Error(err))
		resp := dto.NewErrorResponse(errs.ErrDeletingUser, "Deletion of user account process failed")
		return e.JSON(http.StatusUnauthorized, resp)
	}

	deletionAccountLog.Info("Deleted user account", zap.Any("id", id))
	return e.JSON(http.StatusOK, dto.NewOkReponse("Successfully deleted user's account"))
}
