package controllers

import (
	"car-bond/internals/database"
	"car-bond/internals/models/companyRegistration"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Create a company registration

func CreateCompany(c *fiber.Ctx) error {
	db := database.DB.Db
	company := new(companyRegistration.Company)
	err := c.BodyParser(company)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Something's wrong with your input", "data": err})
	}

	err = db.Create(&company).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed to create company", "data": err})
	}

	return c.Status(201).JSON(fiber.Map{"status": "success", "message": "Company created successfully", "data": company})
}

// Get all companies
func GetAllCompanies(c *fiber.Ctx) error {
	db := database.DB.Db
	var companies []companyRegistration.Company
	db.Find(&companies)
	if len(companies) == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Companies not found"})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Companies fetched successfully", "data": companies})
}

// Get a single company by ID

func GetSingleCompanyById(c *fiber.Ctx) error {
	db := database.DB.Db
	id := c.Params("id")
	var company companyRegistration.Company
	db.First(&company, id)
	if company.ID == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Company not found"})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Company fetched successfully", "data": company})
}

// Update a company
func UpdateCompany(c *fiber.Ctx) error {
	// Struct to parse and validate input
	type UpdateCompany struct {
		Name      string `json:"name"`
		StartDate string `json:"start_date"`
		UpdatedBy string `json:"updated_by"`
	}

	// Initialize database instance
	db := database.DB.Db

	// Fetch company ID from URL params
	id := c.Params("id")
	var company companyRegistration.Company

	// Find the company by ID
	if err := db.First(&company, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Company not found",
		})
	}

	// Parse request body into the UpdateCompany struct
	var updateData UpdateCompany
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input",
			"data":    err.Error(),
		})
	}

	// Update the company fields
	if updateData.Name != "" {
		company.Name = updateData.Name
	}
	if updateData.StartDate != "" {
		company.StartDate = updateData.StartDate
	}
	if updateData.UpdatedBy != "" {
		company.UpdatedBy = updateData.UpdatedBy
	}

	// Save the updated company to the database
	if err := db.Save(&company).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update company",
			"data":    err.Error(),
		})
	}

	// Return success response
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Company updated successfully",
		"data":    company,
	})
}

// Delete a company by ID

func DeleteCompanyById(c *fiber.Ctx) error {
	db := database.DB.Db
	id := c.Params("id")
	var company companyRegistration.Company
	db.First(&company, id)
	if company.ID == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Company not found"})
	}

	err := db.Delete(&company, id).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Failed to delete company", "data": err})
	}
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Company deleted"})
}

// ==================================================================================================================

// Create a company expense

func CreateCompanyExpense(c *fiber.Ctx) error {
	db := database.DB.Db
	companyExpense := new(companyRegistration.CompanyExpense)
	err := c.BodyParser(companyExpense)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Something's wrong with your input", "data": err})
	}

	err = db.Create(&companyExpense).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed to create company", "data": err})
	}

	return c.Status(201).JSON(fiber.Map{"status": "success", "message": "Company created successfully", "data": companyExpense})
}

// Get Company Expenses by ID

func GetCompanyExpenseById(c *fiber.Ctx) error {
	// Initialize database instance
	db := database.DB.Db

	// Retrieve the expense ID and company ID from the request parameters
	id := c.Params("id")
	companyId := c.Params("companyId")

	// Query the database for the expense by its ID and company ID
	var expense companyRegistration.CompanyExpense
	result := db.Preload("Company").Where("id = ? AND company_id = ?", id, companyId).First(&expense)

	// Handle potential database query errors
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{
				"status":  "error",
				"message": "Expense not found for the specified company",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch the expense",
			"error":   result.Error.Error(),
		})
	}

	// Return the fetched expense
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Expense fetched successfully",
		"data":    expense,
	})
}

// Get Company Expenses by Company ID

func GetCompanyExpensesByCompanyId(c *fiber.Ctx) error {
	// Initialize database instance
	db := database.DB.Db

	// Retrieve companyId from the request parameters
	companyId := c.Params("companyId")

	// Query the database for company expenses
	var expenses []companyRegistration.CompanyExpense
	result := db.Preload("Company").Where("company_id = ?", companyId).Find(&expenses)

	// Handle potential database errors
	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve company expenses",
			"error":   result.Error.Error(),
		})
	}

	// Handle the case where no expenses are found
	if len(expenses) == 0 {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "No expenses found for the specified company",
		})
	}

	// Return a success response with the fetched data
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Company expenses fetched successfully",
		"data":    expenses,
	})
}

// Update Company Expenses by

