package repository

import (
	"car-bond/internals/models/carRegistration"
	"car-bond/internals/models/customerRegistration"
	"car-bond/internals/models/saleRegistration"
	"car-bond/internals/utils"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type SaleRepository interface {
	CreateSale(sale *saleRegistration.Sale) error
	GetPaginatedSales(c *fiber.Ctx) (*utils.Pagination, []saleRegistration.Sale, error)
	GetSalePayments(saleID uint) ([]saleRegistration.SalePayment, error)
	GetSalePaymentModes(paymentID uint) ([]saleRegistration.SalePaymentMode, error)
	GetSaleByID(id string) (saleRegistration.Sale, error)
	UpdateSale(sale *saleRegistration.Sale) error
	DeleteByID(id string) error

	// Payment
	CreateInvoice(payment *saleRegistration.SalePayment) error
	GetPaginatedInvoices(c *fiber.Ctx) (*utils.Pagination, []saleRegistration.SalePayment, error)
	FindSalePaymentById(id string) (*saleRegistration.SalePayment, error)
	FindSalePaymentByIdAndSaleId(id, saleId string) (*saleRegistration.SalePayment, error)
	UpdateSalePayment(payment *saleRegistration.SalePayment) error
	DeleteSalePaymentByID(id string) error

	// Payment Mode
	CreatePaymentMode(payment *saleRegistration.SalePaymentMode) error
	GetPaginatedPaymentModes(c *fiber.Ctx) (*utils.Pagination, []saleRegistration.SalePaymentMode, error)
	FindSalePaymentModeByIdAndSalePaymentId(id, salePaymentId string) (*saleRegistration.SalePaymentMode, error)
	FindSalePaymentModeById(id string) (*saleRegistration.SalePaymentMode, error)
	UpdateSalePaymentMode(payment *saleRegistration.SalePaymentMode) error
	DeleteSalePaymentModeByID(id string) error
	GetPaginatedModes(c *fiber.Ctx, mode string) (*utils.Pagination, []saleRegistration.SalePaymentMode, error)

	// Payment Deposit
	CreatePaymentDeposit(payment *saleRegistration.SalePaymentDeposit) error
	GetPaginatedPaymentDeposits(c *fiber.Ctx) (*utils.Pagination, []saleRegistration.SalePaymentDeposit, error)
	FindSalePaymentDepositByIdAndSalePaymentId(id, salePaymentId string) (*saleRegistration.SalePaymentDeposit, error)
	FindSalePaymentDepositById(id string) (*saleRegistration.SalePaymentDeposit, error)
	UpdateSalePaymentDeposit(payment *saleRegistration.SalePaymentDeposit) error
	DeleteSalePaymentDepositByID(id string) error
	GetPaymentDeposits(c *fiber.Ctx, name string) (*utils.Pagination, []saleRegistration.SalePaymentDeposit, error)

	// Customer statement
	GenerateCustomerStatement(customerID uint) (*CustomerStatement, error)
}

type SaleRepositoryImpl struct {
	db *gorm.DB
}

func NewSaleRepository(db *gorm.DB) SaleRepository {
	return &SaleRepositoryImpl{db: db}
}

func (r *SaleRepositoryImpl) CreateSale(sale *saleRegistration.Sale) error {
	return r.db.Create(sale).Error
}

func (r *SaleRepositoryImpl) GetSalePayments(saleID uint) ([]saleRegistration.SalePayment, error) {
	var payments []saleRegistration.SalePayment
	err := r.db.Where("sale_id = ?", saleID).Find(&payments).Error
	return payments, err
}

func (r *SaleRepositoryImpl) GetSalePaymentModes(paymentID uint) ([]saleRegistration.SalePaymentMode, error) {
	var payment_modes []saleRegistration.SalePaymentMode
	err := r.db.Where("sale_payment_id = ?", paymentID).Find(&payment_modes).Error
	return payment_modes, err
}

// ====================

func (r *SaleRepositoryImpl) GetPaginatedSales(c *fiber.Ctx) (*utils.Pagination, []saleRegistration.Sale, error) {
	pagination, sales, err := utils.Paginate(c, r.db.Preload("Car"), saleRegistration.Sale{})
	if err != nil {
		return nil, nil, err
	}
	return &pagination, sales, nil
}

func (r *SaleRepositoryImpl) GetSaleByID(id string) (saleRegistration.Sale, error) {
	var sale saleRegistration.Sale
	err := r.db.Preload("Car").First(&sale, "id = ?", id).Error
	return sale, err
}

func (r *SaleRepositoryImpl) UpdateSale(sale *saleRegistration.Sale) error {
	return r.db.Save(sale).Error
}

// DeleteByID deletes a sale by ID
func (r *SaleRepositoryImpl) DeleteByID(id string) error {
	if err := r.db.Delete(&saleRegistration.Sale{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

// CreateCustomerContact creates a new payment deposit in the database
func (r *SaleRepositoryImpl) CreateInvoice(payment *saleRegistration.SalePayment) error {
	return r.db.Create(payment).Error
}

func (r *SaleRepositoryImpl) GetPaginatedInvoices(c *fiber.Ctx) (*utils.Pagination, []saleRegistration.SalePayment, error) {
	pagination, payments, err := utils.Paginate(c, r.db.Preload("Sale"), saleRegistration.SalePayment{})
	if err != nil {
		return nil, nil, err
	}
	return &pagination, payments, nil
}

func (r *SaleRepositoryImpl) FindSalePaymentByIdAndSaleId(id, saleId string) (*saleRegistration.SalePayment, error) {
	var payment saleRegistration.SalePayment
	result := r.db.Preload("Sale").Where("id = ? AND sale_id = ?", id, saleId).First(&payment)
	if result.Error != nil {
		return nil, result.Error
	}
	return &payment, nil
}

func (r *SaleRepositoryImpl) FindSalePaymentById(id string) (*saleRegistration.SalePayment, error) {
	var payment saleRegistration.SalePayment
	if err := r.db.First(&payment, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *SaleRepositoryImpl) UpdateSalePayment(payment *saleRegistration.SalePayment) error {
	return r.db.Save(payment).Error
}

// Delete salePayment by ID
func (r *SaleRepositoryImpl) DeleteSalePaymentByID(id string) error {
	if err := r.db.Delete(&saleRegistration.SalePayment{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

// CreateCustomerContact creates a new payment mode in the database
func (r *SaleRepositoryImpl) CreatePaymentMode(payment *saleRegistration.SalePaymentMode) error {
	return r.db.Create(payment).Error
}

func (r *SaleRepositoryImpl) GetPaginatedPaymentModes(c *fiber.Ctx) (*utils.Pagination, []saleRegistration.SalePaymentMode, error) {
	pagination, payments, err := utils.Paginate(c, r.db, saleRegistration.SalePaymentMode{})
	if err != nil {
		return nil, nil, err
	}
	return &pagination, payments, nil
}

func (r *SaleRepositoryImpl) FindSalePaymentModeByIdAndSalePaymentId(id, salePaymentId string) (*saleRegistration.SalePaymentMode, error) {
	var paymentMode saleRegistration.SalePaymentMode
	result := r.db.Preload("SalePayment").Where("id = ? AND sale_payment_id = ?", id, salePaymentId).First(&paymentMode)
	if result.Error != nil {
		return nil, result.Error
	}
	return &paymentMode, nil
}

func (r *SaleRepositoryImpl) GetPaginatedModes(c *fiber.Ctx, mode string) (*utils.Pagination, []saleRegistration.SalePaymentMode, error) {
	pagination, modes, err := utils.Paginate(c, r.db.Preload("SalePayment").Where("mode_of_payment = ?", mode), saleRegistration.SalePaymentMode{})
	if err != nil {
		return nil, nil, err
	}
	return &pagination, modes, nil
}

func (r *SaleRepositoryImpl) FindSalePaymentModeById(id string) (*saleRegistration.SalePaymentMode, error) {
	var mode saleRegistration.SalePaymentMode
	if err := r.db.First(&mode, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &mode, nil
}

func (r *SaleRepositoryImpl) UpdateSalePaymentMode(payment *saleRegistration.SalePaymentMode) error {
	return r.db.Save(payment).Error
}

// Delete salePayment by ID
func (r *SaleRepositoryImpl) DeleteSalePaymentModeByID(id string) error {
	if err := r.db.Delete(&saleRegistration.SalePaymentMode{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

// CreateCustomerContact creates a new payment deposit in the database
func (r *SaleRepositoryImpl) CreatePaymentDeposit(deposit *saleRegistration.SalePaymentDeposit) error {
	return r.db.Create(deposit).Error
}

func (r *SaleRepositoryImpl) GetPaginatedPaymentDeposits(c *fiber.Ctx) (*utils.Pagination, []saleRegistration.SalePaymentDeposit, error) {
	pagination, deposits, err := utils.Paginate(c, r.db.Preload("SalePayment"), saleRegistration.SalePaymentDeposit{})
	if err != nil {
		return nil, nil, err
	}
	return &pagination, deposits, nil
}

func (r *SaleRepositoryImpl) FindSalePaymentDepositByIdAndSalePaymentId(id, salePaymentId string) (*saleRegistration.SalePaymentDeposit, error) {
	var paymentDeposit saleRegistration.SalePaymentDeposit
	result := r.db.Preload("SalePayment").Where("id = ? AND sale_payment_id = ?", id, salePaymentId).First(&paymentDeposit)
	if result.Error != nil {
		return nil, result.Error
	}
	return &paymentDeposit, nil
}

func (r *SaleRepositoryImpl) GetPaymentDeposits(c *fiber.Ctx, name string) (*utils.Pagination, []saleRegistration.SalePaymentDeposit, error) {
	pagination, deposits, err := utils.Paginate(c, r.db.Preload("SalePayment").Where("bank_name = ?", name), saleRegistration.SalePaymentDeposit{})
	if err != nil {
		return nil, nil, err
	}
	return &pagination, deposits, nil
}

func (r *SaleRepositoryImpl) FindSalePaymentDepositById(id string) (*saleRegistration.SalePaymentDeposit, error) {
	var deposit saleRegistration.SalePaymentDeposit
	if err := r.db.First(&deposit, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &deposit, nil
}

func (r *SaleRepositoryImpl) UpdateSalePaymentDeposit(deposit *saleRegistration.SalePaymentDeposit) error {
	return r.db.Save(deposit).Error
}

// Delete salePayment by ID
func (r *SaleRepositoryImpl) DeleteSalePaymentDepositByID(id string) error {
	if err := r.db.Delete(&saleRegistration.SalePaymentDeposit{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

func (r *SaleRepositoryImpl) GetOutstandingBalanceByCustomerID(customerID uint) (float64, error) {
	var sales []saleRegistration.Sale
	var totalPaid float64
	var totalOutstanding float64

	// Fetch all sales linked to cars owned by the customer
	if err := r.db.
		Joins("JOIN cars ON sales.car_id = cars.id").
		Where("cars.customer_id = ?", customerID).
		Find(&sales).Error; err != nil {
		return 0, err
	}

	// If no sales found, return zero balance
	if len(sales) == 0 {
		return 0, fmt.Errorf("no sales found for this customer")
	}

	// Iterate over each sale to compute outstanding balances
	for _, sale := range sales {
		totalPaid = 0

		// Sum all payments for the sale
		if err := r.db.Model(&saleRegistration.SalePayment{}).
			Where("sale_id = ?", sale.ID).
			Select("COALESCE(SUM(amount_payed), 0)").
			Scan(&totalPaid).Error; err != nil {
			return 0, err
		}

		// Calculate outstanding balance for this sale
		outstandingBalance := sale.TotalPrice - totalPaid
		if outstandingBalance > 0 {
			totalOutstanding += outstandingBalance
		}
	}

	return totalOutstanding, nil
}

// =========================

type PaymentRecord struct {
	PaymentDate time.Time `json:"payment_date"`
	Amount      float64   `json:"amount"`
}

type SaleStatement struct {
	CarID             uint            `json:"car_id"`
	CarModel          string          `json:"car_model"`
	ChasisNumber      string          `json:"chasis_number"`
	TotalSaleAmount   float64         `json:"total_sale_amount"`
	TotalPaid         float64         `json:"total_paid"`
	OutstandingAmount float64         `json:"outstanding_amount"`
	Payments          []PaymentRecord `json:"payments"`
}

type CustomerStatement struct {
	CustomerID       uint            `json:"customer_id"`
	CustomerName     string          `json:"customer_name"`
	TotalSales       float64         `json:"total_sales"`
	TotalPaid        float64         `json:"total_paid"`
	TotalOutstanding float64         `json:"total_outstanding"`
	Sales            []SaleStatement `json:"sales"`
}

// Helper function to extract car IDs from []Car
func getCarIDs(cars []carRegistration.Car) []uint {
	var ids []uint
	for _, car := range cars {
		ids = append(ids, car.ID)
	}
	return ids
}

// Helper function to get car model by ID
func getCarModelByID(cars []carRegistration.Car, carID uint) string {
	for _, car := range cars {
		if car.ID == carID {
			return car.CarModel
		}
	}
	return ""
}

// Helper function to get chasis number by ID
func getChasisNumberByID(cars []carRegistration.Car, carID uint) string {
	for _, car := range cars {
		if car.ID == carID {
			return car.ChasisNumber
		}
	}
	return ""
}

func (r *SaleRepositoryImpl) GenerateCustomerStatement(customerID uint) (*CustomerStatement, error) {
	var customer customerRegistration.Customer
	var cars []carRegistration.Car
	var sales []saleRegistration.Sale
	var totalSales, totalPaid, totalOutstanding float64

	// Fetch Customer details
	if err := r.db.First(&customer, customerID).Error; err != nil {
		return nil, fmt.Errorf("customer not found")
	}

	// Fetch all cars owned by the customer
	if err := r.db.Where("customer_id = ?", customerID).Find(&cars).Error; err != nil {
		return nil, err
	}

	// Fetch all sales linked to the customer’s cars
	if err := r.db.Where("car_id IN (?)", getCarIDs(cars)).Find(&sales).Error; err != nil {
		return nil, err
	}

	var saleStatements []SaleStatement

	// Process each sale
	for _, sale := range sales {
		var payments []PaymentRecord
		var totalSalePaid float64

		// Fetch all payments for this sale
		rows, err := r.db.Raw(`
			SELECT amount_payed, payment_date 
			FROM sale_payments 
			WHERE sale_id = ? 
			ORDER BY payment_date ASC`, sale.ID).Rows()
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var payment PaymentRecord
			if err := rows.Scan(&payment.Amount, &payment.PaymentDate); err != nil {
				return nil, err
			}
			payments = append(payments, payment)
			totalSalePaid += payment.Amount
		}

		// Calculate outstanding balance
		outstanding := sale.TotalPrice - totalSalePaid

		// Create Sale Statement
		saleStatements = append(saleStatements, SaleStatement{
			CarID:             sale.CarID,
			CarModel:          getCarModelByID(cars, sale.CarID),
			ChasisNumber:      getChasisNumberByID(cars, sale.CarID),
			TotalSaleAmount:   sale.TotalPrice,
			TotalPaid:         totalSalePaid,
			OutstandingAmount: outstanding,
			Payments:          payments,
		})

		// Update totals
		totalSales += sale.TotalPrice
		totalPaid += totalSalePaid
		totalOutstanding += outstanding
	}

	// Create and return the full statement
	return &CustomerStatement{
		CustomerID:       customerID,
		CustomerName:     customer.Surname + " " + customer.Firstname + " " + customer.Othername,
		TotalSales:       totalSales,
		TotalPaid:        totalPaid,
		TotalOutstanding: totalOutstanding,
		Sales:            saleStatements,
	}, nil
}
