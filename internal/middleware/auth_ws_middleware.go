package middleware

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func WebSocketAuth(ctx *fiber.Ctx) error {
	tokenString := ctx.Query("token")
	if tokenString == "" {
		return ctx.Status(fiber.StatusUnauthorized).SendString("Missing token")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})

	if err != nil || !token.Valid {
		return ctx.Status(fiber.StatusUnauthorized).SendString("Invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if ok {
		userId := claims["userId"]
		ctx.Locals("userId", userId)
	}

	if exp, ok := claims["exp"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			return ctx.Status(fiber.StatusUnauthorized).SendString("Token expired")
		}
	}

	return ctx.Next()
}
