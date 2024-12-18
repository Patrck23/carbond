package customerRegistration

import (
	"gorm.io/gorm"
)

type CustomerAddress struct {
	gorm.Model
	ID          int    `json:"id"`
	CustomerID  int    `json:"customer_id"`
	District  	string `json:"district"`
	Subcounty 	string `json:"subcounty"`
	Parish    	string `json:"parish"`
	Village   	string `json:"village"`
	CreatedBy 	string `json:"created_by"`
	UpdatedBy 	int    `json:"updated_by"`
	Customer    Customer `gorm:"references:ID"`
}
