package vehicleevaluation

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type VehicleEvaluation struct {
	gorm.Model
	VehicleUUID uuid.UUID `json:"vehicle_uuid"`
	HSCCode     string    `json:"hsc_code"`
	COO         string    `json:"coo"`
	Description string    `gorm:"size:100;not null" json:"description"`
	CC          string    `gorm:"size:100;not null" json:"cc"`
	CIF         float64   `gorm:"not null" json:"cif"` // Updated to float64
	CreatedBy   string    `gorm:"size:100" json:"created_by"`
	UpdatedBy   string    `gorm:"size:100" json:"updated_by"`
}

func (vehicle *VehicleEvaluation) BeforeCreate(tx *gorm.DB) (err error) {
	vehicle.VehicleUUID = uuid.New()
	return
}
