package saleRegistration

import (
	"car-bond/internals/models/carRegistration"
	"car-bond/internals/models/companyRegistration"
	"gorm.io/gorm"
)

type Sale struct {
	gorm.Model
	// ID           uint                         `gorm:"primaryKey;autoIncrement" json:"id"`
	TotalPrice   float64                     `gorm:"type:numeric;not null" json:"total_price"`
	SaleDate 	 string                       `gorm:"type:date;not null" json:"sale_date"`
	CarID 		 int                          `gorm:"references:ID"`
	Car  		 carRegistration.Car  		  `gorm:"foreignKey:CarID"`
	CompanyID   int                          `gorm:"references:ID"`
	Company     companyRegistration.Company  `gorm:"foreignKey:CompanyID"`
	IsFullPayment bool               		  `json:"is_full_payment"`
	PaymentPeriod int 						  `json:"payment_period`
	CreatedBy    string                       `gorm:"size:100" json:"created_by"`
	UpdatedBy    string                       `gorm:"size:100" json:"updated_by"`
}
