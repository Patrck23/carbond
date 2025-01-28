package routes

import (
	"car-bond/internals/controllers"
	// "car-bond/internals/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// SetupRoutes initializes all routes with their respective controllers
func SetupRoute(app *fiber.App, db *gorm.DB) {

	// Initialize the dbService and other controllers
	dbService := controllers.NewDatabaseService(db)
	groupController := controllers.NewGroupController(dbService)
	roleController := controllers.NewRoleController(dbService)
	resourceController := controllers.NewResourceController(dbService)
	permissionController := controllers.NewPermissionController(dbService)

	api := app.Group("/api")
	// Define routes
	api.Get("/groups", groupController.GetAllGroups)
	groupRoutes := api.Group("/group")
	groupRoutes.Post("/", groupController.CreateGroup)
	groupRoutes.Get("/:code", groupController.GetGroup)
	groupRoutes.Put("/:code", groupController.UpdateGroup)
	groupRoutes.Delete("/:code", groupController.DeleteGroup)

	api.Get("/roles", roleController.GetAllRoles)
	roleRoutes := app.Group("/role")
	roleRoutes.Post("/", roleController.CreateRole)
	roleRoutes.Get("/:code", roleController.GetRole)
	api.Get("/roles", roleController.GetRolesForGroups) //roles?group_codes=group1,group2,group3

	api.Get("/resources", resourceController.GetAllResources)
	api.Get("/resources/:code", resourceController.GetResource)
	resourceRoutes := api.Group("/resourse")
	resourceRoutes.Post("/", resourceController.CreateResource)

	api.Post("/permissions", permissionController.CreatePermission)
	api.Post("/wildcard-permissions", permissionController.CreateWildCardPermission)
	api.Get("/permissions", permissionController.GetGrantedPermissions)
	api.Get("/check-permissions", permissionController.CheckPermissions)
	api.Get("/explict-permissions", permissionController.GetExplicitPermissions)
	api.Get("/wildcard-permissions", permissionController.GetWildCardPermissions)
	api.Get("/exist-permissions", permissionController.ResourceExplicitPermissionsExists)
	api.Get("/group-role-exist", permissionController.GroupsWithRoleExists)
	api.Get("/roles-resource-permisions", permissionController.GetPermissions)

	carDbService := controllers.NewCarRepository(db)
	carController := controllers.NewCarController(carDbService)
	// Create a group for authentication routes
	authGroup := api.Group("/auth")

	// Define the login route and pass the db instance
	authGroup.Post("/login", func(c *fiber.Ctx) error {
		return controllers.Login(c, db)
	})

	// Car
	api.Get("/cars", carController.GetAllCars) // middleware.Protected(),
	car := api.Group("/car")
	car.Post("/", carController.CreateCar)
	car.Get("/id/:id", carController.GetSingleCar)
	car.Get("/vin/:vinNumber", carController.GetSingleCarByVinNumber)
	car.Put("/:id/details", carController.UpdateCar)
	car.Put("/:id/sale", carController.UpdateCar2)
	car.Delete("/:id", carController.DeleteCarByID)
	// Car expense
	api.Get("/carExpenses", carController.GetAllCarExpenses)
	car.Get("/:carId/expenses", carController.GetCarExpensesByCarId)
	car.Get("/:carId/expense/:id", carController.GetCarExpenseById)
	car.Post("/expense", carController.CreateCarExpense)
	car.Put("/expense/:id", carController.UpdateCarExpense)
	car.Delete("/expense/:id", carController.DeleteCarExpenseById)
	// Car port
	car.Get("/ports", carController.GetAllCarPorts)
	car.Get("/:carId/port/:id", carController.GetCarPortById)
	car.Post("/port", carController.CreateCarPort)
	car.Put("/port/:id", carController.UpdateCarPort)
	car.Delete("/:carId/port/:id", carController.DeleteCarPortByID)

	companyDbService := controllers.NewCompanyRepository(db)
	companyController := controllers.NewCompanyController(companyDbService)

	// Company
	api.Get("/companies", companyController.GetAllCompanies)
	company := api.Group("/company")
	company.Get("/:id", companyController.GetSingleCompany)
	company.Post("/", companyController.CreateCompany)
	company.Patch("/:id", companyController.UpdateCompany)
	company.Delete("/:id", companyController.DeleteCompanyByID)
	// Company Expenses
	api.Get("/expenses", companyController.GetAllExpenses)
	company.Get("/expenses/:companyId", companyController.GetCompanyExpensesByCompanyId)
	company.Get("/:companyId/expense/:id", companyController.GetCompanyExpenseById)
	company.Post("/expense", companyController.CreateCompanyExpense)
	company.Put("/expense/:id", companyController.UpdateCompanyExpense)
	company.Delete("/expense/:id", companyController.DeleteCompanyExpenseById)
	// Company Locations
	company.Get("/locations/:companyId", companyController.GetAllCompanyLocations)
	company.Get("/:companyId/location/:id", companyController.GetLocationByCompanyId)
	company.Post("/location", companyController.CreateCompanyLocation)
	company.Put("/location/:id", companyController.UpdateCompanyLocation)
	company.Delete("/location/:id", companyController.DeleteLocationByID)

	customerDbService := controllers.NewCustomerRepository(db)
	customerController := controllers.NewCustomerController(customerDbService)

	// Customer
	api.Get("/customers", customerController.GetAllCustomers)
	customer := api.Group("/customer")
	customer.Get("/:id", customerController.GetSingleCustomer)
	customer.Post("/", customerController.CreateCustomer)
	customer.Put("/:id", customerController.UpdateCustomer)
	customer.Delete("/:id", customerController.DeleteCustomerByID)
	// Upload
	customer.Post("/upload", func(c *fiber.Ctx) error {
		return controllers.UploadCustomerFile(c, db)
	})
	customer.Get("/:id/files", func(c *fiber.Ctx) error {
		return controllers.GetCustomerFiles(c, db)
	})
	customer.Get("/files/:file_id", func(c *fiber.Ctx) error {
		return controllers.GetFile(c, db)
	})
	// Customer contact
	api.Get("/:companyId/contacts", customerController.GetCustomerContactsByCompanyId)
	customer.Get("/contacts/:customerId", customerController.GetCustomerContactsByCustomerId)
	customer.Get("/:customerId/contact/:id", customerController.GetCustomerContactById)
	customer.Post("/contact", customerController.CreateCustomerContact)
	customer.Put("/contact/:id", customerController.UpdateCustomerContact)
	customer.Delete("/:customerId/contact/:id", customerController.DeleteCustomerContactById)
	// Customer address
	api.Get("/:companyId/addresses", customerController.GetCustomerAddressesByCompanyId)
	customer.Get("/addresses/:customerId", customerController.GetCustomerAddressesByCustomerId)
	customer.Get("/:customerId/address/:id", customerController.GetCustomerAddressById)
	customer.Post("/address", customerController.CreateCustomerAddress)
	customer.Put("/address/:id", customerController.UpdateCustomerAddress)
	customer.Delete("/:customerId/address/:id", customerController.DeleteCustomerAddressById)

	userDbService := controllers.NewUserRepository(db)
	userController := controllers.NewUserController(userDbService)

	// User
	// , middleware.CheckPermissionsMiddleware("resource.*", []string{"R", "W"})
	api.Get("/users", userController.GetAllUsers)
	user := api.Group("/user")
	user.Get("/:id", userController.GetUserByID)
	user.Post("/", userController.CreateUser)
	user.Patch("/:id", userController.UpdateUser)
	user.Delete("/:id", userController.DeleteUserByID)

	saleDbService := controllers.NewSaleRepository(db)
	saleController := controllers.NewSaleController(saleDbService)

	// Sale
	api.Get("/sales", saleController.GetAllCarSales)
	sale := api.Group("/sale")
	sale.Get("/:id", saleController.GetCarSale)
	sale.Post("/", saleController.CreateCarSale)
	sale.Put("/:id", saleController.UpdateSale)
	sale.Delete("/:id", saleController.DeleteSaleByID)

	// // Initialize layers
	// repo := &SaleRepositoryImpl{db: db}
	// service := NewSaleService(repo)
	// controller := NewSaleController(service)
	// sale.Get("/searchSales/:criteria",  saleController.SearchByCriteria)

	// Invoice
	api.Get("/invoices", saleController.GetSalePayments)
	invoice := api.Group("/invoice")
	invoice.Get("/:saleId/:id", saleController.FindSalePaymentByIdAndSaleId)
	invoice.Post("/", saleController.CreateInvoice)
	invoice.Put("/:id", saleController.UpdateSalePayment)
	invoice.Delete("/:id", saleController.DeleteSalePaymentByID)
	// Payment
	api.Get("/payments", saleController.GetSalePaymentModes)
	payment := api.Group("/payment")
	payment.Get("/:salePaymentId/:id", saleController.FindSalePaymentModeByIdAndSalePaymentId)
	payment.Post("/", saleController.CreatePaymentMode)
	payment.Get("/:mode", saleController.GetPaymentModesByMode)
	payment.Delete("/:id", saleController.DeleteSalePaymentModeByID)
	payment.Put("/:id", saleController.UpdateSalePaymentMode)

	metaDbService := controllers.NewExcecute(db)
	metaController := controllers.NewMetaController(metaDbService)

	// Meta data
	meta := api.Group("/meta")
	meta.Post("/vehicle-evaluation", metaController.ProcessExcelAndUploadHandler)
	metaGDbService := controllers.NewMetaGetRepository(db)
	metaGController := controllers.NewMetaGetController(metaGDbService)
	meta.Get("/vehicle-evaluation", metaGController.FetchVehicleEvaluationsByDescription)
	meta.Get("/weights", metaGController.GetAllWeightUnits)
	meta.Get("/lengths", metaGController.GetAllLeightUnits)
	meta.Get("/currency", metaGController.GetAllCurrencies)
	NotFoundRoute(app)
}
