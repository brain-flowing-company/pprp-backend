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
	GetAllAgreements(*[]models.AgreementLists) *apperror.AppError
	GetAgreementById(*models.AgreementDetails, string) *apperror.AppError
	GetMyAgreements(*models.MyAgreementResponses, *models.MyAgreementRequests) *apperror.AppError
	CreateAgreement(*models.CreatingAgreements) *apperror.AppError
	DeleteAgreement(string) *apperror.AppError
	UpdateAgreementStatus(*models.UpdatingAgreementStatus, string) *apperror.AppError
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
func (s *serviceImpl) GetAllAgreements(agreements *[]models.AgreementLists) *apperror.AppError {
	err := s.repo.GetAllAgreements(agreements)
	if err != nil {
		s.logger.Error("Could not get all agreements", zap.Error(err))
		return apperror.
			New(apperror.InternalServerError).
			Describe("Could not get all agreements")
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
		s.logger.Error("Could not get agreement by id", zap.Error(err))
		return apperror.
			New(apperror.InternalServerError).
			Describe("Could not get agreement by id")
	}

	return nil
}

func (s *serviceImpl) GetMyAgreements(agreements *models.MyAgreementResponses, agreementRequest *models.MyAgreementRequests) *apperror.AppError {
	if agreementRequest.Order != "ASC" && agreementRequest.Order != "DESC" {
		agreementRequest.Order = "ASC"
	}

	err := s.repo.GetAgreementByUserId(agreements, agreementRequest)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return apperror.
			New(apperror.AppointmentNotFound).
			Describe("Could not find the specified agreement")
	} else if err != nil {
		s.logger.Error("Could not get agreement by user id", zap.Error(err))
		return apperror.
			New(apperror.InternalServerError).
			Describe("Could not get agreement by user id")
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

func (s *serviceImpl) DeleteAgreement(agreementId string) *apperror.AppError {
	if !utils.IsValidUUID(agreementId) {
		return apperror.
			New(apperror.InvalidAgreementId).
			Describe("Invalid agreement id")
	}

	err := s.repo.DeleteAgreement(agreementId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return apperror.
			New(apperror.AgreementNotFound).
			Describe("Could not find the specified agreement")
	} else if err != nil {
		s.logger.Error("Could not delete agreement", zap.Error(err))
		return apperror.
			New(apperror.InternalServerError).
			Describe("Could not delete agreement")
	}

	return nil
}

func (s *serviceImpl) UpdateAgreementStatus(updatingAgreement *models.UpdatingAgreementStatus, agreementId string) *apperror.AppError {
	if !utils.IsValidUUID(agreementId) {
		return apperror.
			New(apperror.InvalidAgreementId).
			Describe("Invalid agreement id")
	}

	err := s.repo.UpdateAgreementStatus(updatingAgreement, agreementId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return apperror.
			New(apperror.AgreementNotFound).
			Describe("Could not find the specified agreement")
	} else if err != nil {
		s.logger.Error("Could not update agreement status", zap.Error(err))
		return apperror.
			New(apperror.InternalServerError).
			Describe("Could not update agreement status")
	}

	return nil
}