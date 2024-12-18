package saleRegistration

import (
	"gorm.io/gorm"
)

type SalePayment struct {
	gorm.Model
	ID           uint                       `gorm:"primaryKey;autoIncrement" json:"id,omitempty"`
	AmountPayed   float64                   `gorm:"type:numeric;not null" json:"amount_payed"`
	PaymentDate   string                    `gorm:"type:date;not null" json:"payment_date"`
	SaleID 		 int                        `gorm:"references:ID"`
	Sale  		 Sale  						`gorm:"foreignKey:SaleID;references:ID" json:"sale"`
	CreatedBy    string                     `gorm:"size:100" json:"created_by"`
	UpdatedBy    string                     `gorm:"size:100" json:"updated_by"`
}