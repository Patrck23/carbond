package repository

import (
	"car-bond/internals/models/metaData"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type MetaGetRepository interface {
	FindEvaluationsByDescription(description string) ([]metaData.VehicleEvaluation, error)
	GetAllWeightUnits(c *fiber.Ctx) ([]metaData.WeightUnit, error)
	GetAllLeightUnits(c *fiber.Ctx) ([]metaData.LeightUnit, error)
	GetAllCurrencies(c *fiber.Ctx) ([]metaData.Currency, error)
	GetAllExpenseCategories(c *fiber.Ctx) ([]metaData.ExpenseCategory, error)
	// FindPortsByName(name string) ([]metaData.Port, error)
	FindPaymentModeBymode(mode string) ([]metaData.PaymentMode, error)
	FindPorts(c *fiber.Ctx) ([]metaData.Port, error)
}

type MetaGetRepositoryImpl struct {
	db *gorm.DB
}

func NewMetaGetRepository(db *gorm.DB) MetaGetRepository {
	return &MetaGetRepositoryImpl{db: db}
}

func (m *MetaGetRepositoryImpl) FindEvaluationsByDescription(description string) ([]metaData.VehicleEvaluation, error) {
	var evaluations []metaData.VehicleEvaluation
	if err := m.db.Where("description LIKE ?", "%"+description+"%").Find(&evaluations).Error; err != nil {
		return nil, err
	}
	return evaluations, nil
}

func (m *MetaGetRepositoryImpl) GetAllWeightUnits(c *fiber.Ctx) ([]metaData.WeightUnit, error) {
	var weights []metaData.WeightUnit
	if err := m.db.Find(&weights).Error; err != nil {
		return nil, err
	}
	return weights, nil
}

func (m *MetaGetRepositoryImpl) GetAllLeightUnits(c *fiber.Ctx) ([]metaData.LeightUnit, error) {
	var lengths []metaData.LeightUnit
	if err := m.db.Find(&lengths).Error; err != nil {
		return nil, err
	}
	return lengths, nil
}

func (m *MetaGetRepositoryImpl) GetAllCurrencies(c *fiber.Ctx) ([]metaData.Currency, error) {
	var currencies []metaData.Currency
	if err := m.db.Find(&currencies).Error; err != nil {
		return nil, err
	}
	return currencies, nil
}

// ExpenseCategory

func (m *MetaGetRepositoryImpl) GetAllExpenseCategories(c *fiber.Ctx) ([]metaData.ExpenseCategory, error) {
	var expenses []metaData.ExpenseCategory
	if err := m.db.Find(&expenses).Error; err != nil {
		return nil, err
	}
	return expenses, nil
}

func (m *MetaGetRepositoryImpl) FindPorts(c *fiber.Ctx) ([]metaData.Port, error) {
	var ports []metaData.Port
	if err := m.db.Find(&ports).Error; err != nil {
		return nil, err
	}
	return ports, nil
}

func (m *MetaGetRepositoryImpl) FindPaymentModeBymode(mode string) ([]metaData.PaymentMode, error) {
	var modes []metaData.PaymentMode
	if err := m.db.Where("mode LIKE ?", "%"+mode+"%").Find(&modes).Error; err != nil {
		return nil, err
	}
	return modes, nil
}
