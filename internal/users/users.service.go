package users

import (
	"errors"

	"github.com/brain-flowing-company/pprp-backend/apperror"
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/brain-flowing-company/pprp-backend/utils"
	"gorm.io/gorm"
)

type Service interface {
	GetAllUsers(*[]models.Users) *apperror.AppError
	GetUserById(*models.Users, string) *apperror.AppError
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

func (service *serviceImpl) GetAllUsers(users *[]models.Users) *apperror.AppError {
	err := service.repo.GetAllUsers(users)

	if err != nil {
		return apperror.InternalServerError
	}

	return nil
}

func (s *serviceImpl) GetUserById(user *models.Users, id string) *apperror.AppError {
	if !utils.IsValidUUID(id) {
		return apperror.InvalidUserId
	}

	err := s.repo.GetUserById(user, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return apperror.UserNotFound
	} else if err != nil {
		return apperror.InternalServerError
	}

	return nil
}

func (s *serviceImpl) CreateUser(user *models.Users) *apperror.AppError {

	err := s.repo.CreateUser(user)

	if err != nil {
		return apperror.InternalServerError
	}

	return nil
}
