package emails

import (
	"fmt"

	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"gorm.io/gorm"
)

type Repository interface {
	CountEmail(*int64, string) error
	CreateEmailVerificationCode(*models.EmailVerificationCodes) error
	GetEmailVerificationCodeByEmail(*models.EmailVerificationCodes, string) error
	CountEmailVerificationCode(*int64, string) error
	DeleteEmailVerificationCode(string) error
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

func (repo *repositoryImpl) CreateEmailVerificationCode(emailVerificationCodes *models.EmailVerificationCodes) error {
	fmt.Println(emailVerificationCodes)
	return repo.db.Create(emailVerificationCodes).Error
}

func (repo *repositoryImpl) GetEmailVerificationCodeByEmail(emailVerificationCodes *models.EmailVerificationCodes, email string) error {
	return repo.db.First(emailVerificationCodes, "email = ?", email).Error
}

func (repo *repositoryImpl) CountEmailVerificationCode(count *int64, email string) error {
	return repo.db.Model(&models.EmailVerificationCodes{}).Where("email = ?", email).Count(count).Error
}

func (repo *repositoryImpl) DeleteEmailVerificationCode(email string) error {
	err := repo.db.First(&models.EmailVerificationCodes{}, "email = ?", email).Error
	if err != nil {
		return err
	}

	return repo.db.Where("email = ?", email).Delete(&models.EmailVerificationCodes{}).Error
}
