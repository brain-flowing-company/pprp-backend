package ratings

import (
	"github.com/brain-flowing-company/pprp-backend/config"
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Service interface {
	CreateRating(*models.Reviews) error
	GetRatingByPropertyId(uuid.UUID, *[]models.RatingResponse) error
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
