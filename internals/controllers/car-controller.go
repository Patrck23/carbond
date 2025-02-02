package controllers

import (
	"car-bond/internals/models/carRegistration"
	"car-bond/internals/utils"

	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type CarRepository interface {
	CreateCar(car *carRegistration.Car) error
	GetPaginatedCars(c *fiber.Ctx) (*utils.Pagination, []carRegistration.Car, error)
	GetCarExpenses(carID uint) ([]carRegistration.CarExpense, error)
	GetCarByID(id string) (carRegistration.Car, error)
	GetCarByVin(vinNumber string) (carRegistration.Car, error)
	UpdateCar(car *carRegistration.Car) error
	DeleteByID(id string) error

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
	GetTotalCarExpenses(carID uint) ([]TotalCarExpense, error)
}

type CarRepositoryImpl struct {
	db *gorm.DB
}

func NewCarRepository(db *gorm.DB) CarRepository {
	return &CarRepositoryImpl{db: db}
}

type CarController struct {
	repo CarRepository
}

func NewCarController(repo CarRepository) *CarController {
	return &CarController{repo: repo}
}

// ==================

func (r *CarRepositoryImpl) CreateCar(car *carRegistration.Car) error {
	return r.db.Create(car).Error
}

func (h *CarController) CreateCar(c *fiber.Ctx) error {
	// Initialize a new Car instance
	car := new(carRegistration.Car)

	// Parse the request body into the Car instance
	if err := c.BodyParser(car); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input provided",
			"data":    err.Error(),
		})
	}

	// Attempt to create the car record using the repository
	if err := h.repo.CreateCar(car); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create car",
			"data":    err.Error(),
		})
	}

	// Return the newly created car record
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Car created successfully",
		"data":    car,
	})
}

// ===================

func (r *CarRepositoryImpl) GetPaginatedCars(c *fiber.Ctx) (*utils.Pagination, []carRegistration.Car, error) {
	pagination, cars, err := utils.Paginate(c, r.db, carRegistration.Car{})
	if err != nil {
		return nil, nil, err
	}
	return &pagination, cars, nil
}

func (h *CarController) GetAllCars(c *fiber.Ctx) error {
	// Fetch paginated cars using the repository
	pagination, cars, err := h.repo.GetPaginatedCars(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve cars",
			"data":    err.Error(),
		})
	}

	// Initialize a response slice to hold cars with their ports and expenses
	var response []fiber.Map

	// Iterate over all cars to fetch associated car ports and expenses
	for _, car := range cars {

		expenses, err := h.repo.GetCarExpenses(car.ID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to retrieve expenses for car ID " + strconv.Itoa(int(car.ID)),
				"data":    err.Error(),
			})
		}

		// Combine car, ports, and expenses into a single response map
		response = append(response, fiber.Map{
			"car":      car,
			"expenses": expenses,
		})
	}

	// Return the paginated response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Cars and associated data retrieved successfully",
		"data":    response,
		"pagination": fiber.Map{
			"total_items":  pagination.TotalItems,
			"total_pages":  pagination.TotalPages,
			"current_page": pagination.CurrentPage,
			"limit":        pagination.ItemsPerPage,
		},
	})
}

// =====================

// ====================

func (r *CarRepositoryImpl) GetCarExpenses(carID uint) ([]carRegistration.CarExpense, error) {
	var expenses []carRegistration.CarExpense
	err := r.db.Where("car_id = ?", carID).Find(&expenses).Error
	return expenses, err
}

// ====================

func (r *CarRepositoryImpl) GetCarByID(id string) (carRegistration.Car, error) {
	var car carRegistration.Car
	err := r.db.First(&car, "id = ?", id).Error
	return car, err
}

// GetSingleCar fetches a car with its associated ports and expenses from the database
func (h *CarController) GetSingleCar(c *fiber.Ctx) error {
	// Get the car ID from the route parameters
	id := c.Params("id")

	// Fetch the car by ID
	car, err := h.repo.GetCarByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Car not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve car",
			"data":    err.Error(),
		})
	}

	// Fetch expenses associated with the car
	expenses, err := h.repo.GetCarExpenses(car.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve expenses",
			"data":    err.Error(),
		})
	}

	// Prepare the response
	response := fiber.Map{
		"car":      car,
		"expenses": expenses,
	}

	// Return the response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Car and associated data retrieved successfully",
		"data":    response,
	})
}

