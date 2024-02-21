package google

import (
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"gorm.io/gorm"
)

type Repository interface {
	CountEmail(count *int64, email string) error
	GetUserByEmail(user *models.Users, email string) error
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

func (repo *repositoryImpl) GetUserByEmail(user *models.Users, email string) error {
	return repo.db.Model(&models.Users{}).First(user, "email = ?", email).Error
}
