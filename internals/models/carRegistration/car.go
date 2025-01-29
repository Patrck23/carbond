package carRegistration

import (
	"car-bond/internals/models/companyRegistration"
	"car-bond/internals/models/customerRegistration"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Car struct {
	gorm.Model
	CarUUID               uuid.UUID                    `json:"car_uuid"`
	VinNumber             string                       `gorm:"size:100;not null;unique" json:"vin_number"`
	EngineNumber          string                       `gorm:"size:100;not null;unique" json:"engine_number"`
	EngineCapacity        string                       `gorm:"size:100;not null" json:"engine_capacity"`
	Make                  string                       `gorm:"size:100;not null" json:"make"`
	CarModel              string                       `gorm:"size:100;not null" json:"model"`
	MaximCarry            int                          `gorm:"size:100;not null" json:"maxim_carry"`
	Weight                int                          `gorm:"size:100;not null" json:"weight"`
	GrossWeight           int                          `gorm:"size:100;not null" json:"gross_weight"`
	Length                int                          `gorm:"size:100;not null" json:"length"`
	Width                 int                          `gorm:"size:100;not null" json:"width"`
	Height                int                          `gorm:"size:100;not null" json:"height"`
	ManufactureYear       int                          `gorm:"not null" json:"maunufacture_year"`
	FirstRegistrationYear int                          `gorm:"not null" json:"first_registration_year"`
	Transmission          string                       `json:"transmission"`
	BodyType              string                       `json:"body_type"`
	Colour                string                       `json:"colour"`
	Auction               string                       `json:"auction"`
	Currency              string                       `json:"currency"`
	CarMillage            int                          `json:"millage"`
	FuelConsumption       string                       `json:"fuel_consumtion"`
	PowerSteering         bool                         `json:"ps"`
	PowerWindow           bool                         `json:"pw"`
	ABS                   bool                         `json:"abs"`
	ADS                   bool                         `json:"ads"`
	AlloyWheel            bool                         `json:"aw"`
	SimpleWheel           bool                         `json:"sw"`
	Navigation            bool                         `json:"navigation"`
	AC                    bool                         `json:"ac"`
	BidPrice              float64                      `gorm:"type:numeric;not null" json:"bid_price"`
	PurchaseDate          string                       `gorm:"type:date;not null" json:"purchase_date"`
	FromCompanyID         *uint                        `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"from_company_id"`
	FromCompany           *companyRegistration.Company `gorm:"foreignKey:FromCompanyID;references:ID" json:"from_company"`
	ToCompanyID           *uint                        `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"to_company_id"`
	ToCompany             *companyRegistration.Company `gorm:"foreignKey:ToCompanyID;references:ID" json:"to_company"`
	Destination           string                       `json:"destination"`
	Port                  string                       `json:"port"`
	// Uganda Edits
	BrokerName       string                         `json:"broker_name"`
	BrokerNumber     string                         `json:"broker_number"`
	VATTax           float64                        `gorm:"type:numeric;not null" json:"vat_tax"`
	NumberPlate      string                         `json:"number_plate"`
	CarTracker       bool                           `json:"car_tracker"`
	CustomerID       *int                           `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"customer_id"`
	Customer         *customerRegistration.Customer `gorm:"foreignKey:CustomerID" json:"customer"`
	CarStatus        string                         `json:"car_status"`
	CarPaymentStatus string                         `json:"car_payment_status"`
	CreatedBy        string                         `gorm:"size:100" json:"created_by"`
	UpdatedBy        string                         `gorm:"size:100" json:"updated_by"`
}

func (car *Car) BeforeCreate(tx *gorm.DB) (err error) {
	car.CarUUID = uuid.New()
	return
}
