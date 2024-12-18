package companyRegistration

import (
	"gorm.io/gorm"
)

type CompanyExpense struct {
	gorm.Model
	CompanyID  	uint      `json:"company_id"`
	Description string   `json:"description"`
	Currency	string   `json:"currency"`
	Amount	 	float64  `json:"amount"`
	ExpenseDate string   `gorm:"type:date" json:"expense_date"`
	CreatedBy 	string   `json:"created_by"`
	UpdatedBy 	string   `json:"updated_by"`
	Company     Company  `gorm:"foreignKey:CompanyID;references:ID" json:"company"`
}