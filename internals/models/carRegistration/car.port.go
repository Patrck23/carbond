package carRegistration

import (
	"gorm.io/gorm"
)

type CarPort struct {
	gorm.Model
	CarID  		uint    `json:"car_id"`
	Name 		string `json:"name"`
	Category 	string `json:"category"`
	CreatedBy 	string `json:"created_by"`
	UpdatedBy 	string `json:"updated_by"`
	Car    		Car    `gorm:"foreignKey:CarID;references:ID" json:"car"`
}