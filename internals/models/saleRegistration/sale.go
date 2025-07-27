package saleRegistration

import (
	"car-bond/internals/models/carRegistration"
	"car-bond/internals/models/companyRegistration"
	"car-bond/internals/models/customerRegistration"

	"gorm.io/gorm"
)

type Sale struct {
	gorm.Model
	TotalPrice    float64                       `gorm:"type:numeric;not null" json:"total_price"`
	DollarRate    float64                       `json:"dollar_rate"`
	SaleDate      string                        `gorm:"type:date;not null" json:"sale_date"`
	CarID         uint                          `gorm:"references:ID" json:"car_id"`
	Car           carRegistration.Car           `gorm:"foreignKey:CarID"`
	CompanyID     int                           `gorm:"references:ID" json:"company_id"`
	Company       companyRegistration.Company   `gorm:"foreignKey:CompanyID"`
	CustomerID    *uint                         `json:"customer_id"`
	Customer      customerRegistration.Customer `gorm:"foreignKey:CustomerID"`
	IsFullPayment bool                          `json:"is_full_payment"`
	InitalPayment float64                       `json:"initial_payment"`
	PaymentPeriod int                           `json:"payment_period"`
	CreatedBy     string                        `gorm:"size:100" json:"created_by"`
	UpdatedBy     string                        `gorm:"size:100" json:"updated_by"`
}
