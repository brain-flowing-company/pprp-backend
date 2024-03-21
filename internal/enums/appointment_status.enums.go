package enums

type AppointmentStatus string

const (
	Pending   AppointmentStatus = "PENDING"
	Confirmed AppointmentStatus = "CONFIRMED"
	Rejected  AppointmentStatus = "REJECTED"
	Cancelled AppointmentStatus = "CANCELLED"
	Archived  AppointmentStatus = "ARCHIVED"
)

var AppointmentStatusMap = map[string]AppointmentStatus{
	"PENDING":   Pending,
	"CONFIRMED": Confirmed,
	"REJECTED":  Rejected,
	"CANCELLED": Cancelled,
	"ARCHIVED":  Archived,
}
