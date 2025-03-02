package controllers

import (
	"car-bond/internals/models/carRegistration"
	"car-bond/internals/models/customerRegistration"
	"car-bond/internals/models/saleRegistration"
	"car-bond/internals/utils"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type SaleRepository interface {
	CreateSale(sale *saleRegistration.Sale) error
	GetPaginatedSales(c *fiber.Ctx) (*utils.Pagination, []saleRegistration.Sale, error)
	GetSalePayments(saleID uint) ([]saleRegistration.SalePayment, error)
	GetSalePaymentModes(paymentID uint) ([]saleRegistration.SalePaymentMode, error)
	GetSaleByID(id string) (saleRegistration.Sale, error)
	UpdateSale(sale *saleRegistration.Sale) error
	DeleteByID(id string) error

	// Payment
	CreateInvoice(payment *saleRegistration.SalePayment) error
	GetPaginatedInvoices(c *fiber.Ctx) (*utils.Pagination, []saleRegistration.SalePayment, error)
	FindSalePaymentById(id string) (*saleRegistration.SalePayment, error)
	FindSalePaymentByIdAndSaleId(id, saleId string) (*saleRegistration.SalePayment, error)
	UpdateSalePayment(payment *saleRegistration.SalePayment) error
	DeleteSalePaymentByID(id string) error

	// Payment Mode
	CreatePaymentMode(payment *saleRegistration.SalePaymentMode) error
	GetPaginatedPaymentModes(c *fiber.Ctx) (*utils.Pagination, []saleRegistration.SalePaymentMode, error)
	FindSalePaymentModeByIdAndSalePaymentId(id, salePaymentId string) (*saleRegistration.SalePaymentMode, error)
	FindSalePaymentModeById(id string) (*saleRegistration.SalePaymentMode, error)
	UpdateSalePaymentMode(payment *saleRegistration.SalePaymentMode) error
	DeleteSalePaymentModeByID(id string) error
	GetPaginatedModes(c *fiber.Ctx, mode string) (*utils.Pagination, []saleRegistration.SalePaymentMode, error)

	// Payment Deposit
	CreatePaymentDeposit(payment *saleRegistration.SalePaymentDeposit) error
	GetPaginatedPaymentDeposits(c *fiber.Ctx) (*utils.Pagination, []saleRegistration.SalePaymentDeposit, error)
	FindSalePaymentDepositByIdAndSalePaymentId(id, salePaymentId string) (*saleRegistration.SalePaymentDeposit, error)
	FindSalePaymentDepositById(id string) (*saleRegistration.SalePaymentDeposit, error)
	UpdateSalePaymentDeposit(payment *saleRegistration.SalePaymentDeposit) error
	DeleteSalePaymentDepositByID(id string) error
	GetPaymentDeposits(c *fiber.Ctx, name string) (*utils.Pagination, []saleRegistration.SalePaymentDeposit, error)

	// Customer statement
	GenerateCustomerStatement(customerID uint) (*CustomerStatement, error)
}

type SaleRepositoryImpl struct {
	db *gorm.DB
}

func NewSaleRepository(db *gorm.DB) SaleRepository {
	return &SaleRepositoryImpl{db: db}
}

type SaleController struct {
	repo SaleRepository
}

func NewSaleController(repo SaleRepository) *SaleController {
	return &SaleController{repo: repo}
}

// ============================================

func (r *SaleRepositoryImpl) CreateSale(sale *saleRegistration.Sale) error {
	return r.db.Create(sale).Error
}

func (h *SaleController) CreateCarSale(c *fiber.Ctx) error {
	// Initialize a new Sale instance
	sale := new(saleRegistration.Sale)

	// Parse the request body into the sale instance
	if err := c.BodyParser(sale); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input provided",
			"data":    err.Error(),
		})
	}

	// Attempt to create the sale record using the repository
	if err := h.repo.CreateSale(sale); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create sale",
			"data":    err.Error(),
		})
	}

	// Return the newly created sale record
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Sale created successfully",
		"data":    sale,
	})
}

// =====================

func (r *SaleRepositoryImpl) GetSalePayments(saleID uint) ([]saleRegistration.SalePayment, error) {
	var payments []saleRegistration.SalePayment
	err := r.db.Where("sale_id = ?", saleID).Find(&payments).Error
	return payments, err
}

func (r *SaleRepositoryImpl) GetSalePaymentModes(paymentID uint) ([]saleRegistration.SalePaymentMode, error) {
	var payment_modes []saleRegistration.SalePaymentMode
	err := r.db.Where("sale_payment_id = ?", paymentID).Find(&payment_modes).Error
	return payment_modes, err
}

