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
	ChasisNumber          string                       `gorm:"size:100;not null;unique" json:"chasis_number"`
	EngineNumber          string                       `gorm:"size:100;" json:"engine_number"`
	EngineCapacity        string                       `gorm:"size:100;" json:"engine_capacity"`
	FrameNumber           string                       `gorm:"size:100;" json:"frame_number"`
	Make                  string                       `gorm:"size:100;" json:"make"`
	CarModel              string                       `gorm:"size:100;" json:"car_model"`
	SeatCapacity          string                       `gorm:"size:50;" json:"seat_capacity"`
	MaximCarry            float64                      `gorm:"size:100;" json:"maxim_carry"`
	Weight                float64                      `gorm:"size:100;" json:"weight"`
	GrossWeight           float64                      `gorm:"size:100;" json:"gross_weight"`
	Length                float64                      `gorm:"size:100;" json:"length"`
	Width                 float64                      `gorm:"size:100;" json:"width"`
	Height                float64                      `gorm:"size:100;" json:"height"`
	ManufactureYear       int                          `json:"manufacture_year"`
	FirstRegistrationYear int                          `json:"first_registration_year"`
	Transmission          string                       `json:"transmission"`
	BodyType              string                       `json:"body_type"`
	Colour                string                       `json:"colour"`
	Auction               string                       `json:"auction"`
	Currency              string                       `json:"currency"`
	CarMillage            int                          `json:"car_millage"`
	FuelConsumption       string                       `json:"fuel_consumption"`
	PowerSteering         bool                         `gorm:"default:false" json:"power_steering"`
	PowerWindow           bool                         `gorm:"default:false" json:"power_window"`
	ABS                   bool                         `gorm:"default:false" json:"abs"`
	ADS                   bool                         `gorm:"default:false" json:"ads"`
	AirBrake              bool                         `gorm:"default:false" json:"air_brake"`
	OilBrake              bool                         `gorm:"default:false" json:"oil_brake"`
	AlloyWheel            bool                         `gorm:"default:false" json:"alloy_wheel"`
	SimpleWheel           bool                         `gorm:"default:false" json:"simple_wheel"`
	Navigation            bool                         `gorm:"default:false" json:"navigation"`
	AC                    bool                         `gorm:"default:false" json:"ac"`
	BidPrice              float64                      `gorm:"type:numeric;" json:"bid_price"`
	VATTax                float64                      `gorm:"type:numeric;" json:"vat_tax"`
	PurchaseDate          string                       `gorm:"type:date;not null" json:"purchase_date"`
	FromCompanyID         *uint                        `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"from_company_id"`
	FromCompany           *companyRegistration.Company `gorm:"foreignKey:FromCompanyID;references:ID" json:"from_company"`
	ToCompanyID           *uint                        `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"to_company_id"`
	ToCompany             *companyRegistration.Company `gorm:"foreignKey:ToCompanyID;references:ID" json:"to_company"`
	OtherEntity           string                       `json:"other_entity"`
	Destination           string                       `json:"destination"`
	Port                  string                       `json:"port"`
	CarShippingInvoiceID  *uint                        `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"car_shipping_invoice_id"`
	CarShippingInvoice    *CarShippingInvoice          `gorm:"foreignKey:CarShippingInvoiceID" json:"car_shipping_invoice"`
	CarStatusJapan        string                       `json:"car_status_japan"` // InStock, Sold, Exported
	// Uganda Edits
	BrokerName       string                         `json:"broker_name"`
	BrokerNumber     string                         `json:"broker_number"`
	NumberPlate      string                         `json:"number_plate"`
	CarTracker       bool                           `json:"car_tracker"`
	CustomerID       *int                           `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"customer_id"`
	Customer         *customerRegistration.Customer `gorm:"foreignKey:CustomerID" json:"customer"`
	CarStatus        string                         `json:"car_status"`         // InTransit, InStock, Sold
	CarPaymentStatus string                         `json:"car_payment_status"` // Fully Payed, Partially Paid, Booked
	CreatedBy        string                         `gorm:"size:100" json:"created_by"`
	UpdatedBy        string                         `gorm:"size:100" json:"updated_by"`

	// Car Photos
	CarPhotos []CarPhoto `gorm:"foreignKey:CarID;constraint:OnDelete:CASCADE;" json:"car_photos"`
}

// CarPhoto struct
type CarPhoto struct {
	gorm.Model
	CarID uint   `gorm:"not null" json:"car_id"`
	URL   string `gorm:"not null" json:"url"`
}

func (car *Car) BeforeCreate(tx *gorm.DB) (err error) {
	car.CarUUID = uuid.New()
	// If ToCompanyID is not nil and not zero, mark as Exported; otherwise, InStock.
	if car.ToCompanyID != nil && *car.ToCompanyID > 0 {
		car.CarStatusJapan = "Exported"
		car.CarStatus = "InTransit"
	} else {
		car.CarStatusJapan = "InStock"
	}

	// Set OtherEntity to ToCompany.Name if it's empty and ToCompanyID is valid
	if car.OtherEntity == "" && car.ToCompanyID != nil && *car.ToCompanyID > 0 {
		var company companyRegistration.Company
		if err := tx.First(&company, *car.ToCompanyID).Error; err == nil {
			car.OtherEntity = company.Name
		}
	}

	return
}
