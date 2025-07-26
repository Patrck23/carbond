package controllers

import (
	"archive/zip"
	"car-bond/internals/models/alertRegistration"
	"car-bond/internals/models/carRegistration"
	"car-bond/internals/repository"
	"car-bond/internals/utils"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CarController struct {
	repo repository.CarRepository
}

func NewCarController(repo repository.CarRepository) *CarController {
	return &CarController{repo: repo}
}

// ==================

// ConvertCarToTransaction converts car updates to a transaction record
func ConvertUpdateCarToTransaction(car *carRegistration.Car) *alertRegistration.Transaction {
	fromCompanyID := uint(0)
	if car.FromCompanyID != nil {
		fromCompanyID = *car.FromCompanyID
	}

	toCompanyID := uint(0)
	if car.ToCompanyID != nil {
		toCompanyID = *car.ToCompanyID
	}

	transactionType := "Storage"
	if toCompanyID > 0 {
		transactionType = "InTransit"
	}

	return &alertRegistration.Transaction{
		CarChasisNumber: car.ChasisNumber,
		FromCompanyId:   fromCompanyID,
		ToCompanyId:     toCompanyID,
		CreatedBy:       car.UpdatedBy,
		UpdatedBy:       car.UpdatedBy,
		ViewStatus:      false,
		TransactionType: transactionType,
	}
}

// ConvertCarToTransaction maps Car struct to Transaction struct
func ConvertCarToTransaction(car *carRegistration.Car) *alertRegistration.Transaction {
	return &alertRegistration.Transaction{
		CarChasisNumber: car.ChasisNumber,
		TransactionType: getTransactionType(car),
		FromCompanyId:   getFromCompanyID(car),
		ToCompanyId:     getToCompanyID(car),
		ViewStatus:      false,
		CreatedBy:       car.CreatedBy,
		UpdatedBy:       car.UpdatedBy,
	}
}

// Determine transaction type based on car status
func getTransactionType(car *carRegistration.Car) string { // InStock, Sold, Exported
	switch car.CarStatusJapan {
	case "Sold":
		return "Sale"
	case "InTransit":
		return "Export"
	case "InStock":
		return "Storage"
	}
	return "Unknown"
}

// Get FromCompanyID (defaults to 0 if nil)
func getFromCompanyID(car *carRegistration.Car) uint {
	if car.FromCompanyID != nil {
		return *car.FromCompanyID
	}
	return 0
}

// Get ToCompanyID (defaults to 0 if nil)
func getToCompanyID(car *carRegistration.Car) uint {
	if car.ToCompanyID != nil {
		return *car.ToCompanyID
	}
	return 0
}

// ========================

// ===================

func (h *CarController) GetAllCars(c *fiber.Ctx) error {
	// Fetch paginated cars using the repository
	pagination, cars, err := h.repo.GetPaginatedCars(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve cars",
			"data":    err.Error(),
		})
	}

	// Initialize a response slice to hold cars with their ports and expenses
	var response []fiber.Map

	// Iterate over all cars to fetch associated car ports and expenses
	for _, car := range cars {

		expenses, err := h.repo.GetCarExpenses(car.ID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to retrieve expenses for car ID " + strconv.Itoa(int(car.ID)),
				"data":    err.Error(),
			})
		}

		// Combine car, ports, and expenses into a single response map
		response = append(response, fiber.Map{
			"car":      car,
			"expenses": expenses,
		})
	}

	// Return the paginated response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Cars and associated data retrieved successfully",
		"data":    response,
		"pagination": fiber.Map{
			"total_items":  pagination.TotalItems,
			"total_pages":  pagination.TotalPages,
			"current_page": pagination.CurrentPage,
			"limit":        pagination.ItemsPerPage,
		},
	})
}

// ====================

// GetSingleCar fetches a car with its associated ports and expenses from the database
func (h *CarController) GetSingleCar(c *fiber.Ctx) error {
	// Get the car ID from the route parameters
	id := c.Params("id")

	// Fetch the car by ID
	car, err := h.repo.GetCarByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Car not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve car",
			"data":    err.Error(),
		})
	}

	// Fetch expenses associated with the car
	expenses, err := h.repo.GetCarExpenses(car.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve expenses",
			"data":    err.Error(),
		})
	}

	// Prepare the response
	response := fiber.Map{
		"car":      car,
		"expenses": expenses,
	}

	// Return the response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Car and associated data retrieved successfully",
		"data":    response,
	})
}

// ====================

