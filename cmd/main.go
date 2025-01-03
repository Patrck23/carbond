package main

import (
	"car-bond/internals/database"
	"car-bond/internals/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	database.Connect()
	app := fiber.New(fiber.Config{
		BodyLimit: 20 * 1024 * 1024, // 20 MB
	})
	app.Use(logger.New())
	app.Use(cors.New())
	routes.SetupRoutes(app)
	app.Listen(":8080")
}