func UpdateCompanyExpense(c *fiber.Ctx) error {
	// Define a struct for input validation
	type UpdateCompanyExpenseInput struct {
		Description string  `json:"description"`
		Currency    string  `json:"currency"`
		Amount      float64 `json:"amount"`
		ExpenseDate string  `gorm:"type:date" json:"expense_date"`
		UpdatedBy   string  `json:"updated_by"`
	}

	db := database.DB.Db
	expenseID := c.Params("id")

	// Find the expense record by ID
	var expense companyRegistration.CompanyExpense
	if err := db.First(&expense, expenseID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Expense not found"})
		}
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Database error", "error": err.Error()})
	}

	// Parse the request body into the input struct
	var input UpdateCompanyExpenseInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Invalid input", "error": err.Error()})
	}

	// Update the fields of the expense record
	expense.Description = input.Description
	expense.Currency = input.Currency
	expense.Amount = input.Amount
	expense.ExpenseDate = input.ExpenseDate
	expense.UpdatedBy = input.UpdatedBy

	// Save the changes to the database
	if err := db.Save(&expense).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed to update expense", "error": err.Error()})
	}

	// Return the updated expense
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Expense updated successfully",
		"data":    expense,
	})
}

// Delete Company Expenses by ID

func DeleteCompanyExpenseById(c *fiber.Ctx) error {
	db := database.DB.Db

	// Parse companyId and expenseId from the URL
	companyId := c.Params("companyId")
	expenseId := c.Params("id")

	// Check if the expense exists and belongs to the specified company
	var expense companyRegistration.CompanyExpense
	if err := db.Where("id = ? AND company_id = ?", expenseId, companyId).First(&expense).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Company expense not found or does not belong to the specified company"})
	}

	// Delete the expense
	if err := db.Delete(&expense).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed to delete company expense", "data": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Company expense deleted successfully"})
}

// Get Company Expenses by Company ID and Expense Date

func GetCompanyExpensesByCompanyIdAndExpenseDate(c *fiber.Ctx) error {
	db := database.DB.Db
	companyId := c.Params("id")
	expenseDate := c.Params("expense_date")
	var expense companyRegistration.CompanyExpense
	db.Where("company_id = ? AND expense_date = ?", companyId, expenseDate).First(&expense)
	if expense.ID == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Company Expenses not found"})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Company Expenses fetched successfully", "data": expense})
}

// Get Company Expenses by Company ID and Expense Description

func GetCompanyExpensesByCompanyIdAndExpenseDescription(c *fiber.Ctx) error {
	db := database.DB.Db
	companyId := c.Params("id")
	expenseDescription := c.Params("expense_description")
	var expense companyRegistration.CompanyExpense
	db.Where("company_id = ? AND description = ?", companyId, expenseDescription).First(&expense)
	if expense.ID == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Company Expenses not found"})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Company Expenses fetched successfully", "data": expense})
}

// Get Company Expenses by Company ID and Currency

func GetCompanyExpensesByCompanyIdAndCurrency(c *fiber.Ctx) error {
	db := database.DB.Db
	companyId := c.Params("id")
	currency := c.Params("currency")
	var expenses []companyRegistration.CompanyExpense
	db.Where("company_id = ? AND currency = ?", companyId, currency).Find(&expenses)
	if len(expenses) == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Company Expenses not found"})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Company Expenses fetched successfully", "data": expenses})
}

// Get Company Expenses by Company ID, Expense Date, and Currency

func GetCompanyExpensesByThree(c *fiber.Ctx) error {
	db := database.DB.Db
	companyId := c.Params("id")
	expenseDate := c.Params("expense_date")
	currency := c.Params("currency")
	var expense companyRegistration.CompanyExpense
	db.Where("company_id = ? AND expense_date = ? AND currency = ?", companyId, expenseDate, currency).First(&expense)
	if expense.ID == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Company Expenses not found"})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Company Expenses fetched successfully", "data": expense})
}

// Get Company Expenses by Company ID, Expense Description, and Currency

func GetCompanyExpensesByThreeDec(c *fiber.Ctx) error {
	db := database.DB.Db
	companyId := c.Params("id")
	expenseDescription := c.Params("expense_description")
	currency := c.Params("currency")
	var expense companyRegistration.CompanyExpense
	db.Where("company_id = ? AND description = ? AND currency = ?", companyId, expenseDescription, currency).First(&expense)
	if expense.ID == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Company Expenses not found"})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Company Expenses fetched successfully", "data": expense})
}

// Get Company Expenses by Company ID, Expense Date, Expense Description, and Currency

