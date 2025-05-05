package config

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/zenitgarden/simple-chat-api/pkg/utils"
)

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found, assuming production environment")
	}
}

func FiberConfig() fiber.Config {
	return fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			// Status code defaults to 500
			code := fiber.StatusInternalServerError
			message := "Internal Server Error"

			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
				message = e.Message
			}

			fmt.Println("Error:", err)

			ctx.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
			return ctx.Status(code).JSON(utils.Res{
				StatusCode: code,
				Message:    message,
				Data:       nil,
			})
		},
	}
}
