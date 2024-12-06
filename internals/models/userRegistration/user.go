package userRegistration

import (
	"car-bond/internals/models/companyRegistration"
	"golang.org/x/crypto/bcrypt"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID          uint      				  `gorm:"primary key;autoIncrement" json:"id"`
	UserUUID 	uuid.UUID 				  `json:"user_uuid"`
	Surname    string                     `gorm:"size:100;not null" json:"surname"`
	Firstname  string                     `gorm:"size:100;not null" json:"firstname"`
	Othername  string                     `gorm:"size:100" json:"othername"`
	Gender     string                     `gorm:"size:10;not null" json:"gender"`
	Title      string                     `gorm:"size:50" json:"title"`
	Password   string                     `gorm:"size:255;not null" json:"password"`
	CompanyID  uint                       `json:"company_id"`
	Company    companyRegistration.Company `gorm:"foreignKey:CompanyID;references:ID" json:"company"`
	CreatedBy  string                     `gorm:"size:100" json:"created_by"`
	UpdatedBy  string                     `gorm:"size:100" json:"updated_by"`
}

// Users struct
type Users struct {
	Users       []User
	CurrentUser int
}

// HashPassword hashes the user's password
func (user *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return nil
}

// BeforeCreate Hook hashes the password before creating a user
func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	// Generate UUID
	user.UserUUID = uuid.New()

	// Hash the password
	if err := user.HashPassword(); err != nil {
		return err
	}
	return nil
}

// BeforeUpdate Hook hashes the password if it has changed
func (user *User) BeforeUpdate(tx *gorm.DB) (err error) {
	// Check if the password field has been updated
	if user.Password != "" {
		if err := user.HashPassword(); err != nil {
			return err
		}
	}
	return nil
}
