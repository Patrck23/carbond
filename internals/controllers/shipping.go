package controllers

import (
	"car-bond/internals/models/carRegistration"
	"car-bond/internals/repository"
	"car-bond/internals/utils"
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ShippingController struct {
	repo repository.ShippingRepository
	DB   *gorm.DB
}

func NewShippingController(repo repository.ShippingRepository, db *gorm.DB) *ShippingController {
	return &ShippingController{repo: repo,
		DB: db}
}

// ============================================

func (h *ShippingController) CreateShippingInvoice(c *fiber.Ctx) error {
	// Initialize a new invoice instance
	invoice := new(carRegistration.CarShippingInvoice)

	// Parse the request body into the invoice instance
	if err := c.BodyParser(invoice); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input provided",
			"data":    err.Error(),
		})
	}

	// Attempt to create the invoice record using the repository
	if err := h.repo.CreateShippingInvoice(invoice); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create shipping invoice",
			"data":    err.Error(),
		})
	}

	// Return the newly created invoice record
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Shipping invoice created successfully",
		"data":    invoice,
	})
}

// ===================

func (h *ShippingController) GetAllShippingInvoices(c *fiber.Ctx) error {
	// Fetch paginated invoices using the repository
	pagination, invoices, err := h.repo.GetPaginatedShippingInvoices(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve invoices",
			"data":    err.Error(),
		})
	}

	var response []fiber.Map

	// Iterate over all invoices to fetch associated invoice cars
	for _, invoice := range invoices {
		// Fetch invoice cars
		cars, err := h.repo.GetCarsByInvoiceId(invoice.ID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to retrieve cars for invoice ID " + strconv.Itoa(int(invoice.ID)),
				"data":    err.Error(),
			})
		}

		response = append(response, fiber.Map{
			"invoice": invoice,
			"cars":    cars,
		})
	}

	// Return the paginated response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "invoices retrieved successfully",
		"data":    response,
		"pagination": fiber.Map{
			"total_items":  pagination.TotalItems,
			"total_pages":  pagination.TotalPages,
			"current_page": pagination.CurrentPage,
			"limit":        pagination.ItemsPerPage,
		},
	})
}

// ================

// GetSingleInvoice fetches a invoice with its associated locations and expenses from the database
func (h *ShippingController) GetSingleInvoice(c *fiber.Ctx) error {
	// Get the invoice ID from the route parameters
	id := c.Params("id")

	// Fetch the invoice by ID
	invoice, err := h.repo.GetShippingInvoiceByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Invoice not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve invoice",
			"data":    err.Error(),
		})
	}

	// Fetch invoice locations associated with the invoice
	invoiceCars, err := h.repo.GetCarsByInvoiceId(invoice.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve cars",
			"data":    err.Error(),
		})
	}

	// Prepare the response
	response := fiber.Map{
		"invoice": invoice,
		"cars":    invoiceCars,
	}

	// Return the response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Invoice and associated data retrieved successfully",
		"data":    response,
	})
}

// ================

// GetSingleInvoice fetches a invoice with its associated locations and expenses from the database
func (h *ShippingController) GetShippingInvoiceByInvoiceNum(c *fiber.Ctx) error {
	// Get the invoice ID from the route parameters
	no := c.Params("no")

	// Fetch the invoice by ID
	invoice, err := h.repo.GetShippingInvoiceByInvoiceNum(no)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Invoice not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve invoice",
			"data":    err.Error(),
		})
	}

	// Fetch invoice cars associated with the invoice
	invoiceCars, err := h.repo.GetCarsByInvoiceId(invoice.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve cars",
			"data":    err.Error(),
		})
	}

	// Prepare the response
	response := fiber.Map{
		"invoice": invoice,
		"cars":    invoiceCars,
	}

	// Return the response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Invoice and associated data retrieved successfully",
		"data":    response,
	})
}

// // ====================