// ====================

func (r *CarRepositoryImpl) GetCarByVin(vinNumber string) (carRegistration.Car, error) {
	var car carRegistration.Car
	err := r.db.First(&car, "vin_number = ?", vinNumber).Error
	return car, err
}

func (h *CarController) GetSingleCarByVinNumber(c *fiber.Ctx) error {
	// Get the car ID from the route parameters
	vinNumber := c.Params("vinNumber")

	// Fetch the car by ID
	car, err := h.repo.GetCarByVin(vinNumber)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Car not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve car",
			"data":    err.Error(),
		})
	}

	// Fetch expenses associated with the car
	expenses, err := h.repo.GetCarExpenses(car.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve expenses",
			"data":    err.Error(),
		})
	}

	// Prepare the response
	response := fiber.Map{
		"car":      car,
		"expenses": expenses,
	}

	// Return the response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Car and associated data retrieved successfully",
		"data":    response,
	})
}

// ====================

func (r *CarRepositoryImpl) UpdateCar(car *carRegistration.Car) error {
	return r.db.Save(car).Error
}

// Define the updateCar struct
type UpdateCarPayload struct {
	VinNumber             string  `json:"vin_number"`
	EngineNumber          string  `json:"engine_number"`
	EngineCapacity        string  `json:"engine_capacity"`
	Make                  string  `json:"make"`
	CarModel              string  `json:"model"`
	MaximCarry            int     `json:"maxim_carry"`
	Weight                int     `json:"weight"`
	GrossWeight           int     `json:"gross_weight"`
	Length                int     `json:"length"`
	Width                 int     `json:"width"`
	Height                int     `json:"height"`
	CarMillage            int     `json:"millage"`
	FuelConsumption       string  `json:"fuel_consumption"`
	ManufactureYear       int     `json:"maunufacture_year"`
	FirstRegistrationYear int     `json:"first_registration_year"`
	Transmission          string  `json:"transmission"`
	BodyType              string  `json:"body_type"`
	Colour                string  `json:"colour"`
	Auction               string  `json:"auction"`
	Currency              string  `json:"currency"`
	PowerSteering         bool    `json:"ps"`
	PowerWindow           bool    `json:"pw"`
	ABS                   bool    `json:"abs"`
	ADS                   bool    `json:"ads"`
	AlloyWheel            bool    `json:"aw"`
	SimpleWheel           bool    `json:"sw"`
	Navigation            bool    `json:"navigation"`
	AC                    bool    `json:"ac"`
	BidPrice              float64 `json:"bid_price"`
	PurchaseDate          string  `json:"purchase_date"`
	FromCompanyID         uint    `json:"from_company_id"`
	ToCompanyID           uint    `json:"to_company_id"`
	Destination           string  `json:"destination"`
	Port                  string  `json:"port"`
	UpdatedBy             string  `json:"updated_by"`
}

type UpdateCarPayload2 struct {
	BrokerName       string  `json:"broker_name"`
	BrokerNumber     string  `json:"broker_number"`
	VATTax           float64 `json:"vat_tax"`
	NumberPlate      string  `json:"number_plate"`
	CarTracker       bool    `json:"car_tracker"`
	CarStatus        string  `json:"car_status"`
	CarPaymentStatus string  `json:"car_payment_status"`
	CustomerID       int     `json:"customer_id"`
	UpdatedBy        string  `json:"updated_by"`
}

// UpdateCar handler function
func (h *CarController) UpdateCar(c *fiber.Ctx) error {
	// Get the car ID from the route parameters
	id := c.Params("id")

	// Find the car in the database
	car, err := h.repo.GetCarByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{
				"status":  "error",
				"message": "Car not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve car",
			"data":    err.Error(),
		})
	}

	// Parse the request body into the UpdateCarPayload struct
	var payload UpdateCarPayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input",
			"data":    err.Error(),
		})
	}

	// Update the car fields using the payload
	updateCarFields(&car, payload) // Pass the parsed payload

	// Save the changes to the database
	if err := h.repo.UpdateCar(&car); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update car",
			"data":    err.Error(),
		})
	}

	// Return the updated car
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Car updated successfully",
		"data":    car,
	})
}

