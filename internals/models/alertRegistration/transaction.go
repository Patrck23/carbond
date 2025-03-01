package alertRegistration

import (
	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	CarChasisNumber string `gorm:"not null;index" json:"car_chasis_number"`
	TransactionType string `gorm:"not null" json:"transaction_type"`
	FromCompanyId   uint   `gorm:"not null;index" json:"from_company_id"`
	ToCompanyId     uint   `json:"to_company_id"`
	CreatedBy       string `gorm:"size:100;not null" json:"created_by"`
	UpdatedBy       string `gorm:"size:100" json:"updated_by"`
	ViewStatus      bool   `json:"view_status"`
}
