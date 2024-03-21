package google

import (
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"gorm.io/gorm"
)

type Repository interface {
	CountEmail(count *int64, email string) error
	GetUserByEmail(user *models.Users, email string) error
	CreateState(state *models.GoogleOAuthStates) error
	GetState(result *models.GoogleOAuthStates, state string) error
	DeleteState(state string) error
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

func (repo *repositoryImpl) CreateState(state *models.GoogleOAuthStates) error {
	return repo.db.Model(&models.GoogleOAuthStates{}).Create(state).Error
}

func (repo *repositoryImpl) GetState(result *models.GoogleOAuthStates, state string) error {
	return repo.db.Model(&models.GoogleOAuthStates{}).Find(result, "code = ? AND expired_at < NOW()", state).Error
}

func (repo *repositoryImpl) DeleteState(state string) error {
	return repo.db.Where("code = ?", state).Delete(&models.GoogleOAuthStates{}).Error
}
