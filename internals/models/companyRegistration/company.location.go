package companyRegistration

import (
	"gorm.io/gorm"
)

type CompanyLocation struct {
	gorm.Model
	ID          uint      `gorm:"primary key;autoIncrement" json:"id,omitempty"`
	CompanyID  	uint      `json:"company_id"`
	Company     Company   `gorm:"foreignKey:CompanyID;references:ID" json:"company"`
	Country   	string    `json:"country"`
	CreatedBy   string    `json:"created_by"`
	UpdatedBy   string    `json:"updated_by"`
}
