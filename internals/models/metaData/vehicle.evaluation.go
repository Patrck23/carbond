package metaData

import (
	"gorm.io/gorm"
)

type VehicleEvaluation struct {
	gorm.Model
	HSCCode     string  `json:"hsc_code"`
	COO         string  `json:"coo"`
	Description string  `gorm:"size:200;not null" json:"description"`
	CC          string  `gorm:"size:200" json:"cc"`
	CIF         float64 `gorm:"not null" json:"cif"`
	CreatedBy   string  `gorm:"size:200" json:"created_by"`
	UpdatedBy   string  `gorm:"size:200" json:"updated_by"`
}
