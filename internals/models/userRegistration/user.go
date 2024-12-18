package userRegistration

import (
	"car-bond/internals/models/companyRegistration"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserUUID 	uuid.UUID 				  `json:"user_uuid"`
	Surname    string                     `gorm:"size:100;not null" json:"surname"`
	Firstname  string                     `gorm:"size:100;not null" json:"firstname"`
	Othername  string                     `gorm:"size:100" json:"othername"`
	Gender     string                     `gorm:"size:10;not null" json:"gender"`
	Title      string                     `gorm:"size:50" json:"title"`
	Username   string 					  `gorm:"uniqueIndex;not null" json:"username"`
	Email      string 					  `gorm:"uniqueIndex;not null" json:"email"`
	Password   string                     `gorm:"size:255;not null" json:"password"`
	CompanyID  uint                       `json:"company_id"`
	Company    companyRegistration.Company `gorm:"foreignKey:CompanyID;references:ID" json:"company"`
	CreatedBy  string                     `gorm:"size:100" json:"created_by"`
	UpdatedBy  string                     `gorm:"size:100" json:"updated_by"`
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	user.UserUUID = uuid.New()
	return nil
}

