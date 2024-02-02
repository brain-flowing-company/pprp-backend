// internal/login/repository.go
package login

import (
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"gorm.io/gorm"
)

type Repository interface {
	GetUserByEmail(email string) (*models.User, error)
}

type repositoryImpl struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repositoryImpl{
		db,
	}
}

func (repo *repositoryImpl) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	err := repo.db.Where("email = ?", email).First(user).Error
	return user, err
}
