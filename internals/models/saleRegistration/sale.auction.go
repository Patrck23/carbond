package saleRegistration

import (
	"car-bond/internals/models/carRegistration"
	"car-bond/internals/models/companyRegistration"

	"gorm.io/gorm"
)

type SaleAuction struct {
	gorm.Model
	CarID              int                         `gorm:"references:ID" json:"car_id"`
	Car                carRegistration.Car         `gorm:"foreignKey:CarID"`
	CompanyID          int                         `gorm:"references:ID" json:"company_id"`
	Company            companyRegistration.Company `gorm:"foreignKey:CompanyID"`
	AuctionUserCompany string                      `gorm:"size:100" json:"auction_user_company"`
	SaleDate           string                      `gorm:"type:date" json:"sale_date"`
	Price              float64                     `json:"price"`
	VATTax             float64                     `json:"vat_tax"`
	RecycleFee         float64                     `json:"recycle_fee"`
	CreatedBy          string                      `gorm:"size:100" json:"created_by"`
	UpdatedBy          string                      `gorm:"size:100" json:"updated_by"`
}