// ====================

func (r *SaleRepositoryImpl) GetPaginatedSales(c *fiber.Ctx) (*utils.Pagination, []saleRegistration.Sale, error) {
	pagination, sales, err := utils.Paginate(c, r.db.Preload("Car"), saleRegistration.Sale{})
	if err != nil {
		return nil, nil, err
	}
	return &pagination, sales, nil
}

func (h *SaleController) GetAllCarSales(c *fiber.Ctx) error {
	pagination, sales, err := h.repo.GetPaginatedSales(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve sales",
			"data":    err.Error(),
		})
	}

	// Initialize a response slice to hold sales with their payments and modes
	var response []fiber.Map

	// Iterate over all sales to fetch associated payments and payment modes
	for _, sale := range sales {
		// Fetch sale payments associated with the sale
		payments, err := h.repo.GetSalePayments(sale.ID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to retrieve sale payments for sale ID " + strconv.Itoa(int(sale.ID)),
				"error":   err.Error(),
			})
		}

		// Iterate over payments to fetch associated payment modes
		var allPaymentModes []fiber.Map
		for _, payment := range payments {
			paymentModes, err := h.repo.GetSalePaymentModes(payment.ID)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status":  "error",
					"message": "Failed to retrieve payment modes for payment ID " + strconv.Itoa(int(payment.ID)),
					"error":   err.Error(),
				})
			}

			allPaymentModes = append(allPaymentModes, fiber.Map{
				"payment":       payment,
				"payment_modes": paymentModes,
			})
		}

		// Combine sale, payments, and payment modes into a single response map
		response = append(response, fiber.Map{
			"sale":          sale,
			"sale_payments": allPaymentModes,
		})
	}

	// Return the paginated response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "sales retrieved successfully",
		"data":    response,
		"pagination": fiber.Map{
			"total_items":  pagination.TotalItems,
			"total_pages":  pagination.TotalPages,
			"current_page": pagination.CurrentPage,
			"limit":        pagination.ItemsPerPage,
		},
	})
}

// ==============

// Get a single car sale by ID

func (r *SaleRepositoryImpl) GetSaleByID(id string) (saleRegistration.Sale, error) {
	var sale saleRegistration.Sale
	err := r.db.Preload("Car").First(&sale, "id = ?", id).Error
	return sale, err
}

// GetCarSale fetches a sale with its associated contacts and addresses from the database
func (h *SaleController) GetCarSale(c *fiber.Ctx) error {
	// Get the sale ID from the route parameters
	id := c.Params("id")

	// Fetch the sale by ID
	sale, err := h.repo.GetSaleByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Sale not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve Sale",
			"data":    err.Error(),
		})
	}

	// Fetch Sale addresses associated with the Sale
	payments, err := h.repo.GetSalePayments(sale.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve sale addresses",
			"data":    err.Error(),
		})
	}

	var allPaymentModes []fiber.Map
	for _, payment := range payments {
		paymentModes, err := h.repo.GetSalePaymentModes(payment.ID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to retrieve payment modes for payment ID " + strconv.Itoa(int(payment.ID)),
				"error":   err.Error(),
			})
		}

		allPaymentModes = append(allPaymentModes, fiber.Map{
			"payment":       payment,
			"payment_modes": paymentModes,
		})
	}

	// Prepare the response
	response := fiber.Map{
		"sale":     sale,
		"payments": allPaymentModes,
	}

	// Return the response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Sale and associated data retrieved successfully",
		"data":    response,
	})
}

// =====================================

func (r *SaleRepositoryImpl) UpdateSale(sale *saleRegistration.Sale) error {
	return r.db.Save(sale).Error
}

// Define the UpdateSale struct
type UpdateSalePayload struct {
	TotalPrice    float64 `json:"total_price"`
	DollarRate    float64 `json:"dollar_rate"`
	SaleDate      string  `json:"sale_date"`
	CarID         uint    `json:"car_id"`
	CompanyID     int     `json:"company_id"`
	IsFullPayment bool    `json:"is_full_payment"`
	InitalPayment float64 `json:"initial_payment"`
	PaymentPeriod int     `json:"payment_period"`
	UpdatedBy     string  `json:"updated_by"`
}

