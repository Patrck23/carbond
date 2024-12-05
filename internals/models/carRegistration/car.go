package carRegistration

import (
	"car-bond/internals/models/companyRegistration"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Car struct {
	gorm.Model
	ID           uint                         `gorm:"primaryKey;autoIncrement" json:"id"`
	CarUUID      uuid.UUID                    `json:"car_uuid"`
	VinNumber    string                       `gorm:"size:100;not null" json:"vin_number"`
	Make         string                       `gorm:"size:100;not null" json:"make"`
	CarModel     string                       `gorm:"size:100;not null" json:"model"`
	Year         int                          `gorm:"not null" json:"year"`
	BidPrice     float64                      `gorm:"type:numeric;not null" json:"bid_price"`
	PurchaseDate string                       `gorm:"type:date;not null" json:"purchase_date"`
	// FromCompanyID int                         `gorm:"references:ID"`
	// FromCompany  companyRegistration.Company  `gorm:"foreignKey:FromCompanyID"`
	// ToCompanyID   int                         `gorm:"references:ID"`
	// ToCompany    companyRegistration.Company  `gorm:"foreignKey:FromCompanyID"`
	FromCompany  companyRegistration.Company  `gorm:"foreignKey:FromCompanyID;references:ID" json:"from_company"`
	ToCompany    companyRegistration.Company  `gorm:"foreignKey:ToCompanyID;references:ID" json:"to_company"`
	FromCompanyID uint                        `json:"from_company_id"`
	ToCompanyID   uint                        `json:"to_company_id"`
	CreatedBy    string                       `gorm:"size:100" json:"created_by"`
	UpdatedBy    string                       `gorm:"size:100" json:"updated_by"`
}

// Cars struct
type Cars struct {
	Cars       []Car
	CurrentCar int
}

func (car *Car) BeforeCreate(tx *gorm.DB) (err error) {
	car.CarUUID = uuid.New()
	return
}