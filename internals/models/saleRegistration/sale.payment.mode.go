package saleRegistration

import (
	"gorm.io/gorm"
)

type SalePaymentMode struct {
	gorm.Model
	ID           	uint                `gorm:"primaryKey;autoIncrement" json:"id"`
	ModeOfPayment   string              `json:"mode_of_payment"`
	transactionId 	string 				`gorm:"unique" json:"transaction_id"`
	SalePaymentID   int                	`gorm:"references:ID"`
	SalePayment  	SalePayment  		`gorm:"foreignKey:SalePaymentID"`
	CreatedBy    	string              `gorm:"size:100" json:"created_by"`
	UpdatedBy    	string              `gorm:"size:100" json:"updated_by"`
}