// UpdateSale handler function
func (h *SaleController) UpdateSale(c *fiber.Ctx) error {
	// Get the sale ID from the route parameters
	id := c.Params("id")

	// Find the sale in the database
	sale, err := h.repo.GetSaleByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{
				"status":  "error",
				"message": "Sale not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve sale",
			"data":    err.Error(),
		})
	}

	// Parse the request body into the UpdateSalePayload struct
	var payload UpdateSalePayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input",
			"data":    err.Error(),
		})
	}

	// Update the sale fields using the payload
	updateSaleFields(&sale, payload) // Pass the parsed payload

	// Save the changes to the database
	if err := h.repo.UpdateSale(&sale); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update sale",
			"data":    err.Error(),
		})
	}

	// Return the updated sale
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "sale updated successfully",
		"data":    sale,
	})
}

// UpdateSaleFields updates the fields of a Sale using the UpdateSale struct
func updateSaleFields(sale *saleRegistration.Sale, updateSaleData UpdateSalePayload) {
	sale.TotalPrice = updateSaleData.TotalPrice
	sale.DollarRate = updateSaleData.DollarRate
	sale.SaleDate = updateSaleData.SaleDate
	sale.CarID = updateSaleData.CarID
	sale.CompanyID = updateSaleData.CompanyID
	sale.IsFullPayment = updateSaleData.IsFullPayment
	sale.InitalPayment = updateSaleData.InitalPayment
	sale.PaymentPeriod = updateSaleData.PaymentPeriod
	sale.UpdatedBy = updateSaleData.UpdatedBy
}

// ============================

// DeleteByID deletes a sale by ID
func (r *SaleRepositoryImpl) DeleteByID(id string) error {
	if err := r.db.Delete(&saleRegistration.Sale{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

// DeleteSaleByID deletes a Sale by its ID
func (h *SaleController) DeleteSaleByID(c *fiber.Ctx) error {
	// Get the Sale ID from the route parameters
	id := c.Params("id")

	// Find the Sale in the database
	sale, err := h.repo.GetSaleByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{
				"status":  "error",
				"message": "Sale not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to find sale",
			"data":    err.Error(),
		})
	}

	// Delete the Sale
	if err := h.repo.DeleteByID(id); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to delete Sale",
			"data":    err.Error(),
		})
	}

	// Return success response
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Sale deleted successfully",
		"data":    sale,
	})
}

// ===============================================================================================

//Create an invoice for a customer

// CreateCustomerContact creates a new payment deposit in the database
func (r *SaleRepositoryImpl) CreateInvoice(payment *saleRegistration.SalePayment) error {
	return r.db.Create(payment).Error
}

// CreateSalePayment handles the creation of a payment deposit
func (h *SaleController) CreateInvoice(c *fiber.Ctx) error {
	// Parse the request body into a salePayment struct
	salePayment := new(saleRegistration.SalePayment)
	if err := c.BodyParser(salePayment); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input data",
			"data":    err.Error(),
		})
	}

	// Create the customer address in the database
	if err := h.repo.CreateInvoice(salePayment); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create payment",
			"data":    err.Error(),
		})
	}

	// Return success response
	return c.Status(201).JSON(fiber.Map{
		"status":  "success",
		"message": "Payment created successfully",
		"data":    salePayment,
	})
}

// =====================

// Get all invoices

func (r *SaleRepositoryImpl) GetPaginatedInvoices(c *fiber.Ctx) (*utils.Pagination, []saleRegistration.SalePayment, error) {
	pagination, payments, err := utils.Paginate(c, r.db.Preload("Sale"), saleRegistration.SalePayment{})
	if err != nil {
		return nil, nil, err
	}
	return &pagination, payments, nil
}

func (h *SaleController) GetSalePayments(c *fiber.Ctx) error {

	// Fetch paginated payments using the repository
	pagination, invoices, err := h.repo.GetPaginatedInvoices(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve payments",
			"data":    err.Error(),
		})
	}

	// Return the paginated response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Payments and associated data retrieved successfully",
		"data":    invoices,
		"pagination": fiber.Map{
			"total_items":  pagination.TotalItems,
			"total_pages":  pagination.TotalPages,
			"current_page": pagination.CurrentPage,
			"limit":        pagination.ItemsPerPage,
		},
	})
}

// ========================

func (r *SaleRepositoryImpl) FindSalePaymentByIdAndSaleId(id, saleId string) (*saleRegistration.SalePayment, error) {
	var payment saleRegistration.SalePayment
	result := r.db.Preload("Sale").Where("id = ? AND sale_id = ?", id, saleId).First(&payment)
	if result.Error != nil {
		return nil, result.Error
	}
	return &payment, nil
}

