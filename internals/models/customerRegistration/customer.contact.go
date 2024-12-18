package customerRegistration

import (
	"gorm.io/gorm"
)

type CustomerContact struct {
	gorm.Model
	ID                 int    `json:"id,omitempty"`
	CustomerID         int    `json:"customer_id"`
	ContactType        int    `json:"contact_type"`
	ContactInformation string `json:"contact_information"`
	CreatedBy          string `json:"created_by"`
	UpdatedBy          int    `json:"updated_by"`
	Customer    	   Customer `gorm:"foreignKey:CustomerID;references:ID" json:"customer"`
}

