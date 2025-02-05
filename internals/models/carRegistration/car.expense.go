package carRegistration

import (
	"gorm.io/gorm"
)

type CarExpense struct {
	gorm.Model
	CarID       uint    `json:"car_id"`
	Description string  `json:"description"`
	Currency    string  `json:"currency"`
	Amount      float64 `json:"amount"`
	DollarRate  float64 `json:"dollar_rate"`
	ExpenseDate string  `gorm:"type:date" json:"expense_date"`
	CreatedBy   string  `json:"created_by"`
	UpdatedBy   string  `json:"updated_by"`
	Car         Car     `gorm:"foreignKey:CarID;references:ID" json:"car"`
}
