package customerRegistration

import (
	"gorm.io/gorm"
)

type CustomerAddress struct {
	gorm.Model
	ID          uint    `json:"id,omitempty"`
	CustomerID  uint    `json:"customer_id"`
	District  	string `json:"district"`
	Subcounty 	string `json:"subcounty"`
	Parish    	string `json:"parish"`
	Village   	string `json:"village"`
	CreatedBy 	string `json:"created_by"`
	UpdatedBy 	string `json:"updated_by"`
	Customer    Customer `gorm:"foreignKey:CustomerID;references:ID" json:"customer"`
}
