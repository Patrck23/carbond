package controllers

import (
	"car-bond/internals/models/saleRegistration"
	"car-bond/internals/utils"
	"errors"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type SaleAuctionRepository interface {
	CreateSale(sale *saleRegistration.SaleAuction) error
	GetPaginatedSales(c *fiber.Ctx) (*utils.Pagination, []saleRegistration.SaleAuction, error)
	GetSaleByID(id string) (saleRegistration.SaleAuction, error)
	UpdateSale(sale *saleRegistration.SaleAuction) error
	DeleteByID(id string) error
}

type SaleAuctionRepositoryImpl struct {
	db *gorm.DB
}

func NewSaleAuctionRepository(db *gorm.DB) SaleAuctionRepository {
	return &SaleAuctionRepositoryImpl{db: db}
}

type SaleAuctionController struct {
	repo SaleAuctionRepository
}

func NewSaleAuctionController(repo SaleAuctionRepository) *SaleAuctionController {
	return &SaleAuctionController{repo: repo}
}

// ============================================

func (r *SaleAuctionRepositoryImpl) CreateSale(sale *saleRegistration.SaleAuction) error {
	return r.db.Create(sale).Error
}

func (h *SaleAuctionController) CreateCarSale(c *fiber.Ctx) error {
	// Initialize a new Sale instance
	sale := new(saleRegistration.SaleAuction)

	// Parse the request body into the sale instance
	if err := c.BodyParser(sale); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input provided",
			"data":    err.Error(),
		})
	}

	// Attempt to create the sale record using the repository
	if err := h.repo.CreateSale(sale); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create sale",
			"data":    err.Error(),
		})
	}

	// Return the newly created sale record
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Sale created successfully",
		"data":    sale,
	})
}

// =====================

func (r *SaleAuctionRepositoryImpl) GetPaginatedSales(c *fiber.Ctx) (*utils.Pagination, []saleRegistration.SaleAuction, error) {
	pagination, sales, err := utils.Paginate(c, r.db.Preload("Car"), saleRegistration.SaleAuction{})
	if err != nil {
		return nil, nil, err
	}
	return &pagination, sales, nil
}

func (h *SaleAuctionController) GetAllCarSales(c *fiber.Ctx) error {
	pagination, sales, err := h.repo.GetPaginatedSales(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve sales",
			"data":    err.Error(),
		})
	}

	// Return the paginated response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "sales retrieved successfully",
		"data":    sales,
		"pagination": fiber.Map{
			"total_items":  pagination.TotalItems,
			"total_pages":  pagination.TotalPages,
			"current_page": pagination.CurrentPage,
			"limit":        pagination.ItemsPerPage,
		},
	})
}

// ==============

// Get a single car sale by ID

func (r *SaleAuctionRepositoryImpl) GetSaleByID(id string) (saleRegistration.SaleAuction, error) {
	var sale saleRegistration.SaleAuction
	err := r.db.Preload("Car").First(&sale, "id = ?", id).Error
	return sale, err
}

// GetCarSale fetches a sale with its associated contacts and addresses from the database
func (h *SaleAuctionController) GetCarSale(c *fiber.Ctx) error {
	// Get the sale ID from the route parameters
	id := c.Params("id")

	// Fetch the sale by ID
	sale, err := h.repo.GetSaleByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Sale not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve Sale",
			"data":    err.Error(),
		})
	}

	// Return the response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Sale and associated data retrieved successfully",
		"data":    sale,
	})
}

// =====================================

func (r *SaleAuctionRepositoryImpl) UpdateSale(sale *saleRegistration.SaleAuction) error {
	return r.db.Save(sale).Error
}

// Define the UpdateSale struct
type UpdateSaleAuctionPayload struct {
	CarID       int     `json:"car_id"`
	CompanyID   int     `json:"company_id"`
	Auction     string  `json:"auction"`
	AuctionDate string  `json:"auction_date"`
	Price       float64 `json:"price"`
	VATTax      float64 `json:"vat_tax"`
	RecycleFee  float64 `json:"recycle_fee"`
	UpdatedBy   string  `json:"updated_by"`
}

// UpdateSale handler function
func (h *SaleAuctionController) UpdateSale(c *fiber.Ctx) error {
	// Get the sale ID from the route parameters
	id := c.Params("id")

	// Find the sale in the database
	sale, err := h.repo.GetSaleByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{
				"status":  "error",
				"message": "Sale not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve sale",
			"data":    err.Error(),
		})
	}

	// Parse the request body into the UpdateSaleAuctionPayload struct
	var payload UpdateSaleAuctionPayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input",
			"data":    err.Error(),
		})
	}

	// Update the sale fields using the payload
	updateSaleAuctionFields(&sale, payload) // Pass the parsed payload

	// Save the changes to the database
	if err := h.repo.UpdateSale(&sale); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update sale",
			"data":    err.Error(),
		})
	}

	// Return the updated sale
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "sale updated successfully",
		"data":    sale,
	})
}

// updateSaleAuctionFields updates the fields of a Sale using the UpdateSale struct
func updateSaleAuctionFields(sale *saleRegistration.SaleAuction, updateSaleData UpdateSaleAuctionPayload) {
	sale.CarID = updateSaleData.CarID
	sale.CompanyID = updateSaleData.CompanyID
	sale.Auction = updateSaleData.Auction
	sale.AuctionDate = updateSaleData.AuctionDate
	sale.Price = updateSaleData.Price
	sale.VATTax = updateSaleData.VATTax
	sale.RecycleFee = updateSaleData.RecycleFee
	sale.UpdatedBy = updateSaleData.UpdatedBy
}

// ============================

// DeleteByID deletes a sale by ID
func (r *SaleAuctionRepositoryImpl) DeleteByID(id string) error {
	if err := r.db.Delete(&saleRegistration.SaleAuction{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

// DeleteSaleByID deletes a Sale by its ID
func (h *SaleAuctionController) DeleteSaleByID(c *fiber.Ctx) error {
	// Get the Sale ID from the route parameters
	id := c.Params("id")

	// Find the Sale in the database
	sale, err := h.repo.GetSaleByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{
				"status":  "error",
				"message": "Sale not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to find sale",
			"data":    err.Error(),
		})
	}

	// Delete the Sale
	if err := h.repo.DeleteByID(id); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to delete Sale",
			"data":    err.Error(),
		})
	}

	// Return success response
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Sale deleted successfully",
		"data":    sale,
	})
}

// ===============================================================================================
