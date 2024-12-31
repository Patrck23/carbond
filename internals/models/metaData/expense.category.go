package metaData

import "gorm.io/gorm"

type ExpenseCategory struct {
	gorm.Model
	Name      string `json:"name"`
	CreatedBy string `json:"created_by"`
	UpdatedBy string `json:"updated_by"`
}
