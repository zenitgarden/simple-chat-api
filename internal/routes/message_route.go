package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zenitgarden/simple-chat-api/internal/handler"
	"github.com/zenitgarden/simple-chat-api/internal/middleware"
)

func MessageRoute(app fiber.Router, messageHandler *handler.MessageHandler) {
	message := app.Group("/messages")
	message.Use(middleware.AuthMiddleware())

	// message.Get("/", messageHandler.FindMessagesByConversationID)
}
