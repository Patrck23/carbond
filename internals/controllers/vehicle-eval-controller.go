package controllers

// import (
// 	"car-bond/internals/database"
// 	"car-bond/internals/models/metaData"
// 	"fmt"
// 	"log"
// 	"os"
// 	"path/filepath"
// 	"strconv"
// 	"strings"

// 	"github.com/gofiber/fiber/v2"
// 	"github.com/xuri/excelize/v2"
// )

// // ProcessExcelAndUpload processes an Excel file and inserts its content into the database
// func ProcessExcelAndUpload(c *fiber.Ctx) error {
// 	db := database.DB.Db
// 	if db == nil {
// 		return fiber.NewError(fiber.StatusInternalServerError, "Database connection not initialized")
// 	}

// 	// Retrieve the uploaded file
// 	file, err := c.FormFile("excel")
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

// 	// Extract table data from the Excel file
// 	data, extractErr := extractTableFromExcel(tempFilePath)
// 	if extractErr != nil {
// 		return fiber.NewError(fiber.StatusInternalServerError, "Failed to extract table data from Excel: "+extractErr.Error())
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

// 		vehicle := &metaData.VehicleEvaluation{
// 			HSCCode:     record[0],
// 			COO:         record[1],
// 			Description: record[2],
// 			CC:          record[3],
// 			CIF:         cif,
// 			CreatedBy:   "system", // Replace with actual user ID
// 			UpdatedBy:   "system",
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
// 		"message": "Excel processed and data uploaded successfully",
// 	})
// }

// // extractTableFromExcel extracts table data from the Excel file
// func extractTableFromExcel(filePath string) ([][]string, error) {
// 	f, err := excelize.OpenFile(filePath)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to open Excel file: %w", err)
// 	}

// 	defer func() {
// 		if closeErr := f.Close(); closeErr != nil {
// 			log.Printf("Warning: Failed to close Excel file: %v", closeErr)
// 		}
// 	}()

// 	// Get the first sheet
// 	sheetName := f.GetSheetName(1)
// 	if sheetName == "" {
// 		return nil, fmt.Errorf("failed to find the first sheet in the Excel file")
// 	}

// 	// Read all rows
// 	rows, err := f.GetRows(sheetName)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to read rows from sheet %s: %w", sheetName, err)
// 	}

// 	var records [][]string
// 	for _, row := range rows {
// 		// Skip empty rows
// 		if len(row) < 5 { // Assuming we need at least 5 columns
// 			continue
// 		}

// 		// Append row as a record
// 		records = append(records, row)
// 	}

// 	if len(records) == 0 {
// 		return nil, fmt.Errorf("no valid data found in the Excel file")
// 	}

// 	return records, nil
// }

// // UploadExcelHandler handles the Excel upload request
// func UploadExcelHandler(c *fiber.Ctx) error {
// 	if err := ProcessExcelAndUpload(c); err != nil {
// 		return err
// 	}

// 	return c.JSON(fiber.Map{
// 		"message": "Excel processed and data uploaded successfully",
// 	})
// }

// =======================================================================

import (
	"car-bond/internals/database"
	"car-bond/internals/models/metaData"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/ledongthuc/pdf"
)

// UploadPDF handles PDF file upload, parsing, and saving entries to the database
func UploadPDF(c *fiber.Ctx) error {
	// Use the global database instance
	db := database.DB.Db

	// Parse PDF file from the request
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to upload file",
		})
	}

	// Create a temporary file to save the uploaded PDF
	tempFile, err := os.CreateTemp("", "uploaded-*.pdf")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create temporary file",
		})
	}
	defer os.Remove(tempFile.Name()) // Ensure the temporary file is deleted after use

	// Save the uploaded file to the temporary file
	if err := c.SaveFile(file, tempFile.Name()); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save file",
		})
	}

	// Open the temporary PDF file
	f, r, err := pdf.Open(tempFile.Name())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to parse PDF",
		})
	}
	defer f.Close()

	// Parse content
	var content string
	totalPage := r.NumPage()
	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
		page := r.Page(pageIndex)
		if !page.V.IsNull() {
			// Extract text from the page
			text, err := page.GetPlainText(nil)
			if err != nil {
				log.Printf("Failed to extract text from page %d: %v", pageIndex, err)
				continue
			}
			content += text
		}
	}

	// Process structured data
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		fields := strings.Fields(line)

		// Ensure there are at least 5 fields (S/N, HSC Code, COO, Description, CC, CIF)
		if len(fields) < 5 {
			continue
		}

		// Extract data
		hscCode := fields[1]
		coo := fields[2]
		description := strings.Join(fields[3:len(fields)-2], " ")
		cc := strings.ReplaceAll(fields[len(fields)-2], ",", "") // Strip commas

		// Parse CIF as a decimal
		cif, err := parseDecimal(fields[len(fields)-1])
		if err != nil {
			// Log invalid CIF entries and continue
			log.Printf("Skipping invalid CIF value: %s", fields[len(fields)-1])
			continue
		}

		// Create and save VehicleEvaluation entry
		evaluation := metaData.VehicleEvaluation{
			HSCCode:     hscCode,
			COO:         coo,
			Description: description,
			CC:          cc,
			CIF:         cif,
			CreatedBy:   "system",
			UpdatedBy:   "system",
		}

		// Save to the database
		if err := db.Create(&evaluation).Error; err != nil {
			log.Printf("Failed to save entry: %v", err)
			continue
		}
	}

	return c.JSON(fiber.Map{
		"message": "Data uploaded successfully",
	})
}

// parseDecimal converts a string to a decimal (float64)
func parseDecimal(input string) (float64, error) {
	// If the value is "nan", return 0.0
	if strings.ToLower(input) == "nan" || input == "" {
		return 0.0, nil
	}
	return strconv.ParseFloat(input, 64)
}
