package models

import (
	"time"

	"github.com/brain-flowing-company/pprp-backend/internal/enums"
	"github.com/google/uuid"
)

type Appointments struct {
	AppointmentId    uuid.UUID               `json:"appointment_id"   example:"123e4567-e89b-12d3-a456-426614174000"`
	PropertyId       uuid.UUID               `json:"property_id"      example:"123e4567-e89b-12d3-a456-426614174000"`
	OwnerUserId      uuid.UUID               `json:"owner_user_id"    example:"123e4567-e89b-12d3-a456-426614174000"`
	DwellerUserId    uuid.UUID               `json:"dweller_user_id"  example:"123e4567-e89b-12d3-a456-426614174000"`
	AppointmentDate  time.Time               `json:"appointment_date" example:"2024-02-18T11:00:00Z"`
	Status           enums.AppointmentStatus `json:"status"           example:"PENDING"`
	Note             string                  `json:"note"             example:"This is a note"`
	CancelledMessage string                  `json:"cancelled_message" example:"This is a cancelled message"`
	CommonModels
}

type CreatingAppointments struct {
	PropertyId      uuid.UUID `json:"property_id"       example:"123e4567-e89b-12d3-a456-426614174000"`
	OwnerUserId     uuid.UUID `json:"owner_user_id"     example:"123e4567-e89b-12d3-a456-426614174000"`
	DwellerUserId   uuid.UUID `json:"dweller_user_id"   example:"123e4567-e89b-12d3-a456-426614174000"`
	AppointmentDate time.Time `json:"appointment_dates" example:"2024-02-18T11:00:00Z"`
	Note            string    `json:"note"             example:"This is a note"`
}

func (a Appointments) TableName() string {
	return "appointments"
}

type UpdatingAppointmentStatus struct {
	Status           enums.AppointmentStatus `json:"status" example:"CANCELLED"`
	CancelledMessage string                  `json:"cancelled_message" example:"This is a cancelled message"`
}

// Data Structure for Apppointment Lists
type AppointmentLists struct {
	AppointmentId    uuid.UUID                `json:"appointment_id"   example:"123e4567-e89b-12d3-a456-426614174000"`
	Property         PropertyAppointmentLists `json:"property" gorm:"foreignKey:AppointmentId; references:AppointmentId; embedded"`
	Owner            OwnerAppointmentLists    `json:"owner" gorm:"foreignKey:AppointmentId; references:AppointmentId; embedded"`
	AppointmentDate  time.Time                `json:"appointment_date" example:"2024-02-18T11:00:00Z"`
	Status           enums.AppointmentStatus  `json:"status"           example:"PENDING"`
	Note             string                   `json:"note"             example:"This is a note"`
	CancelledMessage string                   `json:"cancelled_message" example:"This is a cancelled message"`
	CommonModels
}

type PropertyAppointmentLists struct {
	AppointmentId  uuid.UUID                   `json:"-"`
	PropertyId     uuid.UUID                   `json:"property_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	PropertyName   string                      `json:"property_name" example:"The Base Sukhumvit 77"`
	PropertyType   enums.PropertyTypes         `json:"property_type" example:"CONDO"`
	PropertyImages []PropertyImageAppointments `json:"property_images" gorm:"foreignKey:AppointmentId; references:AppointmentId"`
}

type PropertyImageAppointments struct {
	AppointmentId uuid.UUID `json:"-"`
	PropertyId    uuid.UUID `json:"-"`
	ImageUrl      string    `json:"image_url" example:"https://image_url.com/abcd"`
}

type OwnerAppointmentLists struct {
	AppointmentId        uuid.UUID `json:"-"`
	OwnerUserId          uuid.UUID `json:"owner_user_id"          example:"123e4567-e89b-12d3-a456-426614174000"`
	OwnerFirstName       string    `json:"owner_first_name"       example:"John"`
	OwnerLastName        string    `json:"owner_last_name"        example:"Doe"`
	OwnerProfileImageUrl string    `json:"owner_profile_image_url" example:"https://image_url.com/abcd"`
}

// Data Structure for Appointment Details
type AppointmentDetails struct {
	AppointmentId    uuid.UUID                  `json:"appointment_id"   example:"123e4567-e89b-12d3-a456-426614174000"`
	Property         PropertyAppointmentDetails `json:"property" gorm:"foreignKey:AppointmentId; references:AppointmentId; embedded"`
	Owner            OwnerAppointmentDetails    `json:"owner" gorm:"foreignKey:AppointmentId; references:AppointmentId; embedded"`
	Dweller          DwellerAppointmentDetails  `json:"dweller" gorm:"foreignKey:AppointmentId; references:AppointmentId; embedded"`
	AppointmentDate  time.Time                  `json:"appointment_date" example:"2024-02-18T11:00:00Z"`
	Status           enums.AppointmentStatus    `json:"status"           example:"PENDING"`
	Note             string                     `json:"note"             example:"This is a note"`
	CancelledMessage string                     `json:"cancelled_message" example:"This is a cancelled message"`
	CommonModels
}

type PropertyAppointmentDetails struct {
	AppointmentId  uuid.UUID                   `json:"-"`
	PropertyId     uuid.UUID                   `json:"property_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	PropertyName   string                      `json:"property_name" example:"The Base Sukhumvit 77"`
	PropertyType   enums.PropertyTypes         `json:"property_type" example:"CONDO"`
	Address        string                      `json:"address" example:"123/4"`
	Alley          string                      `json:"alley" example:"Pattaya Nua 78"`
	Street         string                      `json:"street" example:"Pattaya"`
	SubDistrict    string                      `json:"sub_district" example:"Bang Bon"`
	District       string                      `json:"district" example:"Bang Phli"`
	Province       string                      `json:"province" example:"Pattaya"`
	Country        string                      `json:"country" example:"Thailand"`
	PostalCode     string                      `json:"postal_code" example:"69096"`
	PropertyImages []PropertyImageAppointments `json:"property_images" gorm:"foreignKey:PropertyId; references:PropertyId"`
	Price          float64                     `json:"price" example:"12345.67"`
	PricePerMonth  float64                     `json:"price_per_month" example:"12345.67"`
}

type OwnerAppointmentDetails struct {
	AppointmentId        uuid.UUID `json:"-"`
	OwnerUserId          uuid.UUID `json:"owner_user_id"          example:"123e4567-e89b-12d3-a456-426614174000"`
	OwnerFirstName       string    `json:"owner_first_name"       example:"John"`
	OwnerLastName        string    `json:"owner_last_name"        example:"Doe"`
	OwnerProfileImageUrl string    `json:"owner_profile_image_url" example:"https://image_url.com/abcd"`
	OwnerPhoneNumber     string    `json:"owner_phone_number"     example:"0812345678"`
}

type DwellerAppointmentDetails struct {
	AppointmentId          uuid.UUID `json:"-"`
	DwellerUserId          uuid.UUID `json:"dweller_user_id"          example:"123e4567-e89b-12d3-a456-426614174000"`
	DwellerFirstName       string    `json:"dweller_first_name"       example:"John"`
	DwellerLastName        string    `json:"dweller_last_name"        example:"Doe"`
	DwellerProfileImageUrl string    `json:"dweller_profile_image_url" example:"https://image_url.com/abcd"`
	DwellerPhoneNumber     string    `json:"dweller_phone_number"     example:"0812345678"`
}
