package repository

import (
	"car-bond/internals/models/saleRegistration"
	"car-bond/internals/utils"

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

func (r *SaleAuctionRepositoryImpl) CreateSale(sale *saleRegistration.SaleAuction) error {
	return r.db.Create(sale).Error
}

func (r *SaleAuctionRepositoryImpl) GetPaginatedSales(c *fiber.Ctx) (*utils.Pagination, []saleRegistration.SaleAuction, error) {
	pagination, sales, err := utils.Paginate(c, r.db.Preload("Car"), saleRegistration.SaleAuction{})
	if err != nil {
		return nil, nil, err
	}
	return &pagination, sales, nil
}

func (r *SaleAuctionRepositoryImpl) GetSaleByID(id string) (saleRegistration.SaleAuction, error) {
	var sale saleRegistration.SaleAuction
	err := r.db.Preload("Car").First(&sale, "id = ?", id).Error
	return sale, err
}

func (r *SaleAuctionRepositoryImpl) UpdateSale(sale *saleRegistration.SaleAuction) error {
	return r.db.Save(sale).Error
}

// DeleteByID deletes a sale by ID
func (r *SaleAuctionRepositoryImpl) DeleteByID(id string) error {
	if err := r.db.Delete(&saleRegistration.SaleAuction{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}
