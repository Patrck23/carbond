package controllers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"car-bond/internals/models/carRegistration"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Upload CarFile handles uploading either a photo or a PDF for a Car
func UploadCarFile(c *fiber.Ctx, db *gorm.DB) error {

	// Get the Car ID from the request form
	CarID := c.FormValue("car_id")
	if CarID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Car ID is required",
		})
	}

	file, err := c.FormFile("scan")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to parse uploaded file",
			"data":    fmt.Sprintf("Error: %s", err.Error()),
		})
	}

	// Ensure the  Car exists
	var car carRegistration.Car
	if err := db.First(&car, "id = ?", CarID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Car not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to verify  car",
			"data":    err.Error(),
		})
	}

	// Validate the file type (photo or PDF)
	allowedExtensions := []string{".jpg", ".jpeg", ".png", ".pdf"}
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !contains(allowedExtensions, ext) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid file type. Only JPG, JPEG, PNG, and PDF are allowed.",
		})
	}

	// Create a directory for storing files
	uploadDir := "./uploads/car_files"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create upload directory",
			"data":    err.Error(),
		})
	}

	// Save the file to the server
	filename := fmt.Sprintf("%d_%d%s", car.ID, time.Now().Unix(), ext)
	filePath := filepath.Join(uploadDir, filename)
	if err := c.SaveFile(file, filePath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to save file",
			"data":    err.Error(),
		})
	}

	// Save the file path to the database
	carFile := carRegistration.CarScan{
		CarID:  car.ID,
		Scan:   filePath,
		Title:  c.FormValue("title"),
		Remark: c.FormValue("remark"),
	}
	if err := db.Create(&carFile).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to save file information",
			"data":    err.Error(),
		})
	}

	// Return success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "File uploaded successfully",
		"data":    carFile,
	})
}

// contains checks if a slice contains a specific string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// UploadCarFiles handles uploading multiple files for a Car
func UploadCarFiles(c *fiber.Ctx, db *gorm.DB) error {

	// Get the Car ID from the request form
	CarID := c.FormValue("car_id")
	if CarID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Car ID is required",
		})
	}

	// Parse the uploaded files
	files, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to parse uploaded files",
			"data":    err.Error(),
		})
	}

	// Ensure the Car exists
	var car carRegistration.Car
	if err := db.First(&car, "id = ?", CarID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Car not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to verify car",
			"data":    err.Error(),
		})
	}

	// Validate the files
	allowedExtensions := []string{".jpg", ".jpeg", ".png", ".pdf"}
	var uploadedFiles []carRegistration.CarScan

	// Create a directory for storing files
	uploadDir := "./uploads/car_files"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create upload directory",
			"data":    err.Error(),
		})
	}

	// Iterate through the files in the request and save each one
	for _, file := range files.File["file"] {
		ext := strings.ToLower(filepath.Ext(file.Filename))
		if !contains(allowedExtensions, ext) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid file type. Only JPG, JPEG, PNG, and PDF are allowed.",
			})
		}

		// Save the file to the server
		filename := fmt.Sprintf("%d_%d%s", car.ID, time.Now().Unix(), ext)
		filePath := filepath.Join(uploadDir, filename)
		if err := c.SaveFile(file, filePath); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to save file",
				"data":    err.Error(),
			})
		}

		// Save the file path to the database
		carFile := carRegistration.CarScan{
			CarID:  car.ID,
			Scan:   filePath,
			Title:  c.FormValue("title"),
			Remark: c.FormValue("remark"),
		}
		if err := db.Create(&carFile).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to save file information",
				"data":    err.Error(),
			})
		}

		// Add to the list of uploaded files
		uploadedFiles = append(uploadedFiles, carFile)
	}

	// Return success response with uploaded files info
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Files uploaded successfully",
		"data":    uploadedFiles,
	})
}

