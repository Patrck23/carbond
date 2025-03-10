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

	// pdbService := middleware.NewDatabaseService(db)

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
	api.Get("/cars", middleware.Protected(), carController.GetAllCars) // , middleware.Protected(), middleware.PermissionMiddleware(pdbService, "resource.*", []string{"R", "W"})
	car := api.Group("/car")
	car.Post("/", middleware.Protected(), carController.CreateCar)
	car.Get("/id/:id", middleware.Protected(), carController.GetSingleCar)
	car.Get("/vin/:ChasisNumber", middleware.Protected(), carController.GetSingleCarByChasisNumber)
	car.Put("/:id/details", middleware.Protected(), carController.UpdateCar)
	car.Put("/:id/sale", middleware.Protected(), carController.UpdateCar2)
	car.Put("/:id/shipping-invoice", middleware.Protected(), carController.UpdateCar3)
	car.Delete("/:id", middleware.Protected(), carController.DeleteCarByID)
	// Car expense
	api.Get("/carExpenses", middleware.Protected(), carController.GetAllCarExpenses)
	car.Get("/:carId/expenses", carController.GetCarExpensesByCarId)
	car.Get("/:carId/expense/:id", carController.GetCarExpenseById)
	car.Post("/expense", middleware.Protected(), carController.CreateCarExpenses)
	car.Put("/expense/:id", middleware.Protected(), carController.UpdateCarExpense)
	car.Delete("/:carId/expense/:id", middleware.Protected(), carController.DeleteCarExpenseById)
	api.Get("/total-car-expense/:id", middleware.Protected(), carController.GetTotalCarExpenses)
	api.Get("/cars/search", middleware.Protected(), carController.SearchCars)
	car.Get("uploads", middleware.Protected(), carController.FetchCarUploads)
	car.Get("dash", middleware.Protected(), carController.GetDashboardData)
	car.Get("dash/:companyId", middleware.Protected(), carController.GetCompanyDashboardData)

	car.Post("/upload", middleware.Protected(), func(c *fiber.Ctx) error {
		return controllers.UploadCarFile(c, db)
	})
	car.Get("/:id/files", middleware.Protected(), func(c *fiber.Ctx) error {
		return controllers.GetCarFiles(c, db)
	})
	car.Get("/files/:file_id", middleware.Protected(), func(c *fiber.Ctx) error {
		return controllers.GetFile(c, db)
	})

	shippingDbService := controllers.NewShippingRepository(db)
	shippingController := controllers.NewShippingController(shippingDbService)

	shipping := api.Group("/shipping")
	shipping.Get("/invoices", shippingController.GetAllShippingInvoices)
	shipping.Post("/invoices", middleware.Protected(), shippingController.CreateShippingInvoice)
	shipping.Get("/invoice/:id", middleware.Protected(), shippingController.GetSingleInvoice)
	shipping.Get("/invoice/no/:no", middleware.Protected(), shippingController.GetShippingInvoiceByInvoiceNum)
	shipping.Put("/invoice/:id", middleware.Protected(), shippingController.UpdateShippingInvoice)
	shipping.Delete("/invoice/:id", middleware.Protected(), shippingController.DeleteShippingInvoiceByID)

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
	customer.Get("/:id/upload", middleware.Protected(), customerController.FetchCustomerUpload)
	api.Get("/customers/search", middleware.Protected(), customerController.SearchCustomers)

	// Customer contact
	api.Get("/:companyId/contacts", middleware.Protected(), customerController.GetCustomerContactsByCompanyId)
	customer.Get("/contacts/:customerId", middleware.Protected(), customerController.GetCustomerContactsByCustomerId)
	customer.Get("/:customerId/contact/:id", middleware.Protected(), customerController.GetCustomerContactById)
	customer.Post("/contact", middleware.Protected(), customerController.CreateCustomerContact)
	customer.Put("/contact/:id", middleware.Protected(), customerController.UpdateCustomerContact)
	customer.Delete("/:customerId/contact/:id", middleware.Protected(), customerController.DeleteCustomerContactById)
	// Customer upload
	customer.Get("/:id/upload", middleware.Protected(), customerController.FetchCustomerUpload)
	// Customer address
	api.Get("/:companyId/addresses", middleware.Protected(), customerController.GetCustomerAddressesByCompanyId)
	customer.Get("/addresses/:customerId", middleware.Protected(), customerController.GetCustomerAddressesByCustomerId)
	customer.Get("/:customerId/address/:id", middleware.Protected(), customerController.GetCustomerAddressById)
	customer.Post("/address", middleware.Protected(), customerController.CreateCustomerAddress)
	customer.Put("/address/:id", middleware.Protected(), customerController.UpdateCustomerAddress)
	customer.Delete("/:customerId/address/:id", middleware.Protected(), customerController.DeleteCustomerAddressById)

	userDbService := controllers.NewUserRepository(db)
	userController := controllers.NewUserController(userDbService)

	api.Get("/users", middleware.Protected(), userController.GetAllUsers)
	api.Get("/users/:companyId", middleware.Protected(), userController.GetUsersByCompany)
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
	sale.Get("/statement/:customerId", middleware.Protected(), saleController.GenerateCustomerStatement)

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
	// Deposits
	api.Get("/deposits", middleware.Protected(), saleController.GetSalePaymentDeposits)
	deposit := api.Group("/deposit")
	deposit.Get("/:salePaymentId/:id", middleware.Protected(), saleController.FindSalePaymentDepositByIdAndSalePaymentId)
	deposit.Post("/", middleware.Protected(), saleController.CreatePaymentDeposit)
	deposit.Get("/:name", middleware.Protected(), saleController.GetPaymentDepositsByName)
	deposit.Delete("/:id", middleware.Protected(), saleController.DeleteSalePaymentDepositByID)
	deposit.Put("/:id", middleware.Protected(), saleController.UpdateSalePaymentDeposit)

	// Auction Sale
	auctionSaleDbService := controllers.NewSaleAuctionRepository(db)
	auctionSaleController := controllers.NewSaleAuctionController(auctionSaleDbService)

	// Sale
	api.Get("/auction-sales", middleware.Protected(), auctionSaleController.GetAllCarSales)
	SaleAuction := api.Group("/auction-sale")
	SaleAuction.Get("/:id", middleware.Protected(), auctionSaleController.GetCarSale)
	SaleAuction.Post("/", middleware.Protected(), auctionSaleController.CreateCarSale)
	SaleAuction.Put("/:id", middleware.Protected(), auctionSaleController.UpdateSale)
	SaleAuction.Delete("/:id", middleware.Protected(), auctionSaleController.DeleteSaleByID)

	// Alert data
	alertDbService := controllers.NewAlertRepository(db)
	alertController := controllers.NewAlertController(alertDbService)
	api.Get("/alerts/search", middleware.Protected(), alertController.SearchAlerts)
	api.Put("/alert/:id", middleware.Protected(), alertController.UpdateAlert)

	// Meta data
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
	meta.Get("/expenses", middleware.Protected(), metaGController.GetAllExpenseCategories)
	meta.Get("/ports", middleware.Protected(), metaGController.FindPorts)
	meta.Get("/payment-modes", middleware.Protected(), metaGController.FindPaymentModeBymode)
	app.Static("/uploads", "./api/uploads")
	NotFoundRoute(app)
}
