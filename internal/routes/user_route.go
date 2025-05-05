package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zenitgarden/simple-chat-api/internal/handler"
	"github.com/zenitgarden/simple-chat-api/internal/middleware"
)

func UserRoutes(app fiber.Router, userHandler *handler.UserHandler) {
	user := app.Group("/users")
	user.Use(middleware.AuthMiddleware())
	user.Patch("/:id", userHandler.UpdateUser)
	user.Get("/:id", userHandler.FindUserByID)
	user.Get("/", userHandler.FindUsers)
}
