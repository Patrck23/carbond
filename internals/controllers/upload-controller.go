package controllers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"car-bond/internals/models/customerRegistration"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Upload CustomerFile handles uploading either a photo or a PDF for a Customer
func UploadCustomerFile(c *fiber.Ctx, db *gorm.DB) error {

	// Get the Customer ID from the request form
	customerID := c.FormValue("customer_id")
	if customerID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Customer ID is required",
		})
	}

	// Parse the uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to parse uploaded file",
			"data":    err.Error(),
		})
	}

	// Ensure the  Customer exists
	var customer customerRegistration.Customer
	if err := db.First(&customer, "id = ?", customerID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Customer not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to verify  customer",
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
	uploadDir := "./uploads/customer_files"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create upload directory",
			"data":    err.Error(),
		})
	}

	// Save the file to the server
	filename := fmt.Sprintf("%d_%d%s", customer.ID, time.Now().Unix(), ext)
	filePath := filepath.Join(uploadDir, filename)
	if err := c.SaveFile(file, filePath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to save file",
			"data":    err.Error(),
		})
	}

	// Save the file path to the database
	customerFile := customerRegistration.CustomerScan{
		CustomerID: customer.ID,
		Scan:       filePath,
		Title:      c.FormValue("title"),
		Remark:     c.FormValue("remark"),
	}
	if err := db.Create(&customerFile).Error; err != nil {
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
		"data":    customerFile,
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

// GetCustomerFiles retrieves all files (photos or PDFs) associated with a specific customer
func GetCustomerFiles(c *fiber.Ctx, db *gorm.DB) error {

	// Get the  customer ID from the route parameters
	customerID := c.Params("id")
	if customerID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Customer ID is required",
		})
	}

	// Verify if the customer exists
	var customer customerRegistration.Customer
	if err := db.First(&customer, "id = ?", customerID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Customer not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve customer information",
			"data":    err.Error(),
		})
	}

	// Query all files associated with the customer
	var customerFiles []customerRegistration.CustomerScan
	if err := db.Where("customer_id = ?", customerID).Find(&customerFiles).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve customer files",
			"data":    err.Error(),
		})
	}

	// If no files found, return a response indicating so
	if len(customerFiles) == 0 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"message": "No files found for this customer",
			"data":    []string{},
		})
	}

	// Return the list of files
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Files retrieved successfully",
		"data":    customerFiles,
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
	var customerFile customerRegistration.CustomerScan
	if err := db.First(&customerFile, "id = ?", fileID).Error; err != nil {
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
	filePath := filepath.Clean(customerFile.Scan)

	// Serve the file to the client
	// return c.SendFile(filePath, true) // `true` ensures the file is served as an attachment
	return c.SendFile(filePath, false) // Displays inline
}