// UpdateCarFields updates the fields of a car using the updateCar struct
func updateCarFields(car *carRegistration.Car, updateCarData UpdateCarPayload) {
	car.VinNumber = updateCarData.VinNumber
	car.EngineNumber = updateCarData.EngineNumber
	car.EngineCapacity = updateCarData.EngineCapacity
	car.Make = updateCarData.Make
	car.CarModel = updateCarData.CarModel
	car.MaximCarry = updateCarData.MaximCarry
	car.Weight = updateCarData.Weight
	car.GrossWeight = updateCarData.GrossWeight
	car.Length = updateCarData.Length
	car.Width = updateCarData.Width
	car.Height = updateCarData.Height
	car.ManufactureYear = updateCarData.ManufactureYear
	car.FirstRegistrationYear = updateCarData.FirstRegistrationYear
	car.Transmission = updateCarData.Transmission
	car.BodyType = updateCarData.BodyType
	car.Colour = updateCarData.Colour
	car.Auction = updateCarData.Auction
	car.PowerSteering = updateCarData.PowerSteering
	car.PowerWindow = updateCarData.PowerWindow
	car.ABS = updateCarData.ABS
	car.ADS = updateCarData.ADS
	car.AlloyWheel = updateCarData.AlloyWheel
	car.SimpleWheel = updateCarData.SimpleWheel
	car.Navigation = updateCarData.Navigation
	car.AC = updateCarData.AC
	car.Currency = updateCarData.Currency
	car.BidPrice = updateCarData.BidPrice
	car.PurchaseDate = updateCarData.PurchaseDate

	if updateCarData.FromCompanyID != 0 {
		car.FromCompanyID = &updateCarData.FromCompanyID
	} else {
		car.FromCompanyID = nil
	}

	if updateCarData.ToCompanyID != 0 {
		car.ToCompanyID = &updateCarData.ToCompanyID
	} else {
		car.ToCompanyID = nil
	}

	car.Destination = updateCarData.Destination
	car.UpdatedBy = updateCarData.UpdatedBy
}

// UpdateCar handler function
func (h *CarController) UpdateCar2(c *fiber.Ctx) error {
	// Get the car ID from the route parameters
	id := c.Params("id")

	// Find the car in the database
	car, err := h.repo.GetCarByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{
				"status":  "error",
				"message": "Car not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve car",
			"data":    err.Error(),
		})
	}

	// Parse the request body into the UpdateCarPayload struct
	var payload UpdateCarPayload2
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input",
			"data":    err.Error(),
		})
	}

	// Update the car fields using the payload
	updateCar2Fields(&car, payload) // Pass the parsed payload

	// Save the changes to the database
	if err := h.repo.UpdateCar(&car); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update car",
			"data":    err.Error(),
		})
	}

	// Return the updated car
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Car updated successfully",
		"data":    car,
	})
}

// UpdateCarFields updates the fields of a car using the updateCar struct
func updateCar2Fields(car *carRegistration.Car, updateCarData UpdateCarPayload2) {
	// Update the car fields
	car.BrokerName = updateCarData.BrokerName
	car.BrokerNumber = updateCarData.BrokerNumber
	car.NumberPlate = updateCarData.NumberPlate
	car.VATTax = updateCarData.VATTax
	car.CarTracker = updateCarData.CarTracker
	car.CarStatus = updateCarData.CarStatus
	car.CarPaymentStatus = updateCarData.CarPaymentStatus
	// Assign foreign keys if provided
	if updateCarData.CustomerID != 0 {
		car.CustomerID = &updateCarData.CustomerID
	} else {
		car.CustomerID = nil
	}
	car.UpdatedBy = updateCarData.UpdatedBy
}

// ====================

// DeleteByID deletes a car by ID
func (r *CarRepositoryImpl) DeleteByID(id string) error {
	if err := r.db.Delete(&carRegistration.Car{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

// DeleteCarByID deletes a car by its ID
func (h *CarController) DeleteCarByID(c *fiber.Ctx) error {
	// Get the car ID from the route parameters
	id := c.Params("id")

	// Find the car in the database
	car, err := h.repo.GetCarByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{
				"status":  "error",
				"message": "Car not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to find car",
			"data":    err.Error(),
		})
	}

	// Delete the car
	if err := h.repo.DeleteByID(id); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to delete car",
			"data":    err.Error(),
		})
	}

	// Return success response
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Car deleted successfully",
		"data":    car,
	})
}

// Car Expense
// ===================================================================================================
// Create a car expense

