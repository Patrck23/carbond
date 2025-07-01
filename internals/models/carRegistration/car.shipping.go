package carRegistration

import "gorm.io/gorm"

// CustomerScan represents a scan or photo associated with a car
type CarShippingInvoice struct {
	gorm.Model
	InvoiceNo    string `gorm:"unique;not null" json:"invoice_no"`
	ShipDate     string `gorm:"not null" json:"ship_date"`
	VesselName   string `json:"vessel_name"`
	FromLocation string `json:"from_location"`
	ToLocation   string `json:"to_location"`
	CreatedBy    string `gorm:"size:100" json:"created_by"`
	UpdatedBy    string `gorm:"size:100" json:"updated_by"`

	Locked bool `gorm:"default:false" json:"locked"` // ← Add this

	Cars []Car `gorm:"foreignKey:CarShippingInvoiceID" json:"cars"`
}
