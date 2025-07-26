package controllers

import (
	"car-bond/internals/models/companyRegistration"
	"car-bond/internals/utils"

	"car-bond/internals/repository"
	"errors"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type CompanyController struct {
	repo repository.CompanyRepository
}

func NewCompanyController(repo repository.CompanyRepository) *CompanyController {
	return &CompanyController{repo: repo}
}

// ============================================

func (h *CompanyController) CreateCompany(c *fiber.Ctx) error {
	// Initialize a new Company instance
	company := new(companyRegistration.Company)

	// Parse the request body into the company instance
	if err := c.BodyParser(company); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input provided",
			"data":    err.Error(),
		})
	}

	// Attempt to create the company record using the repository
	if err := h.repo.CreateCompany(company); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create company",
			"data":    err.Error(),
		})
	}

	// Return the newly created company record
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Company created successfully",
		"data":    company,
	})
}

// ===================

func (h *CompanyController) GetAllCompanies(c *fiber.Ctx) error {
	// Fetch paginated Companies using the repository
	pagination, companies, err := h.repo.GetPaginatedCompanies(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve Companies",
			"data":    err.Error(),
		})
	}

	// Return the paginated response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Companies retrieved successfully",
		"data":    companies,
		"pagination": fiber.Map{
			"total_items":  pagination.TotalItems,
			"total_pages":  pagination.TotalPages,
			"current_page": pagination.CurrentPage,
			"limit":        pagination.ItemsPerPage,
		},
	})
}

// =====================

// GetSingleCompany fetches a company with its associated locations and expenses from the database
func (h *CompanyController) GetSingleCompany(c *fiber.Ctx) error {
	// Get the company ID from the route parameters
	id := c.Params("id")

	// Fetch the company by ID
	company, err := h.repo.GetCompanyByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Company not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve Company",
			"data":    err.Error(),
		})
	}

	// Fetch company locations associated with the company
	companyLocations, err := h.repo.GetCompanyLocations(company.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve company locations",
			"data":    err.Error(),
		})
	}

	// Prepare the response
	response := fiber.Map{
		"company":          company,
		"company_location": companyLocations,
	}

	// Return the response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Company and associated data retrieved successfully",
		"data":    response,
	})
}

// ====================

// Define the UpdateCompany struct
type UpdateCompanyPayload struct {
	Name      string `json:"name"`
	StartDate string `json:"start_date"`
	UpdatedBy string `json:"updated_by"`
}

// UpdateCompany handler function
func (h *CompanyController) UpdateCompany(c *fiber.Ctx) error {
	// Get the company ID from the route parameters
	id := c.Params("id")

	// Find the company in the database
	company, err := h.repo.GetCompanyByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{
				"status":  "error",
				"message": "Company not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve company",
			"data":    err.Error(),
		})
	}

	// Parse the request body into the UpdateCompanyPayload struct
	var payload UpdateCompanyPayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input",
			"data":    err.Error(),
		})
	}

	// Update the company fields using the payload
	updateCompanyFields(&company, payload) // Pass the parsed payload

	// Save the changes to the database
	if err := h.repo.UpdateCompany(&company); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update company",
			"data":    err.Error(),
		})
	}

	// Return the updated company
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "company updated successfully",
		"data":    company,
	})
}

// UpdateCompanyFields updates the fields of a company using the UpdateCompany struct
func updateCompanyFields(company *companyRegistration.Company, updateCompanyData UpdateCompanyPayload) {
	company.Name = updateCompanyData.Name
	company.StartDate = updateCompanyData.StartDate
	company.UpdatedBy = updateCompanyData.UpdatedBy
}

// ======================

// DeleteCompanyByID deletes a Company by its ID
func (h *CompanyController) DeleteCompanyByID(c *fiber.Ctx) error {
	// Get the Company ID from the route parameters
	id := c.Params("id")

	// Find the Company in the database
	company, err := h.repo.GetCompanyByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{
				"status":  "error",
				"message": "Company not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to find Company",
			"data":    err.Error(),
		})
	}

	// Delete the Company
	if err := h.repo.DeleteByID(id); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to delete Company",
			"data":    err.Error(),
		})
	}

	// Return success response
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Company deleted successfully",
		"data":    company,
	})
}

