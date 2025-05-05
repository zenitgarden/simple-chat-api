package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zenitgarden/simple-chat-api/internal/handler"
	"github.com/zenitgarden/simple-chat-api/internal/middleware"
)

func ParticipantRoutes(app fiber.Router, participanHandler *handler.ParticipantHandler) {
	participant := app.Group("/participants")
	participant.Use(middleware.AuthMiddleware())

	participant.Delete("/:id", participanHandler.DeleteParticipant)
	participant.Post("/", participanHandler.CreateParticipant)
}
