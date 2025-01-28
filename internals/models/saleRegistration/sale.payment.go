package saleRegistration

import (
	"gorm.io/gorm"
)

type SalePayment struct {
	gorm.Model
	AmountPayed float64 `gorm:"type:numeric;not null" json:"amount_payed"`
	PaymentDate string  `gorm:"type:date;not null" json:"payment_date"`
	SaleID      uint    `gorm:"references:ID" json:"sale_id"`
	Sale        Sale    `gorm:"foreignKey:SaleID;references:ID"`
	CreatedBy   string  `gorm:"size:100" json:"created_by"`
	UpdatedBy   string  `gorm:"size:100" json:"updated_by"`
}
