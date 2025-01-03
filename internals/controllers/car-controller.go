package controllers

import (
	"car-bond/internals/database"
	"car-bond/internals/models/carRegistration"
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Create a car
func CreateCar(c *fiber.Ctx) error {
	// Get the database instance
	db := database.DB.Db

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

	// Attempt to create the car record in the database
	if err := db.Create(&car).Error; err != nil {
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

func GetAllCars(c *fiber.Ctx) error {
	// Get the database instance
	db := database.DB.Db

	// Initialize a variable to hold all cars
	var cars []carRegistration.Car

	// Query all cars from the database
	err := db.Find(&cars).Error
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
		// Fetch car ports associated with the car
		var carPorts []carRegistration.CarPort
		err := db.Where("car_id = ?", car.ID).Find(&carPorts).Error
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to retrieve car ports for car ID " + strconv.Itoa(int(car.ID)),
				"data":    err.Error(),
			})
		}

		// Fetch expenses associated with the car
		var expenses []carRegistration.CarExpense
		err = db.Where("car_id = ?", car.ID).Find(&expenses).Error
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to retrieve expenses for car ID " + strconv.Itoa(int(car.ID)),
				"data":    err.Error(),
			})
		}

		// Combine car, ports, and expenses into a single response map
		response = append(response, fiber.Map{
			"car":       car,
			"car_ports": carPorts,
			"expenses":  expenses,
		})
	}

	// Return the response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Cars and associated data retrieved successfully",
		"data":    response,
	})
}

// GetSingleCar fetches a car with its associated ports and expenses from the database
func GetSingleCar(c *fiber.Ctx) error {
	// Get the database instance
	db := database.DB.Db

	// Get the car ID from the route parameters
	id := c.Params("id")

	// Initialize variables for car, ports, and expenses
	var car carRegistration.Car
	var carPorts []carRegistration.CarPort
	var expenses []carRegistration.CarExpense // Assuming you have an Expense model

	// Query the car by ID and preload associated ports and expenses
	err := db.First(&car, "id = ?", id).Error
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

	// Query the CarPorts associated with the car
	err = db.Where("car_id = ?", car.ID).Find(&carPorts).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve car ports",
			"data":    err.Error(),
		})
	}

	// Query the Expenses associated with the car (assuming Expense model has CarID)
	err = db.Where("car_id = ?", car.ID).Find(&expenses).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve expenses",
			"data":    err.Error(),
		})
	}

	// Combine the car, its ports, and expenses in a response
	response := fiber.Map{
		"car":       car,
		"car_ports": carPorts,
		"expenses":  expenses,
	}

	// Return the response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Car and associated data retrieved successfully",
		"data":    response,
	})
}

// GetSingleCar fetches a car with its associated ports and expenses using vin_number
func GetSingleCarByVinNumber(c *fiber.Ctx) error {
	// Get the database instance
	db := database.DB.Db

	// Get the vin_number from the route parameters
	vinNumber := c.Params("vinNumber")

	// Initialize variables for car, ports, and expenses
	var car carRegistration.Car
	var carPorts []carRegistration.CarPort
	var expenses []carRegistration.CarExpense // Assuming you have an Expense model

	// Query the car by vin_number
	err := db.First(&car, "vin_number = ?", vinNumber).Error
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

	// Query the CarPorts associated with the car
	err = db.Where("car_id = ?", car.ID).Find(&carPorts).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve car ports",
			"data":    err.Error(),
		})
	}

	// Query the Expenses associated with the car (assuming Expense model has CarID)
	err = db.Where("car_id = ?", car.ID).Find(&expenses).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve expenses",
			"data":    err.Error(),
		})
	}

	// Combine the car, its ports, and expenses in a response
	response := fiber.Map{
		"car":       car,
		"car_ports": carPorts,
		"expenses":  expenses,
	}

	// Return the response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Car and associated data retrieved successfully",
		"data":    response,
	})
}

