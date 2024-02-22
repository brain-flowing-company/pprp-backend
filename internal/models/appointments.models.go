package models

import (
	"time"

	"github.com/google/uuid"
)

type Appointments struct {
	AppointmentId      uuid.UUID          `json:"appointment_id"   example:"123e4567-e89b-12d3-a456-426614174000"`
	PropertyId         uuid.UUID          `json:"property_id"      example:"123e4567-e89b-12d3-a456-426614174000"`
	OwnerUserId        uuid.UUID          `json:"owner_user_id"    example:"123e4567-e89b-12d3-a456-426614174000"`
	DwellerUserId      uuid.UUID          `json:"dweller_user_id"  example:"123e4567-e89b-12d3-a456-426614174000"`
	AppointmentDate    time.Time          `json:"appointment_date" example:"2024-02-18T11:00:00Z"`
	AppointmentsStatus AppointmentsStatus `json:"status"           example:"PENDING"`
	CommonModels
}

type CreatingAppointments struct {
	PropertyId       uuid.UUID   `json:"property_id"       example:"123e4567-e89b-12d3-a456-426614174000"`
	OwnerUserId      uuid.UUID   `json:"owner_user_id"     example:"123e4567-e89b-12d3-a456-426614174000"`
	DwellerUserId    uuid.UUID   `json:"dweller_user_id"   example:"123e4567-e89b-12d3-a456-426614174000"`
	AppointmentDates []time.Time `json:"appointment_dates" example:"2024-02-18T11:00:00Z"`
}

type DeletingAppointments struct {
	AppointmentIds []string `json:"appointmentIds" example:"123e4567-e89b-12d3-a456-426614174000"`
}

func (a Appointments) TableName() string {
	return "appointments"
}

type AppointmentsStatus string

const (
	Pending       AppointmentsStatus = "PENDING"
	Approved      AppointmentsStatus = "APPROVED"
	Rejected      AppointmentsStatus = "REJECTED"
	RequestChange AppointmentsStatus = "REQUEST_CHANGE"
	Cancelled     AppointmentsStatus = "CANCELLED"
	Completed     AppointmentsStatus = "COMPLETED"
)

var AppointmentStatusMap = map[string]AppointmentsStatus{
	"PENDING":        Pending,
	"APPROVED":       Approved,
	"REJECTED":       Rejected,
	"REQUEST_CHANGE": RequestChange,
	"CANCELLED":      Cancelled,
	"COMPLETED":      Completed,
}

type UpdatingAppointmentStatus struct {
	Status AppointmentsStatus `json:"status"`
}
