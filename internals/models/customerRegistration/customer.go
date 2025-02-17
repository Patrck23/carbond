package customerRegistration

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Customer struct {
	gorm.Model
	CustomerUUID uuid.UUID `json:"customer_uuid"`
	Surname      string    `json:"surname"`
	Firstname    string    `json:"firstname"`
	Othername    string    `json:"othername"`
	Gender       string    `json:"gender"`
	Nationality  string    `json:"nationality"`
	Age          uint      `json:"age"`
	DOB          string    `gorm:"type:date" json:"dob"`
	Telephone    string    `json:"telephone"`
	Email        string    `json:"email"`
	NIN          string    `json:"nin"`
	CreatedBy    string    `json:"created_by"`
	UpdatedBy    string    `json:"updated_by"`
	UploadFile   string    `json:"upload_file"` // Store the file path or URL
}

func (customer *Customer) BeforeCreate(tx *gorm.DB) (err error) {
	customer.CustomerUUID = uuid.New()
	return
}
