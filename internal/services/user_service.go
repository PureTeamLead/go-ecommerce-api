package services

import (
	"errors"
	"eshop/internal/infrastructure/constants"
	"eshop/internal/infrastructure/errs"
	"eshop/internal/infrastructure/hashing"
	"eshop/internal/repositories"
	"eshop/internal/transport/http-server/dto"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

//TODO: add methods updateinfo, updatepassword, getallusers

type UserService interface {
	Register(req *dto.RegisterRequest) (uuid.UUID, error)
	Login(req *dto.LoginRequest) error
	DeleteAccount(req *dto.DeleteRequest) (uuid.UUID, error)
}

type userService struct {
	ur     repositories.UserRepository
	logger *zap.Logger
}

func NewUserService(userRepo repositories.UserRepository, logger *zap.Logger) UserService {
	return &userService{ur: userRepo, logger: logger}
}

func (u *userService) Register(req *dto.RegisterRequest) (uuid.UUID, error) {
	var id uuid.UUID
	regLog := u.logger.With(zap.String("UseCase", "Register Service"))
	hashedPassword, err := hashing.HashPassword(req.Password)
	if err != nil {
		regLog.Error("hashing password:", zap.Error(err))
		return constants.EmptyID, errs.ErrHashing
	}

	newUser := UserFromRequest(req.Username, req.Email, hashedPassword, req.IsAdmin)
	regLog.Info("user created", zap.Any("id", newUser.ID), zap.String("username", newUser.Username))

	id, err = u.ur.Create(newUser)
	if err != nil {
		regLog.Error("DB error", zap.Error(err))
		return constants.EmptyID, errs.ErrDB
	}
	regLog.Info("user successfully added to database", zap.Any("id", newUser.ID))

	return id, nil
}

func (u *userService) Login(req *dto.LoginRequest) error {
	loginLog := u.logger.With(zap.String("UseCase", "Login Service"), zap.Any("id", req.ID))

	user, err := u.ur.GetByID(req.ID)
	if err != nil {
		loginLog.Error("Unable to fetch user from repository", zap.Error(err))
		return errs.ErrDB
	}

	if err = hashing.VerifyPassword(req.Password, user.Password); err != nil {
		if errors.Is(err, errs.ErrWrongPassword) {
			loginLog.Error("User typed in wrong password", zap.String("request password", req.Password))
			return errs.ErrWrongPassword
		}
		loginLog.Error("Server error", zap.Error(err))
		return errs.ErrHashing
	}

	if req.Username != user.Username {
		loginLog.Error("User typed in wrong username", zap.String("request username", req.Username), zap.String("right username", user.Username))
		return errs.ErrWrongUsername
	}

	loginLog.Info("user with successfully logged in", zap.Any("id", req.ID))
	return nil
}

func (u *userService) DeleteAccount(req *dto.DeleteRequest) (uuid.UUID, error) {
	deletionLog := u.logger.With(zap.String("UseCase", "Deletion Service"), zap.Any("id", req.ID))

	user, err := u.ur.GetByID(req.ID)
	if err != nil {
		deletionLog.Error("Unable to fetch user from repository", zap.Error(err))
		return constants.EmptyID, errs.ErrDB
	}

	if err = hashing.VerifyPassword(req.Password, user.Password); err != nil {
		if errors.Is(err, errs.ErrWrongPassword) {
			deletionLog.Error("User typed in wrong password", zap.String("request password", req.Password))
			return constants.EmptyID, errs.ErrWrongPassword
		}
		deletionLog.Error("Server error", zap.Error(err))
		return constants.EmptyID, errs.ErrHashing
	}

	if err = u.ur.Delete(user.ID); err != nil {
		deletionLog.Error("Deleting user error", zap.Error(err))
		return constants.EmptyID, errs.ErrDeletingUser
	}

	deletionLog.Info("User successfully deleted")
	return user.ID, nil
}
