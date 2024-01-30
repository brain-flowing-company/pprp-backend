package property

import (
	"errors"

	"github.com/brain-flowing-company/pprp-backend/apperror"
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/brain-flowing-company/pprp-backend/utils"
	"gorm.io/gorm"
)

type Service interface {
	GetPropertyById(*models.Property, string) *apperror.AppError
}

type serviceImpl struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &serviceImpl{
		repo,
	}
}

func (s *serviceImpl) GetPropertyById(property *models.Property, id string) *apperror.AppError {
	if !utils.IsValidUUID(id) {
		return apperror.InvalidPropertyId
	}

	err := s.repo.GetPropertyById(property, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return apperror.InvalidPropertyId
	} else if err != nil {
		return apperror.InternalServerError
	}

	return nil
}
