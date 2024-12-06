package saleRegistration

import (
	"gorm.io/gorm"
)

type SalePayment struct {
	gorm.Model
	ID           uint                       `gorm:"primaryKey;autoIncrement" json:"id"`
	AmountPayed   float64                   `gorm:"type:numeric;not null" json:"amount_payed"`
	PaymentDate   string                    `gorm:"type:date;not null" json:"payment_date"`
	SaleID 		 int                        `gorm:"references:ID"`
	Sale  		 Sale  		  				`gorm:"foreignKey:SaleID"`
	CreatedBy    string                     `gorm:"size:100" json:"created_by"`
	UpdatedBy    string                     `gorm:"size:100" json:"updated_by"`
}

// SalePayments struct
type SalePayments struct {
	SalePayments       []SalePayment
	CurrentSale int
}