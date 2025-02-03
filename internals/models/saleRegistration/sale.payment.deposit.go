package saleRegistration

import (
	"gorm.io/gorm"
)

type SalePaymentDeposit struct {
	gorm.Model
	BankName        string      `json:"bank_name"`
	BankAccount     string      `json:"bank_account"`
	BankBranch      string      `json:"bank_branch"`
	AmountDeposited float64     `json:"amount_deposited"`
	DateDeposited   string      `gorm:"type:date" json:"date_deposited"`
	DepositScan     string      `json:"deposit_scan"`
	SalePaymentID   uint        `gorm:"references:ID" json:"sale_payment_id"`
	SalePayment     SalePayment `gorm:"foreignKey:SalePaymentID;references:ID"`
	CreatedBy       string      `gorm:"size:100" json:"created_by"`
	UpdatedBy       string      `gorm:"size:100" json:"updated_by"`
}
