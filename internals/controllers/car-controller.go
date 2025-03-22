package controllers

import (
	"archive/zip"
	"car-bond/internals/models/alertRegistration"
	"car-bond/internals/models/carRegistration"
	"car-bond/internals/utils"
	"database/sql"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"

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
	GetCarByVin(ChasisNumber string) (carRegistration.Car, error)
	UpdateCar(car *carRegistration.Car) error
	DeleteByID(id string) error
	SearchPaginatedCars(c *fiber.Ctx) (*utils.Pagination, []carRegistration.Car, error)
	UpdateCarJapan(id string, updates map[string]interface{}) error

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
	GetTotalCarExpenses(carID uint) (CarExpenseResponse, error)

	CreateCarPhotos(photos []carRegistration.CarPhoto) error
	CreateCarExpenses(expenses []carRegistration.CarExpense) error
	// Photos
	CreateCarPhoto(photo *carRegistration.CarPhoto) error
	DeleteCarPhotoByURL(photoURL string) error
	DeleteCarPhotoByID(photoID uint) error
	GetCarPhotosBycarID(carId uint) ([]carRegistration.CarPhoto, error)

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

type CarController struct {
	repo CarRepository
}

func NewCarController(repo CarRepository) *CarController {
	return &CarController{repo: repo}
}

// ==================

func (r *CarRepositoryImpl) CreateAlert(alert *alertRegistration.Transaction) error {
	return r.db.Create(alert).Error
}

// ConvertCarToTransaction converts car updates to a transaction record
func ConvertUpdateCarToTransaction(carID uint, updates map[string]interface{}) *alertRegistration.Transaction {
	// Retrieve from and to company IDs as pointers
	fromCompanyID, okFrom := updates["from_company_id"].(*uint)
	toCompanyID, okTo := updates["to_company_id"].(*uint)

	// If the value is nil or missing, assign zero
	if !okFrom || fromCompanyID == nil {
		fromCompanyID = new(uint) // Create a zero-initialized uint pointer
		*fromCompanyID = 0        // Assign 0
	}

	if !okTo || toCompanyID == nil {
		toCompanyID = new(uint) // Create a zero-initialized uint pointer
		*toCompanyID = 0        // Assign 0
	}

	// Initialize the transaction object
	transaction := &alertRegistration.Transaction{
		CarChasisNumber: updates["chasis_number"].(string),
		FromCompanyId:   *fromCompanyID, // Dereference the pointer
		ToCompanyId:     *toCompanyID,   // Dereference the pointer
		CreatedBy:       updates["updated_by"].(string),
		UpdatedBy:       updates["updated_by"].(string),
		ViewStatus:      false,
	}

	// Set the TransactionType based on the condition
	if *toCompanyID > 0 {
		transaction.TransactionType = "InTransit"
	} else {
		transaction.TransactionType = "Storage" // Set a default value if needed
	}

	return transaction
}

// ConvertCarToTransaction maps Car struct to Transaction struct
func ConvertCarToTransaction(car *carRegistration.Car) *alertRegistration.Transaction {
	return &alertRegistration.Transaction{
		CarChasisNumber: car.ChasisNumber,
		TransactionType: getTransactionType(car),
		FromCompanyId:   getFromCompanyID(car),
		ToCompanyId:     getToCompanyID(car),
		ViewStatus:      false,
		CreatedBy:       car.CreatedBy,
		UpdatedBy:       car.UpdatedBy,
	}
}

// Determine transaction type based on car status
func getTransactionType(car *carRegistration.Car) string { // InStock, Sold, Exported
	if car.CarStatusJapan == "Sold" {
		return "Sale"
	} else if car.CarStatusJapan == "InTransit" {
		return "Export"
	} else if car.CarStatusJapan == "InStock" {
		return "Storage"
	}
	return "Unknown"
}

// Get FromCompanyID (defaults to 0 if nil)
func getFromCompanyID(car *carRegistration.Car) uint {
	if car.FromCompanyID != nil {
		return *car.FromCompanyID
	}
	return 0
}

// Get ToCompanyID (defaults to 0 if nil)
func getToCompanyID(car *carRegistration.Car) uint {
	if car.ToCompanyID != nil {
		return *car.ToCompanyID
	}
	return 0
}

// ========================

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

func (h *CarController) CreateCar(c *fiber.Ctx) error {
	// Parse form data
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input",
			"data":    err.Error(),
		})
	}

	// Extract images
	files := form.File["images"]

	// Create directory if not exists
	uploadDir := "./uploads/car_files"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create upload directory",
			"data":    err.Error(),
		})
	}

	// Store image paths
	var carPhotos []carRegistration.CarPhoto

	for _, file := range files {
		// Sanitize file name (replace spaces with underscores)
		cleanFileName := strings.ReplaceAll(file.Filename, " ", "_")
		filePath := fmt.Sprintf("%s%s", uploadDir, cleanFileName)

		// Save file to disk
		if err := c.SaveFile(file, filePath); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to save image",
				"data":    err.Error(),
			})
		}

		// Append to carPhotos
		carPhotos = append(carPhotos, carRegistration.CarPhoto{URL: fmt.Sprintf("./uploads/car_files/%s", cleanFileName)})
	}

	// Parse other form fields into a Car instance
	car := &carRegistration.Car{
		ChasisNumber:          c.FormValue("chasis_number"),
		EngineNumber:          c.FormValue("engine_number"),
		FrameNumber:           c.FormValue("frame_number"),
		EngineCapacity:        c.FormValue("engine_capacity"),
		Make:                  c.FormValue("make"),
		CarModel:              c.FormValue("car_model"),
		MaximCarry:            utils.StrToInt(c.FormValue("maxim_carry")),
		Weight:                utils.StrToInt(c.FormValue("weight")),
		GrossWeight:           utils.StrToInt(c.FormValue("gross_weight")),
		Length:                utils.StrToInt(c.FormValue("length")),
		Width:                 utils.StrToInt(c.FormValue("width")),
		Height:                utils.StrToInt(c.FormValue("height")),
		ManufactureYear:       utils.StrToInt(c.FormValue("manufacture_year")),
		FirstRegistrationYear: utils.StrToInt(c.FormValue("first_registration_year")),
		Transmission:          c.FormValue("transmission"),
		BodyType:              c.FormValue("body_type"),
		Colour:                c.FormValue("colour"),
		Auction:               c.FormValue("auction"),
		Currency:              c.FormValue("currency"),
		CarMillage:            utils.StrToInt(c.FormValue("car_millage")),
		FuelConsumption:       c.FormValue("fuel_consumption"),
		BidPrice:              utils.StrToFloat(c.FormValue("bid_price")),
		VATTax:                utils.StrToFloat(c.FormValue("vat_tax")),
		// DollarRate:            utils.StrToFloat(c.FormValue("dollar_rate")),
		PurchaseDate: c.FormValue("purchase_date"),

		PowerSteering: utils.StrToBool(c.FormValue("power_steering")),
		PowerWindow:   utils.StrToBool(c.FormValue("power_window")),
		ABS:           utils.StrToBool(c.FormValue("abs")),
		ADS:           utils.StrToBool(c.FormValue("ads")),
		AirBrake:      utils.StrToBool(c.FormValue("air_brake")),
		OilBrake:      utils.StrToBool(c.FormValue("oil_brake")),
		AlloyWheel:    utils.StrToBool(c.FormValue("alloy_wheel")),
		SimpleWheel:   utils.StrToBool(c.FormValue("simple_wheel")),
		Navigation:    utils.StrToBool(c.FormValue("navigation")),
		AC:            utils.StrToBool(c.FormValue("ac")),

		// Convert IDs from string to uint
		FromCompanyID:        utils.StrToUintPointer(c.FormValue("from_company_id")),
		ToCompanyID:          utils.StrToUintPointer(c.FormValue("to_company_id")),
		CarShippingInvoiceID: utils.StrToUintPointer(c.FormValue("car_shipping_invoice_id")),
		BrokerName:           c.FormValue("broker_name"),
		BrokerNumber:         c.FormValue("broker_number"),
		NumberPlate:          c.FormValue("number_plate"),
		CarTracker:           utils.StrToBool(c.FormValue("car_tracker")),
		CustomerID:           utils.StrToIntPointer(c.FormValue("customer_id")),
		CarStatus:            c.FormValue("car_status"),
		CarPaymentStatus:     c.FormValue("car_payment_status"),
		CreatedBy:            c.FormValue("created_by"),
	}

	// Assign images to car
	car.CarPhotos = carPhotos

	// Create car in database
	if err := h.repo.CreateCar(car); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create car",
			"data":    err.Error(),
		})
	}

	// Convert Car to Transaction
	transaction := ConvertCarToTransaction(car)

	// Save to database
	if err := h.repo.CreateAlert(transaction); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create transaction"})
	}

	// Return success response
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Car created successfully",
		"data":    car,
	})
}

