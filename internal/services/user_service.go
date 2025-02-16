package services

import (
	"errors"
	"eshop/internal/domain/entities"
	"eshop/internal/infrastructure/constants"
	"eshop/internal/infrastructure/errs"
	"eshop/internal/infrastructure/hashing"
	"github.com/google/uuid"
)

type userRepository interface {
	Create(user *entities.User) (uuid.UUID, error)
	Delete(id uuid.UUID) error
	GetByID(id uuid.UUID) (*entities.User, error)
	GetAll() ([]entities.User, error)
	Update(user *entities.User) error
}

type UserService struct {
	ur userRepository
}

func NewUserService(userRepo userRepository) *UserService {
	return &UserService{ur: userRepo}
}

func (u *UserService) Register(username, password string, email string, isAdmin bool) (uuid.UUID, error) {
	var id uuid.UUID
	hashedPassword, err := hashing.HashPassword(password)
	if err != nil {
		return constants.EmptyID, errs.ErrHashing
	}

	newUser := entities.NewUser(username, hashedPassword, email, isAdmin)

	id, err = u.ur.Create(newUser)
	if err != nil {
		return constants.EmptyID, errs.ErrDB
	}

	return id, nil
}

func (u *UserService) Login(id uuid.UUID, username, password string) error {
	userDB, err := u.ur.GetByID(id)
	if err != nil {
		return errs.ErrDB
	}

	if err = hashing.VerifyPassword(password, userDB.Password); err != nil {
		if errors.Is(err, errs.ErrWrongPassword) {
			return errs.ErrWrongPassword
		}
		return errs.ErrHashing
	}

	if username != userDB.Username {
		return errs.ErrWrongUsername
	}

	return nil
}

func (u *UserService) DeleteAccount(id uuid.UUID, password string) (uuid.UUID, error) {
	userDB, err := u.ur.GetByID(id)
	if err != nil {
		return constants.EmptyID, errs.ErrDB
	}

	if err = hashing.VerifyPassword(password, userDB.Password); err != nil {
		return constants.EmptyID, err
	}

	if err = u.ur.Delete(id); err != nil {
		return constants.EmptyID, errs.ErrDeletingUser
	}

	return id, nil
}

func (u *UserService) UpdateInfo(id uuid.UUID, username string, oldPassword string, newPassword string, email string, isAdmin bool) (*entities.User, error) {
	userDB, err := u.ur.GetByID(id)
	if err != nil {
		return nil, err
	}

	if err = hashing.VerifyPassword(oldPassword, userDB.Password); err != nil {
		return nil, err
	}

	updatedUser := entities.UpdateUser(id, username, newPassword, email, isAdmin, userDB.CreatedAt)

	if err = u.ur.Update(updatedUser); err != nil {
		return nil, err
	}

	return updatedUser, nil
}

func (u *UserService) GetAll() ([]entities.User, error) {
	users, err := u.ur.GetAll()
	if err != nil {
		return nil, err
	}

	return users, nil
}
