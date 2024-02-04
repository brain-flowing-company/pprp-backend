package register

import (
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"gorm.io/gorm"
)

type Repository interface {
	CreateUser(*models.Users) error
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

func (repo *repositoryImpl) CreateUser(user *models.Users) error {
	return repo.db.Create(user).Error // Create() is a gorm method
}
func (repo *repositoryImpl) GetUserByEmail(email string) (*models.Users, error) {
	user := &models.Users{}
	err := repo.db.Where("email = ?", email).First(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}
