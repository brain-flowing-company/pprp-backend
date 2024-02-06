package property

import (
	"errors"

	"github.com/brain-flowing-company/pprp-backend/apperror"
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/brain-flowing-company/pprp-backend/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Service interface {
	GetPropertyById(*models.Property, string) *apperror.AppError
	GetAllProperties(*[]models.Property) *apperror.AppError
}

type serviceImpl struct {
	repo   Repository
	logger *zap.Logger
}

func NewService(repo Repository, logger *zap.Logger) Service {
	return &serviceImpl{
		repo,
		logger,
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
		s.logger.Error("Could not get property by id", zap.String("id", id), zap.Error(err))
		return apperror.InternalServerError
	}

	return nil
}

func (s *serviceImpl) GetAllProperties(properties *[]models.Property) *apperror.AppError {
	err := s.repo.GetAllProperties(properties)
	if err != nil {
		s.logger.Error("Could not get all properties", zap.Error(err))
		return apperror.InternalServerError
	}

	return nil
}