// ===================

func (r *CarRepositoryImpl) GetPaginatedCars(c *fiber.Ctx) (*utils.Pagination, []carRegistration.Car, error) {
	pagination, cars, err := utils.Paginate(c, r.db.Preload("CarPhotos"), carRegistration.Car{})
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

// ====================

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

func (r *CarRepositoryImpl) GetCarByVin(ChasisNumber string) (carRegistration.Car, error) {
	var car carRegistration.Car
	err := r.db.Preload("CarPhotos").First(&car, "chasis_number = ?", ChasisNumber).Error
	return car, err
}

func (h *CarController) GetSingleCarByChasisNumber(c *fiber.Ctx) error {
	// Get the car ID from the route parameters
	ChasisNumber := c.Params("ChasisNumber")

	// Fetch the car by ID
	car, err := h.repo.GetCarByVin(ChasisNumber)
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

func (r *CarRepositoryImpl) UpdateCarJapan(id string, updates map[string]interface{}) error {
	return r.db.Model(&carRegistration.Car{}).Where("id = ?", id).Updates(updates).Error
}

// Define the UpdateCarPayload struct
type UpdateCarPayload struct {
	ChasisNumber          string `form:"chasis_number"`
	EngineNumber          string `form:"engine_number"`
	FrameNumber           string `form:"frame_number"`
	EngineCapacity        string `form:"engine_capacity"`
	Make                  string `form:"make"`
	CarModel              string `form:"car_model"`
	MaximCarry            string `form:"maxim_carry"`
	Weight                string `form:"weight"`
	GrossWeight           string `form:"gross_weight"`
	Length                string `form:"length"`
	Width                 string `form:"width"`
	Height                string `form:"height"`
	CarMillage            string `form:"car_millage"`
	FuelConsumption       string `form:"fuel_consumption"`
	ManufactureYear       string `form:"manufacture_year"`
	FirstRegistrationYear string `form:"first_registration_year"`
	Transmission          string `form:"transmission"`
	BodyType              string `form:"body_type"`
	Colour                string `form:"colour"`
	Auction               string `form:"auction"`
	Currency              string `form:"currency"`
	PowerSteering         string `form:"power_steering"`
	PowerWindow           string `form:"power_window"`
	ABS                   string `form:"abs"`
	ADS                   string `form:"ads"`
	AirBrake              string `form:"air_brake"`
	OilBrake              string `form:"oil_brake"`
	AlloyWheel            string `form:"alloy_wheel"`
	SimpleWheel           string `form:"simple_wheel"`
	Navigation            string `form:"navigation"`
	AC                    string `form:"ac"`
	BidPrice              string `form:"bid_price"`
	VATTax                string `form:"vat_tax"`
	PurchaseDate          string `form:"purchase_date"`
	FromCompanyID         string `form:"from_company_id"`
	ToCompanyID           string `form:"to_company_id"`
	Destination           string `form:"destination"`
	CarShippingInvoiceID  string `form:"car_shipping_invoice_id"`
	Port                  string `form:"port"`
	UpdatedBy             string `form:"updated_by"`
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

	// Check if the request body is empty
	if (UpdateCarPayload{} == payload) {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Empty request body",
		})
	}

	form, err := c.MultipartForm()

	// Convert struct fields into a map dynamically
	updates := make(map[string]interface{})
	val := reflect.ValueOf(payload)
	typ := reflect.TypeOf(payload)

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldName := typ.Field(i).Tag.Get("form")
		fieldValue := field.String()

		// Convert field values based on expected type
		switch fieldName {
		case "maxim_carry", "weight", "gross_weight", "length", "width", "height", "car_millage", "manufacture_year", "first_registration_year":
			updates[fieldName] = utils.StrToInt(fieldValue)
		case "bid_price", "vat_tax":
			updates[fieldName] = utils.StrToFloat(fieldValue)
		case "power_steering", "power_window", "abs", "ads", "air_brake", "oil_brake", "alloy_wheel", "simple_wheel", "navigation", "ac":
			updates[fieldName] = utils.StrToBool(fieldValue)
		case "from_company_id", "to_company_id", "car_shipping_invoice_id":
			updates[fieldName] = utils.StrToUintPointer(fieldValue)
		default:
			updates[fieldName] = fieldValue
		}
	}

	if err == nil {
		// Extract images
		files := form.File["images"]
		uploadDir := "./uploads/car_files/"

		// Ensure the directory exists
		fmt.Println("Ensuring upload directory exists:", uploadDir)
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			fmt.Println("Error creating upload directory:", err)
			return c.Status(500).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to create upload directory",
				"data":    err.Error(),
			})
		}

		var oldPhotoMap = make(map[string]carRegistration.CarPhoto)
		fmt.Println("Mapping existing car photos...")

		// Store existing photos in a map
		for _, oldPhoto := range car.CarPhotos {
			oldPhotoMap[oldPhoto.URL] = oldPhoto
			fmt.Println("Existing photo - ID:", oldPhoto.ID, "URL:", oldPhoto.URL)
		}

		// List to store updated photos
		var updatedCarPhotos []carRegistration.CarPhoto

		fmt.Println("Processing uploaded images...")
		for _, file := range files {
			// Sanitize filename and construct paths
			cleanFileName := strings.ReplaceAll(file.Filename, " ", "_")
			filePath := fmt.Sprintf("%s%s", uploadDir, cleanFileName)       // Ensure the "/" separator
			fileURL := fmt.Sprintf("./uploads/car_files/%s", cleanFileName) // URL for frontend

			fmt.Println("Processing file:", file.Filename, "->", filePath)

			// Check if this image already exists
			if _, exists := oldPhotoMap[fileURL]; exists {
				fmt.Println("Image already exists, keeping:", fileURL)
				// updatedCarPhotos = append(updatedCarPhotos, existingPhoto)
				delete(oldPhotoMap, fileURL) // Remove from map to track missing ones
			} else {
				// Save new image to disk
				fmt.Println("Saving new image:", filePath)
				if err := c.SaveFile(file, filePath); err != nil {
					fmt.Println("Error saving image:", err)
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
						"status":  "error",
						"message": "Failed to save image",
						"data":    err.Error(),
					})
				}

				// Append new image to updated list
				updatedCarPhotos = append(updatedCarPhotos, carRegistration.CarPhoto{
					URL:   fileURL, // Ensure frontend can retrieve image
					CarID: utils.StrToUint(id),
				})
			}

			// Print all updated car photos so far
			fmt.Println("Current updatedCarPhotos list:")
			for _, photo := range updatedCarPhotos {
				fmt.Println("Photo URL:", photo.URL, "Car ID:", photo.CarID)
			}
		}

		// Remove old images that were not retained
		fmt.Println("Removing old images that are no longer needed...")
		for oldPath, oldPhoto := range oldPhotoMap {
			// Convert stored URL to file path
			oldFileName := strings.TrimPrefix(oldPath, "./uploads/car_files/")
			oldFilePath := fmt.Sprintf("%s%s", uploadDir, oldFileName)

			fmt.Println("Photo URL:", oldPath, "Photo ID:", oldPhoto.ID)

			// Delete old photo record from the database
			if err := h.repo.DeleteCarPhotoByID(oldPhoto.ID); err != nil {
				fmt.Println("Error deleting car photo record:", err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status":  "error",
					"message": fmt.Sprintf("Failed to delete car photo record with ID: %d", oldPhoto.ID),
					"data":    err.Error(),
				})
			}

			// Delete file from disk if it exists
			if _, err := os.Stat(oldFilePath); err == nil {
				fmt.Println("Deleting old image:", oldFilePath)
				if err := os.Remove(oldFilePath); err != nil {
					fmt.Println("Error deleting old image:", err)
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
						"status":  "error",
						"message": "Failed to delete old image",
						"data":    err.Error(),
					})
				}
			} else {
				fmt.Println("Old image not found, skipping:", oldFilePath)
			}
		}

		// Insert new car photos into the database
		if len(updatedCarPhotos) > 0 {
			fmt.Println("Inserting new car photos into database...")
			for _, photo := range updatedCarPhotos {
				if err := h.repo.CreateCarPhoto(&photo); err != nil {
					fmt.Println("Error inserting new car photo:", photo.URL, err)
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
						"status":  "error",
						"message": fmt.Sprintf("Failed to insert car photo: %s", photo.URL),
						"data":    err.Error(),
					})
				}
			}
		}

		// Log final updates
		fmt.Println("Final updatedCarPhotos list:", updatedCarPhotos)
		fmt.Println("Car photos successfully updated.")
	}

	// Convert Car to Transaction
	transaction := ConvertUpdateCarToTransaction(car.ID, updates)

	// Save to database
	if err := h.repo.CreateAlert(transaction); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create transaction"})
	}

	// Proceed with updating other car data
	if err := h.repo.UpdateCarJapan(id, updates); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update car",
			"data":    err.Error(),
		})
	}

	// Return success response
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Car updated successfully",
		"data":    updates,
	})
}

