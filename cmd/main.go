package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/zenitgarden/simple-chat-api/internal/config"
	"github.com/zenitgarden/simple-chat-api/internal/module"
	"github.com/zenitgarden/simple-chat-api/pkg/logger"
)

func main() {
	config.LoadEnv()
	app := fiber.New(config.FiberConfig())
	config.ConnectDB()

	app.Use(cors.New(cors.Config{
		AllowOrigins:  os.Getenv("ALLOWED_ORIGINS"),
		AllowMethods:  "GET,POST,PUT,DELETE,PATCH",
		AllowHeaders:  "Origin, Content-Type, Accept, Authorization",
		ExposeHeaders: "Content-Length",
	}))
	// app.Use(middleware.RateLimitMiddleware())
	app.Use(recover.New())
	module.RegisterModule(app, config.DB)

	logger.StartConnectionLogger()

	log.Fatal(app.Listen(":3000"))
}
