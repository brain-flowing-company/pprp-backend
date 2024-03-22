package appointments

import (
	"database/sql"

	"github.com/brain-flowing-company/pprp-backend/internal/enums"
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"gorm.io/gorm"
)

type Repository interface {
	GetAllAppointments(*[]models.AppointmentLists) error
	GetAppointmentById(*models.AppointmentDetails, string) error
	GetAppointmentByOwnerId([]*models.Appointments, string) error
	GetAppointmentByDwellerId([]*models.Appointments, string) error
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

func (repo *repositoryImpl) GetAllAppointments(appointments *[]models.AppointmentLists) error {
	propertiesQuery := `SELECT property_id, property_name, property_type FROM properties`
	ownersQuery := `SELECT user_id AS owner_user_id,
						first_name AS owner_first_name,
						last_name AS owner_last_name,
						profile_image_url AS owner_profile_image_url
					FROM users`
	return repo.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.Appointments{}).
			Raw(`
				SELECT a.appointment_id,
					   p.*,
					   o.*,
					   a.appointment_date,
					   a.status,
					   a.note,
					   a.cancelled_message,
					   a.created_at
					FROM appointments a
					JOIN (` + propertiesQuery + `) AS p 
						ON a.property_id = p.property_id
					JOIN (` + ownersQuery + `) AS o 
						ON a.owner_user_id = o.owner_user_id
					`).
			Scan(appointments).Error; err != nil {
			return err
		}

		for i, appointment := range *appointments {
			if err := repo.db.Model(&models.PropertyImages{}).
				Raw(`
					SELECT image_url
					FROM property_images
					WHERE property_id = @property_id
					`, sql.Named("property_id", appointment.Property.PropertyId)).
				Pluck("image_url", &(*appointments)[i].Property.PropertyImages).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (repo *repositoryImpl) GetAppointmentById(appointment *models.AppointmentDetails, appointmentId string) error {
	return repo.db.Transaction(func(tx *gorm.DB) error {
		// if err := tx.Model(&models.Appointments{}).
		// 	Raw(`
		// 		SELECT appointments
		// 	`)
		return nil
	})
}

func (repo *repositoryImpl) GetAppointmentByOwnerId(appointments []*models.Appointments, ownerUserId string) error {
	return repo.db.Model(&models.Appointments{}).Find(appointments, "owner_user_id = ?", ownerUserId).Error
}

func (repo *repositoryImpl) GetAppointmentByDwellerId(appointments []*models.Appointments, dwellerUserId string) error {
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
