package module

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zenitgarden/simple-chat-api/internal/handler"
	"github.com/zenitgarden/simple-chat-api/internal/repository"
	"github.com/zenitgarden/simple-chat-api/internal/routes"
	"github.com/zenitgarden/simple-chat-api/internal/service"
	"gorm.io/gorm"
)

func RegisterModule(app *fiber.App, db *gorm.DB) {

	// user
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	// auth
	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService, userService)

	// participant
	conversationRepo := repository.NewConversationRepository(db)
	participantRepo := repository.NewParticipantRepository(db)
	participantService := service.NewParticipantService(participantRepo, conversationRepo)
	participantHandler := handler.NewParticpantHandler(participantService)

	// message
	messageRepo := repository.NewMessageRepository(db)
	messageService := service.NewMessageService(messageRepo)

	// conversation
	conversationService := service.NewConversationService(conversationRepo, participantRepo)
	conversationHandler := handler.NewConversationHandler(conversationService, messageService)

	// ws
	wsHandler := handler.NewWsHandler(messageService, participantService, userService)

	// router
	api := app.Group("/api")
	routes.AuthRoutes(api, authHandler)
	routes.UserRoutes(api, userHandler)
	routes.ConversationRoutes(api, conversationHandler)
	routes.ParticipantRoutes(api, participantHandler)

	routes.WsRoute(app, wsHandler)
}
