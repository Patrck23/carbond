package controllers

import (
	"car-bond/internals/database"
	"car-bond/internals/models/carRegistration"

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

// Get All cars from db
func GetAllCars(c *fiber.Ctx) error {
	// Get the database instance
	db := database.DB.Db

	// Initialize a slice to hold the cars
	var cars []carRegistration.Car

	// Query the database for all cars
	if err := db.Find(&cars).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch cars",
			"data":    err.Error(),
		})
	}

	// Check if no cars are found
	if len(cars) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "No cars found",
		})
	}

	// Return the list of cars
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Cars retrieved successfully",
		"data":    cars,
	})
}

// GetSingleCar from db
func GetSingleCar(c *fiber.Ctx) error {
	// Get the database instance
	db := database.DB.Db

	// Get the car ID from the route parameters
	id := c.Params("id")

	// Initialize a variable to hold the car
	var car carRegistration.Car

	// Query the database for the car by ID
	if err := db.First(&car, "id = ?", id).Error; err != nil {
		// Check if the error is due to the record not being found
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Car not found",
			})
		}
		// Handle other database errors
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve car",
			"data":    err.Error(),
		})
	}

	// Return the found car
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Car found",
		"data":    car,
	})
}

// update a car in db
func UpdateCar(c *fiber.Ctx) error {
	// Define a struct for the update payload
	type updateCar struct {
		VinNumber     string  `json:"vin_number"`
		Make          string  `json:"make"`
		CarModel      string  `json:"model"`
		Year          int     `json:"year"`
		BidPrice      float64 `json:"bid_price"`
		VATTax        float64 `json:"vat_tax"`
		PurchaseDate  string  `json:"purchase_date"`
		Destination   string  `json:"destination"`
		FromCompanyID uint    `json:"from_company_id"`
		ToCompanyID   uint    `json:"to_company_id"`
		CustomerID    int     `json:"customer_id"`
		UpdatedBy     string  `json:"updated_by"`
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
	car.Year = updateCarData.Year
	car.BidPrice = updateCarData.BidPrice
	car.VATTax = updateCarData.VATTax
	car.PurchaseDate = updateCarData.PurchaseDate
	car.Destination = updateCarData.Destination
	car.FromCompanyID = updateCarData.FromCompanyID
	car.ToCompanyID = updateCarData.ToCompanyID
	car.CustomerID = updateCarData.CustomerID
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
