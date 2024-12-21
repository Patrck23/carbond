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
	api.Get("/expenses", controllers.GetAllExpenses)
	company.Get("/expenses/:companyId", controllers.GetCompanyExpensesByCompanyId)
	company.Get("/:companyId/expense/:id", controllers.GetCompanyExpenseById)
	company.Post("/expense", controllers.CreateCompanyExpense)
	company.Put("/expense/:id", controllers.UpdateCompanyExpense)
	company.Delete("/expense/:id", controllers.DeleteCompanyExpenseById)
	// Company Locations
	api.Get("/locations", controllers.GetAllLocations)
	company.Get("/locations/:companyId", controllers.GetAllCompanyLocations)
	company.Get("/:companyId/location/:id", controllers.GetCompanyLocationById)
	company.Post("/location", controllers.CreateCompanyLocation)
	company.Put("/location/:id", controllers.UpdateCompanyLocation)
	company.Delete("/:companyId/location/:id", controllers.DeleteCompanyLocationById)

	// Customer
	api.Get("/customers", controllers.GetAllCustomers)
	customer := api.Group("/customer")
	customer.Get("/:id", controllers.GetSingleCustomer)
	customer.Post("/", controllers.CreateCustomer)
	customer.Put("/:id", controllers.UpdateCustomer)
	customer.Delete("/:id", controllers.DeleteCustomerByID)
	// Customer contact
	api.Get("/contacts", controllers.GetAllContacts)
	customer.Get("/contacts/:customerId", controllers.GetCustomerContacts)
	customer.Get("/:customerId/contact/:id", controllers.GetCustomerContactById)
	customer.Post("/contact", controllers.CreateCustomerContact)
	customer.Put("/contact/:id", controllers.UpdateCustomerContact)
	customer.Delete("/:customerId/contact/:id", controllers.DeleteCustomerContactById)
	// Customer address
	api.Get("/addresses", controllers.GetAllAddresses)
	customer.Get("/addresses/:customerId", controllers.GetCustomerAddresses)
	customer.Get("/:customerId/address/:id", controllers.GetCustomerAddressById)
	customer.Post("/address", controllers.CreateCustomerAddress)
	customer.Put("/address/:id", controllers.UpdateCustomerAddress)
	customer.Delete("/:customerId/address/:id", controllers.DeleteCustomerAddressById)

	// Car
	api.Get("/cars", controllers.GetAllCars)
	car := api.Group("/customer")
	car.Get("/:id", controllers.GetSingleCar)
	car.Post("/", controllers.CreateCar)
	car.Put("/:id", controllers.UpdateCar)
	car.Delete("/:id", controllers.DeleteCarByID)
	// Car expense
	// Car port
}
