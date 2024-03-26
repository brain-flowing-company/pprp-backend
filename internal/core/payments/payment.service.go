package payments

import (
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"go.uber.org/zap"
)

type Service interface {
	CreatePayment(*models.Payments) error
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

func (s *serviceImpl) CreatePayment(payment *models.Payments) error {
	err := s.repo.CreatePayment(payment)
	if err != nil {
		s.logger.Error("Failed to create payment", zap.Error(err))
		return err
	}
	return nil
}
