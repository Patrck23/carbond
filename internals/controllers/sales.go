package controllers

import (
	"car-bond/internals/models/carRegistration"
	"car-bond/internals/models/saleRegistration"
	"car-bond/internals/repository"
	"car-bond/internals/utils"
	"errors"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type SaleController struct {
	repo repository.SaleRepository
	db   *gorm.DB
}

func NewSaleController(repo repository.SaleRepository, db *gorm.DB) *SaleController {
	return &SaleController{repo: repo,
		db: db}
}

// ============================================

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

// ==========================

type SaleWithPaymentsInput struct {
	Sale         saleRegistration.Sale      `json:"sale"`
	SalePayments []SalePaymentWithModeInput `json:"sale_payments"`
}

type SalePaymentWithModeInput struct {
	AmountPayed float64                          `json:"amount_payed"`
	PaymentDate string                           `json:"payment_date"`
	CreatedBy   string                           `json:"created_by"`
	UpdatedBy   string                           `json:"updated_by"`
	PaymentMode saleRegistration.SalePaymentMode `json:"payment_mode"`
}

func (h *SaleController) CreateSaleWithPayments(c *fiber.Ctx) error {
	var input SaleWithPaymentsInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input",
			"data":    err.Error(),
		})
	}

	// Validate sale date
	saleDate, err := time.Parse("2006-01-02", input.Sale.SaleDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid sale_date format (expected YYYY-MM-DD)",
			"data":    err.Error(),
		})
	}
	input.Sale.SaleDate = saleDate.Format("2006-01-02")

	tx := h.db.Begin()
	if tx.Error != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed to start transaction"})
	}

	// Save Sale
	if err := tx.Create(&input.Sale).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed to save sale", "data": err.Error()})
	}

	if input.Sale.CarID != 0 && input.Sale.CustomerID != nil && *input.Sale.CustomerID != 0 {
		if err := tx.Model(&carRegistration.Car{}).
			Where("id = ?", input.Sale.CarID).
			Update("customer_id", input.Sale.CustomerID).Error; err != nil {
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to update car's customer_id after sale",
				"data":    err.Error(),
			})
		}
	}

	var savedPayments []saleRegistration.SalePayment
	var savedModes []saleRegistration.SalePaymentMode

	// Loop through sale payments
	for _, p := range input.SalePayments {
		// Parse payment date
		paymentDate, err := time.Parse("2006-01-02", p.PaymentDate)
		if err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid payment_date format (expected YYYY-MM-DD)",
				"data":    err.Error(),
			})
		}

		// Save payment
		payment := saleRegistration.SalePayment{
			AmountPayed: p.AmountPayed,
			PaymentDate: paymentDate.Format("2006-01-02"),
			SaleID:      input.Sale.ID,
			CreatedBy:   p.CreatedBy,
			UpdatedBy:   p.UpdatedBy,
		}
		if err := tx.Create(&payment).Error; err != nil {
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed to save payment", "data": err.Error()})
		}

		// Save mode
		paymentMode := p.PaymentMode
		paymentMode.SalePaymentID = payment.ID
		if err := tx.Create(&paymentMode).Error; err != nil {
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed to save payment mode", "data": err.Error()})
		}

		savedPayments = append(savedPayments, payment)
		savedModes = append(savedModes, paymentMode)
	}

	if err := tx.Commit().Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Transaction commit failed", "data": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":        "success",
		"message":       "Sale and payments created successfully",
		"sale":          input.Sale,
		"sale_payments": savedPayments,
		"payment_modes": savedModes,
	})
}

func (h *SaleController) UpdateSaleWithPayments(c *fiber.Ctx) error {
	saleIDStr := c.Params("id")
	if saleIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Missing sale ID in URL",
		})
	}

	saleID := utils.StrToUint(saleIDStr)
	if saleID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid sale ID in URL",
		})
	}

	var input SaleWithPaymentsInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input",
			"data":    err.Error(),
		})
	}

	// ✅ Assign sale ID to struct
	input.Sale.ID = saleID

	// Validate and format sale date
	saleDate, err := time.Parse("2006-01-02", input.Sale.SaleDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid sale_date format (expected YYYY-MM-DD)",
			"data":    err.Error(),
		})
	}
	input.Sale.SaleDate = saleDate.Format("2006-01-02")

	tx := h.db.Begin()
	if tx.Error != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed to start transaction"})
	}

	// Update Sale
	if err := tx.Model(&saleRegistration.Sale{}).
		Where("id = ?", saleID).
		Updates(input.Sale).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed to update sale", "data": err.Error()})
	}

	if input.Sale.CarID != 0 && input.Sale.CustomerID != nil && *input.Sale.CustomerID != 0 {
		if err := tx.Model(&carRegistration.Car{}).
			Where("id = ?", input.Sale.CarID).
			Update("customer_id", input.Sale.CustomerID).Error; err != nil {
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to update car's customer_id after sale",
				"data":    err.Error(),
			})
		}
	}

	// Delete existing SalePayments (modes must be deleted first)
	if err := tx.Where("sale_payment_id IN (SELECT id FROM sale_payments WHERE sale_id = ?)", saleID).
		Delete(&saleRegistration.SalePaymentMode{}).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed to delete payment modes", "data": err.Error()})
	}
	if err := tx.Where("sale_id = ?", saleID).
		Delete(&saleRegistration.SalePayment{}).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed to delete payments", "data": err.Error()})
	}

	var savedPayments []saleRegistration.SalePayment
	var savedModes []saleRegistration.SalePaymentMode

	for _, p := range input.SalePayments {
		paymentDate, err := time.Parse("2006-01-02", p.PaymentDate)
		if err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid payment_date format (expected YYYY-MM-DD)",
				"data":    err.Error(),
			})
		}

		payment := saleRegistration.SalePayment{
			AmountPayed: p.AmountPayed,
			PaymentDate: paymentDate.Format("2006-01-02"),
			SaleID:      saleID, // ✅ use correct ID
			CreatedBy:   p.CreatedBy,
			UpdatedBy:   p.UpdatedBy,
		}

		if err := tx.Create(&payment).Error; err != nil {
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed to save payment", "data": err.Error()})
		}

		paymentMode := p.PaymentMode
		paymentMode.SalePaymentID = payment.ID
		if err := tx.Create(&paymentMode).Error; err != nil {
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed to save payment mode", "data": err.Error()})
		}

		savedPayments = append(savedPayments, payment)
		savedModes = append(savedModes, paymentMode)
	}

	if err := tx.Commit().Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Transaction commit failed", "data": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":        "success",
		"message":       "Sale and payments updated successfully",
		"sale":          input.Sale,
		"sale_payments": savedPayments,
		"payment_modes": savedModes,
	})
}
