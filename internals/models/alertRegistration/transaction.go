package alertRegistration

import (
	"fmt"

	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	CarChasisNumber string `gorm:"not null;index" json:"car_chasis_number"`
	TransactionType string `gorm:"not null" json:"transaction_type"`
	Description     string `gorm:"size:100" json:"description"`
	FromCompanyId   uint   `gorm:"not null;index" json:"from_company_id"`
	ToCompanyId     uint   `json:"to_company_id"`
	CreatedBy       string `gorm:"size:100;not null" json:"created_by"`
	UpdatedBy       string `gorm:"size:100" json:"updated_by"`
	ViewStatus      bool   `json:"view_status"`
}

// BeforeCreate hook to set the Description based on the TransactionType
func (t *Transaction) BeforeCreate(tx *gorm.DB) (err error) {
	switch t.TransactionType {
	case "InTransit":
		t.Description = fmt.Sprintf("Car %s is currently in transit.", t.CarChasisNumber)
	case "Storage":
		t.Description = fmt.Sprintf("Car %s is stored at the facility.", t.CarChasisNumber)
	case "Sold":
		t.Description = fmt.Sprintf("Car %s has been sold.", t.CarChasisNumber)
	default:
		t.Description = "Transaction type not recognized."
	}
	return nil
}
