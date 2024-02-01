package users

import (
	"errors"

	"github.com/brain-flowing-company/pprp-backend/apperror"
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"gorm.io/gorm"
)

type Service interface {
	CreateUser(*models.Users) *apperror.AppError
	GetAllUsers(*models.Users) *apperror.AppError
}

type serviceImpl struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &serviceImpl{
		repo,
	}
}

func (service *serviceImpl) GetAllUsers(user *models.Users) *apperror.AppError {
	err := service.repo.GetAllUsers(user)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return apperror.InvalidPropertyId
	} else if err != nil {
		return apperror.InternalServerError
	}

	return nil
}

func (service *serviceImpl) CreateUser(user *models.Users) *apperror.AppError {

	err := service.repo.CreateUser(user)

	if err != nil {
		return apperror.InternalServerError
	}

	return nil
}
