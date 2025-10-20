package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	jwt.RegisteredClaims
	ID int
}

func BuildJWTToken(id int, secretKey string, tokenLifeTime time.Duration) (string, error) {
	var token *jwt.Token

	token = jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenLifeTime)),
		},
		ID: id,
	})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("error signed string from []byte %w", err)
	}

	return tokenString, nil
}
