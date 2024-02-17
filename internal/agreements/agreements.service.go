package agreements

import (
	"errors"

	"github.com/brain-flowing-company/pprp-backend/apperror"
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/brain-flowing-company/pprp-backend/utils"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Service interface {
	GetAllAgreements(*[]models.Agreement) *apperror.AppError
	GetAgreementById(*models.Agreement, string) *apperror.AppError
	GetAgreementsByOwnerId(*[]models.Agreement, string) *apperror.AppError
	GetAgreementsByDwellerId(*[]models.Agreement, string) *apperror.AppError
	CreateAgreement(*models.CreatingAgreement) *apperror.AppError
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
func (s *serviceImpl) GetAllAgreements(results *[]models.Agreement) *apperror.AppError {
	err := s.repo.GetAllAgreements(results)
	if err != nil {
		s.logger.Error("Error getting all agreements", zap.Error(err))
		return apperror.New(apperror.InternalServerError).Describe("Error getting all agreements")
	}
	return nil
}

func (s *serviceImpl) GetAgreementById(result *models.Agreement, id string) *apperror.AppError {
	if !utils.IsValidUUID(id) {
		return apperror.New(apperror.InvalidAgreementId).Describe("Invalid agreement id")
	}
	err := s.repo.GetAgreementById(result, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return apperror.New(apperror.AgreementNotFound).Describe("Agreement not found")
	}
	if err != nil {
		s.logger.Error("Error getting agreement by id", zap.Error(err))
		return apperror.New(apperror.InternalServerError).Describe("Error getting agreement by id")
	}
	return nil
}

func (s *serviceImpl) GetAgreementsByOwnerId(agreements *[]models.Agreement, userId string) *apperror.AppError {
	err := s.repo.GetAgreementsByOwnerId(agreements, userId)
	if err != nil {
		s.logger.Error("Could not get agreements by owner id", zap.Error(err))
		return apperror.
			New(apperror.InternalServerError).
			Describe("Could not get agreements by owner id")
	}

	return nil
}

func (s *serviceImpl) CreateAgreement(creatingAgreement *models.CreatingAgreement) *apperror.AppError {
	agreement := models.Agreement{
		AgreementID:   uuid.New(),
		PropertyID:    creatingAgreement.PropertyID,
		OwnerUserID:   creatingAgreement.OwnerUserID,
		DwellerUserID: creatingAgreement.DwellerUserID,
		AgreementDate: creatingAgreement.AgreementDate,
	}

	err := s.repo.CreateAgreement(&agreement)
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return apperror.
			New(apperror.DuplicateAgreement).
			Describe("Could not create agreement")
	} else if err != nil {
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

func (s *serviceImpl) GetAgreementsByDwellerId(agreements *[]models.Agreement, userId string) *apperror.AppError {
	err := s.repo.GetAgreementsByDwellerId(agreements, userId)
	if err != nil {
		s.logger.Error("Could not get agreements by dweller id", zap.Error(err))
		return apperror.
			New(apperror.InternalServerError).
			Describe("Could not get agreements by dweller id")
	}

	return nil
}
