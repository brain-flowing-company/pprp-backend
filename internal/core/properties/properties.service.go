package properties

import (
	"errors"
	"strings"

	"github.com/brain-flowing-company/pprp-backend/apperror"
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/brain-flowing-company/pprp-backend/internal/utils"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Service interface {
	GetAllProperties(*[]models.Properties) *apperror.AppError
	GetPropertyById(*models.Properties, string) *apperror.AppError
	SearchProperties(*[]models.Properties, string) *apperror.AppError
	AddFavoriteProperty(string, uuid.UUID) *apperror.AppError
	RemoveFavoriteProperty(string, uuid.UUID) *apperror.AppError
	GetFavoritePropertiesByUserId(*[]models.Properties, string) *apperror.AppError
	GetTop10Properties(*[]models.Properties) *apperror.AppError
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

func (s *serviceImpl) AddFavoriteProperty(propertyId string, userId uuid.UUID) *apperror.AppError {
	if !utils.IsValidUUID(propertyId) {
		return apperror.
			New(apperror.InvalidPropertyId).
			Describe("Invalid property id")
	}
	propertyIdUuid, _ := uuid.Parse(propertyId)

	favoriteProperty := models.FavoriteProperties{
		PropertyId: propertyIdUuid,
		UserId:     userId,
	}

	err := s.repo.AddFavoriteProperty(&favoriteProperty)
	if err != nil {
		s.logger.Error("Could not add favorite property", zap.Error(err))
		return apperror.
			New(apperror.InternalServerError).
			Describe("Could not add favorite property. Please try again later.")
	}

	return nil
}

func (s *serviceImpl) RemoveFavoriteProperty(propertyId string, userId uuid.UUID) *apperror.AppError {
	if !utils.IsValidUUID(propertyId) {
		return apperror.
			New(apperror.InvalidPropertyId).
			Describe("Invalid property id")
	}

	err := s.repo.RemoveFavoriteProperty(propertyId, userId.String())
	if err != nil {
		s.logger.Error("Could not remove favorite property", zap.Error(err))
		return apperror.
			New(apperror.InternalServerError).
			Describe("Could not remove favorite property. Please try again later.")
	}

	return nil
}

func (s *serviceImpl) GetFavoritePropertiesByUserId(properties *[]models.Properties, userId string) *apperror.AppError {
	if !utils.IsValidUUID(userId) {
		return apperror.
			New(apperror.InvalidUserId).
			Describe("Invalid user id")
	}

	err := s.repo.GetFavoritePropertiesByUserId(properties, userId)
	if err != nil {
		s.logger.Error("Could not get favorite properties by user id", zap.Error(err))
		return apperror.
			New(apperror.InternalServerError).
			Describe("Could not get favorite properties. Please try again later.")
	}

	return nil
}

func (s *serviceImpl) GetTop10Properties(properties *[]models.Properties) *apperror.AppError {
	err := s.repo.GetTop10Properties(properties)
	if err != nil {
		s.logger.Error("Could not get top 10 properties", zap.Error(err))
		return apperror.
			New(apperror.InternalServerError).
			Describe("Could not get top 10 properties. Please try again later.")
	}

	return nil
}
