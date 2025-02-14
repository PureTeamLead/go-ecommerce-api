package entities

import (
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

// TODO: implement validation on creating and updating

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
