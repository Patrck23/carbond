package routes

import (
	"car-bond/internals/controllers"
	// "car-bond/internals/middleware"

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
	user.Patch("/:id", controllers.UpdateUser)  // middleware.Protected(),
	user.Delete("/:id", controllers.DeleteUser) // middleware.Protected(),

	// Company
	api.Get("/companies", controllers.GetAllCompanies)
	company := api.Group("/company")
	company.Get("/:id", controllers.GetSingleCompanyById)
	company.Post("/", controllers.CreateCompany)
	company.Patch("/:id", controllers.UpdateCompany)
	company.Delete("/:id", controllers.DeleteCompanyById)
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
	car := api.Group("/car")
	car.Get("/id/:id", controllers.GetSingleCar)
	car.Get("/vin/:vinNumber", controllers.GetSingleCarByVinNumber)
	car.Post("/", controllers.CreateCar)
	car.Put("/:id", controllers.UpdateCar)
	car.Delete("/:id", controllers.DeleteCarByID)
	// Car expense
	api.Get("/carExpenses", controllers.GetAllCarExpenses)
	car.Get("/expenses/:carId", controllers.GetCarExpensesByCarId)
	car.Get("/:carId/expense/:id", controllers.GetCarExpenseById)
	car.Post("/expense", controllers.CreateCarExpense)
	car.Put("/expense/:id", controllers.UpdateCarExpense)
	car.Delete("/expense/:id", controllers.DeleteCarExpenseById)
	// Car port
	api.Get("/ports", controllers.GetAllPorts)
	car.Get("/ports/:carId", controllers.GetAllCarPorts)
	car.Get("/:carId/port/:id", controllers.GetCarPortById)
	car.Post("/port", controllers.CreateCarPort)
	car.Put("/port/:id", controllers.UpdateCarPort)
	car.Delete("/:carId/port/:id", controllers.DeleteCarPortById)

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

	// Meta data
	api.Post("/vehicle-evaluation", controllers.UploadPDF)
}