func (h *CarController) GetSingleCarByChasisNumber(c *fiber.Ctx) error {
	// Get the car ID from the route parameters
	ChasisNumber := c.Params("ChasisNumber")

	// Fetch the car by ID
	car, err := h.repo.GetCarByVin(ChasisNumber)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Car not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve car",
			"data":    err.Error(),
		})
	}

	// Fetch expenses associated with the car
	expenses, err := h.repo.GetCarExpenses(car.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve expenses",
			"data":    err.Error(),
		})
	}

	// Prepare the response
	response := fiber.Map{
		"car":      car,
		"expenses": expenses,
	}

	// Return the response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Car and associated data retrieved successfully",
		"data":    response,
	})
}

// ====================

// ===================

type UpdateCarPayload2 struct {
	BrokerName       string `json:"broker_name"`
	BrokerNumber     string `json:"broker_number"`
	NumberPlate      string `json:"number_plate"`
	CarTracker       bool   `json:"car_tracker"`
	CarStatus        string `json:"car_status"`
	CarPaymentStatus string `json:"car_payment_status"`
	CustomerID       int    `json:"customer_id"`
	UpdatedBy        string `json:"updated_by"`
}

// UpdateCar handler function
func (h *CarController) UpdateCar2(c *fiber.Ctx) error {
	// Get the car ID from the route parameters
	id := c.Params("id")

	// Find the car in the database
	car, err := h.repo.GetCarByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{
				"status":  "error",
				"message": "Car not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve car",
			"data":    err.Error(),
		})
	}

	// Parse the request body into the UpdateCarPayload struct
	var payload UpdateCarPayload2
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input",
			"data":    err.Error(),
		})
	}

	// Update the car fields using the payload
	updateCar2Fields(&car, payload) // Pass the parsed payload

	// Save the changes to the database
	if err := h.repo.UpdateCar(&car); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update car",
			"data":    err.Error(),
		})
	}

	// Return the updated car
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Car updated successfully",
		"data":    car,
	})
}

// UpdateCarFields updates the fields of a car using the updateCar struct
func updateCar2Fields(car *carRegistration.Car, updateCarData UpdateCarPayload2) {
	// Update the car fields
	car.BrokerName = updateCarData.BrokerName
	car.BrokerNumber = updateCarData.BrokerNumber
	car.NumberPlate = updateCarData.NumberPlate
	car.CarTracker = updateCarData.CarTracker
	car.CarStatus = updateCarData.CarStatus
	car.CarPaymentStatus = updateCarData.CarPaymentStatus
	// Assign foreign keys if provided
	if updateCarData.CustomerID != 0 {
		car.CustomerID = &updateCarData.CustomerID
	} else {
		car.CustomerID = nil
	}
	car.UpdatedBy = updateCarData.UpdatedBy
}

// ====================

type UpdateCarPayload3 struct {
	CarShippingInvoiceID uint   `json:"car_shipping_invoice_id"`
	UpdatedBy            string `json:"updated_by"`
}

// UpdateCar handler function
func (h *CarController) UpdateCar3(c *fiber.Ctx) error {
	id := c.Params("id")

	car, err := h.repo.GetCarByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{
				"status":  "error",
				"message": "Car not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve car",
			"data":    err.Error(),
		})
	}

	var payloadInv UpdateCarPayload3
	if err := c.BodyParser(&payloadInv); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input",
			"data":    err.Error(),
		})
	}

	targetInvoiceID := payloadInv.CarShippingInvoiceID

	if targetInvoiceID != 0 {
		// Check if the invoice exists and is not locked
		invoice, err := h.repo.GetInvoiceByID(targetInvoiceID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Status(404).JSON(fiber.Map{
					"status":  "error",
					"message": "Target invoice not found",
				})
			}
			return c.Status(500).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to retrieve invoice",
				"data":    err.Error(),
			})
		}

		if invoice.Locked {
			return c.Status(400).JSON(fiber.Map{
				"status":  "error",
				"message": "Invoice is locked and cannot be modified",
			})
		}

		// Ensure car is not already assigned to another invoice
		if car.CarShippingInvoiceID != nil && *car.CarShippingInvoiceID != targetInvoiceID {
			return c.Status(400).JSON(fiber.Map{
				"status":  "error",
				"message": "Car is already assigned to another invoice",
			})
		}
	}

	// Update car fields
	updateCar3Fields(&car, payloadInv)

	// Save the changes
	if err := h.repo.UpdateCar(&car); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update car",
			"data":    err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Car updated successfully",
		"data":    car,
	})
}

// UpdateCarFields updates the fields of a car using the updateCar struct
func updateCar3Fields(car *carRegistration.Car, updateCarData UpdateCarPayload3) {
	// Update the car field
	if updateCarData.CarShippingInvoiceID != 0 {
		car.CarShippingInvoiceID = &updateCarData.CarShippingInvoiceID
	} else {
		car.CarShippingInvoiceID = nil
	}
	car.UpdatedBy = updateCarData.UpdatedBy
}

// ====================

