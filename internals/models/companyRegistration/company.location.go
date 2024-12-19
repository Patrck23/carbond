package companyRegistration

import (
	"gorm.io/gorm"
)

type CompanyLocation struct {
	gorm.Model
	CompanyID  	uint      `json:"company_id"`
	Address 	string	  `json:"address"`
	Telephone	string	  `json:"telephone"`
	Country   	string    `json:"country"`
	CreatedBy   string    `json:"created_by"`
	UpdatedBy   string    `json:"updated_by"`
	Company     Company   `gorm:"foreignKey:CompanyID;references:ID" json:"company"`
}
