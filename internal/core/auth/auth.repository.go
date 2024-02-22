// internal/login/repository.go
package auth

import (
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"gorm.io/gorm"
)

type Repository interface {
	GetUserByEmail(email string) (*models.Users, error)
}

type repositoryImpl struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repositoryImpl{
		db,
	}
}

func (repo *repositoryImpl) GetUserByEmail(email string) (*models.Users, error) {
	user := &models.Users{}
	err := repo.db.Where("email = ?", email).First(user).Error
	return user, err
}
