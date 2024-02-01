package users

import (
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	CreateUser(*models.Users) error
	GetAllUsers(*models.Users) error
}

type repositoryImpl struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repositoryImpl{
		db,
	}
}

func (repo *repositoryImpl) GetAllUsers(user *models.Users) error {
	return repo.db.Find(&user).Error
}

func (repo *repositoryImpl) CreateUser(user *models.Users) error {
	user.UserId = uuid.New().String()

	for repo.db.Find(&models.Users{}, "user_id = ?", user.UserId).RowsAffected != 0 {
		user.UserId = uuid.New().String()
	}

	return repo.db.Create(&user).Error
}