// DeleteCarByID deletes a car by its ID
func (h *CarController) DeleteCarByID(c *fiber.Ctx) error {
	// Get the car ID from the route parameters
	id := c.Params("id")

	// Find the car in the database
	car, err := h.repo.GetCarByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{
				"status":  "error",
				"message": "Car not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to find car",
			"data":    err.Error(),
		})
	}

	// Delete the car
	if err := h.repo.DeleteByID(id); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to delete car",
			"data":    err.Error(),
		})
	}

	// Return success response
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Car deleted successfully",
		"data":    car,
	})
}

// Car Expense
// ===================================================================================================
// Create a car expense

// CreateCarExpense handles the creation of a car expense
func (h *CarController) CreateCarExpense(c *fiber.Ctx) error {
	// Parse the request body into a CarExpense struct
	carExpense := new(carRegistration.CarExpense)
	if err := c.BodyParser(carExpense); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input data",
			"data":    err.Error(),
		})
	}

	// Create the car expense in the database
	if err := h.repo.CreateCarExpense(carExpense); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create car expense",
			"data":    err.Error(),
		})
	}

	// Return success response
	return c.Status(201).JSON(fiber.Map{
		"status":  "success",
		"message": "Car expense created successfully",
		"data":    carExpense,
	})
}

// ========================

// CreateCarExpenses handles the creation of multiple car expenses
func (h *CarController) CreateCarExpenses(c *fiber.Ctx) error {
	// Parse the request body into a slice of CarExpense structs
	var carExpenses []carRegistration.CarExpense
	if err := c.BodyParser(&carExpenses); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input data",
			"data":    err.Error(),
		})
	}

	// Validate if there are any expenses in the request
	if len(carExpenses) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "No car expenses provided",
		})
	}

	// Insert each car expense into the database
	for _, expense := range carExpenses {
		if err := h.repo.CreateCarExpense(&expense); err != nil {
			return c.Status(500).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to create car expense",
				"data":    err.Error(),
			})
		}
	}

	// Return success response
	return c.Status(201).JSON(fiber.Map{
		"status":  "success",
		"message": "Car expenses created successfully",
		"data":    carExpenses,
	})
}

// ========================

func (h *CarController) GetAllCarExpenses(c *fiber.Ctx) error {
	// Fetch paginated expenses using the repository
	pagination, expenses, err := h.repo.GetPaginatedExpenses(c)
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

// =======================

// Get Car Expenses by ID

func (h *CarController) GetCarExpenseById(c *fiber.Ctx) error {
	// Retrieve the expense ID and Car ID from the request parameters
	id := c.Params("id")
	carId := c.Params("carId")

	// Fetch the car expense from the repository
	expense, err := h.repo.FindCarExpenseByIdAndCarId(id, carId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Expense not found for the specified car",
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

// ===============
// Get Car Expenses by Car ID

func (h *CarController) GetCarExpensesByCarId(c *fiber.Ctx) error {
	// Retrieve carId from the request parameters
	carId := c.Params("carId")

	// Fetch paginated car expenses using the repository
	pagination, expenses, err := h.repo.GetPaginatedExpensesByCarId(c, carId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve car expenses",
			"error":   err.Error(),
		})
	}

	// Handle the case where no expenses are found
	if len(expenses) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "No expenses found for the specified car",
		})
	}

	// Return a success response with paginated data
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Car expenses fetched successfully",
		"data":    expenses,
		"pagination": fiber.Map{
			"total_items":  pagination.TotalItems,
			"total_pages":  pagination.TotalPages,
			"current_page": pagination.CurrentPage,
			"limit":        pagination.ItemsPerPage,
		},
	})
}

// ====================

func (h *CarController) UpdateCarExpense(c *fiber.Ctx) error {
	// Define a struct for input validation
	type UpdateCarExpenseInput struct {
		Description   string  `json:"description" validate:"required"`
		Currency      string  `json:"currency" validate:"required"`
		Amount        float64 `json:"amount" validate:"required,gt=0"`
		DollarRate    float64 `json:"dollar_rate"`
		ExpenseDate   string  `json:"expense_date" validate:"required"`
		CompanyName   string  `json:"company_name"`
		Destination   string  `json:"destination"`
		ExpenseVAT    float64 `json:"expense_vat"`
		ExpenseRemark string  `json:"expense_remark"`
		UpdatedBy     string  `json:"updated_by" validate:"required"`
	}

	// Parse the expense ID from the request parameters
	expenseID := c.Params("id")

	// Parse and validate the request body
	var input UpdateCarExpenseInput
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
	expense, err := h.repo.FindCarExpenseById(expenseID)
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
	expense.DollarRate = input.DollarRate
	expense.ExpenseDate = input.ExpenseDate
	expense.CompanyName = input.CompanyName
	expense.ExpenseVAT = input.ExpenseVAT
	expense.Destination = input.Destination
	expense.ExpenseRemark = input.ExpenseRemark
	expense.UpdatedBy = input.UpdatedBy

	// Save the updated expense using the repository
	if err := h.repo.UpdateCarExpense(expense); err != nil {
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

// ============

// Delete car Expenses by ID

func (h *CarController) DeleteCarExpenseById(c *fiber.Ctx) error {
	// Parse carId and expenseId from the request parameters
	carId := c.Params("carId")
	expenseId := c.Params("id")

	// Check if the expense exists and belongs to the specified car
	expense, err := h.repo.FindCarExpenseByCarAndId(carId, expenseId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Car expense not found or does not belong to the specified car",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch car expense",
			"error":   err.Error(),
		})
	}

	// Delete the expense using the repository
	if err := h.repo.DeleteCarExpense(expense); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to delete car expense",
			"error":   err.Error(),
		})
	}

	// Return success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Car expense deleted successfully",
	})
}

