package ratings

import (
	"github.com/brain-flowing-company/pprp-backend/config"
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"go.uber.org/zap"
)

type Service interface {
	CreateRating(*models.Reviews) error
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