// CreateCarExpense creates a new car expense in the database
func (r *CarRepositoryImpl) CreateCarExpense(expense *carRegistration.CarExpense) error {
	return r.db.Create(expense).Error
}

// CreateCarExpense handles the creation of a car expense
func (h *CarController) CreateCarExpense(c *fiber.Ctx) error {
	// Parse the request body into a CarExpense struct
	carExpense := new(carRegistration.CarExpense)
	if err := c.BodyParser(carExpense); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input data",
			"data":    err.Error(),
		})
	}

	// Create the car expense in the database
	if err := h.repo.CreateCarExpense(carExpense); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create car expense",
			"data":    err.Error(),
		})
	}

	// Return success response
	return c.Status(201).JSON(fiber.Map{
		"status":  "success",
		"message": "Car expense created successfully",
		"data":    carExpense,
	})
}

// ========================

func (r *CarRepositoryImpl) GetPaginatedExpenses(c *fiber.Ctx) (*utils.Pagination, []carRegistration.CarExpense, error) {
	pagination, expenses, err := utils.Paginate(c, r.db.Preload("Car"), carRegistration.CarExpense{})
	if err != nil {
		return nil, nil, err
	}
	return &pagination, expenses, nil
}

func (h *CarController) GetAllCarExpenses(c *fiber.Ctx) error {
	// Fetch paginated expenses using the repository
	pagination, expenses, err := h.repo.GetPaginatedExpenses(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve expenses",
			"data":    err.Error(),
		})
	}

	// Return the paginated response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Expenses and associated data retrieved successfully",
		"data":    expenses,
		"pagination": fiber.Map{
			"total_items":  pagination.TotalItems,
			"total_pages":  pagination.TotalPages,
			"current_page": pagination.CurrentPage,
			"limit":        pagination.ItemsPerPage,
		},
	})
}

// =======================

// Get Car Expenses by ID

func (r *CarRepositoryImpl) FindCarExpenseByIdAndCarId(id string, carId string) (*carRegistration.CarExpense, error) {
	var expense carRegistration.CarExpense
	result := r.db.Preload("Car").Where("id = ? AND car_id = ?", id, carId).First(&expense)
	if result.Error != nil {
		return nil, result.Error
	}
	return &expense, nil
}

func (h *CarController) GetCarExpenseById(c *fiber.Ctx) error {
	// Retrieve the expense ID and Car ID from the request parameters
	id := c.Params("id")
	carId := c.Params("carId")

	// Fetch the car expense from the repository
	expense, err := h.repo.FindCarExpenseByIdAndCarId(id, carId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Expense not found for the specified car",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch the expense",
			"error":   err.Error(),
		})
	}

	// Return the fetched expense
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Expense fetched successfully",
		"data":    expense,
	})
}

// ===============
// Get Car Expenses by Car ID
func (r *CarRepositoryImpl) GetPaginatedExpensesByCarId(c *fiber.Ctx, carId string) (*utils.Pagination, []carRegistration.CarExpense, error) {
	pagination, expenses, err := utils.Paginate(c, r.db.Preload("Car").Where("car_id = ?", carId), carRegistration.CarExpense{})
	if err != nil {
		return nil, nil, err
	}
	return &pagination, expenses, nil
}

func (h *CarController) GetCarExpensesByCarId(c *fiber.Ctx) error {
	// Retrieve carId from the request parameters
	carId := c.Params("carId")

	// Fetch paginated car expenses using the repository
	pagination, expenses, err := h.repo.GetPaginatedExpensesByCarId(c, carId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve car expenses",
			"error":   err.Error(),
		})
	}

	// Handle the case where no expenses are found
	if len(expenses) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "No expenses found for the specified car",
		})
	}

	// Return a success response with paginated data
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Car expenses fetched successfully",
		"data":    expenses,
		"pagination": fiber.Map{
			"total_items":  pagination.TotalItems,
			"total_pages":  pagination.TotalPages,
			"current_page": pagination.CurrentPage,
			"limit":        pagination.ItemsPerPage,
		},
	})
}

// ====================

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