// ==================================================================================================================
// Create a company expense

// CreateCompanyExpense handles the creation of a company expense
func (h *CompanyController) CreateCompanyExpense(c *fiber.Ctx) error {
	// Parse the request body into a companyExpense struct
	companyExpense := new(companyRegistration.CompanyExpense)
	if err := c.BodyParser(companyExpense); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input data",
			"data":    err.Error(),
		})
	}

	// Create the company expense in the database
	if err := h.repo.CreateCompanyExpense(companyExpense); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create company expense",
			"data":    err.Error(),
		})
	}

	// Return success response
	return c.Status(201).JSON(fiber.Map{
		"status":  "success",
		"message": "Company expense created successfully",
		"data":    companyExpense,
	})
}

// ========================

func (h *CompanyController) GetCompanyExpensesByCompanyId(c *fiber.Ctx) error {
	companyId := c.Params("companyId")

	// Fetch paginated expenses using the repository
	pagination, expenses, err := h.repo.GetPaginatedExpenses(c, companyId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve expenses",
			"data":    err.Error(),
		})
	}

	// Return the paginated response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Expenses and associated data retrieved successfully",
		"data":    expenses,
		"pagination": fiber.Map{
			"total_items":  pagination.TotalItems,
			"total_pages":  pagination.TotalPages,
			"current_page": pagination.CurrentPage,
			"limit":        pagination.ItemsPerPage,
		},
	})
}

// ========================

func (h *CompanyController) GetCompanyExpenseById(c *fiber.Ctx) error {
	// Retrieve the expense ID and Company ID from the request parameters
	id := c.Params("id")
	companyId := c.Params("companyId")

	// Fetch the company expense from the repository
	expense, err := h.repo.FindCompanyExpenseByIdAndCompanyId(id, companyId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Expense not found for the specified company",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch the expense",
			"error":   err.Error(),
		})
	}

	// Return the fetched expense
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Expense fetched successfully",
		"data":    expense,
	})
}

// ========================

func (h *CompanyController) UpdateCompanyExpense(c *fiber.Ctx) error {
	// Define a struct for input validation
	type UpdateCompanyExpenseInput struct {
		Description string  `json:"description" validate:"required"`
		Currency    string  `json:"currency" validate:"required"`
		Amount      float64 `json:"amount" validate:"required,gt=0"`
		ExpenseDate string  `json:"expense_date" validate:"required"`
		UpdatedBy   string  `json:"updated_by" validate:"required"`
	}

	// Parse the expense ID from the request parameters
	expenseID := c.Params("id")

	// Parse and validate the request body
	var input UpdateCompanyExpenseInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input",
			"error":   err.Error(),
		})
	}

	// Use a validation library to validate the input
	if validationErr := utils.ValidateStruct(input); validationErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Validation failed",
			"errors":  validationErr,
		})
	}

	// Fetch the expense record using the repository
	expense, err := h.repo.FindCompanyExpenseById(expenseID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Expense not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch expense",
			"error":   err.Error(),
		})
	}

	// Update the expense fields
	expense.Description = input.Description
	expense.Currency = input.Currency
	expense.Amount = input.Amount
	expense.ExpenseDate = input.ExpenseDate
	expense.UpdatedBy = input.UpdatedBy

	// Save the updated expense using the repository
	if err := h.repo.UpdateCompanyExpense(expense); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update expense",
			"error":   err.Error(),
		})
	}

	// Return the updated expense
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Expense updated successfully",
		"data":    expense,
	})
}

// ====================

func (h *CompanyController) DeleteCompanyExpenseById(c *fiber.Ctx) error {
	// Parse companyId and expenseId from the request parameters
	companyId := c.Params("companyId")
	expenseId := c.Params("id")

	// Check if the expense exists and belongs to the specified company
	expense, err := h.repo.FindCompanyExpenseByCompanyAndId(companyId, expenseId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Company expense not found or does not belong to the specified company",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch company expense",
			"error":   err.Error(),
		})
	}

	// Delete the expense using the repository
	if err := h.repo.DeleteCompanyExpense(expense); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to delete company expense",
			"error":   err.Error(),
		})
	}

	// Return success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Company expense deleted successfully",
	})
}

