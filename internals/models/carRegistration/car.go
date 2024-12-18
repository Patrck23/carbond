package carRegistration

import (
	"car-bond/internals/models/companyRegistration"
	"car-bond/internals/models/customerRegistration"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Car struct {
	gorm.Model
	CarUUID      uuid.UUID                    `json:"car_uuid"`
	VinNumber    string                       `gorm:"size:100;not null" json:"vin_number"`
	Make         string                       `gorm:"size:100;not null" json:"make"`
	CarModel     string                       `gorm:"size:100;not null" json:"model"`
	Year         int                          `gorm:"not null" json:"year"`
	BidPrice     float64                      `gorm:"type:numeric;not null" json:"bid_price"`
	VATTax     float64                      `gorm:"type:numeric;not null" json:"vat_tax"`
	PurchaseDate string                       `gorm:"type:date;not null" json:"purchase_date"`
	FromCompany  companyRegistration.Company  `gorm:"foreignKey:FromCompanyID;references:ID" json:"from_company"`
	ToCompany    companyRegistration.Company  `gorm:"foreignKey:ToCompanyID;references:ID" json:"to_company"`
	Destination	 string						  `json:"destination"`
	FromCompanyID uint                        `json:"from_company_id"`
	ToCompanyID   uint                        `json:"to_company_id"`
	CustomerID    int                         `gorm:"references:ID"`
	Customer     customerRegistration.Customer  `gorm:"foreignKey:CustomerID"`
	CreatedBy    string                       `gorm:"size:100" json:"created_by"`
	UpdatedBy    string                       `gorm:"size:100" json:"updated_by"`
}

func (car *Car) BeforeCreate(tx *gorm.DB) (err error) {
	car.CarUUID = uuid.New()
	return
}