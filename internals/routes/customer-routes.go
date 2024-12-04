package routes

import (
	"car-bond/internals/controllers"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes func
func SetupRoutes(app *fiber.App) {
	// grouping
	api := app.Group("/api")
	v1 := api.Group("/customer")
	// routes
	v1.Get("/", controllers.GetAllCustomers)
	v1.Get("/:id", controllers.GetSingleCustomer)
	v1.Post("/", controllers.CreateCustomer)
	v1.Put("/:id", controllers.UpdateCustomer)
	v1.Delete("/:id", controllers.DeleteCustomerByID)
}
