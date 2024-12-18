package companyRegistration

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Company struct {
	gorm.Model
	ID          uint      `gorm:"primary key;autoIncrement" json:"id"`
	CompanyUUID uuid.UUID `json:"company_uuid"`
	Name   		string    `json:"name"`
	StartDate   string    `gorm:"type:date" json:"start_date"`
	CreatedBy   string    `json:"created_by"`
	UpdatedBy   string    `json:"updated_by"`
}

func (company *Company) BeforeCreate(tx *gorm.DB) (err error) {
	company.CompanyUUID = uuid.New()
	return
}
