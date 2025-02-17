package hashing

import (
	"errors"
	"eshop/internal/infrastructure/constants"
	"eshop/internal/infrastructure/errs"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), constants.HashCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return string(hashedPassword), nil
}

func VerifyPassword(password string, hashedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return errs.ErrWrongPassword
	}

	if err != nil {
		return fmt.Errorf("failed comparing passwords: %w", err)
	}

	return nil
}
