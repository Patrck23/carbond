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
	api.Get("/users", middleware.Protected(), middleware.CheckPermissionsMiddleware("resource.*", []string{"R", "W"}), controllers.GetAllUsers)
	user := api.Group("/user")
	user.Get("/:id", middleware.Protected(), controllers.GetUser)
	user.Post("/", middleware.Protected(), controllers.CreateUser)
	user.Patch("/:id", middleware.Protected(), controllers.UpdateUser)
	user.Delete("/:id", middleware.Protected(), controllers.DeleteUser)

	// Company
	api.Get("/companies", middleware.Protected(), controllers.GetAllCompanies)
	company := api.Group("/company")
	company.Get("/:id", middleware.Protected(), controllers.GetSingleCompanyById)
	// company.Post("/", controllers.CreateCompany)
	company.Patch("/:id", middleware.Protected(), controllers.UpdateCompany)
	// company.Delete("/:id", controllers.DeleteCompanyById)

	// Company Expenses
	api.Get("/expenses", middleware.Protected(), controllers.GetAllExpenses)
	company.Get("/expenses/:companyId", middleware.Protected(), controllers.GetCompanyExpensesByCompanyId)
	company.Get("/:companyId/expense/:id", middleware.Protected(), controllers.GetCompanyExpenseById)
	company.Post("/expense", middleware.Protected(), controllers.CreateCompanyExpense)
	company.Put("/expense/:id", middleware.Protected(), controllers.UpdateCompanyExpense)
	company.Delete("/expense/:id", middleware.Protected(), controllers.DeleteCompanyExpenseById)

	// Company Locations
	api.Get("/locations", middleware.Protected(), controllers.GetAllLocations)
	company.Get("/locations/:companyId", middleware.Protected(), controllers.GetAllCompanyLocations)
	company.Get("/:companyId/location/:id", middleware.Protected(), controllers.GetCompanyLocationById)
	company.Post("/location", middleware.Protected(), controllers.CreateCompanyLocation)
	company.Put("/location/:id", middleware.Protected(), controllers.UpdateCompanyLocation)
	company.Delete("/:companyId/location/:id", middleware.Protected(), controllers.DeleteCompanyLocationById)

	// Customer
	api.Get("/customers", middleware.Protected(), controllers.GetAllCustomers)
	customer := api.Group("/customer")
	customer.Get("/:id", middleware.Protected(), controllers.GetSingleCustomer)
	customer.Post("/", middleware.Protected(), controllers.CreateCustomer)
	customer.Put("/:id", middleware.Protected(), controllers.UpdateCustomer)
	customer.Delete("/:id", middleware.Protected(), controllers.DeleteCustomerByID)
	// Upload
	customer.Post("/upload", middleware.Protected(), controllers.UploadCustomerFile)
	customer.Get("/:id/files", middleware.Protected(), controllers.GetCustomerFiles)
	customer.Get("/files/:file_id", middleware.Protected(), controllers.GetFile)
	// Customer contact
	api.Get("/contacts", middleware.Protected(), controllers.GetAllContacts)
	customer.Get("/contacts/:customerId", middleware.Protected(), controllers.GetCustomerContacts)
	customer.Get("/:customerId/contact/:id", middleware.Protected(), controllers.GetCustomerContactById)
	customer.Post("/contact", middleware.Protected(), controllers.CreateCustomerContact)
	customer.Put("/contact/:id", middleware.Protected(), controllers.UpdateCustomerContact)
	customer.Delete("/:customerId/contact/:id", middleware.Protected(), controllers.DeleteCustomerContactById)
	// Customer address
	api.Get("/addresses", middleware.Protected(), controllers.GetAllAddresses)
	customer.Get("/addresses/:customerId", middleware.Protected(), controllers.GetCustomerAddresses)
	customer.Get("/:customerId/address/:id", middleware.Protected(), controllers.GetCustomerAddressById)
	customer.Post("/address", middleware.Protected(), controllers.CreateCustomerAddress)
	customer.Put("/address/:id", middleware.Protected(), controllers.UpdateCustomerAddress)
	customer.Delete("/:customerId/address/:id", middleware.Protected(), controllers.DeleteCustomerAddressById)

	// Car
	api.Get("/cars", middleware.Protected(), controllers.GetAllCars)
	car := api.Group("/car")
	car.Get("/id/:id", middleware.Protected(), controllers.GetSingleCar)
	car.Get("/vin/:vinNumber", middleware.Protected(), controllers.GetSingleCarByVinNumber)
	car.Post("/", middleware.Protected(), controllers.CreateCar)
	car.Put("/:id/details", middleware.Protected(), controllers.UpdateCar)
	car.Put("/:id/sale", middleware.Protected(), controllers.UpdateCar2)
	car.Delete("/:id", middleware.Protected(), controllers.DeleteCarByID)

	// Car expense
	api.Get("/carExpenses", middleware.Protected(), controllers.GetAllCarExpenses)
	car.Get("/:carId/expenses", middleware.Protected(), controllers.GetCarExpensesByCarId)
	car.Get("/:carId/expense/:id", middleware.Protected(), controllers.GetCarExpenseById)
	car.Post("/expense", middleware.Protected(), controllers.CreateCarExpense)
	car.Put("/expense/:id", middleware.Protected(), controllers.UpdateCarExpense)
	car.Delete("/expense/:id", middleware.Protected(), controllers.DeleteCarExpenseById)
	// Car port
	api.Get("/ports", middleware.Protected(), controllers.GetAllPorts)
	car.Get("/ports/:carId", middleware.Protected(), controllers.GetAllCarPorts)
	car.Get("/:carId/port/:id", middleware.Protected(), controllers.GetCarPortById)
	car.Post("/port", middleware.Protected(), controllers.CreateCarPort)
	car.Put("/port/:id", middleware.Protected(), controllers.UpdateCarPort)
	car.Delete("/:carId/port/:id", middleware.Protected(), controllers.DeleteCarPortById)

	//Sale
	api.Get("/sales", middleware.Protected(), controllers.GetAllCarSales)
	sale := api.Group("/sale")
	sale.Get("/:id", middleware.Protected(), controllers.GetCarSale)
	sale.Post("/", middleware.Protected(), controllers.CreateCarSale)
	sale.Put("/:id", middleware.Protected(), controllers.UpdateCarSale)
	sale.Delete("/:id", middleware.Protected(), controllers.DeleteCarSale)
	sale.Get("/searchSales/:criteria", middleware.Protected(), controllers.SearchByCriteria)
	// Invoice
	api.Get("/invoices", middleware.Protected(), controllers.GetAllInvoices)
	invoice := api.Group("/invoice")
	invoice.Post("/:id", middleware.Protected(), controllers.CreateInvoice)
	invoice.Post("/", middleware.Protected(), controllers.CreateInvoice)
	invoice.Put("/:id", middleware.Protected(), controllers.UpdateInvoice)
	invoice.Delete("/:id", middleware.Protected(), controllers.DeleteInvoiceByID)
	// Payment
	api.Get("/payments", middleware.Protected(), controllers.GetAllPayments)
	payment := api.Group("/payment")
	payment.Get("/:id", middleware.Protected(), controllers.GetPaymentByID)
	payment.Post("/:id", middleware.Protected(), controllers.CreatePayment)
	payment.Get("/:mode", middleware.Protected(), controllers.GetPaymentByModeOfPayment)
	payment.Delete("/:id", middleware.Protected(), controllers.DeletePaymentByID)
	payment.Put("/:id", middleware.Protected(), controllers.UpdatePayment)

	// Meta data
	meta := api.Group("/meta")
	meta.Post("/vehicle-evaluation", middleware.Protected(), controllers.ProcessExcelAndUpload)
	meta.Get("/vehicle-evaluation", middleware.Protected(), controllers.FetchVehicleEvaluationsByDescription)
	meta.Get("/weights", middleware.Protected(), controllers.GetAllWeightUnits)
	meta.Get("/lengths", middleware.Protected(), controllers.GetAllLengthUnits)
	meta.Get("/currency", middleware.Protected(), controllers.GetAllCurrencies)
	NotFoundRoute(app)

}
