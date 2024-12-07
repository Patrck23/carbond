package carRegistration

import (
	"gorm.io/gorm"
)

type CarPort struct {
	gorm.Model
	ID          int    `json:"id"`
	CarID  		int    `json:"car_id"`
	Name 		string `json:"name"`
	Category 	string `json:"category"`
	CreatedBy 	string `json:"created_by"`
	UpdatedBy 	int    `json:"updated_by"`
	Car    		Car    `gorm:"references:ID"`
}

type CarPorts struct {
	CarPorts []CarPort
	CurrentCar   int
}