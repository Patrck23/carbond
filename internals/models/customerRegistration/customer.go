package customerRegistration

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Customer struct {
	gorm.Model
	ID          uint      `gorm:"primary key;autoIncrement" json:"id"`
	CustomerUUID  uuid.UUID `json:"customer_uuid"`
	Surname     string    `json:"surname"`
	Firstname   string    `json:"firstname"`
	Othername   string    `json:"othername"`
	Gender      string    `json:"gender"`
	Nationality string    `json:"nationality"`
	Age         int       `json:"age"`
	DOB         string    `gorm:"type:date" json:"dob"`
	Telephone   string    `json:"telephone"`
	Email		string	  `json:"email"`
	NIN         string    `json:"nin"`
	CreatedBy   string    `json:"created_by"`
	UpdatedBy   string    `json:"updated_by"`
}

// Customers struct
type Customers struct {
	Customers       []Customer
	CurrentCustomer int
}

func (customer *Customer) BeforeCreate(tx *gorm.DB) (err error) {
	customer.CustomerUUID = uuid.New()
	return
}
