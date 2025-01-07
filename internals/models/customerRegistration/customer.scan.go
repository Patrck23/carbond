package customerRegistration

import "gorm.io/gorm"

// CustomerScan represents a scan or photo associated with a car
type CustomerScan struct {
	gorm.Model
	CustomerID uint     `json:"customer_id"`
	Scan       string   `gorm:"not null" json:"scan"`
	Title      string   `json:"title"`
	Remark     string   `json:"remark"`
	Customer   Customer `gorm:"foreignKey:CustomerID;references:ID" json:"customer"`
}
