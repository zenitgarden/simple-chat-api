package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zenitgarden/simple-chat-api/internal/handler"
	"github.com/zenitgarden/simple-chat-api/internal/middleware"
)

func ConversationRoutes(app fiber.Router, conversationHandler *handler.ConversationHandler) {
	conv := app.Group("/conversations")

	conv.Use(middleware.AuthMiddleware())
	conv.Post("/", conversationHandler.CreateConversation)
	conv.Patch("/:id", conversationHandler.UpdateConversation)
	conv.Get("/latest", conversationHandler.GetLatestConversation)
	conv.Get("/:id", conversationHandler.FindConversationByID)
	conv.Get("/", conversationHandler.FindAllConversations)
	conv.Get("/:id/messages", conversationHandler.GetMessagesByConversationID)
}
