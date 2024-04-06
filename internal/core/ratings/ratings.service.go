package ratings

import (
	"errors"

	"github.com/brain-flowing-company/pprp-backend/apperror"
	"github.com/brain-flowing-company/pprp-backend/config"
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Service interface {
	CreateRating(*models.Reviews) error
	GetRatingByPropertyId(uuid.UUID, *[]models.RatingResponse) error
	GetAllRatings(*[]models.RatingResponse) error
	GetRatingByPropertyIdSortedByRating(propertyId uuid.UUID, ratings *[]models.RatingResponse) error
	GetRatingByPropertyIdSortedByNewest(propertyId uuid.UUID, ratings *[]models.RatingResponse) error
	UpdateRatingStatus(updatingRating *models.UpdateRatingStatus, ratingId uuid.UUID) error
}

type serviceImpl struct {
	repo   Repository
	logger *zap.Logger
	cfg    *config.Config
}

func NewService(repo Repository, logger *zap.Logger, cfg *config.Config) Service {
	return &serviceImpl{
		repo,
		logger,
		cfg,
	}
}

func (s *serviceImpl) CreateRating(reviews *models.Reviews) error {
	err := s.repo.CreateRating(reviews)
	if err != nil {
		s.logger.Error("Failed to create rating", zap.Error(err))
		return err
	}
	return nil
}

func (s *serviceImpl) GetRatingByPropertyId(propertyId uuid.UUID, ratings *[]models.RatingResponse) error {
	err := s.repo.GetRatingByPropertyId(propertyId, ratings)
	if err != nil {
		s.logger.Error("Failed to get rating by property id", zap.Error(err))
		return err
	}
	return nil
}

func (s *serviceImpl) GetAllRatings(ratings *[]models.RatingResponse) error {
	err := s.repo.GetAllRatings(ratings)
	if err != nil {
		s.logger.Error("Failed to get all ratings", zap.Error(err))
		return err
	}
	return nil
}

func (s *serviceImpl) GetRatingByPropertyIdSortedByRating(propertyId uuid.UUID, ratings *[]models.RatingResponse) error {
	err := s.repo.GetRatingByPropertyIdSortedByRating(propertyId, ratings)
	if err != nil {
		s.logger.Error("Failed to get rating by property id sorted by rating", zap.Error(err))
		return err
	}
	return nil
}

func (s *serviceImpl) GetRatingByPropertyIdSortedByNewest(propertyId uuid.UUID, ratings *[]models.RatingResponse) error {
	err := s.repo.GetRatingByPropertyIdSortedByNewest(propertyId, ratings)
	if err != nil {
		s.logger.Error("Failed to get rating by property id sorted by newest", zap.Error(err))
		return err
	}
	return nil
}

func (s *serviceImpl) UpdateRatingStatus(updatingRating *models.UpdateRatingStatus, ratingId uuid.UUID) error {
	err := s.repo.UpdateRatingStatus(updatingRating, ratingId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.Error("Rating not found", zap.Error(err))
		return apperror.New(apperror.RatingNotFound).Describe("Rating not found")
	} else if err != nil {
		s.logger.Error("Failed to update rating status", zap.Error(err))
		return apperror.New(apperror.InternalServerError).Describe("Failed to update rating status")
	}
	return nil
}
