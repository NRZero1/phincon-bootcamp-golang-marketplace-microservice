package utils

import (
	"os"
	"time"
	"user_service/internal/domain"
	"user_service/internal/domain/dto/response"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJwtToken(user response.LoginResponse) (string, error) {
	claims := domain.Claims{
		UserID:   user.UserID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 5)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secretKey := []byte(os.Getenv("SECRET"))
	tokenString, err := token.SignedString(secretKey)

	if err != nil {
		return "", ErrCreateSignature
	}
	return tokenString, nil
}