func (h *SaleController) FindSalePaymentByIdAndSaleId(c *fiber.Ctx) error {
	// Retrieve the payment ID and Company ID from the request parameters
	id := c.Params("id")
	saleId := c.Params("saleId")

	// Fetch the company payment from the repository
	payment, err := h.repo.FindSalePaymentByIdAndSaleId(id, saleId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Payment not found for the specified sale",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch the payment",
			"error":   err.Error(),
		})
	}

	// Return the fetched payment
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "payment fetched successfully",
		"data":    payment,
	})
}

// ========================

// Update an invoice

func (r *SaleRepositoryImpl) FindSalePaymentById(id string) (*saleRegistration.SalePayment, error) {
	var payment saleRegistration.SalePayment
	if err := r.db.First(&payment, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *SaleRepositoryImpl) UpdateSalePayment(payment *saleRegistration.SalePayment) error {
	return r.db.Save(payment).Error
}

func (h *SaleController) UpdateSalePayment(c *fiber.Ctx) error {
	// Define a struct for input validation
	type UpdateSalePaymentInput struct {
		AmountPayed float64 `json:"amount_payed" validate:"required"`
		PaymentDate string  `json:"payment_date" validate:"required"`
		SaleID      uint    `json:"sale_id" validate:"required"`
		UpdatedBy   string  `json:"updated_by" validate:"required"`
	}

	// Parse the payment ID from the request parameters
	paymentID := c.Params("id")

	// Parse and validate the request body
	var input UpdateSalePaymentInput
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

	// Fetch the payment record using the repository
	payment, err := h.repo.FindSalePaymentById(paymentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "payment not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch payment",
			"error":   err.Error(),
		})
	}

	// Update the payment fields
	payment.AmountPayed = input.AmountPayed
	payment.PaymentDate = input.PaymentDate
	payment.SaleID = input.SaleID
	payment.UpdatedBy = input.UpdatedBy

	// Save the updated payment using the repository
	if err := h.repo.UpdateSalePayment(payment); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update payment",
			"error":   err.Error(),
		})
	}

	// Return the updated payment
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "payment updated successfully",
		"data":    payment,
	})
}

// =============================

// Delete salePayment by ID
func (r *SaleRepositoryImpl) DeleteSalePaymentByID(id string) error {
	if err := r.db.Delete(&saleRegistration.SalePayment{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

// DeleteCompanyByID deletes a company by its ID
func (h *SaleController) DeleteSalePaymentByID(c *fiber.Ctx) error {
	// Get the company ID from the route parameters
	id := c.Params("id")

	// Find the SalePayment in the database
	salePayment, err := h.repo.FindSalePaymentById(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{
				"status":  "error",
				"message": "salePayment not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to find salePayment",
			"data":    err.Error(),
		})
	}

	// Delete the salePayment
	if err := h.repo.DeleteSalePaymentByID(id); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to delete port",
			"data":    err.Error(),
		})
	}

	// Return success response
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "salePayment deleted successfully",
		"data":    salePayment,
	})
}

// ===============================================================================================

// Get payment by ModeOfPayment

// CreateCustomerContact creates a new payment mode in the database
func (r *SaleRepositoryImpl) CreatePaymentMode(payment *saleRegistration.SalePaymentMode) error {
	return r.db.Create(payment).Error
}

// CreateSalePayment handles the creation of a payment mode
func (h *SaleController) CreatePaymentMode(c *fiber.Ctx) error {
	// Parse the request body into a salePayment struct
	salePayment := new(saleRegistration.SalePaymentMode)
	if err := c.BodyParser(salePayment); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input data",
			"data":    err.Error(),
		})
	}

	// Create the customer address in the database
	if err := h.repo.CreatePaymentMode(salePayment); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create Payment mode",
			"data":    err.Error(),
		})
	}

	// Return success response
	return c.Status(201).JSON(fiber.Map{
		"status":  "success",
		"message": "Payment mode address created successfully",
		"data":    salePayment,
	})
}

// ==============================

func (r *SaleRepositoryImpl) GetPaginatedPaymentModes(c *fiber.Ctx) (*utils.Pagination, []saleRegistration.SalePaymentMode, error) {
	pagination, payments, err := utils.Paginate(c, r.db, saleRegistration.SalePaymentMode{})
	if err != nil {
		return nil, nil, err
	}
	return &pagination, payments, nil
}

func (h *SaleController) GetSalePaymentModes(c *fiber.Ctx) error {

	// Fetch paginated payments using the repository
	pagination, paymentModes, err := h.repo.GetPaginatedPaymentModes(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve payments modes",
			"data":    err.Error(),
		})
	}

	// Return the paginated response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Payment modes and associated data retrieved successfully",
		"data":    paymentModes,
		"pagination": fiber.Map{
			"total_items":  pagination.TotalItems,
			"total_pages":  pagination.TotalPages,
			"current_page": pagination.CurrentPage,
			"limit":        pagination.ItemsPerPage,
		},
	})
}

