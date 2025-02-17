package jwt

import (
	"eshop/internal/infrastructure/constants"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"time"
)

func GenerateJWT(id uuid.UUID, signingKey string) (string, error) {

	claims := jwt.MapClaims{
		"user_id": id,
		"exp":     time.Now().Add(constants.ExpTime).Unix(),
	}

	// encoded
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(signingKey)
}

func ValidateJWT(tokenGot string) (uuid.UUID, error) {
	parsedToken, err := jwt.Parse(tokenGot, func(token *jwt.Token) (interface{}, error) {
		// Get from the func signing key
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("wrong signing method on JWT token")
		}

		return
	})
}