// ===================

type UpdateCarPayload2 struct {
	BrokerName       string `json:"broker_name"`
	BrokerNumber     string `json:"broker_number"`
	NumberPlate      string `json:"number_plate"`
	CarTracker       bool   `json:"car_tracker"`
	CarStatus        string `json:"car_status"`
	CarPaymentStatus string `json:"car_payment_status"`
	CustomerID       int    `json:"customer_id"`
	UpdatedBy        string `json:"updated_by"`
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

type UpdateCarPayload3 struct {
	CarShippingInvoiceID uint   `json:"car_shipping_invoice_id"`
	UpdatedBy            string `json:"updated_by"`
}

// UpdateCar handler function
func (h *CarController) UpdateCar3(c *fiber.Ctx) error {
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
	var payloadInv UpdateCarPayload3
	if err := c.BodyParser(&payloadInv); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input",
			"data":    err.Error(),
		})
	}

	// Update the car fields using the payload
	updateCar3Fields(&car, payloadInv) // Pass the parsed payload

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
func updateCar3Fields(car *carRegistration.Car, updateCarData UpdateCarPayload3) {
	// Update the car field
	if updateCarData.CarShippingInvoiceID != 0 {
		car.CarShippingInvoiceID = &updateCarData.CarShippingInvoiceID
	} else {
		car.CarShippingInvoiceID = nil
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

// CreateCarExpenses handles the creation of multiple car expenses
func (h *CarController) CreateCarExpenses(c *fiber.Ctx) error {
	// Parse the request body into a slice of CarExpense structs
	var carExpenses []carRegistration.CarExpense
	if err := c.BodyParser(&carExpenses); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input data",
			"data":    err.Error(),
		})
	}

	// Validate if there are any expenses in the request
	if len(carExpenses) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "No car expenses provided",
		})
	}

	// Insert each car expense into the database
	for _, expense := range carExpenses {
		if err := h.repo.CreateCarExpense(&expense); err != nil {
			return c.Status(500).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to create car expense",
				"data":    err.Error(),
			})
		}
	}

	// Return success response
	return c.Status(201).JSON(fiber.Map{
		"status":  "success",
		"message": "Car expenses created successfully",
		"data":    carExpenses,
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
		Description   string  `json:"description" validate:"required"`
		Currency      string  `json:"currency" validate:"required"`
		Amount        float64 `json:"amount" validate:"required,gt=0"`
		DollarRate    float64 `json:"dollar_rate"`
		ExpenseDate   string  `json:"expense_date" validate:"required"`
		CarrierName   string  `json:"carrier_name"` // if description == "Carrier car fee(RISKO)"
		ExpenseRemark string  `json:"expense_remark"`
		UpdatedBy     string  `json:"updated_by" validate:"required"`
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
	expense.DollarRate = input.DollarRate
	expense.ExpenseDate = input.ExpenseDate
	expense.CarrierName = input.CarrierName
	expense.ExpenseRemark = input.ExpenseRemark
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

// TotalCarExpense represents an individual expense for a car
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
						WHEN currency = 'JPY' THEN amount
						ELSE amount / dollar_rate
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

// GetTotalCarExpenses handles the request to get car expenses
func (h *CarController) GetTotalCarExpenses(c *fiber.Ctx) error {
	// Convert carId to uint
	carIDStr := c.Params("id")
	carID := utils.StrToUint(carIDStr)

	// Fetch total car expenses using the repository
	totalExpenses, err := h.repo.GetTotalCarExpenses(carID)
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

// ========================================================================================

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
	from_company_id := c.Query("from_company_id")

	// Start building the query
	query := r.db.Model(&carRegistration.Car{})

	// Apply filters based on provided parameters
	if from_company_id != "" {
		if _, err := strconv.Atoi(from_company_id); err == nil {
			query = query.Where("from_company_id = ?", from_company_id)
		}
	}

	if to_company_id != "" {
		if _, err := strconv.Atoi(to_company_id); err == nil {
			query = query.Where("to_company_id = ?", to_company_id)
		}
	}

	if chasis_number != "" {
		query = query.Where("chasis_number LIKE ?", "%"+chasis_number+"%")
	}
	if make != "" {
		query = query.Where("make LIKE ?", "%"+make+"%")
	}
	if model != "" {
		query = query.Where("car_model = ?", model)
	}
	if colour != "" {
		query = query.Where("colour = ?", colour)
	}
	if bodyType != "" {
		query = query.Where("body_type LIKE ?", "%"+bodyType+"%")
	}
	if auction != "" {
		query = query.Where("auction LIKE ?", "%"+auction+"%")
	}
	if destination != "" {
		query = query.Where("destination = ?", destination)
	}
	if port != "" {
		query = query.Where("port = ?", port)
	}
	if broker_name != "" {
		query = query.Where("broker_name LIKE ?", "%"+broker_name+"%")
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

func (h *CarController) SearchCars(c *fiber.Ctx) error {
	// Call the repository function to get paginated search results
	pagination, cars, err := h.repo.SearchPaginatedCars(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve cars",
			"data":    err.Error(),
		})
	}

	// Return the response with pagination details
	return c.Status(200).JSON(fiber.Map{
		"status":     "success",
		"message":    "Cars retrieved successfully",
		"pagination": pagination,
		"data":       cars,
	})
}

// ===================

func (r *CarRepositoryImpl) GetCarPhotosBycarID(carId uint) ([]carRegistration.CarPhoto, error) {
	var photos []carRegistration.CarPhoto
	if err := r.db.Where("car_id = ?", carId).Find(&photos).Error; err != nil {
		return nil, err
	}
	return photos, nil
}

func (h *CarController) FetchCarUploads(c *fiber.Ctx) error {
	// Get the car ID from the route parameters
	id := c.Query("id")

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

	// Retrieve photos of the car
	photos, err := h.repo.GetCarPhotosBycarID(car.ID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve photos",
			"data":    err.Error(),
		})
	}

	// Check if photos exist
	if len(photos) == 0 {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "No photos found for this car",
		})
	}

	// Collect file paths
	var filePaths []string
	for _, photo := range photos {
		// Convert stored URL to file path
		filePath := strings.TrimPrefix(photo.URL, "./")
		fmt.Println("Photo URL:", filePath, "Photo ID:", photo.ID)
		filePaths = append(filePaths, filePath) // Assuming FilePath stores the image location
	}

	// Create a ZIP file
	zipFileName := fmt.Sprintf("car_%d_photos.zip", car.ID)
	zipFilePath := filepath.Join(os.TempDir(), zipFileName)
	zipFile, err := os.Create(zipFilePath)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create ZIP file",
			"data":    err.Error(),
		})
	}
	defer zipFile.Close()

	// Write images to ZIP
	zipWriter := zip.NewWriter(zipFile)
	for _, filePath := range filePaths {
		err := addFileToZip(zipWriter, filePath)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to add file to ZIP",
				"data":    err.Error(),
			})
		}
	}
	zipWriter.Close()

	// Serve the ZIP file
	return c.SendFile(zipFilePath)
}

