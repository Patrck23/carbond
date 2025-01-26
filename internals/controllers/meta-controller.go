package controllers

import (
	"car-bond/internals/models/metaData"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type Excecute interface {
	Begin() Excecute
	Commit() error
	Rollback()
	Exec(query string, args ...interface{}) *gorm.DB
	Create(value interface{}) *gorm.DB
}

type GormDatabase struct {
	db *gorm.DB
}

func NewExcecute(db *gorm.DB) Excecute {
	return &GormDatabase{db: db}
}

type MetaController struct {
	repo Excecute
}

func NewMetaController(repo Excecute) *MetaController {
	return &MetaController{repo: repo}
}

// ===============

func (g *GormDatabase) Begin() Excecute {
	return &GormDatabase{db: g.db.Begin()}
}

func (g *GormDatabase) Commit() error {
	return g.db.Commit().Error
}

func (g *GormDatabase) Rollback() {
	g.db.Rollback()
}

func (g *GormDatabase) Exec(query string, args ...interface{}) *gorm.DB {
	return g.db.Exec(query, args...)
}

func (g *GormDatabase) Create(value interface{}) *gorm.DB {
	return g.db.Create(value)
}

func (m *MetaController) ProcessExcelAndUpload(c *fiber.Ctx, db Excecute) error {
	// Begin transaction
	db = db.Begin()

	// Retrieve the uploaded file
	file, err := c.FormFile("excel")
	if err != nil {
		db.Rollback()
		return fiber.NewError(fiber.StatusBadRequest, "Failed to retrieve uploaded file: "+err.Error())
	}

	// Save the uploaded file to a temporary location
	tempDir := "./uploads"
	if err := os.MkdirAll(tempDir, os.ModePerm); err != nil {
		db.Rollback()
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create upload directory: "+err.Error())
	}

	tempFilePath := filepath.Join(tempDir, file.Filename)
	if err := c.SaveFile(file, tempFilePath); err != nil {
		db.Rollback()
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
		db.Rollback()
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to extract table data from Excel: "+extractErr.Error())
	}

	// Clear old data
	if err := m.repo.Exec("DELETE FROM vehicle_evaluations").Error; err != nil {
		db.Rollback()
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

		if err := m.repo.Create(vehicle).Error; err != nil {
			db.Rollback()
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to insert new data: "+err.Error())
		}
	}

	// Commit the transaction
	if err := db.Commit(); err != nil {
		db.Rollback()
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

func (m *MetaController) ProcessExcelAndUploadHandler(c *fiber.Ctx) error {
	// Pass the db service to the original method
	return m.ProcessExcelAndUpload(c, m.repo)
}
