package config

import (
	"errors"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/zenitgarden/simple-chat-api/pkg/utils"
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
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
