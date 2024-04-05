package ratings

import "github.com/brain-flowing-company/pprp-backend/internal/models"

type Service interface {
	CreateRating(*models.Reviews) error
}

type serviceImpl struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &serviceImpl{
		repo,
	}
}

func (s *serviceImpl) CreateRating(reviews *models.Reviews) error {
	err := s.repo.CreateRating(reviews)
	if err != nil {
		return err
	}
	return nil
}
