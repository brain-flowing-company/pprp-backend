package enums

type AppointmentStatus string

const (
	PendingAppointment   AppointmentStatus = "PENDING"
	ConfirmedAppointment AppointmentStatus = "CONFIRMED"
	RejectedAppointment  AppointmentStatus = "REJECTED"
	CancelledAppointment AppointmentStatus = "CANCELLED"
	ArchivedAppointment  AppointmentStatus = "ARCHIVED"
)

var AppointmentStatusMap = map[string]AppointmentStatus{
	"PENDING":   PendingAppointment,
	"CONFIRMED": ConfirmedAppointment,
	"REJECTED":  RejectedAppointment,
	"CANCELLED": CancelledAppointment,
	"ARCHIVED":  ArchivedAppointment,
}
