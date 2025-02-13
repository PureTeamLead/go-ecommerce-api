package services

import (
	"eshop/internal/domain/entities"
)

func UserFromRequest(username, email string, hashedPassword string, isAdmin bool) *entities.User {
	user := entities.NewUser(username, hashedPassword, email, isAdmin)
	return user
}
