package registration

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Client struct {
	gorm.Model
	ID          uint      `gorm:"primary key;autoIncrement" json:"id"`
	ClientUUID  uuid.UUID `json:"client_uuid"`
	PatientNo   string    `json:"patient_no"`
	Surname     string    `json:"surname"`
	Firstname   string    `json:"firstname"`
	Othername   string    `json:"othername"`
	Gender      string    `json:"gender"`
	Nationality string    `json:"nationality"`
	Age         int       `json:"age"`
	DOB         string    `gorm:"type:date" json:"dob"`
	Telephone   string    `json:"telephone"`
	NIN         string    `json:"nin"`
	CreatedBy   string    `json:"created_by"`
	UpdatedBy   string    `json:"updated_by"`
}

// Clients struct
type Clients struct {
	Clients       []Client
	CurrentClient int
}

func (client *Client) BeforeCreate(tx *gorm.DB) (err error) {
	client.ClientUUID = uuid.New()
	return
}
