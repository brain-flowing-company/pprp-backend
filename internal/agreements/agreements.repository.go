package agreements

import (
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"gorm.io/gorm"
)

type Repository interface {
	GetAllAgreements(*[]models.Agreement) error
	GetAgreementById(*models.Agreement, string) error
	GetAgreementsByOwnerId(*[]models.Agreement, string) error
	GetAgreementsByDwellerId(*[]models.Agreement, string) error
	CreateAgreement(*models.Agreement) error
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

func (repo *repositoryImpl) GetAllAgreements(results *[]models.Agreement) error {
	return repo.db.Model(&models.Agreement{}).
		Find(results).Error
}

func (repo *repositoryImpl) GetAgreementById(result *models.Agreement, id string) error {
	return repo.db.Model(&models.Agreement{}).
		First(result, "agreement_id = ?", id).Error
}

func (repo *repositoryImpl) GetAgreementsByOwnerId(result *[]models.Agreement, id string) error {
	return repo.db.Model(&models.Agreement{}).
		Where("owner_user_id = ?", id).Find(result).Error
}

func (repo *repositoryImpl) GetAgreementsByDwellerId(result *[]models.Agreement, id string) error {
	return repo.db.Model(&models.Agreement{}).
		Where("dweller_user_id = ?", id).Find(result).Error
}

func (repo *repositoryImpl) CreateAgreement(agreement *models.Agreement) error {
	return repo.db.Model(&models.Agreement{}).
		Create(agreement).Error
}

func (repo *repositoryImpl) DeleteAgreement(id string) error {
	return repo.db.Model(&models.Agreement{}).
		Delete(&models.Agreement{}, "agreement_id = ?", id).Error
}
