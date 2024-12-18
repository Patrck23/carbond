package carRegistration

import (
	"gorm.io/gorm"
)

type CarExpense struct {
	gorm.Model
	ID          int    `json:"id"`
	CarID  		int    `json:"car_id"`
	Description string `json:"description"`
	Currency	string `json:"currency"`
	Amount	 	float64 `json:"amount"`
	ExpenseDate string  `gorm:"type:date" json:"expense_date"`
	CreatedBy 	string `json:"created_by"`
	UpdatedBy 	int    `json:"updated_by"`
	Car    		Car    `gorm:"references:ID"`
}