// ====================

// Get Company Expenses by Company ID and Expense Date

func (h *CompanyController) GetCompanyExpensesByCompanyIdAndExpenseDate(c *fiber.Ctx) error {
	// Retrieve companyId and expenseDate from the request parameters
	companyId := c.Params("id")
	expenseDate := c.Params("expense_date")

	// Fetch expenses using the repository
	expenses, err := h.repo.FindCompanyExpensesByCompanyIdAndExpenseDate(companyId, expenseDate)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Company expenses not found for the specified date",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch company expenses",
			"error":   err.Error(),
		})
	}

	// Return the fetched expenses
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Company expenses fetched successfully",
		"data":    expenses,
	})
}

// =====================

// Get Company Expenses by Company ID and Expense Description

func (h *CompanyController) GetCompanyExpensesByCompanyIdAndExpenseDescription(c *fiber.Ctx) error {
	// Retrieve companyId and expenseDescription from the request parameters
	companyId := c.Params("id")
	expenseDescription := c.Params("expense_description")

	// Fetch expenses using the repository
	expenses, err := h.repo.FindCompanyExpensesByCompanyIdAndExpenseDescription(companyId, expenseDescription)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Company expenses not found for the specified description",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch company expenses",
			"error":   err.Error(),
		})
	}

	// Return the fetched expenses
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Company expenses fetched successfully",
		"data":    expenses,
	})
}

// ======================

// // Get Company Expenses by Company ID and Currency

func (h *CompanyController) GetCompanyExpensesByCompanyIdAndCurrency(c *fiber.Ctx) error {
	// Retrieve companyId and currency from the request parameters
	companyId := c.Params("id")
	currency := c.Params("currency")

	// Fetch expenses using the repository
	expenses, err := h.repo.FindCompanyExpensesByCompanyIdAndCurrency(companyId, currency)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Company expenses not found for the specified company and currency",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch company expenses",
			"error":   err.Error(),
		})
	}

	// Return the fetched expenses
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Company expenses fetched successfully",
		"data":    expenses,
	})
}

// =======================

// Get Company Expenses by Company ID, Expense Date, and Currency

func (h *CompanyController) GetCompanyExpensesByThree(c *fiber.Ctx) error {
	// Retrieve companyId, expenseDate, and currency from the request parameters
	companyId := c.Params("id")
	expenseDate := c.Params("expense_date")
	currency := c.Params("currency")

	// Fetch company expenses using the repository
	expenses, err := h.repo.FindCompanyExpensesByThree(companyId, expenseDate, currency)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Company expenses not found for the specified company, date, and currency",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch company expenses",
			"error":   err.Error(),
		})
	}

	// Return the fetched expenses
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Company expenses fetched successfully",
		"data":    expenses,
	})
}

// ==================

// Get Company Expenses by Company ID, Expense Description, and Currency

// GetCompanyExpensesFilters handles the request to get company expenses by filters.
func (h *CompanyController) GetCompanyExpensesFilters(c *fiber.Ctx) error {
	companyId := c.Params("id")
	expenseDate := c.Params("expense_date")
	expenseDescription := c.Params("expense_description")
	currency := c.Params("currency")

	// Fetch company expenses using the service layer
	expenses, err := h.repo.GetCompanyExpensesByFour(companyId, expenseDate, expenseDescription, currency)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch company expenses",
			"error":   err.Error(),
		})
	}

	// If no expenses found
	if len(expenses) == 0 {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Company expenses not found for the specified filters",
		})
	}

	// Return the fetched company expenses
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Company expenses fetched successfully",
		"data":    expenses,
	})
}

// ==================

