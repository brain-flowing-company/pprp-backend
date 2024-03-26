package agreements

import (
	"errors"

	"github.com/brain-flowing-company/pprp-backend/apperror"
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/brain-flowing-company/pprp-backend/internal/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Service interface {
	GetAllAgreements(*[]models.Agreements) *apperror.AppError
	GetAgreementById(*models.AgreementDetails, string) *apperror.AppError
	GetAgreementsByOwnerId(*[]models.Agreements, string) *apperror.AppError
	GetAgreementsByDwellerId(*[]models.Agreements, string) *apperror.AppError
	CreateAgreement(*models.CreatingAgreements) *apperror.AppError
	DeleteAgreement(string) *apperror.AppError
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
func (s *serviceImpl) GetAllAgreements(results *[]models.Agreements) *apperror.AppError {
	err := s.repo.GetAllAgreements(results)
	if err != nil {
		s.logger.Error("Error getting all agreements", zap.Error(err))
		return apperror.New(apperror.InternalServerError).Describe("Error getting all agreements")
	}
	return nil
}

func (s *serviceImpl) GetAgreementById(agreement *models.AgreementDetails, agreementId string) *apperror.AppError {
	if !utils.IsValidUUID(agreementId) {
		return apperror.
			New(apperror.InvalidAgreementId).
			Describe("Invalid agreement id")
	}

	err := s.repo.GetAgreementById(agreement, agreementId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return apperror.
			New(apperror.AgreementNotFound).
			Describe("Could not find the specified agreement")
	} else if err != nil {
		s.logger.Error("Error getting agreement by id", zap.Error(err))
		return apperror.
			New(apperror.InternalServerError).
			Describe("Could not get agreement by id")
	}

	return nil
}

func (s *serviceImpl) GetAgreementsByOwnerId(agreements *[]models.Agreements, userId string) *apperror.AppError {
	err := s.repo.GetAgreementsByOwnerId(agreements, userId)
	if err != nil {
		s.logger.Error("Could not get agreements by owner id", zap.Error(err))
		return apperror.
			New(apperror.InternalServerError).
			Describe("Could not get agreements by owner id")
	}

	return nil
}

func (s *serviceImpl) CreateAgreement(creatingAgreement *models.CreatingAgreements) *apperror.AppError {
	err := s.repo.CreateAgreement(creatingAgreement)
	if err != nil {
		s.logger.Error("Could not create agreement", zap.Error(err))
		return apperror.
			New(apperror.InternalServerError).
			Describe("Could not create agreement")
	}

	return nil
}

func (s *serviceImpl) DeleteAgreement(id string) *apperror.AppError {
	if !utils.IsValidUUID(id) {
		return apperror.New(apperror.InvalidAgreementId).Describe("Invalid agreement id")
	}
	err := s.repo.DeleteAgreement(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return apperror.New(apperror.AgreementNotFound).Describe("Agreement not found")
	}
	if err != nil {
		s.logger.Error("Error deleting agreement", zap.Error(err))
		return apperror.New(apperror.InternalServerError).Describe("Error deleting agreement")
	}
	return nil
}

func (s *serviceImpl) GetAgreementsByDwellerId(agreements *[]models.Agreements, userId string) *apperror.AppError {
	err := s.repo.GetAgreementsByDwellerId(agreements, userId)
	if err != nil {
		s.logger.Error("Could not get agreements by dweller id", zap.Error(err))
		return apperror.
			New(apperror.InternalServerError).
			Describe("Could not get agreements by dweller id")
	}

	return nil
}