// Helper function to add a file to the ZIP archive
func addFileToZip(zipWriter *zip.Writer, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a zip entry
	zipFileWriter, err := zipWriter.Create(filepath.Base(filePath))
	if err != nil {
		return err
	}

	// Copy the file content into the zip entry
	_, err = io.Copy(zipFileWriter, file)
	return err
}

func (h *CarController) FetchCarUploads64(c *fiber.Ctx) error {
	id := c.Query("id")

	car, err := h.repo.GetCarByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Car not found"})
		}
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed to retrieve car", "data": err.Error()})
	}

	photos, err := h.repo.GetCarPhotosBycarID(car.ID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed to retrieve photos", "data": err.Error()})
	}

	if len(photos) == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No photos found for this car"})
	}

	var images []map[string]string
	for _, photo := range photos {
		// Read image file
		imageData, err := os.ReadFile(photo.URL)
		if err != nil {
			continue
		}

		// Convert image to base64
		base64Image := base64.StdEncoding.EncodeToString(imageData)
		images = append(images, map[string]string{
			"filename": photo.URL,
			"data":     "data:image/jpeg;base64," + base64Image,
		})
	}

	// Return JSON with images
	return c.JSON(fiber.Map{"status": "success", "message": "Car images retrieved", "images": images})
}