// GetCarFiles retrieves all files (photos or PDFs) associated with a specific car
func GetCarFiles(c *fiber.Ctx, db *gorm.DB) error {

	// Get the  car ID from the route parameters
	CarID := c.Params("id")
	if CarID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Car ID is required",
		})
	}

	// Verify if the car exists
	var car carRegistration.Car
	if err := db.First(&car, "id = ?", CarID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Car not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve car information",
			"data":    err.Error(),
		})
	}

	// Query all files associated with the car
	var carFiles []carRegistration.CarScan
	if err := db.Where("car_id = ?", CarID).Find(&carFiles).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve Car files",
			"data":    err.Error(),
		})
	}

	// If no files found, return a response indicating so
	if len(carFiles) == 0 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"message": "No files found for this car",
			"data":    []string{},
		})
	}

	// Return the list of files
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Files retrieved successfully",
		"data":    carFiles,
	})
}

func GetFile(c *fiber.Ctx, db *gorm.DB) error {

	// Get the file ID from the route parameters
	fileID := c.Params("file_id")
	if fileID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "File ID is required",
		})
	}

	// Query the file details from the database
	var carFile carRegistration.CarScan
	if err := db.First(&carFile, "id = ?", fileID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "File not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve file details",
			"data":    err.Error(),
		})
	}

	// Get the absolute path of the file
	filePath := filepath.Clean(carFile.Scan)

	// Serve the file to the client
	// return c.SendFile(filePath, true) // `true` ensures the file is served as an attachment
	return c.SendFile(filePath, false) // Displays inline
}

// =================

// UpdateCarFiles handles deleting old files and uploading new ones for a Car
func UpdateCarFiles(c *fiber.Ctx, db *gorm.DB) error {

	// Get the Car ID from the request form
	CarID := c.FormValue("car_id")
	if CarID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Car ID is required",
		})
	}

	// Ensure the Car exists
	var car carRegistration.Car
	if err := db.First(&car, "id = ?", CarID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Car not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to verify car",
			"data":    err.Error(),
		})
	}

	// Fetch existing car files to delete old ones
	var existingFiles []carRegistration.CarScan
	if err := db.Where("car_id = ?", CarID).Find(&existingFiles).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve existing car files",
			"data":    err.Error(),
		})
	}

	// Delete old files from disk and from the database
	for _, file := range existingFiles {
		// Delete the file from the disk
		if err := os.Remove(file.Scan); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to delete old file",
				"data":    err.Error(),
			})
		}

		// Delete the record from the database
		if err := db.Delete(&file).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to delete old file from database",
				"data":    err.Error(),
			})
		}
	}

	// Parse the new files from the form
	files, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to parse uploaded files",
			"data":    err.Error(),
		})
	}

	// Validate the files
	allowedExtensions := []string{".jpg", ".jpeg", ".png", ".pdf"}
	var uploadedFiles []carRegistration.CarScan

	// Create a directory for storing files
	uploadDir := "./uploads/car_files"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create upload directory",
			"data":    err.Error(),
		})
	}

	// Iterate through the files in the request and save each one
	for _, file := range files.File["scan"] {
		ext := strings.ToLower(filepath.Ext(file.Filename))
		if !contains(allowedExtensions, ext) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid file type. Only JPG, JPEG, PNG, and PDF are allowed.",
			})
		}

		// Save the new file to the server
		filename := fmt.Sprintf("%d_%d%s", car.ID, time.Now().Unix(), ext)
		filePath := filepath.Join(uploadDir, filename)
		if err := c.SaveFile(file, filePath); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to save file",
				"data":    err.Error(),
			})
		}

		// Save the new file path to the database
		carFile := carRegistration.CarScan{
			CarID:  car.ID,
			Scan:   filePath,
			Title:  c.FormValue("title"),
			Remark: c.FormValue("remark"),
		}
		if err := db.Create(&carFile).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to save new file information",
				"data":    err.Error(),
			})
		}

		// Add to the list of uploaded files
		uploadedFiles = append(uploadedFiles, carFile)
	}

	// Return success response with uploaded files info
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Car files updated successfully",
		"data":    uploadedFiles,
	})
}
