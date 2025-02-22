package handlers

import (
	"errors"
	"eshop/internal/infrastructure/constants"
	"eshop/internal/infrastructure/errs"
	token "eshop/internal/infrastructure/jwt-token"
	"eshop/internal/transport/http-server/dto"
	"eshop/internal/transport/http-server/dto/requests"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
)

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
		case errs.ErrWrongPassword:
			logging.Error("wrong password", zap.Error(err))
			return e.JSON(http.StatusUnauthorized, dto.NewErrorResponse(errs.ErrWrongPassword, "Wrong password"))
		case errs.ErrWrongUsername:
			logging.Error("wrong username", zap.Error(err))
			return e.JSON(http.StatusUnauthorized, dto.NewErrorResponse(errs.ErrWrongUsername, "Wrong username"))
		default:
			logging.Error("failed on service operation", zap.Error(err))
			resp := dto.NewErrorResponse(err, "User login failed")
			return e.JSON(http.StatusInternalServerError, resp)
		}
	}

	logging.Info("Logged in user", zap.Any("id", r.ID))

	// JWT token
	newToken, err := token.GenerateJWT(r.ID, h.signingKey)
	if err != nil {
		logging.Error("generate token", zap.Error(err))
		resp := dto.NewErrorResponse(err, "Failed generating authentication token")
		return e.JSON(http.StatusInternalServerError, resp)
	}
	logging.Info("Generated token", zap.String("token", newToken))

	authCookie := &http.Cookie{
		Name:     "JWT token",
		HttpOnly: true,
		Value:    newToken,
	}

	e.SetCookie(authCookie)
	logging.Info("Generated cookie and set")

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
		if errors.Is(err, errs.ErrWrongPassword) {
			logging.Error("wrong password", zap.Any("id", id), zap.String("password typed", r.Password))
			resp := dto.NewErrorResponse(err, "Wrong password")
			return e.JSON(http.StatusBadRequest, resp)
		}

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
		switch err {
		case errs.ErrNoUserFound:
			logging.Error("No user found", zap.Any("id", req.ID))
			resp := dto.NewErrorResponse(err, "No user found with provided ID")
			return e.JSON(http.StatusBadRequest, resp)
		case errs.ErrWrongPassword:
			logging.Error("Wrong password", zap.Any("id", req.ID), zap.String("typed password", req.OldPassword), zap.String("real password", user.Password))
			resp := dto.NewErrorResponse(err, "Wrong password")
			return e.JSON(http.StatusUnauthorized, resp)
		default:
			logging.Error("server-side error", zap.Error(err))
			resp := dto.NewErrorResponse(errs.ErrDB, "Server error")
			return e.JSON(http.StatusInternalServerError, resp)
		}
	}

	logging.Info("Operation success", zap.Any("updated user", user))
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

func (h *Handler) CheckJWT(next echo.HandlerFunc) echo.HandlerFunc {
	return func(e echo.Context) error {
		authCookie, err := e.Cookie(constants.CookieJWT)
		if err != nil {
			h.logger.Error("Couldn't get authentication cookie", zap.Error(err))
			resp := dto.NewErrorResponse(err, "Failed to get needed cookie")
			return e.JSON(http.StatusInternalServerError, resp)
		}

		realToken, err := token.ValidateJWT(authCookie.Value, h.GetSigningKey())
		if err != nil {
			h.logger.Error("Wrong JWT token", zap.Error(err))
			return e.Redirect(http.StatusUnauthorized, "/user/login")
		}

		h.logger.Info("User got access to protected route", zap.String("user_id", realToken.Claims.(jwt.MapClaims)["user_id"].(string)))
		return next(e)
	}
}