func (h *CarController) UpdateCarExpense(c *fiber.Ctx) error {
	// Define a struct for input validation
	type UpdateCarExpenseInput struct {
		Description string  `json:"description" validate:"required"`
		Currency    string  `json:"currency" validate:"required"`
		Amount      float64 `json:"amount" validate:"required,gt=0"`
		ExpenseDate string  `json:"expense_date" validate:"required"`
		UpdatedBy   string  `json:"updated_by" validate:"required"`
	}

	// Parse the expense ID from the request parameters
	expenseID := c.Params("id")

	// Parse and validate the request body
	var input UpdateCarExpenseInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input",
			"error":   err.Error(),
		})
	}

	// Use a validation library to validate the input
	if validationErr := utils.ValidateStruct(input); validationErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Validation failed",
			"errors":  validationErr,
		})
	}

	// Fetch the expense record using the repository
	expense, err := h.repo.FindCarExpenseById(expenseID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Expense not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch expense",
			"error":   err.Error(),
		})
	}

	// Update the expense fields
	expense.Description = input.Description
	expense.Currency = input.Currency
	expense.Amount = input.Amount
	expense.ExpenseDate = input.ExpenseDate
	expense.UpdatedBy = input.UpdatedBy

	// Save the updated expense using the repository
	if err := h.repo.UpdateCarExpense(expense); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update expense",
			"error":   err.Error(),
		})
	}

	// Return the updated expense
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Expense updated successfully",
		"data":    expense,
	})
}

// ============

// Delete car Expenses by ID

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

func (h *CarController) DeleteCarExpenseById(c *fiber.Ctx) error {
	// Parse carId and expenseId from the request parameters
	carId := c.Params("carId")
	expenseId := c.Params("id")

	// Check if the expense exists and belongs to the specified car
	expense, err := h.repo.FindCarExpenseByCarAndId(carId, expenseId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Car expense not found or does not belong to the specified car",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch car expense",
			"error":   err.Error(),
		})
	}

	// Delete the expense using the repository
	if err := h.repo.DeleteCarExpense(expense); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to delete car expense",
			"error":   err.Error(),
		})
	}

	// Return success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Car expense deleted successfully",
	})
}

// =====================
// Get Car Expenses by Car ID and Expense Date

func (r *CarRepositoryImpl) FindCarExpensesByCarIdAndExpenseDate(carId, expenseDate string) ([]carRegistration.CarExpense, error) {
	var expenses []carRegistration.CarExpense
	if err := r.db.Where("car_id = ? AND expense_date = ?", carId, expenseDate).Find(&expenses).Error; err != nil {
		return nil, err
	}
	return expenses, nil
}

func (h *CarController) GetCarExpensesByCarIdAndExpenseDate(c *fiber.Ctx) error {
	// Retrieve carId and expenseDate from the request parameters
	carId := c.Params("id")
	expenseDate := c.Params("expense_date")

	// Fetch expenses using the repository
	expenses, err := h.repo.FindCarExpensesByCarIdAndExpenseDate(carId, expenseDate)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Car expenses not found for the specified date",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch car expenses",
			"error":   err.Error(),
		})
	}

	// Return the fetched expenses
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Car expenses fetched successfully",
		"data":    expenses,
	})
}

// =====================
// Get Car Expenses by Car ID and Expense Description

func (r *CarRepositoryImpl) FindCarExpensesByCarIdAndExpenseDescription(carId, expenseDescription string) ([]carRegistration.CarExpense, error) {
	var expenses []carRegistration.CarExpense
	if err := r.db.Where("car_id = ? AND description = ?", carId, expenseDescription).Find(&expenses).Error; err != nil {
		return nil, err
	}
	return expenses, nil
}

func (h *CarController) GetCarExpensesByCarIdAndExpenseDescription(c *fiber.Ctx) error {
	// Retrieve carId and expenseDescription from the request parameters
	carId := c.Params("id")
	expenseDescription := c.Params("expense_description")

	// Fetch expenses using the repository
	expenses, err := h.repo.FindCarExpensesByCarIdAndExpenseDescription(carId, expenseDescription)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Car expenses not found for the specified description",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch car expenses",
			"error":   err.Error(),
		})
	}

	// Return the fetched expenses
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Car expenses fetched successfully",
		"data":    expenses,
	})
}

// ====================
// Get Car Expenses by Car ID and Currency
func (r *CarRepositoryImpl) FindCarExpensesByCarIdAndCurrency(carId, currency string) ([]carRegistration.CarExpense, error) {
	var expenses []carRegistration.CarExpense
	if err := r.db.Where("car_id = ? AND currency = ?", carId, currency).Find(&expenses).Error; err != nil {
		return nil, err
	}
	return expenses, nil
}

