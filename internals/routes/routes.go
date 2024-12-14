package routes

import (
	"car-bond/internals/controllers"
	"car-bond/internals/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes func
func SetupRoutes(app *fiber.App) {
	// grouping
	api := app.Group("/api")

	// Auth
	auth := api.Group("/auth")
	auth.Post("/login", controllers.Login)

	// User
	api.Get("/users", controllers.GetAllUsers)
	user := api.Group("/user")
	user.Get("/:id", controllers.GetUser)
	user.Post("/", controllers.CreateUser)
	user.Patch("/:id", middleware.Protected(), controllers.UpdateUser)
	user.Delete("/:id", middleware.Protected(), controllers.DeleteUser)
	
	// Customer
	api.Get("/customers", controllers.GetAllCustomers)
	customer := api.Group("/customer")
	customer.Get("/:id", controllers.GetSingleCustomer)
	customer.Post("/", controllers.CreateCustomer)
	customer.Put("/:id", controllers.UpdateCustomer)
	customer.Delete("/:id", controllers.DeleteCustomerByID)
}
