package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zenitgarden/simple-chat-api/internal/handler"
	"github.com/zenitgarden/simple-chat-api/internal/middleware"
)

func WsRoute(app *fiber.App, handler *handler.WsHandler) {
	app.Get("/ws", middleware.WebSocketAuth, handler.ConnectToWs)
}
