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
	Status           enums.AppointmentStatus `json:"status"           example:"PENDING"`
	AppointmentDate  time.Time               `json:"appointment_date" example:"2024-02-18T11:00:00Z"`
	Note             string                  `json:"note"             example:"This is a note"`
	CancelledMessage string                  `json:"cancelled_message" example:"This is a cancelled message"`
	CommonModels
}

type CreatingAppointments struct {
	PropertyId       uuid.UUID `json:"property_id"       example:"123e4567-e89b-12d3-a456-426614174000"`
	OwnerUserId      uuid.UUID `json:"owner_user_id"     example:"123e4567-e89b-12d3-a456-426614174000"`
	DwellerUserId    uuid.UUID `json:"dweller_user_id"   example:"123e4567-e89b-12d3-a456-426614174000"`
	AppointmentDate  time.Time `json:"appointment_dates" example:"2024-02-18T11:00:00Z"`
	Note             string    `json:"note"             example:"This is a note"`
	CancelledMessage string    `json:"cancelled_message" example:"This is a cancelled message"`
}

func (a Appointments) TableName() string {
	return "appointments"
}

type UpdatingAppointmentStatus struct {
	Status enums.AppointmentStatus `json:"status"`
}

// Data Structure for Apppointment Lists
type AppointmentLists struct {
	AppointmentId    uuid.UUID                `json:"appointment_id"   example:"123e4567-e89b-12d3-a456-426614174000"`
	Property         PropertyAppointmentLists `json:"property"`
	Owner            UserAppointmentLists     `json:"owner"`
	Status           enums.AppointmentStatus  `json:"status"           example:"PENDING"`
	AppointmentDate  time.Time                `json:"appointment_date" example:"2024-02-18T11:00:00Z"`
	Note             string                   `json:"note"             example:"This is a note"`
	CancelledMessage string                   `json:"cancelled_message" example:"This is a cancelled message"`
	CommonModels
}

type PropertyAppointmentLists struct {
	PropertyId     uuid.UUID           `json:"property_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	PropertyName   string              `json:"property_name" example:"The Base Sukhumvit 77"`
	PropertyType   enums.PropertyTypes `json:"property_type" example:"CONDO"`
	PropertyImages []PropertyImages    `json:"property_images"`
}

type UserAppointmentLists struct {
	UserId          uuid.UUID `json:"user_id"          example:"123e4567-e89b-12d3-a456-426614174000"`
	FirstName       string    `json:"first_name"       example:"John"`
	LastName        string    `json:"last_name"        example:"Doe"`
	ProfileImageUrl string    `json:"profile_image_url" example:"https://image_url.com/abcd"`
}

// Data Structure for Appointment Details
type AppointmentDetails struct {
	AppointmentId    uuid.UUID                  `json:"appointment_id"   example:"123e4567-e89b-12d3-a456-426614174000"`
	Property         PropertyAppointmentDetails `json:"property"`
	Owner            UserAppointmentDetails     `json:"owner"`
	Dweller          UserAppointmentDetails     `json:"dweller"`
	Status           enums.AppointmentStatus    `json:"status"           example:"PENDING"`
	AppointmentDate  time.Time                  `json:"appointment_date" example:"2024-02-18T11:00:00Z"`
	Note             string                     `json:"note"             example:"This is a note"`
	CancelledMessage string                     `json:"cancelled_message" example:"This is a cancelled message"`
	CommonModels
}

type PropertyAppointmentDetails struct {
	PropertyId     uuid.UUID           `json:"property_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	PropertyName   string              `json:"property_name" example:"The Base Sukhumvit 77"`
	PropertyType   enums.PropertyTypes `json:"property_type" example:"CONDO"`
	Address        string              `json:"address" example:"123/4"`
	Alley          string              `json:"alley" example:"Pattaya Nua 78"`
	Street         string              `json:"street" example:"Pattaya"`
	SubDistrict    string              `json:"sub_district" example:"Bang Bon"`
	District       string              `json:"district" example:"Bang Phli"`
	Province       string              `json:"province" example:"Pattaya"`
	Country        string              `json:"country" example:"Thailand"`
	PostalCode     string              `json:"postal_code" example:"69096"`
	PropertyImages []PropertyImages    `json:"property_images"`
	Price          float64             `json:"price" example:"12345.67"`
	PricePerMonth  float64             `json:"price_per_month" example:"12345.67"`
}

type UserAppointmentDetails struct {
	UserId          uuid.UUID `json:"user_id"          example:"123e4567-e89b-12d3-a456-426614174000"`
	FirstName       string    `json:"first_name"       example:"John"`
	LastName        string    `json:"last_name"        example:"Doe"`
	ProfileImageUrl string    `json:"profile_image_url" example:"https://image_url.com/abcd"`
	PhoneNumber     string    `json:"phone_number"     example:"0812345678"`
}