// =====================

func (r *SaleRepositoryImpl) FindSalePaymentModeByIdAndSalePaymentId(id, salePaymentId string) (*saleRegistration.SalePaymentMode, error) {
	var paymentMode saleRegistration.SalePaymentMode
	result := r.db.Preload("SalePayment").Where("id = ? AND sale_payment_id = ?", id, salePaymentId).First(&paymentMode)
	if result.Error != nil {
		return nil, result.Error
	}
	return &paymentMode, nil
}

func (h *SaleController) FindSalePaymentModeByIdAndSalePaymentId(c *fiber.Ctx) error {
	// Retrieve the payment ID and Company ID from the request parameters
	id := c.Params("id")
	salePaymentId := c.Params("salePaymentId")

	// Fetch the company payment from the repository
	payment, err := h.repo.FindSalePaymentModeByIdAndSalePaymentId(id, salePaymentId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Mode not found for the specified payment",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch the mode",
			"error":   err.Error(),
		})
	}

	// Return the fetched payment
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "payment mode fetched successfully",
		"data":    payment,
	})
}

// =============

func (r *SaleRepositoryImpl) GetPaginatedModes(c *fiber.Ctx, mode string) (*utils.Pagination, []saleRegistration.SalePaymentMode, error) {
	pagination, modes, err := utils.Paginate(c, r.db.Preload("SalePayment").Where("mode_of_payment = ?", mode), saleRegistration.SalePaymentMode{})
	if err != nil {
		return nil, nil, err
	}
	return &pagination, modes, nil
}

func (h *SaleController) GetPaymentModesByMode(c *fiber.Ctx) error {
	mode := c.Params("mode")

	// Fetch paginated contacts using the repository
	pagination, modes, err := h.repo.GetPaginatedModes(c, mode)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve payment modes",
			"data":    err.Error(),
		})
	}

	// Return the paginated response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Payment modes and associated data retrieved successfully",
		"data":    modes,
		"pagination": fiber.Map{
			"total_items":  pagination.TotalItems,
			"total_pages":  pagination.TotalPages,
			"current_page": pagination.CurrentPage,
			"limit":        pagination.ItemsPerPage,
		},
	})
}

// =====================

func (r *SaleRepositoryImpl) FindSalePaymentModeById(id string) (*saleRegistration.SalePaymentMode, error) {
	var mode saleRegistration.SalePaymentMode
	if err := r.db.First(&mode, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &mode, nil
}

func (r *SaleRepositoryImpl) UpdateSalePaymentMode(payment *saleRegistration.SalePaymentMode) error {
	return r.db.Save(payment).Error
}

func (h *SaleController) UpdateSalePaymentMode(c *fiber.Ctx) error {
	// Define a struct for input validation
	type UpdateSalePaymentModeInput struct {
		ModeOfPayment string `json:"mode_of_payment" validate:"required"`
		TransactionID string `json:"transaction_id" validate:"required"`
		SalePaymentID uint   `json:"sale_payment_id" validate:"required"`
		UpdatedBy     string `json:"updated_by" validate:"required"`
	}

	// Parse the payment ID from the request parameters
	paymentID := c.Params("id")

	// Parse and validate the request body
	var input UpdateSalePaymentModeInput
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

	// Fetch the payment record using the repository
	payment, err := h.repo.FindSalePaymentModeById(paymentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "payment not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch payment",
			"error":   err.Error(),
		})
	}

	// Update the payment fields
	payment.ModeOfPayment = input.ModeOfPayment
	payment.TransactionID = input.TransactionID
	payment.SalePaymentID = input.SalePaymentID
	payment.UpdatedBy = input.UpdatedBy

	// Save the updated payment using the repository
	if err := h.repo.UpdateSalePaymentMode(payment); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update payment",
			"error":   err.Error(),
		})
	}

	// Return the updated payment
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "payment updated successfully",
		"data":    payment,
	})
}

// =============================

// Delete salePayment by ID
func (r *SaleRepositoryImpl) DeleteSalePaymentModeByID(id string) error {
	if err := r.db.Delete(&saleRegistration.SalePaymentMode{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

// DeleteCompanyByID deletes a company by its ID
func (h *SaleController) DeleteSalePaymentModeByID(c *fiber.Ctx) error {
	// Get the company ID from the route parameters
	id := c.Params("id")

	// Find the SalePayment in the database
	salePayment, err := h.repo.FindSalePaymentModeById(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{
				"status":  "error",
				"message": "salePayment mode not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to find salePayment mode",
			"data":    err.Error(),
		})
	}

	// Delete the salePayment
	if err := h.repo.DeleteSalePaymentModeByID(id); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to delete payment mode",
			"data":    err.Error(),
		})
	}

	// Return success response
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "salePayment mode deleted successfully",
		"data":    salePayment,
	})
}

