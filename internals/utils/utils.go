package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/go-playground/validator/v10"
)

func ParseBody(r *http.Request, x interface{}) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}
	if err := json.Unmarshal(body, x); err != nil {
		return
	}
}

type Pagination struct {
	Page         int   `json:"page"`
	Limit        int   `json:"limit"`
	TotalItems   int64 `json:"total_items"`
	TotalPages   int   `json:"total_pages"`
	CurrentPage  int   `json:"current_page"`
	ItemsPerPage int   `json:"items_per_page"`
}

// Paginate is a helper function to handle pagination
func Paginate[T any](c *fiber.Ctx, db *gorm.DB, model T) (Pagination, []T, error) {
	// Get pagination parameters from query
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10
	} else if limit > 100 { // Example limit cap
		limit = 100
	}

	// Get the total count of items in the database
	var totalItems int64
	err = db.Model(&model).Count(&totalItems).Error
	if err != nil {
		return Pagination{}, nil, fmt.Errorf("failed to count items: %w", err)
	}

	// Calculate the offset
	offset := (page - 1) * limit

	// Query the items with pagination
	var items []T
	err = db.Offset(offset).Limit(limit).Find(&items).Error
	if err != nil {
		return Pagination{}, nil, fmt.Errorf("failed to retrieve items: %w", err)
	}

	// Calculate total pages
	totalPages := int(totalItems) / limit
	if totalItems%int64(limit) > 0 {
		totalPages++
	}

	// Return pagination info and items
	pagination := Pagination{
		Page:         page,
		Limit:        limit,
		TotalItems:   totalItems,
		TotalPages:   totalPages,
		CurrentPage:  page,
		ItemsPerPage: limit,
	}

	return pagination, items, nil
}

var validate = validator.New()

func ValidateStruct(data interface{}) map[string]string {
	err := validate.Struct(data)
	if err == nil {
		return nil
	}

	errors := make(map[string]string)
	for _, err := range err.(validator.ValidationErrors) {
		errors[err.Field()] = err.Tag()
	}
	return errors
}

func StrToInt(value string) int {
	result, err := strconv.Atoi(value)
	if err != nil {
		log.Printf("Error converting string to int: %s", err)
		return 0
	}
	return result
}

func StrToFloat(value string) float64 {
	result, err := strconv.ParseFloat(value, 64)
	if err != nil {
		log.Printf("Error converting string to float: %s", err)
		return 0.0
	}
	return result
}

func StrToBool(value string) bool {
	result, err := strconv.ParseBool(value)
	if err != nil {
		log.Printf("Error converting string to bool: %s", err)
		return false
	}
	return result
}

func StrToUint(value string) uint {
	result, err := strconv.Atoi(value)
	if err != nil {
		log.Printf("Error converting string to uint: %s", err)
		return 0
	}
	return uint(result)
}

func StrToUintPointer(value string) *uint {
	result, err := strconv.Atoi(value)
	if err != nil {
		log.Printf("Error converting string to uint: %s", err)
		return nil
	}
	resultUint := uint(result)
	return &resultUint
}

// Converts a string to a pointer to an int
func StrToIntPointer(value string) *int {
	result, err := strconv.Atoi(value)
	if err != nil {
		log.Printf("Error converting string to int: %s", err)
		return nil
	}
	return &result
}