func (h *CarController) GetCarExpensesByCarIdAndCurrency(c *fiber.Ctx) error {
	// Retrieve carId and currency from the request parameters
	carId := c.Params("id")
	currency := c.Params("currency")

	// Fetch expenses using the repository
	expenses, err := h.repo.FindCarExpensesByCarIdAndCurrency(carId, currency)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Car expenses not found for the specified car and currency",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch car expenses",
			"error":   err.Error(),
		})
	}

	// Return the fetched expenses
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Car expenses fetched successfully",
		"data":    expenses,
	})
}

// ====================
// Get Car Expenses by Car ID, Expense Date, and Currency
func (r *CarRepositoryImpl) FindCarExpensesByThree(carId, expenseDate, currency string) ([]carRegistration.CarExpense, error) {
	var expenses []carRegistration.CarExpense
	if err := r.db.Where("car_id = ? AND expense_date = ? AND currency = ?", carId, expenseDate, currency).Find(&expenses).Error; err != nil {
		return nil, err
	}
	return expenses, nil
}

func (h *CarController) GetCarExpensesByThree(c *fiber.Ctx) error {
	// Retrieve carId, expenseDate, and currency from the request parameters
	carId := c.Params("id")
	expenseDate := c.Params("expense_date")
	currency := c.Params("currency")

	// Fetch car expenses using the repository
	expenses, err := h.repo.FindCarExpensesByThree(carId, expenseDate, currency)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Car expenses not found for the specified car, date, and currency",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch car expenses",
			"error":   err.Error(),
		})
	}

	// Return the fetched expenses
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Car expenses fetched successfully",
		"data":    expenses,
	})
}

// ==================
// Get Car Expenses by Car ID, Expense Date, Expense Description, and Currency
func (r *CarRepositoryImpl) GetCarExpensesByFour(carId, expenseDate, expenseDescription, currency string) ([]carRegistration.CarExpense, error) {
	var expenses []carRegistration.CarExpense
	if err := r.db.Where("car_id = ? AND expense_date = ? AND description = ? AND currency = ?", carId, expenseDate, expenseDescription, currency).Find(&expenses).Error; err != nil {
		return nil, err
	}
	return expenses, nil
}

// GetCarExpensesFilters handles the request to get car expenses by filters.
func (h *CarController) GetCarExpensesFilters(c *fiber.Ctx) error {
	carId := c.Params("id")
	expenseDate := c.Params("expense_date")
	expenseDescription := c.Params("expense_description")
	currency := c.Params("currency")

	// Fetch car expenses using the service layer
	expenses, err := h.repo.GetCarExpensesByFour(carId, expenseDate, expenseDescription, currency)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch car expenses",
			"error":   err.Error(),
		})
	}

	// If no expenses found
	if len(expenses) == 0 {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Car expenses not found for the specified filters",
		})
	}

	// Return the fetched car expenses
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Car expenses fetched successfully",
		"data":    expenses,
	})
}

// Car Port
// ===================================================================================================

// TotalCarExpense represents the total expenses for a car grouped by currency
type TotalCarExpense struct {
	Currency string  `json:"currency"`
	Total    float64 `json:"total"`
}

// GetTotalCarExpenses calculates the total expenses for a given car by its ID, grouped by currency
func (r *CarRepositoryImpl) GetTotalCarExpenses(carID uint) ([]TotalCarExpense, error) {
	var expenses []TotalCarExpense

	// Sum up all expenses for the given car ID, grouped by currency
	err := r.db.Model(&carRegistration.CarExpense{}).
		Where("car_id = ?", carID).
		Select("currency, COALESCE(SUM(amount), 0) AS total").
		Group("currency").
		Scan(&expenses).Error

	if err != nil {
		return nil, err
	}

	return expenses, nil
}

// GetTotalCarExpenses handles the request to get car expenses
func (h *CarController) GetTotalCarExpenses(c *fiber.Ctx) error {
	// Convert carId to uint
	carIDStr := c.Params("id")
	carID, err := strconv.ParseUint(carIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid car ID",
		})
	}

	// Fetch total car expenses using the repository
	totalExpenses, err := h.repo.GetTotalCarExpenses(uint(carID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch total car expenses",
			"error":   err.Error(),
		})
	}

	// Return the total car expenses
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Car expenses fetched successfully",
		"data":    totalExpenses,
	})
}

// ==============
