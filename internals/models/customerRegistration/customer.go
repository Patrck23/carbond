package customerRegistration

import (
	"car-bond/internals/models/companyRegistration"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Customer struct {
	gorm.Model
	CustomerUUID uuid.UUID                    `json:"customer_uuid"`
	Surname      string                       `json:"surname"`
	Firstname    string                       `json:"firstname"`
	Othername    string                       `json:"othername"`
	Gender       string                       `json:"gender"`
	Nationality  string                       `json:"nationality"`
	Age          uint                         `json:"age"` // Always computed, not from payload
	DOB          string                       `gorm:"type:date" json:"dob"`
	Telephone    string                       `json:"telephone"`
	Email        string                       `json:"email"`
	NIN          string                       `json:"nin"`
	CompanyID    *uint                        `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"company_id"`
	Company      *companyRegistration.Company `gorm:"foreignKey:CompanyID;references:ID" json:"company"`
	CreatedBy    string                       `json:"created_by"`
	UpdatedBy    string                       `json:"updated_by"`
	UploadFile   string                       `json:"upload_file"` // Store the file path or URL
}

// BeforeCreate hook
func (c *Customer) BeforeCreate(tx *gorm.DB) (err error) {
	c.CustomerUUID = uuid.New()
	c.calculateAge()
	return
}

// BeforeUpdate hook (recalculate if DOB changes)
func (c *Customer) BeforeUpdate(tx *gorm.DB) (err error) {
	c.calculateAge()
	return
}

// Helper to compute age from DOB
func (c *Customer) calculateAge() {
	if c.DOB == "" {
		c.Age = 0
		return
	}

	// Parse DOB (assuming format YYYY-MM-DD)
	dob, err := time.Parse("2006-01-02", c.DOB)
	if err != nil {
		c.Age = 0
		return
	}

	today := time.Now()
	age := today.Year() - dob.Year()
	if today.YearDay() < dob.YearDay() {
		age--
	}

	if age < 0 {
		age = 0
	}

	c.Age = uint(age)
}