// =====================
// Get Car Expenses by Car ID and Expense Date

func (h *CarController) GetCarExpensesByCarIdAndExpenseDate(c *fiber.Ctx) error {
	// Retrieve carId and expenseDate from the request parameters
	carId := c.Params("id")
	expenseDate := c.Params("expense_date")

	// Fetch expenses using the repository
	expenses, err := h.repo.FindCarExpensesByCarIdAndExpenseDate(carId, expenseDate)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Car expenses not found for the specified date",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch car expenses",
			"error":   err.Error(),
		})
	}

	// Return the fetched expenses
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Car expenses fetched successfully",
		"data":    expenses,
	})
}

// =====================
// Get Car Expenses by Car ID and Expense Description

func (h *CarController) GetCarExpensesByCarIdAndExpenseDescription(c *fiber.Ctx) error {
	// Retrieve carId and expenseDescription from the request parameters
	carId := c.Params("id")
	expenseDescription := c.Params("expense_description")

	// Fetch expenses using the repository
	expenses, err := h.repo.FindCarExpensesByCarIdAndExpenseDescription(carId, expenseDescription)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Car expenses not found for the specified description",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch car expenses",
			"error":   err.Error(),
		})
	}

	// Return the fetched expenses
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Car expenses fetched successfully",
		"data":    expenses,
	})
}

// ====================
// Get Car Expenses by Car ID and Currency

func (h *CarController) GetCarExpensesByCarIdAndCurrency(c *fiber.Ctx) error {
	// Retrieve carId and currency from the request parameters
	carId := c.Params("id")
	currency := c.Params("currency")

	// Fetch expenses using the repository
	expenses, err := h.repo.FindCarExpensesByCarIdAndCurrency(carId, currency)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Car expenses not found for the specified car and currency",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch car expenses",
			"error":   err.Error(),
		})
	}

	// Return the fetched expenses
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Car expenses fetched successfully",
		"data":    expenses,
	})
}

// ====================
// Get Car Expenses by Car ID, Expense Date, and Currency

func (h *CarController) GetCarExpensesByThree(c *fiber.Ctx) error {
	// Retrieve carId, expenseDate, and currency from the request parameters
	carId := c.Params("id")
	expenseDate := c.Params("expense_date")
	currency := c.Params("currency")

	// Fetch car expenses using the repository
	expenses, err := h.repo.FindCarExpensesByThree(carId, expenseDate, currency)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Car expenses not found for the specified car, date, and currency",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch car expenses",
			"error":   err.Error(),
		})
	}

	// Return the fetched expenses
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Car expenses fetched successfully",
		"data":    expenses,
	})
}

// ==================
// Get Car Expenses by Car ID, Expense Date, Expense Description, and Currency

// GetCarExpensesFilters handles the request to get car expenses by filters.
func (h *CarController) GetCarExpensesFilters(c *fiber.Ctx) error {
	carId := c.Params("id")
	expenseDate := c.Params("expense_date")
	expenseDescription := c.Params("expense_description")
	currency := c.Params("currency")

	// Fetch car expenses using the service layer
	expenses, err := h.repo.GetCarExpensesByFour(carId, expenseDate, expenseDescription, currency)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch car expenses",
			"error":   err.Error(),
		})
	}

	// If no expenses found
	if len(expenses) == 0 {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Car expenses not found for the specified filters",
		})
	}

	// Return the fetched car expenses
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Car expenses fetched successfully",
		"data":    expenses,
	})
}

// Car Port
// ===================================================================================================

// TotalCarExpense represents an individual expense for a car

