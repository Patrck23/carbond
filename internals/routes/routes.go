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

	//Sale
	api.Get("/sales", controllers.GetAllCarSales)
	sale := api.Group("/sale")
	sale.Get("/:id", controllers.GetCarSale)
	sale.Post("/", controllers.CreateCarSale)
	sale.Put("/:id", controllers.UpdateCarSale)
	sale.Delete("/:id", controllers.DeleteCarSale)
	sale.Get("/searchSales/:criteria", controllers.SearchByCriteria)

	// Invoice
	api.Get("/invoices", controllers.GetAllInvoices)
	invoice := api.Group("/invoice")
	invoice.Post("/:id", controllers.CreateInvoice)
	invoice.Post("/", controllers.CreateInvoice)
	invoice.Put("/:id", controllers.UpdateInvoice)
	invoice.Delete("/:id", controllers.DeleteInvoiceByID)

	// Payment
	api.Get("/payments", controllers.GetAllPayments)
	payment := api.Group("/payment")
	payment.Get("/:id", controllers.GetPaymentByID)
	payment.Post("/:id", controllers.CreatePayment)
	payment.Get("/:mode", controllers.GetPaymentByModeOfPayment)
	payment.Delete("/:id", controllers.DeletePaymentByID)
	payment.Put("/:id", controllers.UpdatePayment)
}
