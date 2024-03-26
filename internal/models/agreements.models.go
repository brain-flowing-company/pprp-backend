package models

import (
	"time"

	"github.com/brain-flowing-company/pprp-backend/internal/enums"
	"github.com/google/uuid"
)

type Agreements struct {
	AgreementId   uuid.UUID `json:"agreement_id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	AgreementType enums.AgreementTypes `json:"agreement_type" example:"SELLING"`
	PropertyId    uuid.UUID `json:"property_id" example:"00000000-0000-0000-0000-000000000000"`
	OwnerUserId   uuid.UUID `json:"owner_user_id" example:"00000000-0000-0000-0000-000000000000"`
	DwellerUserId uuid.UUID `json:"dweller_user_id" example:"00000000-0000-0000-0000-000000000000"`
	AgreementDate time.Time `json:"agreement_date" example:"2021-01-01T00:00:00Z"`
	Status enums.AgreementStatus `json:"status" example:"AWAITING_DEPOSIT"`
	DepositAmount float64 `json:"deposit_amount" example:"1000000"`
	PaymentPerMonth float64 `json:"payment_per_month" example:"1000000"`
	PaymentDuration int `json:"payment_duration" example:"12"`
	TotalPayment float64 `json:"total_payment" example:"12000000"`
	CancelledMessage string `json:"cancelled_message" example:"This is cancelled message."`
	CommonModels
}

type CreatingAgreements struct {
	AgreementType enums.AgreementTypes `json:"agreement_type"`
	PropertyId    uuid.UUID `json:"property_id"`
	OwnerUserId   uuid.UUID `json:"owner_user_id"`
	DwellerUserId uuid.UUID `json:"dweller_user_id"`
	AgreementDate time.Time `json:"agreement_date"`
	Status enums.AgreementStatus `json:"status"`
	DepositAmount float64 `json:"deposit_amount"`
	PaymentPerMonth float64 `json:"payment_per_month"`
	PaymentDuration int `json:"payment_duration"`
	TotalPayment float64 `json:"total_payment"`
}

func (a Agreements) TableName() string {
	return "agreements"
}

type UpdatingAgreementStatus struct {
	Status enums.AgreementStatus `json:"status"`
	CancelledMessage string `json:"cancelled_message"`
}

type MyAgreementResponses struct {
	OwnerAgreements   []AgreementLists `json:"owner_agreements"`
	DwellerAgreements []AgreementLists `json:"dweller_agreements"`
}

// Data Structure for Agreement Lists
type AgreementLists struct {
	AgreementId     uuid.UUID `json:"agreement_id" example:"00000000-0000-0000-0000-000000000000"`
	AgreementType   enums.AgreementTypes `json:"agreement_type" example:"SELLING"`
	Property        PropertyAgreementLists `json:"property" gorm:"foreignKey:AgreementId; references:AgreementId; embedded"`
	Owner           OwnerAgreementLists `json:"owner" gorm:"foreignKey:AgreementId; references:AgreementId; embedded"`
	AgreementDate   time.Time `json:"agreement_date" example:"2021-01-01T00:00:00Z"`
	Status          enums.AgreementStatus `json:"status" example:"AWAITING_DEPOSIT"`
	CancelledMessage string `json:"cancelled_message" example:"This is cancelled message."`
	CommonModels
}

type PropertyAgreementLists struct {
	AgreementId     uuid.UUID `json:"-"`
	PropertyId      uuid.UUID `json:"property_id" example:"00000000-0000-0000-0000-000000000000"`
	PropertyName    string `json:"property_name" example:"The Base Sukhumvit 77"`
	PropertyType    enums.PropertyTypes `json:"property_type" example:"CONDO"`
	PropertyImages  []PropertyImageAgreements `json:"property_images" gorm:"foreignKey:AgreementId; references:AgreementId"`
}

type PropertyImageAgreements struct {
	AgreementId uuid.UUID `json:"-"`
	PropertyId  uuid.UUID `json:"-"`
	ImageUrl    string `json:"image_url" example:"https://www.example.com/image.jpg"`
}

type OwnerAgreementLists struct {
	AgreementId uuid.UUID `json:"-"`
	OwnerUserId uuid.UUID `json:"owner_user_id" example:"00000000-0000-0000-0000-000000000000"`
	OwnerFirstName string `json:"owner_first_name" example:"John"`
	OwnerLastName string `json:"owner_last_name" example:"Doe"`
	OwnerProfileImageUrl string `json:"owner_profile_image_url" example:"https://www.example.com/image.jpg"`
}

// Data Structure for Agreement Details
type AgreementDetails struct {
	AgreementId     uuid.UUID `json:"agreement_id" example:"00000000-0000-0000-0000-000000000000"`
	AgreementType   enums.AgreementTypes `json:"agreement_type" example:"SELLING"`
	Property        PropertyAgreementDetails `json:"property" gorm:"foreignKey:AgreementId; references:AgreementId; embedded"`
	Owner           OwnerAgreementDetails `json:"owner" gorm:"foreignKey:AgreementId; references:AgreementId; embedded"`
	Dweller         DwellerAgreementDetails `json:"dweller" gorm:"foreignKey:AgreementId; references:AgreementId; embedded"`
	AgreementDate   time.Time `json:"agreement_date" example:"2021-01-01T00:00:00Z"`
	Status          enums.AgreementStatus `json:"status" example:"AWAITING_DEPOSIT"`
	DepositAmount   float64 `json:"deposit_amount" example:"1000000"`
	PaymentPerMonth float64 `json:"payment_per_month" example:"1000000"`
	PaymentDuration int `json:"payment_duration" example:"12"`
	TotalPayment    float64 `json:"total_payment" example:"12000000"`
	CancelledMessage string `json:"cancelled_message" example:"This is cancelled message."`
	CommonModels
}

type PropertyAgreementDetails struct {
	AgreementId  uuid.UUID `json:"-"`
	PropertyId   uuid.UUID `json:"property_id" example:"00000000-0000-0000-0000-000000000000"`
	PropertyName string `json:"property_name" example:"The Base Sukhumvit 77"`
	PropertyType enums.PropertyTypes `json:"property_type" example:"CONDO"`
	Address      string `json:"address" example:"123/456 The Base Sukhumvit 77"`
	Street       string `json:"street" example:"Sukhumvit"`
	Alley        string `json:"alley" example:"77"`
	SubDistrict  string `json:"sub_district" example:"Phra Khanong"`
	District	 string `json:"district" example:"Bangkok"`
	Province	 string `json:"province" example:"Bangkok"`
	Country		 string `json:"country" example:"Thailand"`
	PostalCode	 string `json:"postal_code" example:"10110"`
	PropertyImages []PropertyImageAgreements `json:"property_images" gorm:"foreignKey:PropertyId; references:PropertyId"`
	Price        float64 `json:"price" example:"1000000"`
	PricePerMonth float64 `json:"price_per_month" example:"1000000"`
}

type OwnerAgreementDetails struct {
	AgreementId        uuid.UUID `json:"-"`
	OwnerUserId        uuid.UUID `json:"owner_user_id" example:"00000000-0000-0000-0000-000000000000"`
	OwnerFirstName     string `json:"owner_first_name" example:"John"`
	OwnerLastName      string `json:"owner_last_name" example:"Doe"`
	OwnerProfileImageUrl string `json:"owner_profile_image_url" example:"https://www.example.com/image.jpg"`
	OwnerPhoneNumber   string `json:"owner_phone_number" example:"0812345678"`
}

type DwellerAgreementDetails struct {
	AgreementId        uuid.UUID `json:"-"`
	DwellerUserId      uuid.UUID `json:"dweller_user_id" example:"00000000-0000-0000-0000-000000000000"`
	DwellerFirstName   string `json:"dweller_first_name" example:"John"`
	DwellerLastName    string `json:"dweller_last_name" example:"Doe"`
	DwellerProfileImageUrl string `json:"dweller_profile_image_url" example:"https://www.example.com/image.jpg"`
	DwellerPhoneNumber string `json:"dweller_phone_number" example:"0812345678"`
}