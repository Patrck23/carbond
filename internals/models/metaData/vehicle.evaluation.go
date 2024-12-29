package metaData

import (
	"gorm.io/gorm"
)

type VehicleEvaluation struct {
	gorm.Model
	HSCCode     string  `json:"hsc_code"`
	COO         string  `json:"coo"`
	Description string  `gorm:"size:100;not null" json:"description"`
	CC          string  `gorm:"size:50" json:"cc"`
	CIF         float64 `gorm:"not null" json:"cif"`
	CreatedBy   string  `json:"created_by"`
	UpdatedBy   string  `json:"updated_by"`
}
