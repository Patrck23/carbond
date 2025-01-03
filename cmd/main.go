package main

import (
	"car-bond/internals/database"
	"car-bond/internals/routes"

	_ "car-bond/cmd/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/gofiber/swagger"
)

// @title			Carbond Inventory app
// @version			1.0
// @description		This is a sample swagger for the system
// @termsOfService	http://swagger.io/terms/
func main() {
	database.Connect()
	app := fiber.New(fiber.Config{
		BodyLimit: 20 * 1024 * 1024, // 20 MB
	})
	app.Get("/swagger/*", swagger.HandlerDefault)
	app.Use(logger.New())
	app.Use(cors.New())
	routes.SetupRoutes(app)
	app.Listen(":8080")
}
