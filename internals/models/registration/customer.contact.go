package registration

import (
	"gorm.io/gorm"
)

type CustomerContact struct {
	gorm.Model
	CurrentUser        int    `json:"current_user"`
	ID                 int    `json:"id"`
	CustomerID           int    `json:"customer_id"`
	ContactType        int    `json:"contact_type"`
	ContactInformation string `json:"contact_information"`
	CreatedBy          string `json:"created_by"`
	UpdatedBy          int    `json:"updated_by"`
	Customer             Customer `gorm:"references:ID"`
}
type CustomerContacts struct {
	CustomerContacts []CustomerContact
	CurrentCustomer  int
}
