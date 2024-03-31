package payments

import (
	"fmt"

	"github.com/brain-flowing-company/pprp-backend/config"
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Service interface {
	CreatePayment(*models.Payments) error
	GetPaymentByUserId(*models.MyPaymentsResponse, uuid.UUID) error
}

type serviceImpl struct {
	repo   Repository
	logger *zap.Logger
	cfg    *config.Config
}

func NewService(logger *zap.Logger, repo Repository, cfg *config.Config) Service {
	return &serviceImpl{
		repo,
		logger,
		cfg,
	}
}

func (s *serviceImpl) CreatePayment(payment *models.Payments) error {
	fmt.Println("service payment = ", payment)
	err := s.repo.CreatePayment(payment)
	if err != nil {
		s.logger.Error("Failed to create payment", zap.Error(err))
		return err
	}
	return nil
}

func (s *serviceImpl) GetPaymentByUserId(payments *models.MyPaymentsResponse, userId uuid.UUID) error {
	err := s.repo.GetPaymentByUserId(payments, userId)
	if err != nil {
		s.logger.Error("Failed to get payment by user id", zap.Error(err))
		return err
	}
	return nil
}
