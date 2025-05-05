package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/zenitgarden/simple-chat-api/internal/entity"
)

func GenerateJWT(user *entity.User) (*string, error) {
	secret := os.Getenv("JWT_SECRET_KEY")

	claims := jwt.MapClaims{
		"userId": user.ID,
		"name":   user.Name,
		"email":  user.Email,
		"exp":    time.Now().Add(time.Hour * 48).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedString, err := token.SignedString([]byte(secret))
	if err != nil {
		return nil, err
	}
	return &signedString, nil
}
