package vehicleevaluation

import (
	"gorm.io/gorm"
)

type VehicleEvaluation struct {
	gorm.Model
	HSCCode     string `json:"hsc_code"`
	COO         string `json:"coo"`
	Description string `gorm:"size:100;not null" json:"description"`
	CC          string `gorm:"size:100;not null" json:"cc"`
	CIF         int    `gorm:"not null" json:"cif"`
	CreatedBy   string `gorm:"size:100" json:"created_by"`
	UpdatedBy   string `gorm:"size:100" json:"updated_by"`
}
