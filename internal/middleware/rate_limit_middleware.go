package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/zenitgarden/simple-chat-api/pkg/utils"
)

func RateLimitMiddleware() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        30, // max 30 requests
		Expiration: 1 * time.Minute,
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(utils.Res{
				StatusCode: fiber.StatusTooManyRequests,
				Message:    "Too many requests",
				Data:       nil,
			})
		},
	})
}
