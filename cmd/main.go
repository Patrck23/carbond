package main

import (
	"car-bond/internals/database"
	"car-bond/internals/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func main() {
	database.Connect()
	app := fiber.New(fiber.Config{
		BodyLimit: 20 * 1024 * 1024, // 20 MB
	})
	// Add session middleware
	var sessionStore = session.New()
	app.Use(func(c *fiber.Ctx) error {
		sess, err := sessionStore.Get(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Session error")
		}
		// Save session for reuse in handlers
		c.Locals("session", sess)
		return c.Next()
	})
	app.Use(logger.New())
	app.Use(cors.New())
	routes.SetupRoutes(app)
	app.Listen(":8080")
}