// GetTotalCarExpenses handles the request to get car expenses
func (h *CarController) GetTotalCarExpenses(c *fiber.Ctx) error {
	// Convert carId to uint
	carIDStr := c.Params("id")
	carID := utils.StrToUint(carIDStr)

	// Fetch total car expenses using the repository
	totalExpenses, err := h.repo.GetTotalCarExpenses(carID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch total car expenses",
			"error":   err.Error(),
		})
	}

	// Return the total car expenses
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Car expenses fetched successfully",
		"data":    totalExpenses,
	})
}

// ========================================================================================

func (h *CarController) SearchCars(c *fiber.Ctx) error {
	// Call the repository function to get paginated search results
	pagination, cars, err := h.repo.SearchPaginatedCars(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve cars",
			"data":    err.Error(),
		})
	}

	// Initialize a response slice to hold cars with their ports and expenses
	var response []fiber.Map

	// Iterate over all cars to fetch associated car ports and expenses
	for _, car := range cars {

		expenses, err := h.repo.GetCarExpenses(car.ID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to retrieve expenses for car ID " + strconv.Itoa(int(car.ID)),
				"data":    err.Error(),
			})
		}

		// Combine car, ports, and expenses into a single response map
		response = append(response, fiber.Map{
			"car":      car,
			"expenses": expenses,
		})
	}

	// Return the paginated response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Cars and associated data retrieved successfully",
		"data":    response,
		"pagination": fiber.Map{
			"total_items":  pagination.TotalItems,
			"total_pages":  pagination.TotalPages,
			"current_page": pagination.CurrentPage,
			"limit":        pagination.ItemsPerPage,
		},
	})
}

// ===================

func (h *CarController) FetchCarUploads(c *fiber.Ctx) error {
	// Get the car ID from the route parameters
	id := c.Query("id")

	// Find the car in the database
	car, err := h.repo.GetCarByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{
				"status":  "error",
				"message": "Car not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve car",
			"data":    err.Error(),
		})
	}

	// Retrieve photos of the car
	photos, err := h.repo.GetCarPhotosBycarID(car.ID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve photos",
			"data":    err.Error(),
		})
	}

	// Check if photos exist
	if len(photos) == 0 {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "No photos found for this car",
		})
	}

	// Collect file paths
	var filePaths []string
	for _, photo := range photos {
		// Convert stored URL to file path
		filePath := strings.TrimPrefix(photo.URL, "./")
		fmt.Println("Photo URL:", filePath, "Photo ID:", photo.ID)
		filePaths = append(filePaths, filePath) // Assuming FilePath stores the image location
	}

	// Create a ZIP file
	zipFileName := fmt.Sprintf("car_%d_photos.zip", car.ID)
	zipFilePath := filepath.Join(os.TempDir(), zipFileName)
	zipFile, err := os.Create(zipFilePath)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create ZIP file",
			"data":    err.Error(),
		})
	}
	defer zipFile.Close()

	// Write images to ZIP
	zipWriter := zip.NewWriter(zipFile)
	for _, filePath := range filePaths {
		err := addFileToZip(zipWriter, filePath)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to add file to ZIP",
				"data":    err.Error(),
			})
		}
	}
	zipWriter.Close()

	// Serve the ZIP file
	return c.SendFile(zipFilePath)
}

// Helper function to add a file to the ZIP archive
func addFileToZip(zipWriter *zip.Writer, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a zip entry
	zipFileWriter, err := zipWriter.Create(filepath.Base(filePath))
	if err != nil {
		return err
	}

	// Copy the file content into the zip entry
	_, err = io.Copy(zipFileWriter, file)
	return err
}

func (h *CarController) FetchCarUploads64(c *fiber.Ctx) error {
	id := c.Query("id")

	car, err := h.repo.GetCarByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Car not found"})
		}
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed to retrieve car", "data": err.Error()})
	}

	photos, err := h.repo.GetCarPhotosBycarID(car.ID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed to retrieve photos", "data": err.Error()})
	}

	if len(photos) == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No photos found for this car"})
	}

	var images []map[string]string
	for _, photo := range photos {
		// Read image file
		imageData, err := os.ReadFile(photo.URL)
		if err != nil {
			continue
		}

		// Convert image to base64
		base64Image := base64.StdEncoding.EncodeToString(imageData)
		images = append(images, map[string]string{
			"filename": photo.URL,
			"data":     "data:image/jpeg;base64," + base64Image,
		})
	}

	// Return JSON with images
	return c.JSON(fiber.Map{"status": "success", "message": "Car images retrieved", "images": images})
}

// ==========================================

