package routes

import (
	"car-bond/internals/controllers"
	"car-bond/internals/middleware"

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
	api.Get("/cars", middleware.Protected(), carController.GetAllCars)
	car := api.Group("/car")
	car.Post("/", middleware.Protected(), carController.CreateCar)
	car.Get("/id/:id", middleware.Protected(), carController.GetSingleCar)
	car.Get("/vin/:vinNumber", middleware.Protected(), carController.GetSingleCarByVinNumber)
	car.Put("/:id/details", middleware.Protected(), carController.UpdateCar)
	car.Put("/:id/sale", middleware.Protected(), carController.UpdateCar2)
	car.Delete("/:id", middleware.Protected(), carController.DeleteCarByID)
	// Car expense
	api.Get("/carExpenses", middleware.Protected(), carController.GetAllCarExpenses)
	car.Get("/:carId/expenses", carController.GetCarExpensesByCarId)
	car.Get("/:carId/expense/:id", carController.GetCarExpenseById)
	car.Post("/expense", middleware.Protected(), carController.CreateCarExpense)
	car.Put("/expense/:id", middleware.Protected(), carController.UpdateCarExpense)
	car.Delete("/expense/:id", middleware.Protected(), carController.DeleteCarExpenseById)
	// Car port
	car.Get("/ports", middleware.Protected(), carController.GetAllCarPorts)
	car.Get("/:carId/port/:id", middleware.Protected(), carController.GetCarPortById)
	car.Post("/port", middleware.Protected(), carController.CreateCarPort)
	car.Put("/port/:id", middleware.Protected(), carController.UpdateCarPort)
	car.Delete("/:carId/port/:id", middleware.Protected(), carController.DeleteCarPortByID)

	companyDbService := controllers.NewCompanyRepository(db)
	companyController := controllers.NewCompanyController(companyDbService)

	// Company
	api.Get("/companies", middleware.Protected(), companyController.GetAllCompanies)
	company := api.Group("/company")
	company.Get("/:id", middleware.Protected(), companyController.GetSingleCompany)
	company.Post("/", middleware.Protected(), companyController.CreateCompany)
	company.Patch("/:id", middleware.Protected(), companyController.UpdateCompany)
	company.Delete("/:id", middleware.Protected(), companyController.DeleteCompanyByID)
	// Company Expenses
	api.Get("/expenses", middleware.Protected(), companyController.GetAllExpenses)
	company.Get("/expenses/:companyId", middleware.Protected(), companyController.GetCompanyExpensesByCompanyId)
	company.Get("/:companyId/expense/:id", middleware.Protected(), companyController.GetCompanyExpenseById)
	company.Post("/expense", middleware.Protected(), companyController.CreateCompanyExpense)
	company.Put("/expense/:id", middleware.Protected(), companyController.UpdateCompanyExpense)
	company.Delete("/expense/:id", middleware.Protected(), companyController.DeleteCompanyExpenseById)
	// Company Locations
	company.Get("/locations/:companyId", middleware.Protected(), companyController.GetAllCompanyLocations)
	company.Get("/:companyId/location/:id", middleware.Protected(), companyController.GetLocationByCompanyId)
	company.Post("/location", middleware.Protected(), companyController.CreateCompanyLocation)
	company.Put("/location/:id", middleware.Protected(), companyController.UpdateCompanyLocation)
	company.Delete("/location/:id", middleware.Protected(), companyController.DeleteLocationByID)

	customerDbService := controllers.NewCustomerRepository(db)
	customerController := controllers.NewCustomerController(customerDbService)

	// Customer
	api.Get("/customers", middleware.Protected(), customerController.GetAllCustomers)
	customer := api.Group("/customer")
	customer.Get("/:id", middleware.Protected(), customerController.GetSingleCustomer)
	customer.Post("/", middleware.Protected(), customerController.CreateCustomer)
	customer.Put("/:id", middleware.Protected(), customerController.UpdateCustomer)
	customer.Delete("/:id", middleware.Protected(), customerController.DeleteCustomerByID)
	// Upload
	customer.Post("/upload", middleware.Protected(), func(c *fiber.Ctx) error {
		return controllers.UploadCustomerFile(c, db)
	})
	customer.Get("/:id/files", middleware.Protected(), func(c *fiber.Ctx) error {
		return controllers.GetCustomerFiles(c, db)
	})
	customer.Get("/files/:file_id", middleware.Protected(), func(c *fiber.Ctx) error {
		return controllers.GetFile(c, db)
	})
	// Customer contact
	api.Get("/:companyId/contacts", middleware.Protected(), customerController.GetCustomerContactsByCompanyId)
	customer.Get("/contacts/:customerId", middleware.Protected(), customerController.GetCustomerContactsByCustomerId)
	customer.Get("/:customerId/contact/:id", middleware.Protected(), customerController.GetCustomerContactById)
	customer.Post("/contact", middleware.Protected(), customerController.CreateCustomerContact)
	customer.Put("/contact/:id", middleware.Protected(), customerController.UpdateCustomerContact)
	customer.Delete("/:customerId/contact/:id", middleware.Protected(), customerController.DeleteCustomerContactById)
	// Customer address
	api.Get("/:companyId/addresses", middleware.Protected(), customerController.GetCustomerAddressesByCompanyId)
	customer.Get("/addresses/:customerId", middleware.Protected(), customerController.GetCustomerAddressesByCustomerId)
	customer.Get("/:customerId/address/:id", middleware.Protected(), customerController.GetCustomerAddressById)
	customer.Post("/address", middleware.Protected(), customerController.CreateCustomerAddress)
	customer.Put("/address/:id", middleware.Protected(), customerController.UpdateCustomerAddress)
	customer.Delete("/:customerId/address/:id", middleware.Protected(), customerController.DeleteCustomerAddressById)

	userDbService := controllers.NewUserRepository(db)
	userController := controllers.NewUserController(userDbService)

	// User
	// , middleware.CheckPermissionsMiddleware("resource.*", []string{"R", "W"})
	api.Get("/users", middleware.Protected(), userController.GetAllUsers)
	user := api.Group("/user")
	user.Get("/:id", middleware.Protected(), userController.GetUserByID)
	user.Post("/", middleware.Protected(), userController.CreateUser)
	user.Patch("/:id", middleware.Protected(), userController.UpdateUser)
	user.Delete("/:id", middleware.Protected(), userController.DeleteUserByID)

	saleDbService := controllers.NewSaleRepository(db)
	saleController := controllers.NewSaleController(saleDbService)

	// Sale
	api.Get("/sales", middleware.Protected(), saleController.GetAllCarSales)
	sale := api.Group("/sale")
	sale.Get("/:id", middleware.Protected(), saleController.GetCarSale)
	sale.Post("/", middleware.Protected(), saleController.CreateCarSale)
	sale.Put("/:id", middleware.Protected(), saleController.UpdateSale)
	sale.Delete("/:id", middleware.Protected(), saleController.DeleteSaleByID)

	// // Initialize layers
	// repo := &SaleRepositoryImpl{db: db}
	// service := NewSaleService(repo)
	// controller := NewSaleController(service)
	// sale.Get("/searchSales/:criteria", middleware.Protected(), saleController.SearchByCriteria)

	// Invoice
	api.Get("/invoices", middleware.Protected(), saleController.GetSalePayments)
	invoice := api.Group("/invoice")
	invoice.Get("/:saleId/:id", middleware.Protected(), saleController.FindSalePaymentByIdAndSaleId)
	invoice.Post("/", middleware.Protected(), saleController.CreateInvoice)
	invoice.Put("/:id", middleware.Protected(), saleController.UpdateSalePayment)
	invoice.Delete("/:id", middleware.Protected(), saleController.DeleteSalePaymentByID)
	// Payment
	api.Get("/payments", middleware.Protected(), saleController.GetSalePaymentModes)
	payment := api.Group("/payment")
	payment.Get("/:salePaymentId/:id", middleware.Protected(), saleController.FindSalePaymentModeByIdAndSalePaymentId)
	payment.Post("/", middleware.Protected(), saleController.CreatePaymentMode)
	payment.Get("/:mode", middleware.Protected(), saleController.GetPaymentModesByMode)
	payment.Delete("/:id", middleware.Protected(), saleController.DeleteSalePaymentModeByID)
	payment.Put("/:id", middleware.Protected(), saleController.UpdateSalePaymentMode)

	metaDbService := controllers.NewExcecute(db)
	metaController := controllers.NewMetaController(metaDbService)

	// Meta data
	meta := api.Group("/meta")
	meta.Post("/vehicle-evaluation", middleware.Protected(), metaController.ProcessExcelAndUploadHandler)
	metaGDbService := controllers.NewMetaGetRepository(db)
	metaGController := controllers.NewMetaGetController(metaGDbService)
	meta.Get("/vehicle-evaluation", middleware.Protected(), metaGController.FetchVehicleEvaluationsByDescription)
	meta.Get("/weights", middleware.Protected(), metaGController.GetAllWeightUnits)
	meta.Get("/lengths", middleware.Protected(), metaGController.GetAllLeightUnits)
	meta.Get("/currency", middleware.Protected(), metaGController.GetAllCurrencies)
	NotFoundRoute(app)
}
