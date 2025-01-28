package controllers

import (
	"car-bond/internals/models/metaData"
	"errors"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type MetaGetRepository interface {
	FindEvaluationsByDescription(description string) ([]metaData.VehicleEvaluation, error)
	GetAllWeightUnits(c *fiber.Ctx) ([]metaData.WeightUnit, error)
	GetAllLeightUnits(c *fiber.Ctx) ([]metaData.LeightUnit, error)
	GetAllCurrencies(c *fiber.Ctx) ([]metaData.Currency, error)
	GetAllExpenseCategories(c *fiber.Ctx) ([]metaData.ExpenseCategory, error)
	FindPortsByName(name string) ([]metaData.Port, error)
}

type MetaGetRepositoryImpl struct {
	db *gorm.DB
}

func NewMetaGetRepository(db *gorm.DB) MetaGetRepository {
	return &MetaGetRepositoryImpl{db: db}
}

type MetaGetController struct {
	repo MetaGetRepository
}

func NewMetaGetController(repo MetaGetRepository) *MetaGetController {
	return &MetaGetController{repo: repo}
}

// ======================================

func (m *MetaGetRepositoryImpl) FindEvaluationsByDescription(description string) ([]metaData.VehicleEvaluation, error) {
	var evaluations []metaData.VehicleEvaluation
	if err := m.db.Where("description LIKE ?", "%"+description+"%").Find(&evaluations).Error; err != nil {
		return nil, err
	}
	return evaluations, nil
}

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

func (m *MetaGetRepositoryImpl) GetAllWeightUnits(c *fiber.Ctx) ([]metaData.WeightUnit, error) {
	var weights []metaData.WeightUnit
	if err := m.db.Find(&weights).Error; err != nil {
		return nil, err
	}
	return weights, nil
}

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

func (m *MetaGetRepositoryImpl) GetAllLeightUnits(c *fiber.Ctx) ([]metaData.LeightUnit, error) {
	var lengths []metaData.LeightUnit
	if err := m.db.Find(&lengths).Error; err != nil {
		return nil, err
	}
	return lengths, nil
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

func (m *MetaGetRepositoryImpl) GetAllCurrencies(c *fiber.Ctx) ([]metaData.Currency, error) {
	var currencies []metaData.Currency
	if err := m.db.Find(&currencies).Error; err != nil {
		return nil, err
	}
	return currencies, nil
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

// ExpenseCategory

func (m *MetaGetRepositoryImpl) GetAllExpenseCategories(c *fiber.Ctx) ([]metaData.ExpenseCategory, error) {
	var expenses []metaData.ExpenseCategory
	if err := m.db.Find(&expenses).Error; err != nil {
		return nil, err
	}
	return expenses, nil
}

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

func (m *MetaGetRepositoryImpl) FindPortsByName(name string) ([]metaData.Port, error) {
	var ports []metaData.Port
	if err := m.db.Where("name LIKE ?", "%"+name+"%").Find(&ports).Error; err != nil {
		return nil, err
	}
	return ports, nil
}

func (h *MetaGetController) FindPortsByName(c *fiber.Ctx) error {
	// Retrieve companyId and expenseDate from the request parameters
	name := c.Query("name")
	if name == "" {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Name query parameter is required",
		})
	}

	// Fetch ports using the repository
	ports, err := h.repo.FindPortsByName(name)
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
