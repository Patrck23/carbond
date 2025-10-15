package repository

import (
	"car-bond/internals/middleware"
	"car-bond/internals/models/carRegistration"
	"car-bond/internals/models/customerRegistration"
	"car-bond/internals/models/saleRegistration"
	"car-bond/internals/utils"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	GetSalesSummary(companyID uint) (map[string]float64, error)
	CheckPaymentNotifications(c *fiber.Ctx) ([]Notification, error)
}

type SaleRepositoryImpl struct {
	db *gorm.DB
}

func NewSaleRepository(db *gorm.DB) SaleRepository {
	return &SaleRepositoryImpl{db: db}
}

func (r *SaleRepositoryImpl) GetSalesSummary(companyID uint) (map[string]float64, error) {
	summary := make(map[string]float64)

	// Total sales for the company
	var totalSales float64
	if err := r.db.Model(&saleRegistration.Sale{}).
		Where("company_id = ?", companyID).
		Select("COALESCE(SUM(total_price), 0)").Scan(&totalSales).Error; err != nil {
		return nil, err
	}
	summary["total_sales"] = totalSales

	// Total payments for the company
	var totalPayments float64
	if err := r.db.Model(&saleRegistration.SalePayment{}).
		Joins("JOIN sales ON sales.id = sale_payments.sale_id").
		Where("sales.company_id = ?", companyID).
		Select("COALESCE(SUM(sale_payments.amount_payed), 0)").Scan(&totalPayments).Error; err != nil {
		return nil, err
	}
	summary["total_payments"] = totalPayments
	summary["money_in_mrkt"] = totalSales - totalPayments

	// Total deposits for the company
	var totalDeposits float64
	if err := r.db.Model(&saleRegistration.SalePaymentDeposit{}).
		Joins("JOIN sale_payments ON sale_payments.id = sale_payment_deposits.sale_payment_id").
		Joins("JOIN sales ON sales.id = sale_payments.sale_id").
		Where("sales.company_id = ?", companyID).
		Select("COALESCE(SUM(sale_payment_deposits.amount_deposited), 0)").Scan(&totalDeposits).Error; err != nil {
		return nil, err
	}
	summary["total_deposits"] = totalDeposits

	return summary, nil
}

const (
	CarStatusInTransit = "InTransit"
	CarStatusSold      = "Sold"
)

var ErrCarAlreadySold = errors.New("car already sold")

func (r *SaleRepositoryImpl) CreateSale(sale *saleRegistration.Sale) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Lock the car row to prevent concurrent sales
		var car carRegistration.Car
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&car, sale.CarID).Error; err != nil {
			return fmt.Errorf("failed to find car with ID %d: %w", sale.CarID, err)
		}

		// Check if car is already sold
		if strings.EqualFold(car.CarStatus, CarStatusSold) {
			return fmt.Errorf("car with ID %d is already sold", sale.CarID)
		}

		// Create the sale
		if err := tx.Create(sale).Error; err != nil {
			return fmt.Errorf("failed to create sale: %w", err)
		}

		// Update car status to "Sold"
		result := tx.Model(&carRegistration.Car{}).
			Where("id = ?", sale.CarID).
			Update("car_status", CarStatusSold)

		if result.Error != nil {
			return fmt.Errorf("failed to update car status: %w", result.Error)
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("failed to update car status: no rows affected (possible race condition)")
		}

		return nil
	})
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

	_, company_id, err := middleware.GetUserAndCompanyFromJWT(c)
	if err != nil {
		return nil, nil, err
	}
	// Start building the query
	query := r.db.Preload("Car").Model(&saleRegistration.Sale{})
	// Apply filters based on provided parameters

	query = query.Where("company_id = ?", company_id)

	pagination, sales, err := utils.Paginate(c, query, saleRegistration.Sale{})
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

type Notification struct {
	SaleID       uint
	CustomerName string
	PhoneNumber  string
	Message      string
	DueDate      time.Time
	AmountDue    float64
	CreatedAt    time.Time
}