// GetDashboardData returns aggregated statistics for the dashboard
func (h *CarController) GetDashboardData(c *fiber.Ctx) error {
	totalCars, err := h.repo.GetTotalCars()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve total cars",
			"data":    err.Error(),
		})
	}

	disbandedCars, err := h.repo.GetDisbandedCars()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve disbanded cars",
			"data":    err.Error(),
		})
	}

	carsInStock, err := h.repo.GetCarsInStock()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve cars in stock",
			"data":    err.Error(),
		})
	}

	totalMoneySpent, err := h.repo.GetTotalMoneySpent()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve total money spent",
			"data":    err.Error(),
		})
	}

	totalCarExpenses, err := h.repo.GetTotalCarsExpenses()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve total car expenses",
			"data":    err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Dashboard data retrieved successfully",
		"data": fiber.Map{
			"total_cars":         totalCars,
			"disbanded_cars":     disbandedCars,
			"cars_in_stock":      carsInStock,
			"total_money_spent":  totalMoneySpent,
			"total_car_expenses": totalCarExpenses,
		},
	})
}

// ====================================================

// GetDashboardData returns aggregated statistics for the dashboard
func (h *CarController) GetCompanyDashboardData(c *fiber.Ctx) error {
	companyIdStr := c.Params("companyId")
	companyId := utils.StrToUint(companyIdStr)
	if companyId == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid company ID")
	}

	totalCars, err := h.repo.GetComTotalCars(companyId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve total cars",
			"data":    err.Error(),
		})
	}

	carsSold, err := h.repo.GetComCarsSold(companyId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve cars sold",
			"data":    err.Error(),
		})
	}

	carsInStock, err := h.repo.GetComCarsInStock(companyId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve cars in stock",
			"data":    err.Error(),
		})
	}

	totalMoneySpent, err := h.repo.GetComTotalMoneySpent(companyId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve total money spent",
			"data":    err.Error(),
		})
	}

	totalCarExpenses, err := h.repo.GetComTotalCarsExpenses(companyId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve total car expenses",
			"data":    err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Dashboard data retrieved successfully",
		"data": fiber.Map{
			"total_cars":         totalCars,
			"cars_in_stock":      carsInStock,
			"cars_sold":          carsSold,
			"total_money_spent":  totalMoneySpent,
			"total_car_expenses": totalCarExpenses,
		},
	})
}

// ===================================

type CarFormPayload struct {
	ChasisNumber          string `form:"chasis_number"`
	EngineNumber          string `form:"engine_number"`
	EngineCapacity        string `form:"engine_capacity"`
	FrameNumber           string `form:"frame_number"`
	Make                  string `form:"make"`
	CarModel              string `form:"car_model"`
	SeatCapacity          string `form:"seat_capacity"`
	MaximCarry            string `form:"maxim_carry"`
	Weight                string `form:"weight"`
	GrossWeight           string `form:"gross_weight"`
	Length                string `form:"length"`
	Width                 string `form:"width"`
	Height                string `form:"height"`
	ManufactureYear       string `form:"manufacture_year"`
	FirstRegistrationYear string `form:"first_registration_year"`
	Transmission          string `form:"transmission"`
	BodyType              string `form:"body_type"`
	Colour                string `form:"colour"`
	Auction               string `form:"auction"`
	Currency              string `form:"currency"`
	CarMillage            string `form:"car_millage"`
	FuelConsumption       string `form:"fuel_consumption"`
	PowerSteering         string `form:"power_steering"`
	PowerWindow           string `form:"power_window"`
	ABS                   string `form:"abs"`
	ADS                   string `form:"ads"`
	AirBrake              string `form:"air_brake"`
	OilBrake              string `form:"oil_brake"`
	AlloyWheel            string `form:"alloy_wheel"`
	SimpleWheel           string `form:"simple_wheel"`
	Navigation            string `form:"navigation"`
	AC                    string `form:"ac"`
	BidPrice              string `form:"bid_price"`
	VATTax                string `form:"vat_tax"`
	PurchaseDate          string `form:"purchase_date"`
	FromCompanyID         string `form:"from_company_id"`
	ToCompanyID           string `form:"to_company_id"`
	CarShippingInvoiceID  string `form:"car_shipping_invoice_id"`
	OtherEntity           string `form:"other_entity"`
	Destination           string `form:"destination"`
	Port                  string `form:"port"`
	CarStatusJapan        string `form:"car_status_japan"`
	BrokerName            string `form:"broker_name"`
	BrokerNumber          string `form:"broker_number"`
	NumberPlate           string `form:"number_plate"`
	CarTracker            string `form:"car_tracker"`
	CustomerID            string `form:"customer_id"`
	CarStatus             string `form:"car_status"`
	CarPaymentStatus      string `form:"car_payment_status"`
	CreatedBy             string `form:"created_by"`
	UpdatedBy             string `form:"updated_by"`
}

type CreateCarInput struct {
	Car         carRegistration.Car          `json:"car"`
	CarPhotos   []carRegistration.CarPhoto   `json:"car_photos"`
	CarExpenses []carRegistration.CarExpense `json:"car_expenses"`
}

