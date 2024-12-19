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

	// Company
	api.Get("/companies", controllers.GetAllCompanies)
	company := api.Group("/company")
	company.Get("/:id", controllers.GetSingleCompanyById)
	company.Post("/", controllers.CreateCompany)
	company.Patch("/:id", middleware.Protected(), controllers.UpdateCompany)
	company.Delete("/:id", middleware.Protected(), controllers.DeleteCompanyById)
	// Company Expenses
	company.Get("/expenses/:companyId", controllers.GetCompanyExpensesByCompanyId)
	company.Get("/expense/:companyId/:id", controllers.GetCompanyExpenseById)
	company.Post("/expense", controllers.CreateCompanyExpense)
	company.Put("/expense/:id", controllers.UpdateCompanyExpense)
	company.Delete("/expense/:id", controllers.DeleteCompanyExpenseById)
	// Company Locations
	company.Get("/location/:companyId/:id", controllers.GetCompanyLocationById)
	company.Post("/location", controllers.CreateCompanyLocation)
	company.Put("/location/:id", controllers.UpdateCompanyLocation)
	company.Delete("/location/:companyId/:id", controllers.DeleteCompanyLocationById)

	// Customer
	api.Get("/customers", controllers.GetAllCustomers)
	customer := api.Group("/customer")
	customer.Get("/:id", controllers.GetSingleCustomer)
	customer.Post("/", controllers.CreateCustomer)
	customer.Put("/:id", controllers.UpdateCustomer)
	customer.Delete("/:id", controllers.DeleteCustomerByID)
}
