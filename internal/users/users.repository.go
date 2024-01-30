package users

import (
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"gorm.io/gorm"
)

type Repository interface {
	CreateUser(*models.Users) error
}

type repositoryImpl struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repositoryImpl{
		db,
	}
}

func (repo *repositoryImpl) CreateUser(user *models.Users) error {
	return repo.db.Create(user).Error
}