func (h *CarController) CreateCarWithDetails(c *fiber.Ctx) error {
	// Parse form fields
	carPayload := CarFormPayload{}
	if err := c.BodyParser(&carPayload); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status": "error", "message": "Invalid form input", "data": err.Error(),
		})
	}

	car := carRegistration.Car{}
	val := reflect.ValueOf(carPayload)
	typ := reflect.TypeOf(carPayload)

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldName := typ.Field(i).Tag.Get("form")
		fieldValue := field.String()

		switch fieldName {
		case "car_millage", "manufacture_year", "first_registration_year":
			reflect.ValueOf(&car).Elem().FieldByName(typ.Field(i).Name).SetInt(int64(utils.StrToInt(fieldValue)))

		case "maxim_carry", "weight", "gross_weight", "length", "width", "height", "bid_price", "vat_tax":
			reflect.ValueOf(&car).Elem().FieldByName(typ.Field(i).Name).SetFloat(utils.StrToFloat(fieldValue))

		case "power_steering", "power_window", "abs", "ads", "air_brake", "oil_brake", "alloy_wheel", "simple_wheel", "navigation", "ac", "car_tracker":
			reflect.ValueOf(&car).Elem().FieldByName(typ.Field(i).Name).SetBool(utils.StrToBool(fieldValue))

		case "from_company_id", "to_company_id", "car_shipping_invoice_id":
			reflect.ValueOf(&car).Elem().FieldByName(typ.Field(i).Name).Set(reflect.ValueOf(utils.StrToUintPointer(fieldValue)))

		case "customer_id":
			if fieldValue != "" {
				id, err := strconv.Atoi(fieldValue)
				if err == nil {
					reflect.ValueOf(&car).Elem().FieldByName("CustomerID").Set(reflect.ValueOf(&id))
				}
			}

		default:
			reflect.ValueOf(&car).Elem().FieldByName(typ.Field(i).Name).SetString(fieldValue)
		}
	}

	// Save the car
	if err := h.repo.CreateCar(&car); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to save car",
			"data":    err.Error(),
		})
	}

	// Parse car expenses JSON
	expenseJSON := c.FormValue("car_expenses")
	var expenses []carRegistration.CarExpense
	if expenseJSON != "" {
		if err := json.Unmarshal([]byte(expenseJSON), &expenses); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid expenses JSON",
				"data":    err.Error(),
			})
		}
		for i := range expenses {
			expenses[i].CarID = car.ID
		}
		if err := h.repo.CreateCarExpenses(expenses); err != nil {
			return c.Status(500).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to save expenses",
				"data":    err.Error(),
			})
		}
	}

	// Handle photo uploads
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to parse uploaded files",
			"data":    err.Error(),
		})
	}

	uploadDir := "./uploads/car_files"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create upload directory",
			"data":    err.Error(),
		})
	}

	files := form.File["car_photos"]
	var savedPhotos []carRegistration.CarPhoto

	for _, file := range files {
		uniqueName := fmt.Sprintf("%s_%s", uuid.New().String(), file.Filename)
		cleanFileName := strings.ReplaceAll(uniqueName, " ", "_")
		savePath := filepath.Join(uploadDir, cleanFileName)

		// Save file to disk
		if err := c.SaveFile(file, savePath); err != nil {
			return c.Status(500).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to save file",
				"data":    err.Error(),
			})
		}

		// Save only filename/relative path to DB
		savedPhotos = append(savedPhotos, carRegistration.CarPhoto{
			CarID: car.ID,
			URL:   fmt.Sprintf("./uploads/car_files/%s", cleanFileName),
		})
	}

	// Save photo metadata to DB
	if len(savedPhotos) > 0 {
		if err := h.repo.CreateCarPhotos(savedPhotos); err != nil {
			return c.Status(500).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to save photo metadata",
				"data":    err.Error(),
			})
		}
	}

	// Convert Car to Transaction
	transaction := ConvertCarToTransaction(&car)

	// Save to database
	if err := h.repo.CreateAlert(transaction); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create transaction"})
	}

	// Success response
	return c.Status(201).JSON(fiber.Map{
		"status":       "success",
		"message":      "Car created with photo metadata and expenses",
		"car":          car,
		"car_expenses": expenses,
		"car_photos":   savedPhotos,
	})
}

// ===========================================================