// update a car in db
func UpdateCar(c *fiber.Ctx) error {
	// Define a struct for the update payload
	type updateCar struct {
		VinNumber       string  `json:"vin_number"`
		Make            string  `json:"make"`
		CarModel        string  `json:"model"`
		ManufactureYear int     `json:"maunufacture_year"`
		Currency        string  `json:"currency"`
		BidPrice        float64 `json:"bid_price"`
		VATTax          float64 `json:"vat_tax"`
		PurchaseDate    string  `json:"purchase_date"`
		Destination     string  `json:"destination"`
		FromCompanyID   uint    `json:"from_company_id"`
		ToCompanyID     uint    `json:"to_company_id"`
		CustomerID      int     `json:"customer_id"`
		UpdatedBy       string  `json:"updated_by"`
	}

	// Get the database instance
	db := database.DB.Db

	// Get the car ID from the route parameters
	id := c.Params("id")

	// Find the car in the database by ID
	var car carRegistration.Car
	if err := db.First(&car, "id = ?", id).Error; err != nil {
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

	// Parse the request body into the updateCar struct
	var updateCarData updateCar
	if err := c.BodyParser(&updateCarData); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input",
			"data":    err.Error(),
		})
	}

	// Update the car fields
	car.VinNumber = updateCarData.VinNumber
	car.Make = updateCarData.Make
	car.CarModel = updateCarData.CarModel
	car.ManufactureYear = updateCarData.ManufactureYear
	car.Currency = updateCarData.Currency
	car.BidPrice = updateCarData.BidPrice
	car.VATTax = updateCarData.VATTax
	car.PurchaseDate = updateCarData.PurchaseDate
	car.Destination = updateCarData.Destination
	// Assign foreign keys if provided
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

	if updateCarData.CustomerID != 0 {
		car.CustomerID = &updateCarData.CustomerID
	} else {
		car.CustomerID = nil
	}
	car.UpdatedBy = updateCarData.UpdatedBy

	// Save the changes to the database
	if err := db.Save(&car).Error; err != nil {
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

// delete car in db by ID
func DeleteCarByID(c *fiber.Ctx) error {
	// Get the database instance
	db := database.DB.Db

	// Get the car ID from the route parameters
	id := c.Params("id")

	// Find the car in the database by ID
	var car carRegistration.Car
	if err := db.First(&car, "id = ?", id).Error; err != nil {
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
	if err := db.Delete(&car, "id = ?", id).Error; err != nil {
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
	})
}

// Car Expense
// ===================================================================================================
// Create a car expense

func CreateCarExpense(c *fiber.Ctx) error {
	db := database.DB.Db
	carExpense := new(carRegistration.CarExpense)
	err := c.BodyParser(carExpense)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Something's wrong with your input", "data": err})
	}

	err = db.Create(&carExpense).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed to create car expense", "data": err})
	}

	return c.Status(201).JSON(fiber.Map{"status": "success", "message": "Car expense created successfully", "data": carExpense})
}

func GetAllCarExpenses(c *fiber.Ctx) error {
	// Initialize database instance
	db := database.DB.Db
	var expenses []carRegistration.CarExpense

	// Fetch all expenses without requiring a Car association
	if err := db.Preload("Car").Find(&expenses).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch expenses",
		})
	}

	return c.Status(fiber.StatusOK).JSON(expenses)
}

// Get Car Expenses by ID

func GetCarExpenseById(c *fiber.Ctx) error {
	// Initialize database instance
	db := database.DB.Db

	// Retrieve the expense ID and Car ID from the request parameters
	id := c.Params("id")
	carId := c.Params("carId")

	// Query the database for the expense by its ID and car ID
	var expense carRegistration.CarExpense
	result := db.Preload("Car").Where("id = ? AND car_id = ?", id, carId).First(&expense)

	// Handle potential database query errors
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{
				"status":  "error",
				"message": "Expense not found for the specified car",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch the expense",
			"error":   result.Error.Error(),
		})
	}

	// Return the fetched expense
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Expense fetched successfully",
		"data":    expense,
	})
}

// Get Car Expenses by Car ID

func GetCarExpensesByCarId(c *fiber.Ctx) error {
	// Initialize database instance
	db := database.DB.Db

	// Retrieve carId from the request parameters
	carId := c.Params("carId")

	// Query the database for car expenses
	var expenses []carRegistration.CarExpense
	result := db.Preload("Car").Where("car_id = ?", carId).Find(&expenses)

	// Handle potential database errors
	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve car expenses",
			"error":   result.Error.Error(),
		})
	}

	// Handle the case where no expenses are found
	if len(expenses) == 0 {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "No expenses found for the specified car",
		})
	}

	// Return a success response with the fetched data
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Car expenses fetched successfully",
		"data":    expenses,
	})
}

