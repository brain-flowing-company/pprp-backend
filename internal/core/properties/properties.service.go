package properties

import (
	"errors"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/brain-flowing-company/pprp-backend/apperror"
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/brain-flowing-company/pprp-backend/internal/utils"
	"github.com/brain-flowing-company/pprp-backend/storage"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Service interface {
	GetAllProperties(*models.AllPropertiesResponses, string, string, *utils.PaginatedQuery, *utils.SortedQuery) *apperror.AppError
	GetPropertyById(*models.Properties, string) *apperror.AppError
	GetPropertyByOwnerId(*models.MyPropertiesResponses, string, *utils.PaginatedQuery) *apperror.AppError
	CreateProperty(*models.PropertyInfos, *multipart.FileHeader) *apperror.AppError
	UpdatePropertyById(*models.PropertyInfos, string) *apperror.AppError
	DeletePropertyById(string) *apperror.AppError
	AddFavoriteProperty(string, uuid.UUID) *apperror.AppError
	RemoveFavoriteProperty(string, uuid.UUID) *apperror.AppError
	GetFavoritePropertiesByUserId(*models.MyFavoritePropertiesResponses, string, *utils.PaginatedQuery) *apperror.AppError
	GetTop10Properties(*[]models.Properties, string) *apperror.AppError
}

type serviceImpl struct {
	repo    Repository
	logger  *zap.Logger
	storage storage.Storage
}

func NewService(logger *zap.Logger, repo Repository, storage storage.Storage) Service {
	return &serviceImpl{
		repo,
		logger,
		storage,
	}
}

func (s *serviceImpl) GetAllProperties(properties *models.AllPropertiesResponses, query string, userId string, paginated *utils.PaginatedQuery, sorted *utils.SortedQuery) *apperror.AppError {
	if !utils.IsValidUUID(userId) {
		return apperror.
			New(apperror.InvalidUserId).
			Describe("Invalid user id")
	}

	query = strings.ToLower(strings.TrimSpace(query))
	err := s.repo.GetAllProperties(properties, query, userId, paginated, sorted)
	if err != nil {
		s.logger.Error("Could not search properties", zap.Error(err))
		return apperror.
			New(apperror.InternalServerError).
			Describe("Could not search properties. Please try again later.")
	}

	return nil
}

func (s *serviceImpl) GetPropertyById(property *models.Properties, propertyId string) *apperror.AppError {
	if !utils.IsValidUUID(propertyId) {
		return apperror.
			New(apperror.InvalidPropertyId).
			Describe("Invalid property id")
	}

	var countProperty int64
	countErr := s.repo.CountProperty(&countProperty, propertyId)
	if countErr != nil {
		s.logger.Error("Could not count property by id", zap.Error(countErr))
		return apperror.
			New(apperror.InternalServerError).
			Describe("Could not update property. Please try again later.")
	} else if countProperty == 0 {
		return apperror.
			New(apperror.PropertyNotFound).
			Describe("Could not find the specified property")
	}

	err := s.repo.GetPropertyById(property, propertyId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return apperror.
			New(apperror.PropertyNotFound).
			Describe("Could not find the specified property")
	} else if err != nil {
		s.logger.Error("Could not get property by id", zap.String("id", propertyId), zap.Error(err))
		return apperror.
			New(apperror.InternalServerError).
			Describe("Could not get property. Please try again later.")
	}

	return nil
}

func (s *serviceImpl) GetPropertyByOwnerId(properties *models.MyPropertiesResponses, ownerId string, paginated *utils.PaginatedQuery) *apperror.AppError {
	if !utils.IsValidUUID(ownerId) {
		return apperror.
			New(apperror.InvalidUserId).
			Describe("Invalid user id")
	}

	err := s.repo.GetPropertyByOwnerId(properties, ownerId, paginated)
	if err != nil {
		s.logger.Error("Could not get property by owner id", zap.Error(err))
		return apperror.
			New(apperror.InternalServerError).
			Describe("Could not get property. Please try again later.")
	}

	return nil
}

func (s *serviceImpl) CreateProperty(property *models.PropertyInfos, propertyImages *multipart.FileHeader) *apperror.AppError {

	err := s.repo.CreateProperty(property)
	if err != nil {
		s.logger.Error("Could not create property", zap.Error(err))
		return apperror.
			New(apperror.InternalServerError).
			Describe("Could not create property. Please try again later.")
	}

	return nil
}

func (s *serviceImpl) UpdatePropertyById(property *models.PropertyInfos, propertyId string) *apperror.AppError {
	if !utils.IsValidUUID(propertyId) {
		return apperror.
			New(apperror.InvalidPropertyId).
			Describe("Invalid property id")
	}

	err := s.repo.UpdatePropertyById(property, propertyId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return apperror.
			New(apperror.PropertyNotFound).
			Describe("Could not find the specified property")
	} else if err != nil {
		s.logger.Error("Could not update property by id", zap.String("id", propertyId), zap.Error(err))
		return apperror.
			New(apperror.InternalServerError).
			Describe("Could not update property. Please try again later.")
	}

	return nil
}

func (s *serviceImpl) DeletePropertyById(propertyId string) *apperror.AppError {
	if !utils.IsValidUUID(propertyId) {
		return apperror.
			New(apperror.InvalidPropertyId).
			Describe("Invalid property id")
	}

	var countProperty int64
	countErr := s.repo.CountProperty(&countProperty, propertyId)
	if countErr != nil {
		s.logger.Error("Could not count property by id", zap.String("id", propertyId), zap.Error(countErr))
		return apperror.
			New(apperror.InternalServerError).
			Describe("Could not update property. Please try again later.")
	} else if countProperty == 0 {
		return apperror.
			New(apperror.PropertyNotFound).
			Describe("Could not find the specified property")
	}

	err := s.repo.DeletePropertyById(propertyId)
	if err != nil {
		s.logger.Error("Could not delete property by id", zap.String("id", propertyId), zap.Error(err))
		return apperror.
			New(apperror.InternalServerError).
			Describe("Could not delete property. Please try again later.")
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

func (s *serviceImpl) GetFavoritePropertiesByUserId(properties *models.MyFavoritePropertiesResponses, userId string, paginated *utils.PaginatedQuery) *apperror.AppError {
	if !utils.IsValidUUID(userId) {
		return apperror.
			New(apperror.InvalidUserId).
			Describe("Invalid user id")
	}

	err := s.repo.GetFavoritePropertiesByUserId(properties, userId, paginated)
	if err != nil {
		s.logger.Error("Could not get favorite properties by user id", zap.Error(err))
		return apperror.
			New(apperror.InternalServerError).
			Describe("Could not get favorite properties. Please try again later.")
	}

	return nil
}

func (s *serviceImpl) GetTop10Properties(properties *[]models.Properties, userId string) *apperror.AppError {
	if !utils.IsValidUUID(userId) {
		return apperror.
			New(apperror.InvalidUserId).
			Describe("Invalid user id")
	}

	err := s.repo.GetTop10Properties(properties, userId)
	if err != nil {
		s.logger.Error("Could not get top 10 properties", zap.Error(err))
		return apperror.
			New(apperror.InternalServerError).
			Describe("Could not get top 10 properties. Please try again later.")
	}

	return nil
}

func (s *serviceImpl) uploadPropertyImage(propertyId uuid.UUID, propertyImage *multipart.FileHeader) (string, *apperror.AppError) {
	if propertyImage == nil {
		return "", apperror.
			New(apperror.BadRequest).
			Describe("No citizen card found")
	}

	file, err := propertyImage.Open()
	if err != nil {
		return "", apperror.
			New(apperror.InternalServerError).
			Describe("Could not upload profile image")
	}

	ext := filepath.Ext(propertyImage.Filename)
	ip := utils.NewImageProcessor()

	switch strings.ToLower(ext) {
	case ".jpg":
		fallthrough

	case ".jpeg":
		err = ip.LoadJPEG(file)

	case ".png":
		err = ip.LoadPNG(file)

	default:
		return "", apperror.
			New(apperror.InvalidPropertyImageExtension).
			Describe(fmt.Sprintf("App does not support %v extension", ext))
	}

	if err != nil {
		s.logger.Error("Could not load image", zap.Error(err))
		return "", apperror.
			New(apperror.InternalServerError).
			Describe("Could not process image")
	}

	processedFile, err := ip.Save()
	if err != nil {
		s.logger.Error("Could not create new image", zap.Error(err))
		return "", apperror.
			New(apperror.InternalServerError).
			Describe("Could not process image")
	}

	url, err := s.storage.Upload(fmt.Sprintf("verifications/%v.jpeg", propertyId.String()), processedFile, types.ObjectCannedACLPrivate)

	if err != nil {
		return "", apperror.
			New(apperror.InternalServerError).
			Describe("Could not upload profile image")
	}

	return url, nil
}
