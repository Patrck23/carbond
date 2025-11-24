package repository

import (
	"car-bond/internals/models/carRegistration"
	"car-bond/internals/utils"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ShippingRepository interface {
	CreateShippingInvoice(invoice *carRegistration.CarShippingInvoice) error
	GetPaginatedShippingInvoices(c *fiber.Ctx) (*utils.Pagination, []carRegistration.CarShippingInvoice, error)
	// GetCarsByInvoiceId(invoiceID uint) ([]carRegistration.Car, error)
	GetCarsByInvoiceId(invoiceID uint) (map[string][]carRegistration.Car, error)
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

func (r *ShippingRepositoryImpl) CreateShippingInvoice(invoice *carRegistration.CarShippingInvoice) error {
	return r.db.Create(invoice).Error
}

// func (r *ShippingRepositoryImpl) GetCarsByInvoiceId(invoiceID uint) ([]carRegistration.Car, error) {
// 	var cars []carRegistration.Car
// 	err := r.db.Where("car_shipping_invoice_id = ?", invoiceID).Find(&cars).Error
// 	return cars, err
// }

func (r *ShippingRepositoryImpl) GetCarsByInvoiceId(invoiceID uint) (map[string][]carRegistration.Car, error) {
	var cars []carRegistration.Car
	err := r.db.Where("car_shipping_invoice_id = ?", invoiceID).Find(&cars).Error
	if err != nil {
		return nil, err
	}

	groupedCars := make(map[string][]carRegistration.Car)
	for _, car := range cars {
		key := car.OtherEntity
		groupedCars[key] = append(groupedCars[key], car)
	}

	return groupedCars, nil
}

func (r *ShippingRepositoryImpl) GetPaginatedShippingInvoices(c *fiber.Ctx) (*utils.Pagination, []carRegistration.CarShippingInvoice, error) {
	// — parse the exclude_locked param (defaults to “false” if absent) —
	excludeLocked := false
	if q := strings.TrimSpace(c.Query("exclude_locked")); q != "" {
		b, err := strconv.ParseBool(q)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid exclude_locked: %w", err)
		}
		excludeLocked = b
	}
	// — build the base query, including your Cars preload —
	query := r.db.
		Model(&carRegistration.CarShippingInvoice{})

	// — if exclude_locked=true, add WHERE locked = false —
	if excludeLocked {
		query = query.Where("locked = ?", false)
	}

	pagination, invoices, err := utils.Paginate(c, query, carRegistration.CarShippingInvoice{})
	if err != nil {
		return nil, nil, err
	}

	return &pagination, invoices, nil
}

func (r *ShippingRepositoryImpl) GetShippingInvoiceByID(invoiceID string) (carRegistration.CarShippingInvoice, error) {
	var invoice carRegistration.CarShippingInvoice
	err := r.db.First(&invoice, "id = ?", invoiceID).Error
	return invoice, err
}

func (r *ShippingRepositoryImpl) GetShippingInvoiceByInvoiceNum(invoiceNo string) (carRegistration.CarShippingInvoice, error) {
	var invoice carRegistration.CarShippingInvoice
	err := r.db.First(&invoice, "invoice_no = ?", invoiceNo).Error
	return invoice, err
}

func (r *ShippingRepositoryImpl) UpdateShippingInvoice(invoice *carRegistration.CarShippingInvoice) error {
	return r.db.Save(invoice).Error
}

// DeleteByID deletes a Invoice by ID
func (r *ShippingRepositoryImpl) DeleteShippingInvoiceByID(id string) error {
	if err := r.db.Delete(&carRegistration.CarShippingInvoice{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
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

func (r *ShippingRepositoryImpl) UnlockInvoice(id uint, updatedBy string) error {
	return r.db.Model(&carRegistration.CarShippingInvoice{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"locked":     false,
			"updated_by": updatedBy,
		}).Error
}
