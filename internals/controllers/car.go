package controllers

import (
	"archive/zip"
	"car-bond/internals/models/alertRegistration"
	"car-bond/internals/models/carRegistration"
	"car-bond/internals/repository"
	"car-bond/internals/utils"
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

type CarController struct {
	repo repository.CarRepository
}

func NewCarController(repo repository.CarRepository) *CarController {
	return &CarController{repo: repo}
}

// ==================

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
		SeatCapacity:          c.FormValue("seat_capacity"),
		MaximCarry:            utils.StrToFloat(c.FormValue("maxim_carry")),
		Weight:                utils.StrToFloat(c.FormValue("weight")),
		GrossWeight:           utils.StrToFloat(c.FormValue("gross_weight")),
		Length:                utils.StrToFloat(c.FormValue("length")),
		Width:                 utils.StrToFloat(c.FormValue("width")),
		Height:                utils.StrToFloat(c.FormValue("height")),
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
		OtherEntity:          c.FormValue("other_entity"),
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

// Define the UpdateCarPayload struct
type UpdateCarPayload struct {
	ChasisNumber          string `form:"chasis_number"`
	EngineNumber          string `form:"engine_number"`
	FrameNumber           string `form:"frame_number"`
	EngineCapacity        string `form:"engine_capacity"`
	Make                  string `form:"make"`
	CarModel              string `form:"car_model"`
	SeatCapacity          string `form:"seat_capacity"`
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
	OtherEntity           string `json:"other_entity"`
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
		case "car_millage", "manufacture_year", "first_registration_year":
			updates[fieldName] = utils.StrToInt(fieldValue)
		case "maxim_carry", "weight", "gross_weight", "length", "width", "height", "bid_price", "vat_tax":
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

	if val, ok := updates["to_company_id"]; ok {
		if ptr, ok := val.(*uint); ok && ptr != nil && *ptr != 0 {
			// 1) mark status
			updates["car_status_japan"] = "Exported"

			// 2) fetch company name & add to updates
			if name, err := h.repo.GetCompanyNameByID(*ptr); err == nil {
				updates["other_entity"] = name // <-- new field
			} else if err != gorm.ErrRecordNotFound {
				// surface DB error (missing company simply skips the field)
				return c.Status(500).JSON(fiber.Map{
					"status":  "error",
					"message": "Failed to fetch destination company",
					"data":    err.Error(),
				})
			}
		}
	}

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
	id := c.Params("id")

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

	var payloadInv UpdateCarPayload3
	if err := c.BodyParser(&payloadInv); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input",
			"data":    err.Error(),
		})
	}

	targetInvoiceID := payloadInv.CarShippingInvoiceID

	if targetInvoiceID != 0 {
		// Check if the invoice exists and is not locked
		invoice, err := h.repo.GetInvoiceByID(targetInvoiceID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Status(404).JSON(fiber.Map{
					"status":  "error",
					"message": "Target invoice not found",
				})
			}
			return c.Status(500).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to retrieve invoice",
				"data":    err.Error(),
			})
		}

		if invoice.Locked {
			return c.Status(400).JSON(fiber.Map{
				"status":  "error",
				"message": "Invoice is locked and cannot be modified",
			})
		}

		// Ensure car is not already assigned to another invoice
		if car.CarShippingInvoiceID != nil && *car.CarShippingInvoiceID != targetInvoiceID {
			return c.Status(400).JSON(fiber.Map{
				"status":  "error",
				"message": "Car is already assigned to another invoice",
			})
		}
	}

	// Update car fields
	updateCar3Fields(&car, payloadInv)

	// Save the changes
	if err := h.repo.UpdateCar(&car); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update car",
			"data":    err.Error(),
		})
	}

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

func (h *CarController) UpdateCarExpense(c *fiber.Ctx) error {
	// Define a struct for input validation
	type UpdateCarExpenseInput struct {
		Description   string  `json:"description" validate:"required"`
		Currency      string  `json:"currency" validate:"required"`
		Amount        float64 `json:"amount" validate:"required,gt=0"`
		DollarRate    float64 `json:"dollar_rate"`
		ExpenseDate   string  `json:"expense_date" validate:"required"`
		CompanyName   string  `json:"company_name"`
		Destination   string  `json:"destination"`
		ExpenseVAT    float64 `json:"expense_vat"`
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
	expense.CompanyName = input.CompanyName
	expense.ExpenseVAT = input.ExpenseVAT
	expense.Destination = input.Destination
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

// ===================

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

func (h *CarController) UpdateCarWithDetails(c *fiber.Ctx) error {
	carIDStr := c.Params("id")
	carID := utils.StrToUint(carIDStr)

	var input CreateCarInput

	// Parse request body
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input",
			"data":    err.Error(),
		})
	}

	// Ensure the Car ID matches the route param
	input.Car.ID = carID

	// -------------------------------------------------------------------
	// OPTIONAL EXPORT + OTHER_ENTITY UPDATE
	// -------------------------------------------------------------------
	if input.Car.ToCompanyID != nil && *input.Car.ToCompanyID != 0 {
		// 1) Mark as Exported
		input.Car.CarStatusJapan = "Exported"

		// 2) Fetch and assign company name
		companyName, err := h.repo.GetCompanyNameByID(*input.Car.ToCompanyID)
		if err != nil && err != gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to fetch destination company",
				"data":    err.Error(),
			})
		}
		input.Car.OtherEntity = companyName
	}

	// Update car and its expenses
	if err := h.repo.UpdateCarWithExpenses(&input.Car, input.CarExpenses); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update car and expenses",
			"data":    err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Car and expenses updated successfully",
		"data":    input.Car,
	})
}
