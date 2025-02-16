package handlers

import (
	"eshop/internal/infrastructure/errs"
	"eshop/internal/transport/http-server/dto"
	"eshop/internal/transport/http-server/dto/requests"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
)

// TODO: handle errors like in login

func (h *Handler) UserLogin(e echo.Context) error {
	logging := h.logger.With(zap.String("Use Case", "Login"))
	var r requests.Login
	if err := e.Bind(&r); err != nil {
		logging.Error("fail on binding request", zap.Error(err))
		resp := dto.NewErrorResponse(errs.ErrBadRequest, "Failed accepting request")
		return e.JSON(http.StatusBadRequest, resp)
	}

	if err := h.usrs.Login(r.ID, r.Username, r.Password); err != nil {
		switch err {
		case errs.ErrHashing:
			logging.Error("hash error", zap.Error(err))
			return e.JSON(http.StatusInternalServerError, dto.NewErrorResponse(errs.ErrHashing, "Failed to work with password"))
		case errs.ErrWrongPassword:
			logging.Error("wrong password", zap.Error(err))
			return e.JSON(http.StatusBadRequest, dto.NewErrorResponse(errs.ErrWrongPassword, "Wrong password"))
		case errs.ErrWrongUsername:
			logging.Error("wrong username", zap.Error(err))
			return e.JSON(http.StatusBadRequest, dto.NewErrorResponse(errs.ErrWrongUsername, "Wrong username"))
		default:
			logging.Error("failed on service operation", zap.Error(err))
			resp := dto.NewErrorResponse(err, "User login failed")
			return e.JSON(http.StatusUnauthorized, resp)
		}
	}

	logging.Info("Logged in user", zap.Any("id", r.ID))
	return e.JSON(http.StatusOK, dto.NewOkReponse("Successfully logged in, id", r.ID))
}

func (h *Handler) UserRegister(e echo.Context) error {
	logging := h.logger.With(zap.String("Use Case", "Register"))

	var r requests.Register
	if err := e.Bind(&r); err != nil {
		logging.Error("fail on binding request", zap.Error(err))
		resp := dto.NewErrorResponse(errs.ErrBadRequest, "Failed accepting request")
		return e.JSON(http.StatusBadRequest, resp)
	}

	id, err := h.usrs.Register(r.Username, r.Password, r.Email, r.IsAdmin)
	if err != nil {
		logging.Error("failed on service operation", zap.Error(err))
		resp := dto.NewErrorResponse(errs.ErrCreatingUser, "User register process failed")
		return e.JSON(http.StatusUnauthorized, resp)
	}

	logging.Info("Registered user", zap.Any("id", id))
	return e.JSON(http.StatusOK, dto.NewOkReponse("Successfully registered with id", id))
}

func (h *Handler) UserDeleteAccount(e echo.Context) error {
	logging := h.logger.With(zap.String("Use Case", "Delete Account"))

	var r requests.DeleteUser
	if err := e.Bind(&r); err != nil {
		logging.Error("fail on binding request", zap.Error(err))
		resp := dto.NewErrorResponse(errs.ErrBadRequest, "Failed accepting request")
		return e.JSON(http.StatusBadRequest, resp)
	}

	id, err := h.usrs.DeleteAccount(r.ID, r.Password)
	if err != nil {
		logging.Error("failed on service operation", zap.Error(err))
		resp := dto.NewErrorResponse(errs.ErrDeletingUser, "Deletion of user account process failed")
		return e.JSON(http.StatusUnauthorized, resp)
	}

	logging.Info("Deleted user account", zap.Any("id", id))
	return e.JSON(http.StatusOK, dto.NewOkReponse("Successfully deleted user's account, id:", id))
}

func (h *Handler) UserUpdate(e echo.Context) error {
	logging := h.logger.With(zap.String("Use Case", "Update User"))

	var req requests.UpdateUserInfo
	if err := e.Bind(&req); err != nil {
		logging.Error("fail on binding request", zap.Error(err))
		resp := dto.NewErrorResponse(errs.ErrBadRequest, "Failed accepting request")
		return e.JSON(http.StatusBadRequest, resp)
	}

	user, err := h.usrs.UpdateInfo(req.ID, req.Username, req.OldPassword, req.NewPassword, req.Email, req.IsAdmin)
	if err != nil {
		logging.Error("DB fail", zap.Error(err))
		return e.JSON(http.StatusInternalServerError, dto.NewErrorResponse(errs.ErrDB, "Failed to update user info"))
	}

	logging.Info("Operation success", zap.Any("new user", user))
	return e.JSON(http.StatusOK, dto.NewOkReponse("user info updated", user))
}

func (h *Handler) GetAllUsers(e echo.Context) error {
	logging := h.logger.With(zap.String("Use Case", "Get All Users"))

	users, err := h.usrs.GetAll()
	if err != nil {
		logging.Error("DB fail", zap.Error(err))
		return e.JSON(http.StatusInternalServerError, dto.NewErrorResponse(errs.ErrDB, "Failed to fetch all users"))
	}

	logging.Info("Operation success", zap.Any("users", users))
	return e.JSON(http.StatusOK, dto.NewOkReponse("all users fetched", users))
}
