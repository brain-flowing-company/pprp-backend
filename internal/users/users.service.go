package users

import (
	"github.com/brain-flowing-company/pprp-backend/apperror"
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/google/uuid"
)

type Service interface {
	CreateUser(*models.Users) *apperror.AppError
}

type serviceImpl struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &serviceImpl{
		repo,
	}
}

func (service *serviceImpl) CreateUser(user *models.Users) *apperror.AppError {
	user.UserId = uuid.New().String()

	err := service.repo.CreateUser(user)

	if err != nil {
		return apperror.InternalServerError
	}

	return nil
}
