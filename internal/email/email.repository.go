package email

import (
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"gorm.io/gorm"
)

type Repository interface {
	CountEmail(*int64, string) error
	CreateEmailVerificationData(*models.EmailVerificationData) error
	GetEmailVerificationDataByEmail(*models.EmailVerificationData, string) error
	CountEmailVerificationData(*int64, string) error
	DeleteEmailVerificationData(string) error
}

type repositoryImpl struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repositoryImpl{
		db,
	}
}

func (repo *repositoryImpl) CountEmail(count *int64, email string) error {
	return repo.db.Model(&models.Users{}).Where("email = ?", email).Count(count).Error
}

func (repo *repositoryImpl) CreateEmailVerificationData(emailVerificationData *models.EmailVerificationData) error {
	return repo.db.Create(emailVerificationData).Error
}

func (repo *repositoryImpl) GetEmailVerificationDataByEmail(emailVerificationData *models.EmailVerificationData, email string) error {
	return repo.db.First(emailVerificationData, "email = ?", email).Error
}

func (repo *repositoryImpl) CountEmailVerificationData(count *int64, email string) error {
	return repo.db.Model(&models.EmailVerificationData{}).Where("email = ?", email).Count(count).Error
}

func (repo *repositoryImpl) DeleteEmailVerificationData(email string) error {
	err := repo.db.First(&models.EmailVerificationData{}, "email = ?", email).Error
	if err != nil {
		return err
	}

	return repo.db.Where("email = ?", email).Delete(&models.EmailVerificationData{}).Error
}
