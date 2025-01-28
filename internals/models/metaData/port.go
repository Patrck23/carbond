package metaData

import (
	"gorm.io/gorm"
)

type Port struct {
	gorm.Model
	Name      string `json:"name"`
	Location  string `json:"location"`
	Category  string `json:"category"`
	Function  string `json:"function"`
	CreatedBy string `json:"created_by"`
	UpdatedBy string `json:"updated_by"`
}
