package dto

import "github.com/google/uuid"

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	IsAdmin  bool   `json:"isAdmin"`
}

type LoginRequest struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Password string    `json:"password"`
}

type DeleteRequest struct {
	ID       uuid.UUID `json:"id"`
	Password string    `json:"password"`
}