// 1️⃣ Extend your payload
type UpdateShippingInvoicePayload struct {
	InvoiceNo    string `json:"invoice_no"`
	ShipDate     string `json:"ship_date"` // "YYYY-MM-DD"
	VesselName   string `json:"vessel_name"`
	FromLocation string `json:"from_location"`
	ToLocation   string `json:"to_location"`
	UpdatedBy    string `json:"updated_by"`
	CarIDs       []uint `json:"car_ids,omitempty"` // ← new: full list of IDs you want linked
}

func (h *ShippingController) UpdateShippingInvoice(c *fiber.Ctx) error {
	// … load existing invoice as before …
	invoice, err := h.repo.GetShippingInvoiceByID(c.Params("id"))
	if err != nil {
		// handle 404 or 500
	}

	// 2a️⃣ parse payload
	var payload UpdateShippingInvoicePayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input",
			"data":    err.Error(),
		})
	}

	// 2b️⃣ update scalar fields
	invoice.InvoiceNo = payload.InvoiceNo
	invoice.ShipDate = payload.ShipDate
	invoice.VesselName = payload.VesselName
	invoice.FromLocation = payload.FromLocation
	invoice.ToLocation = payload.ToLocation
	invoice.UpdatedBy = payload.UpdatedBy

	// 2c️⃣ persist the scalar fields first
	if err := h.repo.UpdateShippingInvoice(&invoice); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update invoice",
			"data":    err.Error(),
		})
	}

	// 2d️⃣ if CarIDs was provided, sync the many→many
	if payload.CarIDs != nil {
		// Step 1: Fetch all requested cars
		var carsToAssign []carRegistration.Car
		if err := h.DB.Where("id IN ?", payload.CarIDs).Find(&carsToAssign).Error; err != nil {
			return c.Status(400).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid car_ids provided",
				"data":    err.Error(),
			})
		}

		if err := h.DB.Model(&carRegistration.Car{}).
			Where("id IN ? AND car_shipping_invoice_id != ?", payload.CarIDs, invoice.ID).
			Update("car_shipping_invoice_id", nil).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to unassign cars from other invoices",
				"data":    err.Error(),
			})
		}

		// Step 3: Assign the cars to this invoice
		if err := h.DB.Model(&invoice).Association("Cars").Replace(carsToAssign); err != nil {
			return c.Status(500).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to reassign cars to invoice",
				"data":    err.Error(),
			})
		}
		invoice.Cars = carsToAssign
	}

	// 2e️⃣ return the updated invoice (with Cars)
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Invoice updated successfully",
		"data":    invoice,
	})
}

// ======================

// DeleteShippingInvoiceByID deletes a Invoice by its ID
func (h *ShippingController) DeleteShippingInvoiceByID(c *fiber.Ctx) error {
	// Get the Invoice ID from the route parameters
	id := c.Params("id")

	// Find the Invoice in the database
	invoice, err := h.repo.GetShippingInvoiceByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{
				"status":  "error",
				"message": "invoice not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to find invoice",
			"data":    err.Error(),
		})
	}

	// Delete the invoice
	if err := h.repo.DeleteShippingInvoiceByID(id); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to delete invoice",
			"data":    err.Error(),
		})
	}

	// Return success response
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "invoice deleted successfully",
		"data":    invoice,
	})
}

// ---------------------------------------------------------------------
// 3) Unlock invoice  (optional convenience helper)
// ---------------------------------------------------------------------

func getUsernameOrDefault(c *fiber.Ctx, def string) string {
	if v := c.Locals("username"); v != nil {
		if s, ok := v.(string); ok && s != "" {
			return s
		}
	}
	return def
}

func (ic *ShippingController) LockInvoice(c *fiber.Ctx) error {
	id := utils.StrToUint(c.Params("id"))
	updatedBy := getUsernameOrDefault(c, "system")

	if err := ic.repo.LockInvoice(id, updatedBy); err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Invoice not found"})
		case repository.ErrAlreadyLocked:
			return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Invoice is already locked"})
		default:
			return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed to lock invoice", "data": err.Error()})
		}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Invoice locked successfully",
	})
}
