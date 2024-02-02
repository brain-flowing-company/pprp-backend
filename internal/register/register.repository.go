package register

import (
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"gorm.io/gorm"
)

type Repository interface {
	CreateUser(*models.User) error
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

func (repo *repositoryImpl) CreateUser(user *models.User) error {
	return repo.db.Create(user).Error // Create() is a gorm method
}
func (repo *repositoryImpl) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	err := repo.db.Where("email = ?", email).First(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}
