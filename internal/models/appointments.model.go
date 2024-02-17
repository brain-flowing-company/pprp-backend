package models

import (
	"time"

	"github.com/google/uuid"
)

type Appointments struct {
	AppointmentId   uuid.UUID  `json:"appointment_id"   example:"123e4567-e89b-12d3-a456-426614174000"`
	PropertyId      uuid.UUID  `json:"property_id"      example:"123e4567-e89b-12d3-a456-426614174000"`
	OwnerUserId     uuid.UUID  `json:"owner_user_id"    example:"123e4567-e89b-12d3-a456-426614174000"`
	DwellerUserId   uuid.UUID  `json:"dweller_user_id"  example:"123e4567-e89b-12d3-a456-426614174000"`
	AppointmentDate time.Time  `json:"appointment_date" example:"123e4567-e89b-12d3-a456-426614174000"`
	CreatedAt       *time.Time `gorm:"autoCreateTime"`
	UpdatedAt       *time.Time `gorm:"autoUpdateTime"`
	DeletedAt       *time.Time `gorm:"default:null"`
}

func (a Appointments) TableName() string {
	return "appointments"
}
