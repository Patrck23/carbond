package controllers

import (
	"car-bond/internals/database"
	"car-bond/internals/models/metaData"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/xuri/excelize/v2"
)

// Vehicle Evaluation

// ProcessExcelAndUpload processes an Excel file and inserts its content into the database
func ProcessExcelAndUpload(c *fiber.Ctx) error {
	db := database.DB.Db
	if db == nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Database connection not initialized")
	}

	// Retrieve the uploaded file
	file, err := c.FormFile("excel")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Failed to retrieve uploaded file: "+err.Error())
	}

	// Save the uploaded file to a temporary location
	tempDir := "./uploads"
	if err := os.MkdirAll(tempDir, os.ModePerm); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create upload directory: "+err.Error())
	}

	tempFilePath := filepath.Join(tempDir, file.Filename)
	if err := c.SaveFile(file, tempFilePath); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to save uploaded file: "+err.Error())
	}
	defer func() {
		if removeErr := os.Remove(tempFilePath); removeErr != nil {
			log.Printf("Warning: Failed to remove temporary file: %v", removeErr)
		}
	}()

	// Extract table data from the Excel file
	data, extractErr := extractTableFromExcel(tempFilePath)
	if extractErr != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to extract table data from Excel: "+extractErr.Error())
	}

	// Start a database transaction
	tx := db.Begin()

	// Clear old data
	if err := tx.Exec("DELETE FROM vehicle_evaluations").Error; err != nil {
		tx.Rollback()
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to clear old data: "+err.Error())
	}

	// Insert new data
	for _, record := range data {

		// Parse CIF as a float64
		cif, parseErr := strconv.ParseFloat(strings.ReplaceAll(record[4], " ", ""), 64)
		if parseErr != nil {
			log.Printf("Failed to parse CIF for record: %v, error: %v", record, parseErr)
			continue
		}

		vehicle := &metaData.VehicleEvaluation{
			HSCCode:     record[0],
			COO:         record[1],
			Description: record[2],
			CC:          record[3],
			CIF:         cif,
			CreatedBy:   "system", // Replace with actual user ID
			UpdatedBy:   "system",
		}

		if err := tx.Create(vehicle).Error; err != nil {
			tx.Rollback()
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to insert new data: "+err.Error())
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to commit transaction: "+err.Error())
	}

	return c.JSON(fiber.Map{
		"message": "Excel processed and data uploaded successfully",
	})
}

func extractTableFromExcel(filePath string) ([][]string, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open Excel file: %w", err)
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			log.Printf("Warning: Failed to close Excel file: %v", closeErr)
		}
	}()

	// Fetch the list of all sheet names
	sheetList := f.GetSheetList()
	if len(sheetList) == 0 {
		return nil, fmt.Errorf("no sheets found in the Excel file")
	}

	// Use the first sheet as default
	sheetName := sheetList[0]
	log.Printf("Using sheet: %s", sheetName)

	// Read all rows from the sheet
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("failed to read rows from sheet %s: %w", sheetName, err)
	}

	var records [][]string
	for _, row := range rows {
		// Skip rows with fewer columns than the minimum required
		if len(row) < 5 {
			continue
		}

		// Trim trailing empty cells (optional)
		for len(row) > 0 && row[len(row)-1] == "" {
			row = row[:len(row)-1]
		}

		records = append(records, row)
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("no valid data found in sheet %s", sheetName)
	}

	return records, nil
}

// UploadExcelHandler handles the Excel upload request
func UploadExcelHandler(c *fiber.Ctx) error {
	if err := ProcessExcelAndUpload(c); err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message": "Excel processed and data uploaded successfully",
	})
}

// FetchVehicleEvaluationsByDescription fetches items from VehicleEvaluation that match a given description
func FetchVehicleEvaluationsByDescription(c *fiber.Ctx) error {

	// Get the database instance
	db := database.DB.Db

	// Extract the description from the query parameters
	description := c.Query("description")
	if description == "" {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Description query parameter is required",
		})
	}

	var evaluations []metaData.VehicleEvaluation

	// Fetch records that match the description
	if err := db.Where("description LIKE ?", "%"+description+"%").Find(&evaluations).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Error fetching records",
			"error":   err.Error(),
		})
	}

	// Return the results
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Records fetched successfully",
		"data":    evaluations,
	})
}

// =======================================================================

// Meta Units

func GetAllWeightUnits(c *fiber.Ctx) error {
	db := database.DB.Db
	var units []metaData.WeightUnit
	// find all users in the database
	db.Find(&units)
	// If no customer found, return an error
	if len(units) == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Units not found"})
	}
	// return users
	return c.Status(200).JSON(fiber.Map{"status": "sucess", "message": "units Found", "data": units})
}

func GetAllLengthUnits(c *fiber.Ctx) error {
	db := database.DB.Db
	var units []metaData.LeightUnit
	// find all users in the database
	db.Find(&units)
	// If no customer found, return an error
	if len(units) == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Units not found"})
	}
	// return users
	return c.Status(200).JSON(fiber.Map{"status": "sucess", "message": "units Found", "data": units})
}

func GetAllCurrencies(c *fiber.Ctx) error {
	db := database.DB.Db
	var currencies []metaData.Currency
	// find all users in the database
	db.Find(&currencies)
	// If no customer found, return an error
	if len(currencies) == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Units not found"})
	}
	// return users
	return c.Status(200).JSON(fiber.Map{"status": "sucess", "message": "units Found", "data": currencies})
}
