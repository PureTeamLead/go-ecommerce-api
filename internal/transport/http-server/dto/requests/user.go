package requests

import "github.com/google/uuid"

type Register struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	IsAdmin  bool   `json:"is_admin"`
}

type Login struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Password string    `json:"password"`
}

type UpdateUserInfo struct {
	ID          uuid.UUID `json:"id"`
	Username    string    `json:"username"`
	OldPassword string    `json:"old_password"`
	NewPassword string    `json:"new_password"`
	Email       string    `json:"email"`
	IsAdmin     bool      `json:"is_admin"`
}

type DeleteUser struct {
	ID       uuid.UUID `json:"id"`
	Password string    `json:"password"`
}