// Update car Expenses by

func UpdateCarExpense(c *fiber.Ctx) error {
	// Define a struct for input validation
	type UpdateCarExpenseInput struct {
		Description string  `json:"description"`
		Currency    string  `json:"currency"`
		Amount      float64 `json:"amount"`
		ExpenseDate string  `gorm:"type:date" json:"expense_date"`
		UpdatedBy   string  `json:"updated_by"`
	}

	db := database.DB.Db
	expenseID := c.Params("id")

	// Find the expense record by ID
	var expense carRegistration.CarExpense
	if err := db.First(&expense, expenseID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Expense not found"})
		}
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Database error", "error": err.Error()})
	}

	// Parse the request body into the input struct
	var input UpdateCarExpenseInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Invalid input", "error": err.Error()})
	}

	// Update the fields of the expense record
	expense.Description = input.Description
	expense.Currency = input.Currency
	expense.Amount = input.Amount
	expense.ExpenseDate = input.ExpenseDate
	expense.UpdatedBy = input.UpdatedBy

	// Save the changes to the database
	if err := db.Save(&expense).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed to update expense", "error": err.Error()})
	}

	// Return the updated expense
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Expense updated successfully",
		"data":    expense,
	})
}

// Delete car Expenses by ID

func DeleteCarExpenseById(c *fiber.Ctx) error {
	db := database.DB.Db

	// Parse carId and expenseId from the URL
	carId := c.Params("carId")
	expenseId := c.Params("id")

	// Check if the expense exists and belongs to the specified car
	var expense carRegistration.CarExpense
	if err := db.Where("id = ? AND car_id = ?", expenseId, carId).First(&expense).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Car expense not found or does not belong to the specified car"})
	}

	// Delete the expense
	if err := db.Delete(&expense).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed to delete car expense", "data": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Car expense deleted successfully"})
}

// Get Car Expenses by Car ID and Expense Date

func GetCarExpensesByCarIdAndExpenseDate(c *fiber.Ctx) error {
	db := database.DB.Db
	carId := c.Params("id")
	expenseDate := c.Params("expense_date")
	var expense carRegistration.CarExpense
	db.Where("car_id = ? AND expense_date = ?", carId, expenseDate).First(&expense)
	if expense.ID == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Car Expenses not found"})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Car Expenses fetched successfully", "data": expense})
}

// Get Car Expenses by Car ID and Expense Description

func GetCarExpensesByCarIdAndExpenseDescription(c *fiber.Ctx) error {
	db := database.DB.Db
	carId := c.Params("id")
	expenseDescription := c.Params("expense_description")
	var expense carRegistration.CarExpense
	db.Where("car_id = ? AND description = ?", carId, expenseDescription).First(&expense)
	if expense.ID == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Car Expenses not found"})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Car Expenses fetched successfully", "data": expense})
}

// Get Car Expenses by Car ID and Currency

func GetCarExpensesByCarIdAndCurrency(c *fiber.Ctx) error {
	db := database.DB.Db
	carId := c.Params("id")
	currency := c.Params("currency")
	var expenses []carRegistration.CarExpense
	db.Where("car_id = ? AND currency = ?", carId, currency).Find(&expenses)
	if len(expenses) == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Car Expenses not found"})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Car Expenses fetched successfully", "data": expenses})
}

// Get Car Expenses by Car ID, Expense Date, and Currency

func GetCarExpensesByThree(c *fiber.Ctx) error {
	db := database.DB.Db
	carId := c.Params("id")
	expenseDate := c.Params("expense_date")
	currency := c.Params("currency")
	var expense carRegistration.CarExpense
	db.Where("car_id = ? AND expense_date = ? AND currency = ?", carId, expenseDate, currency).First(&expense)
	if expense.ID == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Car Expenses not found"})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Car Expenses fetched successfully", "data": expense})
}

// Get Car Expenses by Car ID, Expense Description, and Currency

func GetCarExpensesByThreeDec(c *fiber.Ctx) error {
	db := database.DB.Db
	carId := c.Params("id")
	expenseDescription := c.Params("expense_description")
	currency := c.Params("currency")
	var expense carRegistration.CarExpense
	db.Where("car_id = ? AND description = ? AND currency = ?", carId, expenseDescription, currency).First(&expense)
	if expense.ID == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Car Expenses not found"})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Car Expenses fetched successfully", "data": expense})
}

