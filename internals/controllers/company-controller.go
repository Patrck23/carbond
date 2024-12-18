package controllers

import (
	"car-bond/internals/database"
	"car-bond/internals/models/companyRegistration"

	"github.com/gofiber/fiber/v2"
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
	type UpdateCompany struct {
		Name      string `json:"name"`
		StartDate string `gorm:"type:date" json:"start_date"`
		CreatedBy string `json:"created_by"`
		UpdatedBy string `json:"updated_by"`
	}
	db := database.DB.Db
	id := c.Params("id")
	var company companyRegistration.Company
	db.First(&company, id)
	if company.ID == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Company not found"})
	}

	err := c.BodyParser(&company)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Something's wrong with your input", "data": err})
	}
	company.Name = company.Name
	company.StartDate = company.StartDate
	company.UpdatedBy = company.UpdatedBy

	db.Save(&company)
	return c.JSON(fiber.Map{"status": "success", "message": "Company updated successfully", "data": company})
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

// Get Company Expenses by ID

func GetCompanyExpenseById(c *fiber.Ctx) error {
	// Initialize database instance
	db := database.DB.Db

	// Retrieve the expense ID and company ID from the request parameters
	id := c.Params("id")
	companyId := c.Params("companyId")

	// Query the database for the expense by its ID and company ID
	var expense companyRegistration.CompanyExpenses
	result := db.Where("id = ? AND company_id = ?", id, companyId).First(&expense)

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
	var expenses []companyRegistration.CompanyExpenses
	result := db.Where("company_id = ?", companyId).Find(&expenses)

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
	type UpdateCompanyExpense struct {
		Description string  `json:"description"`
		Currency    string  `json:"currency"`
		Amount      float64 `json:"amount"`
		ExpenseDate string  `gorm:"type:date" json:"expense_date"`
		UpdatedBy   int     `json:"updated_by"`
	}

	db := database.DB.Db
	id := c.Params("id")
	var company companyRegistration.Company
	db.First(&company, id)
	if company.ID == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Company not found"})
	}

	var UpdateCompanyExpense UpdateCompanyExpense
	err := c.BodyParser(&UpdateCompanyExpense)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Something's wrong with your input", "data": err})
	}
	UpdateCompanyExpense.Description = UpdateCompanyExpense.Description
	UpdateCompanyExpense.Currency = UpdateCompanyExpense.Currency
	UpdateCompanyExpense.Amount = UpdateCompanyExpense.Amount
	UpdateCompanyExpense.ExpenseDate = UpdateCompanyExpense.ExpenseDate
	UpdateCompanyExpense.UpdatedBy = UpdateCompanyExpense.UpdatedBy
	db.Save(&company)
	return c.JSON(fiber.Map{"status": "success", "message": "Company Expenses updated successfully", "data": company})
}

// Delete Company Expenses by ID

func DeleteCompanyExpenseById(c *fiber.Ctx) error {
	db := database.DB.Db
	id := c.Params("id")
	var company companyRegistration.Company
	db.First(&company, id)
	if company.ID == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Company not found"})
	}

	var expense companyRegistration.CompanyExpense
	db.First(&expense, id)
	if expense.ID == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Company Expenses not found"})
	}

	err := db.Delete(&expense, id).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Failed to delete company expenses", "data": err})
	}
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Company Expenses deleted"})
}

// Get Company Expenses by Company ID and Expense Date

func GetCompanyExpensesByCompanyIdAndExpenseDate(c *fiber.Ctx) error {
	db := database.DB.Db
	companyId := c.Params("id")
	expenseDate := c.Params("expense_date")
	var expense companyRegistration.CompanyExpenses
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
	var expense companyRegistration.CompanyExpenses
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
	var expenses []companyRegistration.CompanyExpenses
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
	var expense companyRegistration.CompanyExpenses
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
	var expense companyRegistration.CompanyExpenses
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
	var expense companyRegistration.CompanyExpenses
	db.Where("company_id = ? AND expense_date = ? AND description = ? AND currency = ?", companyId, expenseDate, expenseDescription, currency).First(&expense)
	if expense.ID == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Company Expenses not found"})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Company Expenses fetched successfully", "data": expense})
}
