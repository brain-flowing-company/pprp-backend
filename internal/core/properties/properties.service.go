package properties

import (
	"errors"
	"strings"

	"github.com/brain-flowing-company/pprp-backend/apperror"
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/brain-flowing-company/pprp-backend/internal/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Service interface {
	GetPropertyById(*models.Properties, string) *apperror.AppError
	GetAllProperties(*[]models.Properties) *apperror.AppError
	SearchProperties(*[]models.Properties, string) *apperror.AppError
}

type serviceImpl struct {
	repo   Repository
	logger *zap.Logger
}

func NewService(logger *zap.Logger, repo Repository) Service {
	return &serviceImpl{
		repo,
		logger,
	}
}

func (s *serviceImpl) GetPropertyById(property *models.Properties, id string) *apperror.AppError {
	if !utils.IsValidUUID(id) {
		return apperror.
			New(apperror.InvalidPropertyId).
			Describe("Invalid property id")
	}

	err := s.repo.GetPropertyById(property, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return apperror.
			New(apperror.PropertyNotFound).
			Describe("Could not find the specified property")
	} else if err != nil {
		s.logger.Error("Could not get property by id", zap.String("id", id), zap.Error(err))
		return apperror.
			New(apperror.InternalServerError).
			Describe("Could not get property. Please try again later.")
	}

	return nil
}

func (s *serviceImpl) GetAllProperties(properties *[]models.Properties) *apperror.AppError {
	err := s.repo.GetAllProperties(properties)
	if err != nil {
		s.logger.Error("Could not get all properties", zap.Error(err))
		return apperror.
			New(apperror.InternalServerError).
			Describe("Could not get all properties. Please try again later.")
	}

	return nil
}

func (s *serviceImpl) SearchProperties(properties *[]models.Properties, query string) *apperror.AppError {
	query = strings.ToLower(strings.TrimSpace(query))
	err := s.repo.SearchProperties(properties, query)
	if err != nil {
		s.logger.Error("Could not search properties", zap.Error(err))
		return apperror.
			New(apperror.InternalServerError).
			Describe("Could not search properties. Please try again later.")
	}

	return nil
}
