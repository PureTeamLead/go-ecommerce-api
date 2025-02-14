package services

import (
	"errors"
	"eshop/internal/domain/entities"
	"eshop/internal/infrastructure/constants"
	"eshop/internal/infrastructure/errs"
	"eshop/internal/infrastructure/hashing"
	"eshop/internal/repositories"
	"github.com/google/uuid"
)

//TODO: add methods updateinfo, updatepassword, getallusers

type UserService interface {
	Register(user *entities.User) (uuid.UUID, error)
	Login(user *entities.User) error
	DeleteAccount(user *entities.User) (uuid.UUID, error)
	UpdateInfo(user *entities.User) (*entities.User, error)
	GetAll() ([]entities.User, error)
}

type userService struct {
	ur repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{ur: userRepo}
}

func (u *userService) Register(user *entities.User) (uuid.UUID, error) {
	var id uuid.UUID
	hashedPassword, err := hashing.HashPassword(user.Password)
	if err != nil {
		return constants.EmptyID, errs.ErrHashing
	}

	user.Password = hashedPassword

	id, err = u.ur.Create(user)
	if err != nil {
		return constants.EmptyID, errs.ErrDB
	}

	return id, nil
}

func (u *userService) Login(user *entities.User) error {
	userDB, err := u.ur.GetByID(user.ID)
	if err != nil {
		return errs.ErrDB
	}

	if err = hashing.VerifyPassword(user.Password, userDB.Password); err != nil {
		if errors.Is(err, errs.ErrWrongPassword) {
			return errs.ErrWrongPassword
		}
		return errs.ErrHashing
	}

	if user.Username != userDB.Username {
		return errs.ErrWrongUsername
	}

	return nil
}

func (u *userService) DeleteAccount(user *entities.User) (uuid.UUID, error) {
	userDB, err := u.ur.GetByID(user.ID)
	if err != nil {
		return constants.EmptyID, errs.ErrDB
	}

	if err = hashing.VerifyPassword(user.Password, userDB.Password); err != nil {
		if errors.Is(err, errs.ErrWrongPassword) {
			return constants.EmptyID, errs.ErrWrongPassword
		}
		return constants.EmptyID, errs.ErrHashing
	}

	if err = u.ur.Delete(user.ID); err != nil {
		return constants.EmptyID, errs.ErrDeletingUser
	}

	return user.ID, nil
}

func (u *userService) UpdateInfo(user *entities.User) (*entities.User, error) {
	updatedUser := entities.UpdateUser(user.ID, user.Username, user.Password, user.Email, user.IsAdmin, user.CreatedAt)

	if err := u.ur.Update(updatedUser); err != nil {
		return nil, err
	}

	return updatedUser, nil
}

func (u *userService) GetAll() ([]entities.User, error) {
	users, err := u.ur.GetAll()
	if err != nil {
		return nil, err
	}

	return users, nil
}