func (r *SaleRepositoryImpl) CheckPaymentNotifications(c *fiber.Ctx) ([]Notification, error) {
	var sales []saleRegistration.Sale
	query := r.db.
		Preload("Customer").
		Preload("Car").
		Preload("Company")

	// Apply company filter if provided
	companyID := c.Query("company_id")
	if companyID != "" {
		if id, err := strconv.Atoi(companyID); err == nil {
			query = query.Where("company_id = ?", id)
		}
	}

	if err := query.Find(&sales).Error; err != nil {
		return nil, err
	}

	var notifications []Notification
	now := time.Now()

	for _, sale := range sales {
		// Calculate total paid
		var totalPaid float64
		var payments []saleRegistration.SalePayment
		if err := r.db.Where("sale_id = ?", sale.ID).Find(&payments).Error; err != nil {
			return nil, err
		}
		for _, p := range payments {
			totalPaid += p.AmountPayed
		}

		// Skip if fully paid
		if totalPaid >= sale.TotalPrice {
			continue
		}

		// Car number plate
		carPlate := sale.Car.NumberPlate
		if carPlate == "" {
			carPlate = "N/A"
		}

		// Parse sale date
		saleDate := parseDate(sale.SaleDate)
		fmt.Println("Parsed sale date:", saleDate)

		var dueDate time.Time
		if sale.IsFullPayment {
			// Full payment: set due date to sale date (or could be nil)
			dueDate = saleDate
			fmt.Println(dueDate)
			fmt.Println(saleDate)
		} else {
			// Installments: calculate due date based on last payment or sale date
			if len(payments) == 0 {
				dueDate = saleDate.AddDate(0, sale.PaymentPeriod, 0) // sale date + total periods
			} else {
				// Last payment date
				lastPaymentDate := parseDate(payments[len(payments)-1].PaymentDate)
				dueDate = lastPaymentDate.AddDate(0, 1, 0) // next installment due next month
			}
		}

		// Message
		message := ""
		if sale.IsFullPayment {
			message = fmt.Sprintf("Full payment not received for sale #%d (Car: %s)", sale.ID, carPlate)
		} else {
			message = fmt.Sprintf("Installment payment overdue for sale #%d (Car: %s)", sale.ID, carPlate)
		}

		// Append notification
		notifications = append(notifications, Notification{
			SaleID:       sale.ID,
			CustomerName: sale.Customer.Firstname,
			PhoneNumber:  sale.Customer.Telephone,
			Message:      message,
			DueDate:      dueDate,
			AmountDue:    sale.TotalPrice - totalPaid,
			CreatedAt:    now,
		})
	}

	return notifications, nil
}

func parseDate(dateStr string) time.Time {
	t, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		fmt.Println("Warning: failed to parse date:", dateStr, err)
		return time.Time{}
	}
	return t
}

func (r *SaleRepositoryImpl) DeleteByID(id string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Fetch the sale first to get CarID
		var saleRecord saleRegistration.Sale
		if err := tx.First(&saleRecord, "id = ?", id).Error; err != nil {
			return fmt.Errorf("failed to fetch sale to delete: %w", err)
		}

		// Get all payments for this sale
		var payments []saleRegistration.SalePayment
		if err := tx.Where("sale_id = ?", id).Find(&payments).Error; err != nil {
			return fmt.Errorf("failed to fetch sale payments: %w", err)
		}

		// Collect payment IDs
		var paymentIDs []uint
		for _, p := range payments {
			paymentIDs = append(paymentIDs, p.ID)
		}

		// Delete deposits tied to those payments
		if len(paymentIDs) > 0 {
			if err := tx.Where("sale_payment_id IN ?", paymentIDs).
				Delete(&saleRegistration.SalePaymentDeposit{}).Error; err != nil {
				return fmt.Errorf("failed to delete sale payment deposits: %w", err)
			}
		}

		// Delete the payments
		if err := tx.Where("sale_id = ?", id).
			Delete(&saleRegistration.SalePayment{}).Error; err != nil {
			return fmt.Errorf("failed to delete sale payments: %w", err)
		}

		// Delete the sale
		if err := tx.Delete(&saleRegistration.Sale{}, "id = ?", id).Error; err != nil {
			return fmt.Errorf("failed to delete sale: %w", err)
		}

		// Update the car status back to "InTransit"
		result := tx.Model(&carRegistration.Car{}).
			Where("id = ?", saleRecord.CarID).
			Update("car_status", "InStock")

		if result.Error != nil {
			return fmt.Errorf("failed to update car status back to InTransit: %w", result.Error)
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("failed to update car status: no rows affected (possible race condition)")
		}

		return nil
	})
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
