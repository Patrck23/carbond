package controllers

import (
	"car-bond/internals/models/alertRegistration"
	"car-bond/internals/repository"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type AlertController struct {
	repo     repository.AlertRepository
	saleRepo repository.SaleRepository
}

func NewAlertController(repo repository.AlertRepository, saleRepo repository.SaleRepository) *AlertController {
	return &AlertController{
		repo:     repo,
		saleRepo: saleRepo,
	}
}

// ============================================

func (h *AlertController) CreateAlert(c *fiber.Ctx) error {
	// Initialize a new Alert instance
	alert := new(alertRegistration.Transaction)

	// Parse the request body into the alert instance
	if err := c.BodyParser(alert); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input provided",
			"data":    err.Error(),
		})
	}

	// Attempt to create the alert record using the repository
	if err := h.repo.CreateAlert(alert); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create alert",
			"data":    err.Error(),
		})
	}

	// Return the newly created alert record
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Alert created successfully",
		"data":    alert,
	})
}

// ========================

func (h *AlertController) SearchAlerts(c *fiber.Ctx) error {
	// Call the repository function to get paginated search results
	pagination, alerts, err := h.repo.SearchPaginatedAlerts(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve alerts",
			"data":    err.Error(),
		})
	}

	companyNotifications, err := h.saleRepo.CheckPaymentNotifications(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve sales notifications",
			"data":    err.Error(),
		})
	}

	// Return the response with pagination details
	return c.Status(200).JSON(fiber.Map{
		"status":              "success",
		"message":             "Alerts retrieved successfully",
		"pagination":          pagination,
		"data":                alerts,
		"sales_notifications": companyNotifications,
	})
}

// ======================

// Define the UpdateAlert struct
type UpdateAlertPayload struct {
	ViewStatus bool   `json:"view_status"`
	UpdatedBy  string `json:"updated_by"`
}

// UpdateAlert handler function
func (h *AlertController) UpdateAlert(c *fiber.Ctx) error {
	// Get the Alert ID from the route parameters
	id := c.Params("id")

	// Find the alert in the database
	alert, err := h.repo.GetAlertByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{
				"status":  "error",
				"message": "Alert not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve Alert",
			"data":    err.Error(),
		})
	}

	// Parse the request body into the UpdateAlertPayload struct
	var payload UpdateAlertPayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input",
			"data":    err.Error(),
		})
	}

	// Update the Alert fields using the payload
	updateAlertFields(&alert, payload) // Pass the parsed payload

	// Save the changes to the database
	if err := h.repo.UpdateAlert(&alert); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update alert",
			"data":    err.Error(),
		})
	}

	// Return the updated alert
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "alert updated successfully",
		"data":    alert,
	})
}

// UpdateAlertFields updates the fields of a alert using the UpdateAlert struct
func updateAlertFields(alert *alertRegistration.Transaction, updateAlertData UpdateAlertPayload) {
	alert.ViewStatus = updateAlertData.ViewStatus
	alert.UpdatedBy = updateAlertData.UpdatedBy
}
