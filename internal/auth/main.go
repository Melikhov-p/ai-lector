package auth

import (
	"fmt"
	"time"

	"errors"

	"github.com/Melikhov-p/ai-lector/internal/domain/user"
	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID int
	Phone  string
}

func BuildJWTToken(usr *user.User, secretKey string, tokenLifeTime time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenLifeTime)),
		},
		UserID: int(usr.ID()),
		Phone:  usr.Phone(),
	})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("error signed string from []byte %w", err)
	}

	return tokenString, nil
}

// ValidateJWTToken валидирует токен - возвращает userID
func ValidateJWTToken(tokenString, secretKey string) (int64, error) {
	claims := Claims{}

	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected token singing method: %s", token.Method)
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return -1, errors.New("token expired")
		}
		return -1, fmt.Errorf("error parsing token with claims %w", err)
	}

	if !token.Valid {
		return -1, errors.New("invalid token")
	}

	return int64(claims.UserID), nil
}
