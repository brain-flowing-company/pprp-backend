package appointments

import (
	"database/sql"

	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"gorm.io/gorm"
)

type Repository interface {
	GetAllAppointments(*[]models.AppointmentLists) error
	GetAppointmentById(*models.AppointmentDetails, string) error
	GetAppointmentByUserId(*models.MyAppointmentResponses, *models.MyAppointmentRequests) error
	CreateAppointment(*models.CreatingAppointments) error
	DeleteAppointment(string) error
	UpdateAppointmentStatus(*models.UpdatingAppointmentStatus, string) error
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
	propertyQuery := `SELECT p.property_id, p.property_name, p.property_type, p.address,
							p.alley, p.street, p.sub_district, p.district, p.province, p.country,
							p.postal_code, s.price, r.price_per_month 
						FROM properties p
						LEFT JOIN selling_properties s ON p.property_id = s.property_id
						LEFT JOIN renting_properties r ON p.property_id = r.property_id
						WHERE p.property_id = @property_id`

	ownerQuery := `SELECT user_id AS owner_user_id,
						first_name AS owner_first_name,
						last_name AS owner_last_name,
						profile_image_url AS owner_profile_image_url,
						phone_number AS owner_phone_number
					FROM users
					WHERE user_id = @owner_user_id`

	dwellerQuery := `SELECT user_id AS dweller_user_id,
						first_name AS dweller_first_name,
						last_name AS dweller_last_name,
						profile_image_url AS dweller_profile_image_url,
						phone_number AS dweller_phone_number
					FROM users
					WHERE user_id = @dweller_user_id`

	return repo.db.Transaction(func(tx *gorm.DB) error {
		var existingAppointment models.Appointments
		if err := tx.Model(&models.Appointments{}).First(&existingAppointment, "appointment_id = ?", appointmentId).Error; err != nil {
			return err
		}

		if err := tx.Model(&models.Appointments{}).
			Raw(`
				SELECT a.appointment_id, 
					   p.*, 
					   o.*, 
					   d.*,
					   a.appointment_date, 
					   a.status,
					   a.note,
					   a.cancelled_message,
					   a.created_at
					FROM appointments a
					JOIN (`+propertyQuery+`) p ON a.property_id = p.property_id
					JOIN (`+ownerQuery+`) o ON a.owner_user_id = o.owner_user_id
					JOIN (`+dwellerQuery+`) d ON a.dweller_user_id = d.dweller_user_id
			`, sql.Named("property_id", existingAppointment.PropertyId),
				sql.Named("owner_user_id", existingAppointment.OwnerUserId),
				sql.Named("dweller_user_id", existingAppointment.DwellerUserId)).
			Scan(appointment).Error; err != nil {
			return err
		}

		if err := tx.Model(&models.PropertyImages{}).
			Raw(`
				SELECT image_url
				FROM property_images
				WHERE property_id = @property_id
				`, sql.Named("property_id", appointment.Property.PropertyId)).
			Pluck("image_url", &appointment.Property.PropertyImages).Error; err != nil {
			return err
		}

		return nil
	})
}

func (repo *repositoryImpl) GetAppointmentByUserId(appointmentResponse *models.MyAppointmentResponses, appointmentRequest *models.MyAppointmentRequests) error {
	propertiesQuery := `SELECT property_id, property_name, property_type FROM properties`

	ownersQuery := `SELECT user_id AS owner_user_id,
						first_name AS owner_first_name,
						last_name AS owner_last_name,
						profile_image_url AS owner_profile_image_url
					FROM users`

	appointmentListsQuery := `SELECT a.appointment_id,
								p.*,
								o.*,
								a.appointment_date,
								a.status,
								a.note,
								a.cancelled_message,
								a.created_at
							FROM appointments a
							JOIN (` + propertiesQuery + `) AS p ON a.property_id = p.property_id
							JOIN (` + ownersQuery + `) AS o ON a.owner_user_id = o.owner_user_id`

	return repo.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.Appointments{}).
			Raw(appointmentListsQuery+`
				WHERE a.owner_user_id = @userId
				ORDER BY a.appointment_date `+appointmentRequest.Order+`
				`, sql.Named("userId", appointmentRequest.UserId)).
			Scan(&appointmentResponse.OwnerAppointments).Error; err != nil {
			return err
		}

		if err := tx.Model(&models.Appointments{}).
			Raw(appointmentListsQuery+`
				WHERE a.dweller_user_id = @userId
				ORDER BY a.appointment_date `+appointmentRequest.Order+`
				`, sql.Named("userId", appointmentRequest.UserId)).
			Scan(&appointmentResponse.DwellerAppointments).Error; err != nil {
			return err
		}

		for i, appointment := range appointmentResponse.OwnerAppointments {
			if err := repo.db.Model(&models.PropertyImages{}).
				Raw(`
					SELECT image_url
					FROM property_images
					WHERE property_id = @property_id
					`, sql.Named("property_id", appointment.Property.PropertyId)).
				Pluck("image_url", &appointmentResponse.OwnerAppointments[i].Property.PropertyImages).Error; err != nil {
				return err
			}
		}

		for i, appointment := range appointmentResponse.DwellerAppointments {
			if err := repo.db.Model(&models.PropertyImages{}).
				Raw(`
					SELECT image_url
					FROM property_images
					WHERE property_id = @property_id
					`, sql.Named("property_id", appointment.Property.PropertyId)).
				Pluck("image_url", &appointmentResponse.DwellerAppointments[i].Property.PropertyImages).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (repo *repositoryImpl) CreateAppointment(appointment *models.CreatingAppointments) error {
	return repo.db.Exec(`INSERT INTO appointments (appointment_id, property_id, owner_user_id, dweller_user_id, appointment_date, note) VALUES (?, ?, ?, ?, ?, ?)`,
		appointment.AppointmentId, appointment.PropertyId, appointment.OwnerUserId, appointment.DwellerUserId, appointment.AppointmentDate, appointment.Note).Error
}

func (repo *repositoryImpl) DeleteAppointment(appointmentId string) error {
	if err := repo.db.Model(&models.Appointments{}).First(&models.Agreements{}, "appointment_id = ?", appointmentId).Error; err != nil {
		return err
	}

	return repo.db.Where("appointment_id = ?", appointmentId).Delete(&models.Appointments{}).Error
}

func (repo *repositoryImpl) UpdateAppointmentStatus(updatingAppointment *models.UpdatingAppointmentStatus, appointmentId string) error {
	if err := repo.db.Model(&models.Appointments{}).First(&models.Appointments{}, "appointment_id = ?", appointmentId).Error; err != nil {
		return err
	}

	return repo.db.Model(&models.Appointments{}).Where("appointment_id = ?", appointmentId).Updates(updatingAppointment).Error
}
