package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zenitgarden/simple-chat-api/internal/handler"
	"github.com/zenitgarden/simple-chat-api/internal/middleware"
)

func AuthRoutes(app fiber.Router, authHandler *handler.AuthHandler) {
	auth := app.Group("/auth")
	auth.Get("/profile", middleware.AuthMiddleware(), authHandler.Profile)
	auth.Post("/login", authHandler.Login)
	auth.Post("/register", authHandler.RegisterUser)
}
