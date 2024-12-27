package controllers

// import (
// 	"car-bond/internals/database"
// 	vehicleevaluation "car-bond/internals/models/vehicleEvaluation"
// 	"fmt"
// 	"log"
// 	"os"
// 	"path/filepath"
// 	"strconv"
// 	"strings"

// 	"github.com/gofiber/fiber/v2"
// 	"rsc.io/pdf"
// )

// // ProcessPDFAndUpload processes a PDF file and inserts its content into the database
// func ProcessPDFAndUpload(c *fiber.Ctx) error {
// 	db := database.DB.Db
// 	if db == nil {
// 		return fiber.NewError(fiber.StatusInternalServerError, "Database connection not initialized")
// 	}

// 	// Retrieve the uploaded file
// 	file, err := c.FormFile("pdf")
// 	if err != nil {
// 		return fiber.NewError(fiber.StatusBadRequest, "Failed to retrieve uploaded file: "+err.Error())
// 	}

// 	// Save the uploaded file to a temporary location
// 	tempDir := "./uploads"
// 	if err := os.MkdirAll(tempDir, os.ModePerm); err != nil {
// 		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create upload directory: "+err.Error())
// 	}

// 	tempFilePath := filepath.Join(tempDir, file.Filename)
// 	if err := c.SaveFile(file, tempFilePath); err != nil {
// 		return fiber.NewError(fiber.StatusInternalServerError, "Failed to save uploaded file: "+err.Error())
// 	}
// 	defer func() {
// 		if removeErr := os.Remove(tempFilePath); removeErr != nil {
// 			log.Printf("Warning: Failed to remove temporary file: %v", removeErr)
// 		}
// 	}()

// 	// Extract table data from the PDF
// 	data, extractErr := extractTableFromPDF(tempFilePath)
// 	if extractErr != nil {
// 		return fiber.NewError(fiber.StatusInternalServerError, "Failed to extract table data from PDF: "+extractErr.Error())
// 	}

// 	// Start a database transaction
// 	tx := db.Begin()

// 	// Optionally clear old data
// 	if err := tx.Exec("DELETE FROM vehicle_evaluations").Error; err != nil {
// 		tx.Rollback()
// 		return fiber.NewError(fiber.StatusInternalServerError, "Failed to clear old data: "+err.Error())
// 	}

// 	// Insert new data
// 	for _, record := range data {
// 		cif, parseErr := strconv.Atoi(strings.ReplaceAll(record[4], " ", ""))
// 		if parseErr != nil {
// 			log.Printf("Failed to parse CIF for record: %v, error: %v", record, parseErr)
// 			continue
// 		}

// 		vehicle := &vehicleevaluation.VehicleEvaluation{
// 			HSCCode:     record[0],
// 			COO:         record[1],
// 			Description: record[2],
// 			CC:          record[3],
// 			CIF:         cif,
// 			CreatedBy:   "admin", // Replace with actual user ID
// 			UpdatedBy:   "admin",
// 		}

// 		if err := tx.Create(vehicle).Error; err != nil {
// 			tx.Rollback()
// 			return fiber.NewError(fiber.StatusInternalServerError, "Failed to insert new data: "+err.Error())
// 		}
// 	}

// 	// Commit the transaction
// 	if err := tx.Commit().Error; err != nil {
// 		return fiber.NewError(fiber.StatusInternalServerError, "Failed to commit transaction: "+err.Error())
// 	}

// 	return c.JSON(fiber.Map{
// 		"message": "PDF processed and data uploaded successfully",
// 	})
// }

// // extractTableFromPDF extracts table data from the PDF file using a free library
// func extractTableFromPDF(pdfPath string) ([][]string, error) {
// 	// Open the PDF file
// 	file, err := os.Open(pdfPath)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to open PDF file: %v", err)
// 	}
// 	defer file.Close()

// 	// Load the PDF
// 	doc, err := pdf.Open(file.Name())
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to open PDF document: %v", err)
// 	}

// 	var records [][]string

// 	// Extract text from all pages
// 	for i := 0; i < doc.NumPage(); i++ {
// 		page := doc.Page(i + 1)
// 		if page == nil {
// 			continue
// 		}

// 		text := page.Text()
// 		lines := strings.Split(text, "\n")

// 		// Process each line as a potential row in the table
// 		for _, line := range lines {
// 			line = strings.TrimSpace(line)
// 			if line == "" {
// 				continue
// 			}
// 			records = append(records, strings.Split(line, ","))
// 		}
// 	}

// 	return records, nil
// }

// // UploadPDFHandler handles the PDF upload request
// func UploadPDFHandler(c *fiber.Ctx) error {
// 	if err := ProcessPDFAndUpload(c); err != nil {
// 		return err
// 	}

// 	return c.JSON(fiber.Map{
// 		"message": "PDF processed and data uploaded successfully",
// 	})
// }
