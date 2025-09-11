package repository

import (
	"car-bond/internals/models/alertRegistration"
	"car-bond/internals/models/carRegistration"
	"car-bond/internals/models/companyRegistration"
	"car-bond/internals/utils"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type CarRepository interface {
	CreateCar(car *carRegistration.Car) error
	GetPaginatedCars(c *fiber.Ctx) (*utils.Pagination, []carRegistration.Car, error)
	GetCarExpenses(carID uint) ([]carRegistration.CarExpense, error)
	GetCarByID(id string) (carRegistration.Car, error)
	GetCarByVin(ChasisNumber string) (carRegistration.Car, error)
	UpdateCar(car *carRegistration.Car) error
	DeleteByID(id string) error
	SearchPaginatedCars(c *fiber.Ctx) (*utils.Pagination, []carRegistration.Car, error)
	UpdateCarJapan(id string, updates map[string]interface{}) error
	CountCarsByInvoiceExcludingID(invoiceID uint, excludeCarID uint) (int64, error)
	GetCompanyNameByID(id uint) (string, error)
	GetInvoiceByID(id uint) (carRegistration.CarShippingInvoice, error)
	UpdateCarStatusByID(carID uint, status string) error

	// Expense
	CreateCarExpense(expense *carRegistration.CarExpense) error
	GetPaginatedExpenses(c *fiber.Ctx) (*utils.Pagination, []carRegistration.CarExpense, error)
	FindCarExpenseByIdAndCarId(id string, carId string) (*carRegistration.CarExpense, error)
	GetPaginatedExpensesByCarId(c *fiber.Ctx, carId string) (*utils.Pagination, []carRegistration.CarExpense, error)
	FindCarExpenseById(id string) (*carRegistration.CarExpense, error)
	UpdateCarExpense(expense *carRegistration.CarExpense) error
	DeleteCarExpense(expense *carRegistration.CarExpense) error
	FindCarExpenseByCarAndId(carId, expenseId string) (*carRegistration.CarExpense, error)
	FindCarExpensesByCarIdAndExpenseDate(carId, expenseDate string) ([]carRegistration.CarExpense, error)
	FindCarExpensesByCarIdAndExpenseDescription(carId, expenseDescription string) ([]carRegistration.CarExpense, error)
	FindCarExpensesByCarIdAndCurrency(carId, currency string) ([]carRegistration.CarExpense, error)
	FindCarExpensesByThree(carId, expenseDate, currency string) ([]carRegistration.CarExpense, error)
	GetCarExpensesByFour(carId, expenseDate, expenseDescription, currency string) ([]carRegistration.CarExpense, error)
	UpdateCarWithExpenses(car *carRegistration.Car, expenses []carRegistration.CarExpense) error
	GetTotalCarExpenses(carID uint) (CarExpenseResponse, error)

	CreateCarPhotos(photos []carRegistration.CarPhoto) error
	CreateCarExpenses(expenses []carRegistration.CarExpense) error
	// Photos
	CreateCarPhoto(photo *carRegistration.CarPhoto) error
	DeleteCarPhotoByURL(photoURL string) error
	DeleteCarPhotoByID(photoID uint) error
	GetCarPhotosBycarID(carId uint) ([]carRegistration.CarPhoto, error)
	DeleteCarPhotos(carID uint) error

	CreateAlert(alert *alertRegistration.Transaction) error

	GetTotalCars() (int64, error)
	GetDisbandedCars() (int64, error)
	GetCarsInStock() (int64, error)
	GetTotalMoneySpent() (float64, error)
	GetTotalCarsExpenses() (map[string]float64, error)

	GetComTotalCars(companyID uint) (int64, error)
	GetComCarsInStock(companyID uint) (int64, error)
	GetComCarsSold(companyID uint) (int64, error)
	GetComTotalMoneySpent(companyID uint) (float64, error)
	GetComTotalCarsExpenses(companyID uint) (map[string]float64, error)
}

type CarRepositoryImpl struct {
	db *gorm.DB
}

func NewCarRepository(db *gorm.DB) CarRepository {
	return &CarRepositoryImpl{db: db}
}

func (r *CarRepositoryImpl) CountCarsByInvoiceExcludingID(invoiceID uint, excludeCarID uint) (int64, error) {
	var count int64
	err := r.db.Model(&carRegistration.Car{}).
		Where("car_shipping_invoice_id = ? AND id != ?", invoiceID, excludeCarID).
		Count(&count).Error
	return count, err
}

func (r *CarRepositoryImpl) CreateAlert(alert *alertRegistration.Transaction) error {
	return r.db.Create(alert).Error
}

func (r *CarRepositoryImpl) CreateCarPhoto(photo *carRegistration.CarPhoto) error {
	return r.db.Create(photo).Error
}

func (r *CarRepositoryImpl) DeleteCarPhotoByURL(photoURL string) error {
	return r.db.Where("url = ?", photoURL).Delete(&carRegistration.CarPhoto{}).Error
}

func (r *CarRepositoryImpl) DeleteCarPhotoByID(photoID uint) error {
	return r.db.Where("id = ?", photoID).Delete(&carRegistration.CarPhoto{}).Error
}

func (r *CarRepositoryImpl) CreateCar(car *carRegistration.Car) error {
	return r.db.Create(car).Error
}

func (r *CarRepositoryImpl) GetPaginatedCars(c *fiber.Ctx) (*utils.Pagination, []carRegistration.Car, error) {
	pagination, cars, err := utils.Paginate(c, r.db.Preload("CarPhotos"), carRegistration.Car{})
	if err != nil {
		return nil, nil, err
	}
	return &pagination, cars, nil
}

func (r *CarRepositoryImpl) GetCarExpenses(carID uint) ([]carRegistration.CarExpense, error) {
	var expenses []carRegistration.CarExpense
	err := r.db.Where("car_id = ?", carID).Find(&expenses).Error
	return expenses, err
}

// ====================

func (r *CarRepositoryImpl) GetCarByID(id string) (carRegistration.Car, error) {
	var car carRegistration.Car
	err := r.db.Preload("CarPhotos").First(&car, "id = ?", id).Error
	return car, err
}

func (r *CarRepositoryImpl) GetCarByVin(ChasisNumber string) (carRegistration.Car, error) {
	var car carRegistration.Car
	err := r.db.Preload("CarPhotos").First(&car, "chasis_number = ?", ChasisNumber).Error
	return car, err
}

func (r *CarRepositoryImpl) UpdateCar(car *carRegistration.Car) error {
	return r.db.Save(car).Error
}

func (r *CarRepositoryImpl) UpdateCarJapan(id string, updates map[string]interface{}) error {
	return r.db.Model(&carRegistration.Car{}).Where("id = ?", id).Updates(updates).Error
}

func (r *CarRepositoryImpl) GetCompanyNameByID(id uint) (string, error) {
	var company companyRegistration.Company
	if err := r.db.Select("name").First(&company, id).Error; err != nil {
		return "", err
	}
	return company.Name, nil
}

func (r *CarRepositoryImpl) GetInvoiceByID(id uint) (carRegistration.CarShippingInvoice, error) {
	var invoice carRegistration.CarShippingInvoice
	err := r.db.First(&invoice, id).Error
	return invoice, err
}

// DeleteByID deletes a car by ID
func (r *CarRepositoryImpl) DeleteByID(id string) error {
	if err := r.db.Delete(&carRegistration.Car{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

// CreateCarExpense creates a new car expense in the database
func (r *CarRepositoryImpl) CreateCarExpense(expense *carRegistration.CarExpense) error {
	return r.db.Create(expense).Error
}

func (r *CarRepositoryImpl) GetPaginatedExpenses(c *fiber.Ctx) (*utils.Pagination, []carRegistration.CarExpense, error) {
	pagination, expenses, err := utils.Paginate(c, r.db.Preload("Car"), carRegistration.CarExpense{})
	if err != nil {
		return nil, nil, err
	}
	return &pagination, expenses, nil
}

func (r *CarRepositoryImpl) FindCarExpenseByIdAndCarId(id string, carId string) (*carRegistration.CarExpense, error) {
	var expense carRegistration.CarExpense
	result := r.db.Preload("Car").Where("id = ? AND car_id = ?", id, carId).First(&expense)
	if result.Error != nil {
		return nil, result.Error
	}
	return &expense, nil
}

func (r *CarRepositoryImpl) GetPaginatedExpensesByCarId(c *fiber.Ctx, carId string) (*utils.Pagination, []carRegistration.CarExpense, error) {
	pagination, expenses, err := utils.Paginate(c, r.db.Preload("Car").Where("car_id = ?", carId), carRegistration.CarExpense{})
	if err != nil {
		return nil, nil, err
	}
	return &pagination, expenses, nil
}

func (r *CarRepositoryImpl) FindCarExpenseById(id string) (*carRegistration.CarExpense, error) {
	var expense carRegistration.CarExpense
	if err := r.db.First(&expense, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &expense, nil
}

func (r *CarRepositoryImpl) UpdateCarExpense(expense *carRegistration.CarExpense) error {
	return r.db.Save(expense).Error
}

func (r *CarRepositoryImpl) FindCarExpenseByCarAndId(carId, expenseId string) (*carRegistration.CarExpense, error) {
	var expense carRegistration.CarExpense
	if err := r.db.Where("id = ? AND car_id = ?", expenseId, carId).First(&expense).Error; err != nil {
		return nil, err
	}
	return &expense, nil
}

func (r *CarRepositoryImpl) DeleteCarExpense(expense *carRegistration.CarExpense) error {
	return r.db.Delete(expense).Error
}

func (r *CarRepositoryImpl) FindCarExpensesByCarIdAndExpenseDate(carId, expenseDate string) ([]carRegistration.CarExpense, error) {
	var expenses []carRegistration.CarExpense
	if err := r.db.Where("car_id = ? AND expense_date = ?", carId, expenseDate).Find(&expenses).Error; err != nil {
		return nil, err
	}
	return expenses, nil
}

func (r *CarRepositoryImpl) FindCarExpensesByCarIdAndExpenseDescription(carId, expenseDescription string) ([]carRegistration.CarExpense, error) {
	var expenses []carRegistration.CarExpense
	if err := r.db.Where("car_id = ? AND description = ?", carId, expenseDescription).Find(&expenses).Error; err != nil {
		return nil, err
	}
	return expenses, nil
}

func (r *CarRepositoryImpl) FindCarExpensesByCarIdAndCurrency(carId, currency string) ([]carRegistration.CarExpense, error) {
	var expenses []carRegistration.CarExpense
	if err := r.db.Where("car_id = ? AND currency = ?", carId, currency).Find(&expenses).Error; err != nil {
		return nil, err
	}
	return expenses, nil
}

func (r *CarRepositoryImpl) FindCarExpensesByThree(carId, expenseDate, currency string) ([]carRegistration.CarExpense, error) {
	var expenses []carRegistration.CarExpense
	if err := r.db.Where("car_id = ? AND expense_date = ? AND currency = ?", carId, expenseDate, currency).Find(&expenses).Error; err != nil {
		return nil, err
	}
	return expenses, nil
}

func (r *CarRepositoryImpl) GetCarExpensesByFour(carId, expenseDate, expenseDescription, currency string) ([]carRegistration.CarExpense, error) {
	var expenses []carRegistration.CarExpense
	if err := r.db.Where("car_id = ? AND expense_date = ? AND description = ? AND currency = ?", carId, expenseDate, expenseDescription, currency).Find(&expenses).Error; err != nil {
		return nil, err
	}
	return expenses, nil
}

type TotalCarExpense struct {
	Currency   string  `json:"currency"`
	DollarRate float64 `json:"dollar_rate"`
	Total      float64 `json:"total"`
}

// CarExpenseResponse represents the response structure for car expenses
type CarExpenseResponse struct {
	TotalCarPriceJapan            float64           `json:"total_car_price_japan"`              // Total car price in dollars
	BidPrice                      float64           `json:"bid_price"`                          // Bid price in dollars
	VATPrice                      float64           `json:"vat_tax"`                            // VAT price in dollars
	TotalExpenseJapan             float64           `json:"total_expense_japan"`                // Total expenses in dollars
	Expenses                      []TotalCarExpense `json:"expenses"`                           // List of expenses
	TotalCarPriceAndExpensesJapan float64           `json:"total_car_price_and_expenses_japan"` // Total car price and expenses
	TotalExpenseOther             float64           `json:"total_expense_other"`                // Total expenses in dollars
}

// GetTotalCarExpenses calculates the total expenses for a given car by its ID
func (r *CarRepositoryImpl) GetTotalCarExpenses(carID uint) (CarExpenseResponse, error) {
	var expenses []TotalCarExpense
	var carDetails struct {
		Currency   string  `json:"currency"`
		DollarRate float64 `json:"dollar_rate"`
		VATTax     float64 `json:"vat_tax"`
		BidPrice   float64 `json:"bid_price"`
	}

	// Fetch car details
	err := r.db.Model(&carRegistration.Car{}).
		Select("currency, vat_tax, bid_price").
		Where("id = ?", carID).
		Scan(&carDetails).Error

	if err != nil {
		return CarExpenseResponse{}, err
	}

	// Debugging output
	fmt.Printf("Car Details: %+v\n", carDetails)

	// Sum up all expenses for the given car ID, grouped by currency
	err = r.db.Model(&carRegistration.CarExpense{}).
		Where("car_id = ?", carID).
		Select(`
			currency, 
			dollar_rate, 
			COALESCE(
				SUM(
					CASE 
						WHEN currency = 'JPY' THEN amount * (1 + expense_vat/100)
						ELSE (amount * (1 + expense_vat/100)) / dollar_rate
					END
				), 
				0
			) AS total
		`).
		Group("currency, dollar_rate").
		Scan(&expenses).Error

	if err != nil {
		return CarExpenseResponse{}, err
	}

	totalExpenseYen := 0.0
	totalExpenseOther := 0.0

	for _, expense := range expenses {
		if expense.Currency == "JPY" {
			totalExpenseYen += expense.Total // Sum Yen separately
		} else {
			totalExpenseOther += expense.Total // Sum other currencies converted to dollars
		}
	}

	// Calculate VAT price in dollars
	vatPrice := (carDetails.BidPrice * carDetails.VATTax / 100) // / carDetails.DollarRate

	// Prepare the response
	response := CarExpenseResponse{
		TotalCarPriceJapan:            (carDetails.BidPrice) + vatPrice,                   // / carDetails.DollarRate // Total car price in dollars
		BidPrice:                      carDetails.BidPrice,                                // / carDetails.DollarRate,                                      // Bid price in dollars
		VATPrice:                      vatPrice,                                           // VAT price in dollars
		TotalExpenseJapan:             totalExpenseYen,                                    // Total expenses in dollars
		Expenses:                      expenses,                                           // List of expenses
		TotalCarPriceAndExpensesJapan: (carDetails.BidPrice) + vatPrice + totalExpenseYen, // Total car price and expenses / carDetails.DollarRate
		TotalExpenseOther:             totalExpenseOther,
	}

	// If there are no expenses, set expenses to an empty slice
	if len(expenses) == 0 {
		response.Expenses = []TotalCarExpense{}
	}

	return response, nil
}

func (r *CarRepositoryImpl) SearchPaginatedCars(c *fiber.Ctx) (*utils.Pagination, []carRegistration.Car, error) {
	// Get query parameters from request
	chasis_number := c.Query("chasisNumber")
	make := c.Query("make")
	model := c.Query("car_model")
	colour := c.Query("colour")
	bodyType := c.Query("bodyType")
	auction := c.Query("auction")
	destination := c.Query("destination")
	port := c.Query("port")
	minMillage := c.Query("min_millage")
	maxMillage := c.Query("max_millage")
	fuel_consumption := c.Query("fuelConsumption")
	broker_name := c.Query("brokerName")
	car_tracker := c.Query("carTracker")
	car_status := c.Query("carStatus")
	car_payment_status := c.Query("carPaymentStatus")
	manufacture_year := c.Query("maunufactureYear")
	bid_price := c.Query("bidPrice")
	maxBidPrice := c.Query("max_bid_price")
	minBidPrice := c.Query("min_bid_price")
	to_company_id := c.Query("to_company_id")
	to_company := c.Query("to_company")
	from_company_id := c.Query("from_company_id")

	// Start building the query
	query := r.db.Model(&carRegistration.Car{})

	// Apply filters based on provided parameters
	if from_company_id != "" {
		if _, err := strconv.Atoi(from_company_id); err == nil {
			query = query.Where("from_company_id = ?", from_company_id)
		}
	}

	if to_company != "" {
		query = query.Where("other_entity ILIKE ?", "%"+to_company+"%")
	} else if to_company_id != "" {
		if _, err := strconv.Atoi(to_company_id); err == nil {
			query = query.Where("to_company_id = ?", to_company_id)
		}
	}

	if chasis_number != "" {
		query = query.Where("chasis_number ILIKE ?", "%"+chasis_number+"%")
	}
	if make != "" {
		query = query.Where("make ILIKE ?", "%"+make+"%")
	}
	if model != "" {
		query = query.Where("car_model = ?", model)
	}
	if colour != "" {
		query = query.Where("colour = ?", colour)
	}
	if bodyType != "" {
		query = query.Where("body_type ILIKE ?", "%"+bodyType+"%")
	}
	if auction != "" {
		query = query.Where("auction ILIKE ?", "%"+auction+"%")
	}
	if destination != "" {
		query = query.Where("destination = ?", destination)
	}
	if port != "" {
		query = query.Where("port = ?", port)
	}
	if broker_name != "" {
		query = query.Where("broker_name ILIKE ?", "%"+broker_name+"%")
	}
	if car_status != "" {
		query = query.Where("car_status = ?", car_status)
	}
	if car_payment_status != "" {
		query = query.Where("car_payment_status = ?", car_payment_status)
	}
	if fuel_consumption != "" {
		query = query.Where("car_payment_status = ?", fuel_consumption)
	}
	if car_tracker != "" {
		boolValue, err := strconv.ParseBool(car_tracker)
		if err == nil {
			query = query.Where("car_tracker = ?", boolValue)
		}
	}

	if minMillage != "" {
		minMillageValue, err := strconv.ParseFloat(minMillage, 64)
		if err == nil {
			query = query.Where("car_millage >= ?", minMillageValue)
		}
	}

	if maxMillage != "" {
		maxMillageValue, err := strconv.ParseFloat(maxMillage, 64)
		if err == nil {
			query = query.Where("car_millage <= ?", maxMillageValue)
		}
	}

	if minBidPrice != "" {
		minBidPriceValue, err := strconv.ParseFloat(minBidPrice, 64)
		if err == nil {
			query = query.Where("bid_price >= ?", minBidPriceValue)
		}
	}

	if maxBidPrice != "" {
		maxBidPriceValue, err := strconv.ParseFloat(maxBidPrice, 64)
		if err == nil {
			query = query.Where("bid_price <= ?", maxBidPriceValue)
		}
	}

	if manufacture_year != "" {
		if _, err := strconv.Atoi(manufacture_year); err == nil {
			query = query.Where("manufacture_year = ?", manufacture_year)
		}
	}

	if bid_price != "" {
		if _, err := strconv.Atoi(bid_price); err == nil {
			query = query.Where("bid_price = ?", bid_price)
		}
	}

	// Call the pagination helper
	pagination, cars, err := utils.Paginate(c, query, carRegistration.Car{})
	if err != nil {
		return nil, nil, err
	}

	return &pagination, cars, nil
}

func (r *CarRepositoryImpl) GetCarPhotosBycarID(carId uint) ([]carRegistration.CarPhoto, error) {
	var photos []carRegistration.CarPhoto
	if err := r.db.Where("car_id = ?", carId).Find(&photos).Error; err != nil {
		return nil, err
	}
	return photos, nil
}

func (r *CarRepositoryImpl) GetTotalCars() (int64, error) {
	var count int64
	err := r.db.Model(&carRegistration.Car{}).Count(&count).Error
	return count, err
}

func (r *CarRepositoryImpl) GetDisbandedCars() (int64, error) {
	var count int64
	err := r.db.Model(&carRegistration.Car{}).Where("car_status_japan != ?", "InStock").Count(&count).Error
	return count, err
}

func (r *CarRepositoryImpl) GetCarsInStock() (int64, error) {
	var count int64
	err := r.db.Model(&carRegistration.Car{}).Where("car_status_japan = ?", "InStock").Count(&count).Error
	return count, err
}

func (r *CarRepositoryImpl) GetTotalMoneySpent() (float64, error) {
	var total float64
	err := r.db.Model(&carRegistration.Car{}).Select("SUM(bid_price + ((vat_tax * bid_price)/100))").Scan(&total).Error
	return total, err
}

func (r *CarRepositoryImpl) GetTotalCarsExpenses() (map[string]float64, error) {
	var results []struct {
		Currency   string
		DollarRate float64
		Total      float64
	}

	err := r.db.Model(&carRegistration.CarExpense{}).
		Select("currency, dollar_rate, SUM(amount * (1 + (expense_vat/100.0))) as total").
		Group("currency, dollar_rate").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	// 🔎 Print raw query results
	for _, summary := range results {
		fmt.Printf("Currency: %s | Dollar Rate: %.4f | Total: %.2f\n",
			summary.Currency, summary.DollarRate, summary.Total)
	}

	// Dynamically build totals
	totalExpenses := make(map[string]float64)

	for _, res := range results {
		// Always add original currency
		totalExpenses[res.Currency] += res.Total

		// Convert to USD only if a valid dollar rate exists (> 0)
		if res.DollarRate > 0 {
			if res.Currency == "USD" {
				totalExpenses["USD"] += res.Total
			} else {
				totalExpenses["USD"] += res.Total / res.DollarRate
			}
		}
	}

	return totalExpenses, nil
}

// GetTotalCars returns the total count of cars that belong to the given company
func (r *CarRepositoryImpl) GetComTotalCars(companyID uint) (int64, error) {
	var count int64
	// A car is considered for the company if it is either sold from or bought by the company.
	err := r.db.Model(&carRegistration.Car{}).
		Where("to_company_id = ?", companyID).
		Count(&count).Error
	return count, err
}

func (r *CarRepositoryImpl) GetComCarsInStock(companyID uint) (int64, error) {
	var count int64
	err := r.db.Model(&carRegistration.Car{}).
		Where("to_company_id = ? AND LOWER(car_status) = LOWER(?)", companyID, "In stock").
		Count(&count).Error
	return count, err
}

func (r *CarRepositoryImpl) GetComCarsSold(companyID uint) (int64, error) {
	var count int64
	err := r.db.Model(&carRegistration.Car{}).
		Where("to_company_id = ? AND LOWER(car_status) = LOWER(?)", companyID, "Sold").
		Count(&count).Error
	return count, err
}

func (r *CarRepositoryImpl) GetComTotalMoneySpent(companyID uint) (float64, error) {
	var total sql.NullFloat64 // Use sql.NullFloat64 to handle NULL
	err := r.db.Model(&carRegistration.Car{}).
		Where("to_company_id = ?", companyID).
		Where("cars.deleted_at IS NULL").
		Select("SUM(bid_price + ((vat_tax * bid_price)/100))").
		Scan(&total).Error

	if err != nil {
		return 0, err
	}

	if !total.Valid {
		return 0, nil // Return 0 if there's no data to sum
	}

	return total.Float64, nil // Return the actual value
}

func (r *CarRepositoryImpl) GetComTotalCarsExpenses(companyID uint) (map[string]float64, error) {
	var results []struct {
		Currency   string
		DollarRate float64
		Total      float64
	}

	err := r.db.Model(&carRegistration.CarExpense{}).
		Joins("JOIN cars ON car_expenses.car_id = cars.id").
		Where("cars.to_company_id = ?", companyID).
		Select(`
			car_expenses.currency, 
			car_expenses.dollar_rate, 
			SUM(car_expenses.amount * (1 + (car_expenses.expense_vat/100.0))) as total
		`).
		Group("car_expenses.currency, car_expenses.dollar_rate").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}
	// 🔎 Debug print of raw results
	fmt.Println("Raw results from DB:")
	for _, res := range results {
		fmt.Printf("Currency: %s | DollarRate: %.4f | Total: %.2f\n", res.Currency, res.DollarRate, res.Total)
	}

	// Build totals dynamically for all currencies
	totalExpenses := make(map[string]float64)
	for _, res := range results {
		// If it's USD, no conversion needed
		if res.Currency == "USD" {
			totalExpenses["USD"] += res.Total
			continue
		}

		// If no dollar rate provided, skip to avoid division by zero
		if res.DollarRate == 100 {
			continue
		}

		// Always convert to USD equivalent and store original currency too
		totalExpenses[res.Currency] += res.Total
		totalExpenses["USD"] += res.Total / res.DollarRate
	}

	return totalExpenses, nil
}

// CreateCarPhotos creates car photos in the database
func (r *CarRepositoryImpl) CreateCarPhotos(photos []carRegistration.CarPhoto) error {
	if len(photos) == 0 {
		return nil
	}
	return r.db.Create(&photos).Error
}

// CreateCarExpenses creates car expenses in the database
func (r *CarRepositoryImpl) CreateCarExpenses(expenses []carRegistration.CarExpense) error {
	if len(expenses) == 0 {
		return nil
	}
	return r.db.Create(&expenses).Error
}

func (r *CarRepositoryImpl) UpdateCarWithExpenses(car *carRegistration.Car, expenses []carRegistration.CarExpense) error {
	tx := r.db.Begin()

	// Omit fields that should not be updated
	if err := tx.Model(&carRegistration.Car{}).
		Where("id = ?", car.ID).
		// Omit("car_uuid", "created_at", "created_by", "broker_name", "broker_number",
		// 	"number_plate", "car_tracker", "customer_id", "car_status",
		// 	"car_payment_status").
		Updates(car).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Delete existing expenses
	if err := tx.Where("car_id = ?", car.ID).Delete(&carRegistration.CarExpense{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Insert new expenses
	for i := range expenses {
		expenses[i].CarID = car.ID
	}
	if len(expenses) > 0 {
		if err := tx.Create(&expenses).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func (r *CarRepositoryImpl) DeleteCarPhotos(carID uint) error {
	var photos []carRegistration.CarPhoto

	// Fetch all photo records
	if err := r.db.Where("car_id = ?", carID).Find(&photos).Error; err != nil {
		return err
	}

	// Delete physical files
	for _, photo := range photos {
		filePath := photo.URL
		if strings.HasPrefix(filePath, "./uploads/car_files/") {
			filePath = strings.TrimPrefix(filePath, "./")
		}
		os.Remove(filePath)

	}

	// Delete from DB
	return r.db.Where("car_id = ?", carID).Delete(&carRegistration.CarPhoto{}).Error
}

func (r *CarRepositoryImpl) UpdateCarStatusByID(carID uint, status string) error {
	return r.db.Model(&carRegistration.Car{}).
		Where("id = ?", carID).
		Update("car_status", status).Error
}