// =========================================================================================================

// Get deposit by SalePaymentDeposit

// CreateCustomerContact creates a new payment deposit in the database
func (r *SaleRepositoryImpl) CreatePaymentDeposit(deposit *saleRegistration.SalePaymentDeposit) error {
	return r.db.Create(deposit).Error
}

// CreateSalePayment handles the creation of a payment deposit
func (h *SaleController) CreatePaymentDeposit(c *fiber.Ctx) error {
	// Parse the request body into a salePayment struct
	saleDeposit := new(saleRegistration.SalePaymentDeposit)
	if err := c.BodyParser(saleDeposit); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input data",
			"data":    err.Error(),
		})
	}

	// Create the customer address in the database
	if err := h.repo.CreatePaymentDeposit(saleDeposit); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create Payment deposit",
			"data":    err.Error(),
		})
	}

	// Return success response
	return c.Status(201).JSON(fiber.Map{
		"status":  "success",
		"message": "Payment deposit created successfully",
		"data":    saleDeposit,
	})
}

// ==============================

func (r *SaleRepositoryImpl) GetPaginatedPaymentDeposits(c *fiber.Ctx) (*utils.Pagination, []saleRegistration.SalePaymentDeposit, error) {
	pagination, deposits, err := utils.Paginate(c, r.db.Preload("SalePayment"), saleRegistration.SalePaymentDeposit{})
	if err != nil {
		return nil, nil, err
	}
	return &pagination, deposits, nil
}

func (h *SaleController) GetSalePaymentDeposits(c *fiber.Ctx) error {

	// Fetch paginated payments using the repository
	pagination, paymentDeposits, err := h.repo.GetPaginatedPaymentDeposits(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve payments deposits",
			"data":    err.Error(),
		})
	}

	// Return the paginated response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Payment deposits and associated data retrieved successfully",
		"data":    paymentDeposits,
		"pagination": fiber.Map{
			"total_items":  pagination.TotalItems,
			"total_pages":  pagination.TotalPages,
			"current_page": pagination.CurrentPage,
			"limit":        pagination.ItemsPerPage,
		},
	})
}

// =====================

func (r *SaleRepositoryImpl) FindSalePaymentDepositByIdAndSalePaymentId(id, salePaymentId string) (*saleRegistration.SalePaymentDeposit, error) {
	var paymentDeposit saleRegistration.SalePaymentDeposit
	result := r.db.Preload("SalePayment").Where("id = ? AND sale_payment_id = ?", id, salePaymentId).First(&paymentDeposit)
	if result.Error != nil {
		return nil, result.Error
	}
	return &paymentDeposit, nil
}

func (h *SaleController) FindSalePaymentDepositByIdAndSalePaymentId(c *fiber.Ctx) error {
	// Retrieve the payment ID and Company ID from the request parameters
	id := c.Params("id")
	salePaymentId := c.Params("salePaymentId")

	// Fetch the company deposit from the repository
	deposit, err := h.repo.FindSalePaymentDepositByIdAndSalePaymentId(id, salePaymentId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Deposit not found for the specified payment",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch the mode",
			"error":   err.Error(),
		})
	}

	// Return the fetched payment
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "deposit fetched successfully",
		"data":    deposit,
	})
}

// =============

func (r *SaleRepositoryImpl) GetPaymentDeposits(c *fiber.Ctx, name string) (*utils.Pagination, []saleRegistration.SalePaymentDeposit, error) {
	pagination, deposits, err := utils.Paginate(c, r.db.Preload("SalePayment").Where("bank_name = ?", name), saleRegistration.SalePaymentDeposit{})
	if err != nil {
		return nil, nil, err
	}
	return &pagination, deposits, nil
}

func (h *SaleController) GetPaymentDepositsByName(c *fiber.Ctx) error {
	name := c.Params("name")

	// Fetch paginated contacts using the repository
	pagination, deposits, err := h.repo.GetPaymentDeposits(c, name)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve payment modes",
			"data":    err.Error(),
		})
	}

	// Return the paginated response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Payment deposits and associated data retrieved successfully",
		"data":    deposits,
		"pagination": fiber.Map{
			"total_items":  pagination.TotalItems,
			"total_pages":  pagination.TotalPages,
			"current_page": pagination.CurrentPage,
			"limit":        pagination.ItemsPerPage,
		},
	})
}

