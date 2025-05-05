package middleware

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/zenitgarden/simple-chat-api/pkg/utils"
)

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(utils.Res{
				StatusCode: fiber.StatusUnauthorized,
				Message:    "Missing or invalid token",
				Data:       nil,
			})
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse and verify the token
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
			return []byte(os.Getenv("JWT_SECRET_KEY")), nil
		})
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(utils.Res{
				StatusCode: fiber.StatusUnauthorized,
				Message:    "Invalid or expired token",
				Data:       nil,
			})
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			userId := claims["userId"]
			c.Locals("userId", userId)
		}

		return c.Next()
	}
}