// ==========================================

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
		Select("currency, dollar_rate, SUM(amount) as total").
		Group("currency, dollar_rate").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	totalExpenses := map[string]float64{
		"USD": 0,
		"JPY": 0,
		// "TotalUSD": 0,
	}

	for _, res := range results {
		if res.Currency == "JPY" {
			totalExpenses["JPY"] += res.Total
			// totalExpenses["TotalUSD"] += res.Total * res.DollarRate
		} else {
			totalExpenses["USD"] += res.Total * res.DollarRate
			// totalExpenses["TotalUSD"] += res.Total * res.DollarRate
		}
	}

	return totalExpenses, nil
}

// GetDashboardData returns aggregated statistics for the dashboard
func (h *CarController) GetDashboardData(c *fiber.Ctx) error {
	totalCars, err := h.repo.GetTotalCars()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve total cars",
			"data":    err.Error(),
		})
	}

	disbandedCars, err := h.repo.GetDisbandedCars()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve disbanded cars",
			"data":    err.Error(),
		})
	}

	carsInStock, err := h.repo.GetCarsInStock()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve cars in stock",
			"data":    err.Error(),
		})
	}

	totalMoneySpent, err := h.repo.GetTotalMoneySpent()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve total money spent",
			"data":    err.Error(),
		})
	}

	totalCarExpenses, err := h.repo.GetTotalCarsExpenses()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve total car expenses",
			"data":    err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Dashboard data retrieved successfully",
		"data": fiber.Map{
			"total_cars":         totalCars,
			"disbanded_cars":     disbandedCars,
			"cars_in_stock":      carsInStock,
			"total_money_spent":  totalMoneySpent,
			"total_car_expenses": totalCarExpenses,
		},
	})
}