func (h *CarController) UpdateCarWithDetails(c *fiber.Ctx) error {
	carIDStr := c.Params("id")
	carID := utils.StrToUint(carIDStr)

	// Parse form fields into payload struct
	var payload CarFormPayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid form input",
			"data":    err.Error(),
		})
	}

	// Build Car struct via reflection
	car := carRegistration.Car{}

	car.ID = carID

	val := reflect.ValueOf(payload)
	typ := reflect.TypeOf(payload)

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldName := typ.Field(i).Tag.Get("form")
		fieldValue := field.String()

		switch fieldName {
		case "car_millage", "manufacture_year", "first_registration_year":
			reflect.ValueOf(&car).Elem().FieldByName(typ.Field(i).Name).SetInt(int64(utils.StrToInt(fieldValue)))

		case "maxim_carry", "weight", "gross_weight", "length", "width", "height", "bid_price", "vat_tax":
			reflect.ValueOf(&car).Elem().FieldByName(typ.Field(i).Name).SetFloat(utils.StrToFloat(fieldValue))

		case "power_steering", "power_window", "abs", "ads", "air_brake", "oil_brake", "alloy_wheel", "simple_wheel", "navigation", "ac", "car_tracker":
			reflect.ValueOf(&car).Elem().FieldByName(typ.Field(i).Name).SetBool(utils.StrToBool(fieldValue))

		case "from_company_id", "to_company_id", "car_shipping_invoice_id":
			reflect.ValueOf(&car).Elem().FieldByName(typ.Field(i).Name).Set(reflect.ValueOf(utils.StrToUintPointer(fieldValue)))

		case "customer_id":
			if fieldValue != "" {
				id, err := strconv.Atoi(fieldValue)
				if err == nil {
					reflect.ValueOf(&car).Elem().FieldByName("CustomerID").Set(reflect.ValueOf(&id))
				}
			}

		default:
			reflect.ValueOf(&car).Elem().FieldByName(typ.Field(i).Name).SetString(fieldValue)
		}
	}

	// Optional logic for export
	if car.ToCompanyID != nil && *car.ToCompanyID != 0 {
		car.CarStatusJapan = "Exported"
		companyName, err := h.repo.GetCompanyNameByID(*car.ToCompanyID)
		if err != nil && err != gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to fetch destination company",
				"data":    err.Error(),
			})
		}
		car.OtherEntity = companyName
	}

	// Parse and update car expenses
	expenseJSON := c.FormValue("car_expenses")
	var expenses []carRegistration.CarExpense
	if expenseJSON != "" {
		if err := json.Unmarshal([]byte(expenseJSON), &expenses); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid car_expenses JSON",
				"data":    err.Error(),
			})
		}
		for i := range expenses {
			expenses[i].CarID = car.ID
		}
	}

	// Update car and its expenses
	if err := h.repo.UpdateCarWithExpenses(&car, expenses); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update car and expenses",
			"data":    err.Error(),
		})
	}

	// Handle uploaded photos
	form, err := c.MultipartForm()
	if err == nil && form.File != nil && len(form.File["car_photos"]) > 0 {
		uploadDir := "./uploads/car_files"
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			return c.Status(500).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to create upload directory",
				"data":    err.Error(),
			})
		}

		// Remove old photos (optional)
		if err := h.repo.DeleteCarPhotos(car.ID); err != nil {
			return c.Status(500).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to delete old photos",
				"data":    err.Error(),
			})
		}

		// Save new photo metadata
		var newPhotos []carRegistration.CarPhoto
		for _, file := range form.File["car_photos"] {
			uniqueName := fmt.Sprintf("%s_%s", uuid.New().String(), file.Filename)
			cleanName := strings.ReplaceAll(uniqueName, " ", "_")
			savePath := filepath.Join(uploadDir, cleanName)

			if err := c.SaveFile(file, savePath); err != nil {
				return c.Status(500).JSON(fiber.Map{
					"status":  "error",
					"message": "Failed to save file",
					"data":    err.Error(),
				})
			}

			newPhotos = append(newPhotos, carRegistration.CarPhoto{
				CarID: car.ID,
				URL:   fmt.Sprintf("./uploads/car_files/%s", cleanName),
			})
		}

		if err := h.repo.CreateCarPhotos(newPhotos); err != nil {
			return c.Status(500).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to save new photo metadata",
				"data":    err.Error(),
			})
		}
	}

	// Convert Car to Transaction
	transaction := ConvertUpdateCarToTransaction(&car)

	if car.ToCompanyID != nil && *car.ToCompanyID != 0 {
		// 1) mark status
		car.CarStatusJapan = "Exported"

		// 2) fetch company name & add to car
		if name, err := h.repo.GetCompanyNameByID(*car.ToCompanyID); err == nil {
			car.OtherEntity = name
		} else if err != gorm.ErrRecordNotFound {
			return c.Status(500).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to fetch destination company",
				"data":    err.Error(),
			})
		}
	}

	// Save to database
	if err := h.repo.CreateAlert(transaction); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create transaction"})
	}

	// Success
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Car and related data updated successfully",
		"car":     car,
	})
}