// =====================

func (r *SaleRepositoryImpl) FindSalePaymentDepositById(id string) (*saleRegistration.SalePaymentDeposit, error) {
	var deposit saleRegistration.SalePaymentDeposit
	if err := r.db.First(&deposit, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &deposit, nil
}

func (r *SaleRepositoryImpl) UpdateSalePaymentDeposit(deposit *saleRegistration.SalePaymentDeposit) error {
	return r.db.Save(deposit).Error
}

func (h *SaleController) UpdateSalePaymentDeposit(c *fiber.Ctx) error {
	// Define a struct for input validation
	type UpdateSalePaymentDepositInput struct {
		BankName        string  `json:"bank_name" validate:"required"`
		BankAccount     string  `json:"bank_account" validate:"required"`
		BankBranch      string  `json:"bank_branch" validate:"required"`
		AmountDeposited float64 `json:"amount_deposited" validate:"required"`
		DateDeposited   string  `json:"date_deposited" validate:"required"`
		DepositScan     string  `json:"deposit_scan" validate:"required"`
		SalePaymentID   uint    `json:"sale_payment_id" validate:"required"`
		UpdatedBy       string  `json:"updated_by" validate:"required"`
	}

	// Parse the payment ID from the request parameters
	paymentID := c.Params("id")

	// Parse and validate the request body
	var input UpdateSalePaymentDepositInput
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

	// Fetch the payment record using the repository
	deposit, err := h.repo.FindSalePaymentDepositById(paymentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "payment not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch payment",
			"error":   err.Error(),
		})
	}

	// Update the deposit fields
	deposit.BankName = input.BankName
	deposit.BankAccount = input.BankAccount
	deposit.BankBranch = input.BankBranch
	deposit.AmountDeposited = input.AmountDeposited
	deposit.DateDeposited = input.DateDeposited
	deposit.DepositScan = input.DepositScan
	deposit.SalePaymentID = input.SalePaymentID
	deposit.UpdatedBy = input.UpdatedBy

	// Save the updated payment using the repository
	if err := h.repo.UpdateSalePaymentDeposit(deposit); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update deposit",
			"error":   err.Error(),
		})
	}

	// Return the updated payment
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Deposit updated successfully",
		"data":    deposit,
	})
}

// =============================

// Delete salePayment by ID
func (r *SaleRepositoryImpl) DeleteSalePaymentDepositByID(id string) error {
	if err := r.db.Delete(&saleRegistration.SalePaymentDeposit{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

// DeleteCompanyByID deletes a company by its ID
func (h *SaleController) DeleteSalePaymentDepositByID(c *fiber.Ctx) error {
	// Get the company ID from the route parameters
	id := c.Params("id")

	// Find the SalePaymentDeposit in the database
	deposit, err := h.repo.FindSalePaymentDepositById(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{
				"status":  "error",
				"message": "salePayment deposit not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to find salePayment deposit",
			"data":    err.Error(),
		})
	}

	// Delete the salePayment deposit
	if err := h.repo.DeleteSalePaymentDepositByID(id); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to delete payment mode",
			"data":    err.Error(),
		})
	}

	// Return success response
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "payment deposit deleted successfully",
		"data":    deposit,
	})
}

// ================================

func (r *SaleRepositoryImpl) GetOutstandingBalanceByCustomerID(customerID uint) (float64, error) {
	var sales []saleRegistration.Sale
	var totalPaid float64
	var totalOutstanding float64

	// Fetch all sales linked to cars owned by the customer
	if err := r.db.
		Joins("JOIN cars ON sales.car_id = cars.id").
		Where("cars.customer_id = ?", customerID).
		Find(&sales).Error; err != nil {
		return 0, err
	}

	// If no sales found, return zero balance
	if len(sales) == 0 {
		return 0, fmt.Errorf("no sales found for this customer")
	}

	// Iterate over each sale to compute outstanding balances
	for _, sale := range sales {
		totalPaid = 0

		// Sum all payments for the sale
		if err := r.db.Model(&saleRegistration.SalePayment{}).
			Where("sale_id = ?", sale.ID).
			Select("COALESCE(SUM(amount_payed), 0)").
			Scan(&totalPaid).Error; err != nil {
			return 0, err
		}

		// Calculate outstanding balance for this sale
		outstandingBalance := sale.TotalPrice - totalPaid
		if outstandingBalance > 0 {
			totalOutstanding += outstandingBalance
		}
	}

	return totalOutstanding, nil
}