// Get Car Expenses by Car ID, Expense Date, Expense Description, and Currency

func GetCarExpensesFilters(c *fiber.Ctx) error {
	db := database.DB.Db
	carId := c.Params("id")
	expenseDate := c.Params("expense_date")
	expenseDescription := c.Params("expense_description")
	currency := c.Params("currency")
	var expense carRegistration.CarExpense
	db.Where("car_id = ? AND expense_date = ? AND description = ? AND currency = ?", carId, expenseDate, expenseDescription, currency).First(&expense)
	if expense.ID == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Car Expenses not found"})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Car Expenses fetched successfully", "data": expense})
}

// Car Port
// ===================================================================================================

// Create a car port

func CreateCarPort(c *fiber.Ctx) error {
	db := database.DB.Db
	port := new(carRegistration.CarPort)
	err := c.BodyParser(port)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Something's wrong with your input", "data": err})
	}

	err = db.Create(&port).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed to create car port", "data": err})
	}

	return c.Status(201).JSON(fiber.Map{"status": "success", "message": "Car port created successfully", "data": port})
}

func GetAllCarPorts(c *fiber.Ctx) error {
	// Initialize database instance
	db := database.DB.Db
	carId := c.Params("carId")

	var ports []carRegistration.CarPort

	// Fetch all ports with associated car details
	if err := db.Preload("Car").Where("car_id = ?", carId).Find(&ports).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch car ports",
		})
	}

	return c.Status(fiber.StatusOK).JSON(ports)
}

// Get car port by ID

func GetCarPortById(c *fiber.Ctx) error {
	// Initialize database instance
	db := database.DB.Db

	// Retrieve the port ID and Car ID from the request parameters
	id := c.Params("id")
	carId := c.Params("carId")

	// Query the database for the port by its ID and Car ID
	var port carRegistration.CarPort
	result := db.Preload("Car").Where("id = ? AND car_id = ?", id, carId).First(&port)

	// Handle potential database query errors
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{
				"status":  "error",
				"message": "Port not found for the specified car",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch the port",
			"error":   result.Error.Error(),
		})
	}

	// Return the fetched port
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Port fetched successfully",
		"data":    port,
	})
}

// Update car port

func UpdateCarPort(c *fiber.Ctx) error {
	// Define a struct for input validation
	type UpdateCarPortInput struct {
		Name      string `json:"name"`
		Category  string `json:"category"`
		UpdatedBy string `json:"updated_by"`
	}

	db := database.DB.Db
	portID := c.Params("id")

	// Find the port record by ID
	var port carRegistration.CarPort
	if err := db.First(&port, portID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Port not found"})
		}
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Database error", "error": err.Error()})
	}

	// Parse the request body into the input struct
	var input UpdateCarPortInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Invalid input", "error": err.Error()})
	}

	// Update the fields of the Port record
	port.Name = input.Name
	port.Category = input.Category
	port.UpdatedBy = input.UpdatedBy

	// Save the changes to the database
	if err := db.Save(&port).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed to update port", "error": err.Error()})
	}

	// Return the updated port
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Port updated successfully",
		"data":    port,
	})
}

// Delete car port by ID

func DeleteCarPortById(c *fiber.Ctx) error {
	// Initialize database instance
	db := database.DB.Db

	// Parse parameters
	id := c.Params("id")
	carId := c.Params("carId")

	// Check if the car exists
	var car carRegistration.Car
	if err := db.First(&car, carId).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Car not found"})
	}

	// Check if the car port exists and belongs to the car
	var port carRegistration.CarPort
	if err := db.Where("id = ? AND car_id = ?", id, carId).First(&port).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Car port not found or does not belong to the specified car"})
	}

	// Delete the Car port
	if err := db.Delete(&port).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed to delete car port", "data": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Car port deleted successfully"})
}

func GetAllPorts(c *fiber.Ctx) error {
	// Initialize database instance
	db := database.DB.Db
	var locations []carRegistration.CarPort

	// Fetch all locations without requiring a Car association
	if err := db.Preload("Car").Find(&locations).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch locations",
		})
	}

	return c.Status(fiber.StatusOK).JSON(locations)
}
