package controllers

import (
	"car-bond/internals/models/alertRegistration"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type AlertRepository interface {
	CreateAlert(alert *alertRegistration.Transaction) error
}

type AlertRepositoryImpl struct {
	db *gorm.DB
}

func NewAlertRepository(db *gorm.DB) AlertRepository {
	return &AlertRepositoryImpl{db: db}
}

type AlertController struct {
	repo AlertRepository
}

func NewAlertController(repo AlertRepository) *AlertController {
	return &AlertController{repo: repo}
}

// ============================================

func (r *AlertRepositoryImpl) CreateAlert(alert *alertRegistration.Transaction) error {
	return r.db.Create(alert).Error
}

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
