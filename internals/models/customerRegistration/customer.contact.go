package customerRegistration

import (
	"gorm.io/gorm"
)

type CustomerContact struct {
	gorm.Model
	ID                 uint    `json:"id,omitempty"`
	CustomerID         uint    `json:"customer_id"`
	ContactType        string  `json:"contact_type"`
	ContactInformation string `json:"contact_information"`
	CreatedBy          string `json:"created_by"`
	UpdatedBy          string   `json:"updated_by"`
	Customer    	   Customer `gorm:"foreignKey:CustomerID;references:ID" json:"customer"`
}

