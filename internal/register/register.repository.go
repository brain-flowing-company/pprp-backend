package register

import (
	"fmt"

	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"gorm.io/gorm"
)

type Repository interface {
	CreateUser(*models.User) error
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
	// bcrytpt password before saving
	fmt.Println("user", user)
	if err := user.HashPassword(); err != nil {
		return err
	}

	return repo.db.Create(user).Error // Create() is a gorm method
}
