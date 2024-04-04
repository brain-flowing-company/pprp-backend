package agreements

import (
	"database/sql"

	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"gorm.io/gorm"
)

type Repository interface {
	GetAllAgreements(*[]models.AgreementLists) error
	GetAgreementById(*models.AgreementDetails, string) error
	GetAgreementByUserId(*models.MyAgreementResponses, *models.MyAgreementRequests) error
	CreateAgreement(*models.CreatingAgreements) error
	DeleteAgreement(string) error
	UpdateAgreementStatus(*models.UpdatingAgreementStatus, string) error
}

type repositoryImpl struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repositoryImpl{
		db,
	}
}

func (repo *repositoryImpl) GetAllAgreements(agreements *[]models.AgreementLists) error {
	propertiesQuery := `SELECT property_id, property_name, property_type FROM properties`

	ownersQuery := `SELECT user_id AS owner_user_id,
						first_name AS owner_first_name,
						last_name AS owner_last_name,
						profile_image_url AS owner_profile_image_url
					FROM users`

	return repo.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.Agreements{}).
			Raw(`
				SELECT a.agreement_id,
					   p.*,
					   o.*,
					   a.agreement_date, 
					   a.status, 
					   a.deposit_amount,
					   a.payment_per_month,
					   a.payment_duration,
					   a.total_payment,
					   a.cancelled_message,
					   a.created_at
					FROM agreements a
					JOIN (` + propertiesQuery + `) AS p ON a.property_id = p.property_id
					JOIN (` + ownersQuery + `) AS o ON a.owner_user_id = o.owner_user_id
					`).
			Scan(agreements).Error; err != nil {
			return err
		}

		for i, agreement := range *agreements {
			if err := repo.db.Model(&models.PropertyImages{}).
				Raw(`
					SELECT image_url
					FROM property_images
					WHERE property_id = @property_id
					`, sql.Named("property_id", agreement.Property.PropertyId)).
				Pluck("image_url", &(*agreements)[i].Property.PropertyImages).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (repo *repositoryImpl) GetAgreementById(agreement *models.AgreementDetails, agreementId string) error {
	propertyQuery := `SELECT p.property_id, p.property_name, p.property_type, p.address,
							p.alley, p.street, p.sub_district, p.district, p.province, 
							p.country, p.postal_code, s.price, r.price_per_month 
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
		var existingAgreement models.Agreements
		if err := tx.Model(&models.Agreements{}).First(&existingAgreement, "agreement_id = ?", agreementId).Error; err != nil {
			return err
		}

		if err := tx.Model(&models.Agreements{}).
			Raw(`
				SELECT a.agreement_id, 
					   a.agreement_type,
					   p.*, 
					   o.*, 
					   d.*,
					   a.agreement_date, 
					   a.status,
					   a.deposit_amount,
					   a.payment_per_month,
					   a.payment_duration,
					   a.total_payment,
					   a.cancelled_message,
					   a.created_at
					FROM agreements a
					JOIN (`+propertyQuery+`) p ON a.property_id = p.property_id
					JOIN (`+ownerQuery+`) o ON a.owner_user_id = o.owner_user_id
					JOIN (`+dwellerQuery+`) d ON a.dweller_user_id = d.dweller_user_id
			`, sql.Named("property_id", existingAgreement.PropertyId),
				sql.Named("owner_user_id", existingAgreement.OwnerUserId),
				sql.Named("dweller_user_id", existingAgreement.DwellerUserId)).
			Scan(agreement).Error; err != nil {
			return err
		}

		if err := tx.Model(&models.PropertyImages{}).
			Raw(`
				SELECT image_url
				FROM property_images
				WHERE property_id = @property_id
				`, sql.Named("property_id", agreement.Property.PropertyId)).
			Pluck("image_url", &agreement.Property.PropertyImages).Error; err != nil {
			return err
		}

		return nil
	})
}

func (repo *repositoryImpl) GetAgreementByUserId(agreementResponse *models.MyAgreementResponses, agreementRequest *models.MyAgreementRequests) error {
	propertiesQuery := `SELECT property_id, property_name, property_type FROM properties`

	ownersQuery := `SELECT user_id AS owner_user_id,
						first_name AS owner_first_name,
						last_name AS owner_last_name,
						profile_image_url AS owner_profile_image_url
					FROM users`

	agreementListsQuery := `SELECT a.agreement_id, 
								a.agreement_type, 
								p.*, 
								o.*, 
								a.agreement_date, 
								a.status, 
								a.deposit_amount,
								a.payment_per_month,
								a.payment_duration,
								a.total_payment,
								a.cancelled_message,
								a.created_at
							FROM agreements a
							JOIN (` + propertiesQuery + `) AS p ON a.property_id = p.property_id
							JOIN (` + ownersQuery + `) AS o ON a.owner_user_id = o.owner_user_id`

	return repo.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.Agreements{}).
			Raw(agreementListsQuery+`
				WHERE a.owner_user_id = @user_id
				ORDER BY a.created_at `+agreementRequest.Order+`
				`, sql.Named("user_id", agreementRequest.UserId)).
			Scan(&agreementResponse.OwnerAgreements).Error; err != nil {
			return err
		}

		if err := tx.Model(&models.Agreements{}).
			Raw(agreementListsQuery+`
				WHERE a.dweller_user_id = @user_id
				ORDER BY a.created_at `+agreementRequest.Order+`
				`, sql.Named("user_id", agreementRequest.UserId)).
			Scan(&agreementResponse.DwellerAgreements).Error; err != nil {
			return err
		}

		for i, agreement := range agreementResponse.OwnerAgreements {
			if err := tx.Model(&models.PropertyImages{}).
				Raw(`
					SELECT image_url
					FROM property_images
					WHERE property_id = @property_id
					`, sql.Named("property_id", agreement.Property.PropertyId)).
				Pluck("image_url", &agreementResponse.OwnerAgreements[i].Property.PropertyImages).Error; err != nil {
				return err
			}
		}

		for i, agreement := range agreementResponse.DwellerAgreements {
			if err := tx.Model(&models.PropertyImages{}).
				Raw(`
					SELECT image_url
					FROM property_images
					WHERE property_id = @property_id
					`, sql.Named("property_id", agreement.Property.PropertyId)).
				Pluck("image_url", &agreementResponse.DwellerAgreements[i].Property.PropertyImages).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (repo *repositoryImpl) CreateAgreement(agreement *models.CreatingAgreements) error {
	return repo.db.Exec(`INSERT INTO agreements (agreement_id, agreement_type, property_id, owner_user_id, dweller_user_id, agreement_date, 
		status, deposit_amount, payment_per_month, payment_duration, total_payment) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		agreement.AgreementId, agreement.AgreementType, agreement.PropertyId, agreement.OwnerUserId, agreement.DwellerUserId, agreement.AgreementDate,
		agreement.Status, agreement.DepositAmount, agreement.PaymentPerMonth, agreement.PaymentDuration, agreement.TotalPayment).Error
}

func (repo *repositoryImpl) DeleteAgreement(agreementId string) error {
	if err := repo.db.Model(&models.Agreements{}).First(&models.Agreements{}, "agreement_id = ?", agreementId).Error; err != nil {
		return err
	}

	return repo.db.Where("agreement_id = ?", agreementId).Delete(&models.Agreements{}).Error
}

func (repo *repositoryImpl) UpdateAgreementStatus(updatingAgreement *models.UpdatingAgreementStatus, agreementId string) error {
	if err := repo.db.Model(&models.Agreements{}).First(&models.Agreements{}, "agreement_id = ?", agreementId).Error; err != nil {
		return err
	}

	return repo.db.Model(&models.Agreements{}).Where("agreement_id = ?", agreementId).Updates(updatingAgreement).Error
}
