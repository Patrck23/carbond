package carRegistration

import "gorm.io/gorm"

// CustomerScan represents a scan or photo associated with a car
type CarScan struct {
	gorm.Model
	CarID  uint   `json:"car_id"`
	Scan   string `gorm:"not null" json:"scan"`
	Title  string `json:"title"`
	Remark string `json:"remark"`
	Car    Car    `gorm:"foreignKey:CarID;references:ID"`
}
