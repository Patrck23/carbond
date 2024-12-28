package controllers

import (
	vehicleevaluation "car-bond/internals/models/vehicleEvaluation"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/ledongthuc/pdf"
	"gorm.io/gorm"
)

func UploadPDF(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Parse PDF file from the request
		file, err := c.FormFile("file")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Failed to upload file",
			})
		}

		filePath := fmt.Sprintf("./%s", file.Filename)
		if err := c.SaveFile(file, filePath); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to save file",
			})
		}

		// Open PDF file
		f, r, err := pdf.Open(filePath)
		defer f.Close()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to parse PDF",
			})
		}

		// Parse content
		var content string
		totalPage := r.NumPage()
		for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
			page := r.Page(pageIndex)
			if !page.V.IsNull() {
				// Use the Content() method to extract text
				text, err := page.GetPlainText(nil)
				if err != nil {
					log.Printf("Failed to extract text from page %d: %v", pageIndex, err)
					continue
				}
				content += text
			}
		}

		// Extract structured data
		lines := strings.Split(content, "\n")
		for _, line := range lines {
			fields := strings.Fields(line)

			// Ensure at least 5 fields (S/N, HSC Code, COO, Description, CC, CIF)
			if len(fields) < 5 {
				continue
			}

			// Skip the first column (S/N)
			hscCode := fields[1]
			coo := fields[2]
			description := strings.Join(fields[3:len(fields)-2], " ")
			cc := fields[len(fields)-2]

			// Parse CIF as decimal
			cif, err := parseDecimal(fields[len(fields)-1])
			if err != nil {
				log.Println("Skipping invalid CIF value:", fields[len(fields)-1])
				continue
			}

			// Create VehicleEvaluation entry
			evaluation := vehicleevaluation.VehicleEvaluation{
				HSCCode:     hscCode,
				COO:         coo,
				Description: description,
				CC:          cc,
				CIF:         cif,     // Decimal value
				CreatedBy:   "admin", // Default CreatedBy
			}

			// Save to the database
			if err := db.Create(&evaluation).Error; err != nil {
				log.Println("Failed to save entry:", err)
			}
		}

		return c.JSON(fiber.Map{
			"message": "Data uploaded successfully",
		})
	}
}

func parseDecimal(input string) (float64, error) {
	return strconv.ParseFloat(input, 64)
}