// =========================

type PaymentRecord struct {
	PaymentDate time.Time `json:"payment_date"`
	Amount      float64   `json:"amount"`
}

type SaleStatement struct {
	CarID             uint            `json:"car_id"`
	CarModel          string          `json:"car_model"`
	ChasisNumber      string          `json:"chasis_number"`
	TotalSaleAmount   float64         `json:"total_sale_amount"`
	TotalPaid         float64         `json:"total_paid"`
	OutstandingAmount float64         `json:"outstanding_amount"`
	Payments          []PaymentRecord `json:"payments"`
}

type CustomerStatement struct {
	CustomerID       uint            `json:"customer_id"`
	CustomerName     string          `json:"customer_name"`
	TotalSales       float64         `json:"total_sales"`
	TotalPaid        float64         `json:"total_paid"`
	TotalOutstanding float64         `json:"total_outstanding"`
	Sales            []SaleStatement `json:"sales"`
}

func (r *SaleRepositoryImpl) GenerateCustomerStatement(customerID uint) (*CustomerStatement, error) {
	var customer customerRegistration.Customer
	var cars []carRegistration.Car
	var sales []saleRegistration.Sale
	var totalSales, totalPaid, totalOutstanding float64

	// Fetch Customer details
	if err := r.db.First(&customer, customerID).Error; err != nil {
		return nil, fmt.Errorf("customer not found")
	}

	// Fetch all cars owned by the customer
	if err := r.db.Where("customer_id = ?", customerID).Find(&cars).Error; err != nil {
		return nil, err
	}

	// Fetch all sales linked to the customer’s cars
	if err := r.db.Where("car_id IN (?)", getCarIDs(cars)).Find(&sales).Error; err != nil {
		return nil, err
	}

	var saleStatements []SaleStatement

	// Process each sale
	for _, sale := range sales {
		var payments []PaymentRecord
		var totalSalePaid float64

		// Fetch all payments for this sale
		rows, err := r.db.Raw(`
			SELECT amount_payed, payment_date 
			FROM sale_payments 
			WHERE sale_id = ? 
			ORDER BY payment_date ASC`, sale.ID).Rows()
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var payment PaymentRecord
			if err := rows.Scan(&payment.Amount, &payment.PaymentDate); err != nil {
				return nil, err
			}
			payments = append(payments, payment)
			totalSalePaid += payment.Amount
		}

		// Calculate outstanding balance
		outstanding := sale.TotalPrice - totalSalePaid

		// Create Sale Statement
		saleStatements = append(saleStatements, SaleStatement{
			CarID:             sale.CarID,
			CarModel:          getCarModelByID(cars, sale.CarID),
			ChasisNumber:      getChasisNumberByID(cars, sale.CarID),
			TotalSaleAmount:   sale.TotalPrice,
			TotalPaid:         totalSalePaid,
			OutstandingAmount: outstanding,
			Payments:          payments,
		})

		// Update totals
		totalSales += sale.TotalPrice
		totalPaid += totalSalePaid
		totalOutstanding += outstanding
	}

	// Create and return the full statement
	return &CustomerStatement{
		CustomerID:       customerID,
		CustomerName:     customer.Surname + " " + customer.Firstname + " " + customer.Othername,
		TotalSales:       totalSales,
		TotalPaid:        totalPaid,
		TotalOutstanding: totalOutstanding,
		Sales:            saleStatements,
	}, nil
}

// Helper function to extract car IDs from []Car
func getCarIDs(cars []carRegistration.Car) []uint {
	var ids []uint
	for _, car := range cars {
		ids = append(ids, car.ID)
	}
	return ids
}

// Helper function to get car model by ID
func getCarModelByID(cars []carRegistration.Car, carID uint) string {
	for _, car := range cars {
		if car.ID == carID {
			return car.CarModel
		}
	}
	return ""
}

// Helper function to get chasis number by ID
func getChasisNumberByID(cars []carRegistration.Car, carID uint) string {
	for _, car := range cars {
		if car.ID == carID {
			return car.ChasisNumber
		}
	}
	return ""
}

func (h *SaleController) GenerateCustomerStatement(c *fiber.Ctx) error {
	customerIdstr := c.Params("customerId")
	customerID := utils.StrToUint(customerIdstr)

	// Fetch the customer statement from the repository
	statement, err := h.repo.GenerateCustomerStatement(customerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Deposit not found for the specified payment",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch the statement",
			"error":   err.Error(),
		})
	}

	// Return the fetched payment
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "statement fetched successfully",
		"data":    statement,
	})
}