func GetCompanyExpensesFilters(c *fiber.Ctx) error {
	db := database.DB.Db
	companyId := c.Params("id")
	expenseDate := c.Params("expense_date")
	expenseDescription := c.Params("expense_description")
	currency := c.Params("currency")
	var expense companyRegistration.CompanyExpense
	db.Where("company_id = ? AND expense_date = ? AND description = ? AND currency = ?", companyId, expenseDate, expenseDescription, currency).First(&expense)
	if expense.ID == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Company Expenses not found"})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Company Expenses fetched successfully", "data": expense})
}

func GetAllExpenses(c *fiber.Ctx) error {
	// Initialize database instance
	db := database.DB.Db
	var expenses []companyRegistration.CompanyExpense

	// Fetch all locations without requiring a Company association
	if err := db.Preload("Company").Find(&expenses).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch locations",
		})
	}

	return c.Status(fiber.StatusOK).JSON(expenses)
}

// ==================================================================================================================

// Create a company location

func CreateCompanyLocation(c *fiber.Ctx) error {
	db := database.DB.Db
	companyLocation := new(companyRegistration.CompanyLocation)
	err := c.BodyParser(companyLocation)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Something's wrong with your input", "data": err})
	}

	err = db.Create(&companyLocation).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed to create company location", "data": err})
	}

	return c.Status(201).JSON(fiber.Map{"status": "success", "message": "Company created successfully", "data": companyLocation})
}

func GetAllCompanyLocations(c *fiber.Ctx) error {
	// Initialize database instance
	db := database.DB.Db
	companyId := c.Params("companyId")

	var locations []companyRegistration.CompanyLocation

	// Fetch all locations with associated company details
	if err := db.Preload("Company").Where("company_id = ?", companyId).Find(&locations).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch company locations",
		})
	}

	return c.Status(fiber.StatusOK).JSON(locations)
}

// Get Company Location by ID

func GetCompanyLocationById(c *fiber.Ctx) error {
	// Initialize database instance
	db := database.DB.Db

	// Retrieve the location ID and company ID from the request parameters
	id := c.Params("id")
	companyId := c.Params("companyId")

	// Query the database for the location by its ID and company ID
	var location companyRegistration.CompanyLocation
	result := db.Preload("Company").Where("id = ? AND company_id = ?", id, companyId).First(&location)

	// Handle potential database query errors
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{
				"status":  "error",
				"message": "Location not found for the specified company",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch the location",
			"error":   result.Error.Error(),
		})
	}

	// Return the fetched location
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Location fetched successfully",
		"data":    location,
	})
}

// Update Company Location

func UpdateCompanyLocation(c *fiber.Ctx) error {
	// Define a struct for input validation
	type UpdateCompanyLocationInput struct {
		Address   string `json:"address"`
		Telephone string `json:"telephone"`
		Country   string `json:"country"`
		UpdatedBy string `json:"updated_by"`
	}

	db := database.DB.Db
	locationID := c.Params("id")

	// Find the location record by ID
	var location companyRegistration.CompanyLocation
	if err := db.First(&location, locationID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Location not found"})
		}
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Database error", "error": err.Error()})
	}

	// Parse the request body into the input struct
	var input UpdateCompanyLocationInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Invalid input", "error": err.Error()})
	}

	// Update the fields of the location record
	location.Address = input.Address
	location.Telephone = input.Telephone
	location.Country = input.Country
	location.UpdatedBy = input.UpdatedBy

	// Save the changes to the database
	if err := db.Save(&location).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed to update location", "error": err.Error()})
	}

	// Return the updated location
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Location updated successfully",
		"data":    location,
	})
}

// Delete Company Location by ID

func DeleteCompanyLocationById(c *fiber.Ctx) error {
	// Initialize database instance
	db := database.DB.Db

	// Parse parameters
	id := c.Params("id")
	companyId := c.Params("companyId")

	// Check if the company exists
	var company companyRegistration.Company
	if err := db.First(&company, companyId).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Company not found"})
	}

	// Check if the company location exists and belongs to the company
	var companyLocation companyRegistration.CompanyLocation
	if err := db.Where("id = ? AND company_id = ?", id, companyId).First(&companyLocation).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Company location not found or does not belong to the specified company"})
	}

	// Delete the company location
	if err := db.Delete(&companyLocation).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed to delete company location", "data": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Company location deleted successfully"})
}

func GetAllLocations(c *fiber.Ctx) error {
	// Initialize database instance
	db := database.DB.Db
	var locations []companyRegistration.CompanyLocation

	// Fetch all locations without requiring a Company association
	if err := db.Preload("Company").Find(&locations).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch locations",
		})
	}

	return c.Status(fiber.StatusOK).JSON(locations)
}
