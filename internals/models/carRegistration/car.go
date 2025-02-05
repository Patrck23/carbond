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
	FuelConsumption       string                       `json:"fuel_consumption"`
	PowerSteering         bool                         `gorm:"default:false" json:"ps"`
	PowerWindow           bool                         `gorm:"default:false" json:"pw"`
	ABS                   bool                         `gorm:"default:false" json:"abs"`
	ADS                   bool                         `gorm:"default:false" json:"ads"`
	AlloyWheel            bool                         `gorm:"default:false" json:"aw"`
	SimpleWheel           bool                         `gorm:"default:false" json:"sw"`
	Navigation            bool                         `gorm:"default:false" json:"navigation"`
	AC                    bool                         `gorm:"default:false" json:"ac"`
	BidPrice              float64                      `gorm:"type:numeric;not null" json:"bid_price"`
	VATTax                float64                      `gorm:"type:numeric;not null" json:"vat_tax"`
	DollarRate            float64                      `json:"dollar_rate"`
	PurchaseDate          string                       `gorm:"type:date;not null" json:"purchase_date"`
	FromCompanyID         *uint                        `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"from_company_id"`
	FromCompany           *companyRegistration.Company `gorm:"foreignKey:FromCompanyID;references:ID" json:"from_company"`
	ToCompanyID           *uint                        `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"to_company_id"`
	ToCompany             *companyRegistration.Company `gorm:"foreignKey:ToCompanyID;references:ID" json:"to_company"`
	Destination           string                       `json:"destination"`
	Port                  string                       `json:"port"`
	CarShippingInvoiceID  *uint                        `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"car_shipping_invoice_id"`
	CarShippingInvoice    *CarShippingInvoice          `gorm:"foreignKey:CarShippingInvoiceID" json:"car_shipping_invoice"`
	// Uganda Edits
	BrokerName       string                         `json:"broker_name"`
	BrokerNumber     string                         `json:"broker_number"`
	NumberPlate      string                         `json:"number_plate"`
	CarTracker       bool                           `json:"car_tracker"`
	CustomerID       *int                           `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"customer_id"`
	Customer         *customerRegistration.Customer `gorm:"foreignKey:CustomerID" json:"customer"`
	CarStatus        string                         `json:"car_status"`
	CarPaymentStatus string                         `json:"car_payment_status"`
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
	return
}
