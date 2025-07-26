package controllers

import (
	"car-bond/internals/repository"
	"errors"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type MetaGetController struct {
	repo repository.MetaGetRepository
}

func NewMetaGetController(repo repository.MetaGetRepository) *MetaGetController {
	return &MetaGetController{repo: repo}
}

// ======================================

func (h *MetaGetController) FetchVehicleEvaluationsByDescription(c *fiber.Ctx) error {
	// Retrieve companyId and expenseDate from the request parameters
	description := c.Query("description")
	if description == "" {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Description query parameter is required",
		})
	}

	// Fetch evaluations using the repository
	evaluations, err := h.repo.FindEvaluationsByDescription(description)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "evaluations not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch records",
			"error":   err.Error(),
		})
	}

	// Return the fetched evaluations
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Records fetched successfully",
		"data":    evaluations,
	})
}

// =======================================================================

// Meta Units

func (h *MetaGetController) GetAllWeightUnits(c *fiber.Ctx) error {
	// Fetch evaluations using the repository
	weights, err := h.repo.GetAllWeightUnits(c)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Weights not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch records",
			"error":   err.Error(),
		})
	}

	// Return the fetched evaluations
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Records fetched successfully",
		"data":    weights,
	})
}

func (h *MetaGetController) GetAllLeightUnits(c *fiber.Ctx) error {
	lengths, err := h.repo.GetAllLeightUnits(c)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Length not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch records",
			"error":   err.Error(),
		})
	}

	// Return the fetched evaluations
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Records fetched successfully",
		"data":    lengths,
	})
}
func (h *MetaGetController) GetAllCurrencies(c *fiber.Ctx) error {
	currencies, err := h.repo.GetAllCurrencies(c)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "currencies not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch records",
			"error":   err.Error(),
		})
	}

	// Return the fetched evaluations
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Records fetched successfully",
		"data":    currencies,
	})
}

// =================

func (h *MetaGetController) GetAllExpenseCategories(c *fiber.Ctx) error {
	currencies, err := h.repo.GetAllExpenseCategories(c)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Expense categories not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch records",
			"error":   err.Error(),
		})
	}

	// Return the fetched evaluations
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Records fetched successfully",
		"data":    currencies,
	})
}

// ======================================

func (h *MetaGetController) FindPorts(c *fiber.Ctx) error {
	// Retrieve companyId and expenseDate from the request parameter

	// Fetch ports using the repository
	ports, err := h.repo.FindPorts(c)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Ports not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch records",
			"error":   err.Error(),
		})
	}

	// Return the fetched evaluations
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Records fetched successfully",
		"data":    ports,
	})
}

// ======================================

func (h *MetaGetController) FindPaymentModeBymode(c *fiber.Ctx) error {
	// Retrieve companyId and expenseDate from the request parameters
	mode := c.Query("mode")
	if mode == "" {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Mode query parameter is required",
		})
	}

	// Fetch ports using the repository
	ports, err := h.repo.FindPaymentModeBymode(mode)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Payment modes not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch records",
			"error":   err.Error(),
		})
	}

	// Return the fetched evaluations
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Records fetched successfully",
		"data":    ports,
	})
}
