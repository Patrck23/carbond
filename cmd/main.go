package main

import (
	"car-bond/internals/database"
	"car-bond/internals/routes"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func main() {
	// Create a new Fiber instance
	app := fiber.New(fiber.Config{
		BodyLimit: 20 * 1024 * 1024, // 20 MB
	})

	// Initialize and connect to the database
	db := database.NewDatabase()
	db.Connect()
	defer db.Close() // Ensure database connection is closed when the app shuts down

	// Run migrations and seed data
	db.Migrate()
	db.Seed()

	var sessionStore = session.New()
	// // Configure Redis session storage
	// var sessionStore = session.New(session.Config{
	// 	Storage: redis.New(redis.Config{
	// 		Host: "localhost", // Redis host
	// 		Port: 6379,        // Redis port
	// 	}),
	// })
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
	// Use CORS middleware with specific configuration
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // Allow all origins
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// Setup routes
	routes.SetupRoute(app, db.GetDB())

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-quit
		log.Println("Shutting down server...")
		if err := app.Shutdown(); err != nil {
			log.Fatalf("Error during shutdown: %v", err)
		}
	}()

	// Start the server
	app.Listen(":8080")
}
