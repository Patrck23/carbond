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

//	@title			Carbond Inventory app
//	@version		1.0
//	@description	This is a sample swagger for Fiber
//	@termsOfService	http://swagger.io/terms/
//	@BasePath		/
func main() {
	database.Connect()
	app := fiber.New(fiber.Config{
		BodyLimit: 20 * 1024 * 1024, // 20 MB
	})
	app.Get("/swagger/*", swagger.HandlerDefault)
	app.Use(logger.New())
	app.Use(cors.New())
	routes.SetupRoutes(app)
	// handle unavailable route
	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(404) // => 404 "Not Found"
	})
	app.Listen(":8080")
}
