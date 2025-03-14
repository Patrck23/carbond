package saleRegistration

import (
	"gorm.io/gorm"
)

type SalePaymentMode struct {
	gorm.Model
	ModeOfPayment string      `json:"mode_of_payment"`
	TransactionID string      `gorm:"unique" json:"transaction_id"`
	SalePaymentID uint        `gorm:"references:ID" json:"sale_payment_id"`
	SalePayment   SalePayment `gorm:"foreignKey:SalePaymentID;references:ID"`
	CreatedBy     string      `gorm:"size:100" json:"created_by"`
	UpdatedBy     string      `gorm:"size:100" json:"updated_by"`
}
