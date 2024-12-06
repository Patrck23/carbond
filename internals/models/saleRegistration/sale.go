package saleRegistration

import (
	"car-bond/internals/models/carRegistration"
	"car-bond/internals/models/customerRegistration"
	"gorm.io/gorm"
)

type Sale struct {
	gorm.Model
	ID           uint                         `gorm:"primaryKey;autoIncrement" json:"id"`
	TotalPrice   float64                     `gorm:"type:numeric;not null" json:"total_price"`
	SaleDate 	 string                       `gorm:"type:date;not null" json:"sale_date"`
	CarID 		 int                          `gorm:"references:ID"`
	Car  		 carRegistration.Car  		  `gorm:"foreignKey:CarID"`
	CustomerID   int                          `gorm:"references:ID"`
	Customer     customerRegistration.Customer `gorm:"foreignKey:CustomerID"`
	IsFullPayment bool               		  `json:"is_full_payment"`
	PaymentPeriod int 						  `json:"payment_period`
	CreatedBy    string                       `gorm:"size:100" json:"created_by"`
	UpdatedBy    string                       `gorm:"size:100" json:"updated_by"`
}

// Sales struct
type Sales struct {
	Sales       []Sale
	CurrentSale int
}