func (h *CompanyController) GetAllExpenses(c *fiber.Ctx) error {
	// Fetch paginated Companies using the repository
	pagination, companies, err := h.repo.GetPaginatedAllExpenses(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve Companies",
			"data":    err.Error(),
		})
	}

	// Return the paginated response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Companies retrieved successfully",
		"data":    companies,
		"pagination": fiber.Map{
			"total_items":  pagination.TotalItems,
			"total_pages":  pagination.TotalPages,
			"current_page": pagination.CurrentPage,
			"limit":        pagination.ItemsPerPage,
		},
	})
}

// ==================================================================================================================

// Create a company location

// CreateCompanyLocation handles the creation of a company location
func (h *CompanyController) CreateCompanyLocation(c *fiber.Ctx) error {
	// Parse the request body into a CompanyLocation struct
	CompanyLocation := new(companyRegistration.CompanyLocation)
	if err := c.BodyParser(CompanyLocation); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input data",
			"data":    err.Error(),
		})
	}

	// Create the company location in the database
	if err := h.repo.CreateCompanyLocation(CompanyLocation); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create location",
			"data":    err.Error(),
		})
	}

	// Return success response
	return c.Status(201).JSON(fiber.Map{
		"status":  "success",
		"message": "Location created successfully",
		"data":    CompanyLocation,
	})
}

// ===============================

// Get Company Expenses by Company ID, Expense Description, and Currency

// GetCompanyExpensesFilters handles the request to get company expenses by filters.
func (h *CompanyController) GetAllCompanyLocations(c *fiber.Ctx) error {
	companyId := c.Params("id")

	// Fetch company expenses using the service layer
	expenses, err := h.repo.GetAllCompanyLocations(companyId)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch company locations",
			"error":   err.Error(),
		})
	}

	// If no expenses found
	if len(expenses) == 0 {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Company locations not found for the specified filters",
		})
	}

	// Return the fetched company expenses
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Company locations fetched successfully",
		"data":    expenses,
	})
}

// ===============================

func (h *CompanyController) GetLocationByCompanyId(c *fiber.Ctx) error {
	companyId := c.Params("companyId")
	id := c.Params("id")

	// Fetch the company expense from the repository
	expense, err := h.repo.GetLocationByCompanyId(id, companyId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Location not found for the specified company",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch the location",
			"error":   err.Error(),
		})
	}

	// Return the fetched expense
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Location fetched successfully",
		"data":    expense,
	})
}

// ===============================

func (h *CompanyController) UpdateCompanyLocation(c *fiber.Ctx) error {
	// Define a struct for input validation
	type UpdateCompanyLocationInput struct {
		Address   string `json:"address" validate:"required"`
		Telephone string `json:"telephone" validate:"required"`
		Country   string `json:"country" validate:"required"`
		UpdatedBy string `json:"updated_by" validate:"required"`
	}

	// Parse the Location ID from the request parameters
	locationID := c.Params("id")

	// Parse and validate the request body
	var input UpdateCompanyLocationInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input",
			"error":   err.Error(),
		})
	}

	// Use a validation library to validate the input
	if validationErr := utils.ValidateStruct(input); validationErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Validation failed",
			"errors":  validationErr,
		})
	}

	// Fetch the location record using the repository
	location, err := h.repo.FindLocationById(locationID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "location not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch location",
			"error":   err.Error(),
		})
	}

	// Update the location fields
	location.Address = input.Address
	location.Country = input.Country
	location.Telephone = input.Telephone
	location.UpdatedBy = input.UpdatedBy

	// Save the updated location using the repository
	if err := h.repo.UpdateCompanyLocation(location); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update location",
			"error":   err.Error(),
		})
	}

	// Return the updated location
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "location updated successfully",
		"data":    location,
	})
}

// ===============================

// DeleteCompanyByID deletes a company by its ID
func (h *CompanyController) DeleteLocationByID(c *fiber.Ctx) error {
	// Get the company ID from the route parameters
	id := c.Params("id")

	// Find the location in the database
	location, err := h.repo.FindLocationById(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{
				"status":  "error",
				"message": "location not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to find location",
			"data":    err.Error(),
		})
	}

	// Delete the location
	if err := h.repo.DeleteLocationByID(id); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to delete port",
			"data":    err.Error(),
		})
	}

	// Return success response
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "location deleted successfully",
		"data":    location,
	})
}
