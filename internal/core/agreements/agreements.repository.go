package agreements

import (
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"gorm.io/gorm"
)

type Repository interface {
	GetAllAgreements(*[]models.Agreements) error
	GetAgreementById(*models.Agreements, string) error
	GetAgreementsByOwnerId(*[]models.Agreements, string) error
	GetAgreementsByDwellerId(*[]models.Agreements, string) error
	CreateAgreement(*models.CreatingAgreements) error
	DeleteAgreement(string) error
}

type repositoryImpl struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repositoryImpl{
		db,
	}
}

func (repo *repositoryImpl) GetAllAgreements(results *[]models.Agreements) error {
	return repo.db.Model(&models.Agreements{}).
		Find(results).Error
}

func (repo *repositoryImpl) GetAgreementById(result *models.Agreements, id string) error {
	return repo.db.Model(&models.Agreements{}).
		First(result, "agreement_id = ?", id).Error
}

func (repo *repositoryImpl) GetAgreementsByOwnerId(result *[]models.Agreements, id string) error {
	return repo.db.Model(&models.Agreements{}).
		Where("owner_user_id = ?", id).Find(result).Error
}

func (repo *repositoryImpl) GetAgreementsByDwellerId(result *[]models.Agreements, id string) error {
	return repo.db.Model(&models.Agreements{}).
		Where("dweller_user_id = ?", id).Find(result).Error
}

func (repo *repositoryImpl) CreateAgreement(agreement *models.CreatingAgreements) error {
	return repo.db.Exec(`INSERT INTO agreements (agreement_type, property_id, owner_user_id, dweller_user_id, agreement_date, 
		status, deposit_amount, payment_per_month, payment_duration, total_payment) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		agreement.AgreementType, agreement.PropertyId, agreement.OwnerUserId, agreement.DwellerUserId, agreement.AgreementDate, 
		agreement.Status, agreement.DepositAmount, agreement.PaymentPerMonth, agreement.PaymentDuration, agreement.TotalPayment ).Error
}

func (repo *repositoryImpl) DeleteAgreement(id string) error {
	return repo.db.Model(&models.Agreements{}).
		Delete(&models.Agreements{}, "agreement_id = ?", id).Error
}
