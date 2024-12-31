package metaData

import "gorm.io/gorm"

type WeightUnit struct {
	gorm.Model
	Name      string `json:"name"`
	Symbol    string `json:"symbol"`
	CreatedBy string `json:"created_by"`
	UpdatedBy string `json:"updated_by"`
}
