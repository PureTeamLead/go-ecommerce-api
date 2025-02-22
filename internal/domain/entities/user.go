package entities

import (
	"eshop/internal/infrastructure/errs"
	"github.com/badoux/checkmail"
	"github.com/go-passwd/validator"
	"github.com/google/uuid"
	"time"
)

//Getters, setters, validating, constructor
// Sometimes need to create validated_product and constructor for it

type User struct {
	ID        uuid.UUID
	Username  string
	Password  string
	Email     string
	IsAdmin   bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(username string, password string, email string, isAdmin bool) *User {
	return &User{
		ID:        uuid.New(),
		Username:  username,
		Password:  password,
		Email:     email,
		IsAdmin:   isAdmin,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func UpdateUser(id uuid.UUID, username string, password string, email string, isAdmin bool, createdAt time.Time) *User {
	return &User{
		ID:        id,
		Username:  username,
		Password:  password,
		Email:     email,
		IsAdmin:   isAdmin,
		CreatedAt: createdAt,
		UpdatedAt: time.Now(),
	}
}

func ValidateUser(email string, password string) error {
	if err := checkmail.ValidateFormat(email); err != nil {
		return errs.ErrInvalidEmail
	}

	pswdValidator := validator.New(validator.MinLength(8, nil))
	if err := pswdValidator.Validate(password); err != nil {
		return errs.ErrBadPassword
	}

	return nil
}