// ====================================================

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
		Where("to_company_id = ? AND car_status = ?", companyID, "Instock").
		Count(&count).Error
	return count, err
}

func (r *CarRepositoryImpl) GetComCarsSold(companyID uint) (int64, error) {
	var count int64
	err := r.db.Model(&carRegistration.Car{}).
		Where("to_company_id = ? AND car_status = ?", companyID, "Sold").
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
		Select("car_expenses.currency, car_expenses.dollar_rate, SUM(car_expenses.amount) as total, cars.to_company_id").
		Group("car_expenses.currency, car_expenses.dollar_rate, cars.to_company_id").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	totalExpenses := map[string]float64{
		"USD": 0,
		"JPY": 0,
	}

	for _, res := range results {
		if res.Currency == "JPY" {
			totalExpenses["JPY"] += res.Total
		} else {
			totalExpenses["USD"] += res.Total * res.DollarRate
		}
	}

	return totalExpenses, nil
}

// GetDashboardData returns aggregated statistics for the dashboard
func (h *CarController) GetCompanyDashboardData(c *fiber.Ctx) error {
	companyIdStr := c.Params("companyId")
	companyId := utils.StrToUint(companyIdStr)
	if companyId == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid company ID")
	}

	totalCars, err := h.repo.GetComTotalCars(companyId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve total cars",
			"data":    err.Error(),
		})
	}

	carsSold, err := h.repo.GetComCarsSold(companyId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve cars sold",
			"data":    err.Error(),
		})
	}

	carsInStock, err := h.repo.GetComCarsInStock(companyId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve cars in stock",
			"data":    err.Error(),
		})
	}

	totalMoneySpent, err := h.repo.GetComTotalMoneySpent(companyId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve total money spent",
			"data":    err.Error(),
		})
	}

	totalCarExpenses, err := h.repo.GetComTotalCarsExpenses(companyId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve total car expenses",
			"data":    err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Dashboard data retrieved successfully",
		"data": fiber.Map{
			"total_cars":         totalCars,
			"cars_in_stock":      carsInStock,
			"cars_sold":          carsSold,
			"total_money_spent":  totalMoneySpent,
			"total_car_expenses": totalCarExpenses,
		},
	})
}

