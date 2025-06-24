package controllers

import (
	"car-bond/internals/models/carRegistration"
	"car-bond/internals/utils"
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ShippingRepository interface {
	CreateShippingInvoice(invoice *carRegistration.CarShippingInvoice) error
	GetPaginatedShippingInvoices(c *fiber.Ctx) (*utils.Pagination, []carRegistration.CarShippingInvoice, error)
	GetCarsByInvoiceId(invoiceID uint) ([]carRegistration.Car, error)
	GetShippingInvoiceByID(invoiceID string) (carRegistration.CarShippingInvoice, error)
	GetShippingInvoiceByInvoiceNum(invoiceNo string) (carRegistration.CarShippingInvoice, error)
	UpdateShippingInvoice(invoice *carRegistration.CarShippingInvoice) error
	DeleteShippingInvoiceByID(id string) error

	UnlockInvoice(id uint, updatedBy string) error
	LockInvoice(id uint, updatedBy string) error
}

type ShippingRepositoryImpl struct {
	db *gorm.DB
}

func NewShippingRepository(db *gorm.DB) ShippingRepository {
	return &ShippingRepositoryImpl{db: db}
}

type ShippingController struct {
	repo ShippingRepository
}

func NewShippingController(repo ShippingRepository) *ShippingController {
	return &ShippingController{repo: repo}
}

// ============================================

func (r *ShippingRepositoryImpl) CreateShippingInvoice(invoice *carRegistration.CarShippingInvoice) error {
	return r.db.Create(invoice).Error
}

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

func (r *ShippingRepositoryImpl) GetCarsByInvoiceId(invoiceID uint) ([]carRegistration.Car, error) {
	var cars []carRegistration.Car
	err := r.db.Where("car_shipping_invoice_id = ?", invoiceID).Find(&cars).Error
	return cars, err
}

func (r *ShippingRepositoryImpl) GetPaginatedShippingInvoices(c *fiber.Ctx) (*utils.Pagination, []carRegistration.CarShippingInvoice, error) {
	pagination, invoices, err := utils.Paginate(c, r.db, carRegistration.CarShippingInvoice{})
	if err != nil {
		return nil, nil, err
	}

	return &pagination, invoices, nil
}

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

func (r *ShippingRepositoryImpl) GetShippingInvoiceByID(invoiceID string) (carRegistration.CarShippingInvoice, error) {
	var invoice carRegistration.CarShippingInvoice
	err := r.db.First(&invoice, "id = ?", invoiceID).Error
	return invoice, err
}

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

func (r *ShippingRepositoryImpl) GetShippingInvoiceByInvoiceNum(invoiceNo string) (carRegistration.CarShippingInvoice, error) {
	var invoice carRegistration.CarShippingInvoice
	err := r.db.First(&invoice, "invoice_no = ?", invoiceNo).Error
	return invoice, err
}

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

func (r *ShippingRepositoryImpl) UpdateShippingInvoice(invoice *carRegistration.CarShippingInvoice) error {
	return r.db.Save(invoice).Error
}

// Define the UpdateShippingInvoice struct
type UpdateShippingInvoicePayload struct {
	InvoiceNo    string `json:"invoice_no"`
	ShipDate     string `json:"ship_date"`
	VesselName   string `json:"vessel_name"`
	FromLocation string `json:"from_location"`
	ToLocation   string `json:"to_location"`
	UpdatedBy    string `json:"updated_by"`
}

// UpdateShippingInvoice handler function
func (h *ShippingController) UpdateShippingInvoice(c *fiber.Ctx) error {
	// Get the invoice ID from the route parameters
	id := c.Params("id")

	// Find the invoice in the database
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
			"message": "Failed to retrieve invoice",
			"data":    err.Error(),
		})
	}

	// Parse the request body into the UpdateShippingInvoicePayload struct
	var payload UpdateShippingInvoicePayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input",
			"data":    err.Error(),
		})
	}

	// Update the invoice fields using the payload
	updateShippingInvoiceFields(&invoice, payload) // Pass the parsed payload

	// Save the changes to the database
	if err := h.repo.UpdateShippingInvoice(&invoice); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update invoice",
			"data":    err.Error(),
		})
	}

	// Return the updated invoice
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "invoice updated successfully",
		"data":    invoice,
	})
}

// UpdateinvoiceFields updates the fields of a invoice using the Updateinvoice struct
func updateShippingInvoiceFields(invoice *carRegistration.CarShippingInvoice, updateShippingInvoiceData UpdateShippingInvoicePayload) {
	invoice.InvoiceNo = updateShippingInvoiceData.InvoiceNo
	invoice.ShipDate = updateShippingInvoiceData.ShipDate
	invoice.VesselName = updateShippingInvoiceData.VesselName
	invoice.FromLocation = updateShippingInvoiceData.FromLocation
	invoice.ToLocation = updateShippingInvoiceData.ToLocation
	invoice.UpdatedBy = updateShippingInvoiceData.UpdatedBy
}

// ======================

// DeleteByID deletes a Invoice by ID
func (r *ShippingRepositoryImpl) DeleteShippingInvoiceByID(id string) error {
	if err := r.db.Delete(&carRegistration.CarShippingInvoice{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

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

var ErrAlreadyLocked = errors.New("invoice already locked")

func (r *ShippingRepositoryImpl) LockInvoice(id uint, updatedBy string) error {
	tx := r.db.Model(&carRegistration.CarShippingInvoice{}).
		Where("id = ? AND locked = ?", id, false).
		Updates(map[string]interface{}{
			"locked":     true,
			"updated_by": updatedBy,
		})

	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		// Either not found or already locked; check which.
		var tmp carRegistration.CarShippingInvoice
		if err := r.db.Select("id", "locked").First(&tmp, id).Error; err != nil {
			return err // not found
		}
		return ErrAlreadyLocked
	}
	return nil
}

// ---------------------------------------------------------------------
// 3) Unlock invoice  (optional convenience helper)
// ---------------------------------------------------------------------
func (r *ShippingRepositoryImpl) UnlockInvoice(id uint, updatedBy string) error {
	return r.db.Model(&carRegistration.CarShippingInvoice{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"locked":     false,
			"updated_by": updatedBy,
		}).Error
}

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
		case ErrAlreadyLocked:
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
