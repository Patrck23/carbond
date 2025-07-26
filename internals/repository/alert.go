package repository

import (
	"car-bond/internals/models/alertRegistration"
	"car-bond/internals/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type AlertRepository interface {
	CreateAlert(alert *alertRegistration.Transaction) error
	SearchPaginatedAlerts(c *fiber.Ctx) (*utils.Pagination, []alertRegistration.Transaction, error)
	GetAlertByID(id string) (alertRegistration.Transaction, error)
	UpdateAlert(alert *alertRegistration.Transaction) error
}

type AlertRepositoryImpl struct {
	db *gorm.DB
}

func NewAlertRepository(db *gorm.DB) AlertRepository {
	return &AlertRepositoryImpl{db: db}
}

func (r *AlertRepositoryImpl) CreateAlert(alert *alertRegistration.Transaction) error {
	return r.db.Create(alert).Error
}

func (r *AlertRepositoryImpl) SearchPaginatedAlerts(c *fiber.Ctx) (*utils.Pagination, []alertRegistration.Transaction, error) {
	// Get query parameters from request
	chasisNumber := c.Query("car_chasis_number")
	transactionType := c.Query("transaction_type")
	viewStatus := c.Query("view_status")
	toCompanyId := c.Query("to_company_id")
	fromCompanyId := c.Query("from_company_id")

	// Start building the query
	query := r.db.Model(&alertRegistration.Transaction{})

	// Apply filters based on provided parameters
	if fromCompanyId != "" {
		if _, err := strconv.Atoi(fromCompanyId); err == nil {
			query = query.Where("from_company_id = ?", fromCompanyId)
		}
	}

	if toCompanyId != "" {
		if _, err := strconv.Atoi(toCompanyId); err == nil {
			query = query.Where("to_company_id = ?", toCompanyId)
		}
	}

	if chasisNumber != "" {
		query = query.Where("chasis_number LIKE ?", "%"+chasisNumber+"%")
	}
	if transactionType != "" {
		query = query.Where("transaction_type LIKE ?", "%"+transactionType+"%")
	}
	if viewStatus != "" {
		query = query.Where("view_status = ?", viewStatus)
	}
	// Call the pagination helper
	pagination, alerts, err := utils.Paginate(c, query, alertRegistration.Transaction{})
	if err != nil {
		return nil, nil, err
	}

	return &pagination, alerts, nil
}

func (r *AlertRepositoryImpl) GetAlertByID(id string) (alertRegistration.Transaction, error) {
	var alert alertRegistration.Transaction
	err := r.db.First(&alert, "id = ?", id).Error
	return alert, err
}

func (r *AlertRepositoryImpl) UpdateAlert(alert *alertRegistration.Transaction) error {
	return r.db.Save(alert).Error
}
