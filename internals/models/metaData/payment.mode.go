package metaData

import (
	"gorm.io/gorm"
)

type PaymentMode struct {
	gorm.Model
	Mode        string `json:"mode"`
	Description string `json:"description"`
	Category    string `json:"category"`
	CreatedBy   string `json:"created_by"`
	UpdatedBy   string `json:"updated_by"`
}
