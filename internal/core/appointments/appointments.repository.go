package appointments

import (
	"github.com/brain-flowing-company/pprp-backend/internal/enums"
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"gorm.io/gorm"
)

type Repository interface {
	GetAllAppointments(*[]models.Appointments) error
	GetAppointmentById(*models.Appointments, string) error
	GetAppointmentByOwnerId(*[]models.Appointments, string) error
	GetAppointmentByDwellerId(*[]models.Appointments, string) error
	CreateAppointment(*models.Appointments) error
	DeleteAppointment(string) error
	UpdateAppointmentStatus(string, enums.AppointmentStatus) error
}

type repositoryImpl struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repositoryImpl{
		db,
	}
}

func (repo *repositoryImpl) GetAllAppointments(appointments *[]models.Appointments) error {
	return repo.db.Model(&models.Appointments{}).Find(appointments).Error
}

func (repo *repositoryImpl) GetAppointmentById(appointment *models.Appointments, appointmentId string) error {
	return repo.db.Model(&models.Appointments{}).First(appointment, "appointment_id = ?", appointmentId).Error
}

func (repo *repositoryImpl) GetAppointmentByOwnerId(appointments *[]models.Appointments, ownerUserId string) error {
	return repo.db.Model(&models.Appointments{}).Find(appointments, "owner_user_id = ?", ownerUserId).Error
}

func (repo *repositoryImpl) GetAppointmentByDwellerId(appointments *[]models.Appointments, dwellerUserId string) error {
	return repo.db.Model(&models.Appointments{}).Find(appointments, "dweller_user_id = ?", dwellerUserId).Error
}

func (repo *repositoryImpl) CreateAppointment(appointment *models.Appointments) error {
	return repo.db.Create(appointment).Error
}

func (repo *repositoryImpl) DeleteAppointment(appointmentId string) error {
	if err := repo.db.Model(&models.Appointments{}).First("appointment_id = ?", appointmentId).Error; err != nil {
		return err
	}

	return repo.db.Where("appointment_id = ?", appointmentId).Delete(&models.Appointments{}).Error
}

func (repo *repositoryImpl) UpdateAppointmentStatus(appointmentId string, status enums.AppointmentStatus) error {
	return repo.db.Model(&models.Appointments{}).Where("appointment_id = ?", appointmentId).Update("appointments_status", status).Error
}