// ===================================

type CreateCarInput struct {
	Car carRegistration.Car `json:"car"`
	// CarPhotos   []carRegistration.CarPhoto   `json:"car_photos"`
	CarExpenses []carRegistration.CarExpense `json:"car_expenses"`
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

// CreateCarWithDetails handles the creation of a car with photos and expenses
func (h *CarController) CreateCarWithDetails(c *fiber.Ctx) error {
	// Parse the request body into a CreateCarInput struct
	input := new(CreateCarInput)
	if err := c.BodyParser(input); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input data",
			"data":    err.Error(),
		})
	}

	// Create the car in the database
	if err := h.repo.CreateCar(&input.Car); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create car",
			"data":    err.Error(),
		})
	}

	// // Associate photos with the car
	// for i := range input.CarPhotos {
	// 	input.CarPhotos[i].CarID = input.Car.ID
	// }
	// if err := h.repo.CreateCarPhotos(input.CarPhotos); err != nil {
	// 	return c.Status(500).JSON(fiber.Map{
	// 		"status":  "error",
	// 		"message": "Failed to add car photos",
	// 		"data":    err.Error(),
	// 	})
	// }

	// Associate expenses with the car
	for i := range input.CarExpenses {
		input.CarExpenses[i].CarID = input.Car.ID
	}
	if err := h.repo.CreateCarExpenses(input.CarExpenses); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to add car expenses",
			"data":    err.Error(),
		})
	}

	// Return success response
	return c.Status(201).JSON(fiber.Map{
		"status":  "success",
		"message": "Car with details created successfully",
		"data":    input.Car,
	})
}

// ===========================================================
