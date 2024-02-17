package appointments

import (
	"errors"

	"github.com/brain-flowing-company/pprp-backend/apperror"
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/brain-flowing-company/pprp-backend/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Service interface {
	GetAllAppointments(*[]models.Appointments) *apperror.AppError
	GetAppointmentsById(*models.Appointments, string) *apperror.AppError
	CreateAppointments(*[]models.Appointments) *apperror.AppError
	UpdateAppointments(*models.Appointments, string) *apperror.AppError
	DeleteAppointments(string) *apperror.AppError
}

type serviceImpl struct {
	logger *zap.Logger
	repo   Repository
}

func NewService(logger *zap.Logger, repo Repository) Service {
	return &serviceImpl{
		logger,
		repo,
	}
}

func (s *serviceImpl) GetAllAppointments(apps *[]models.Appointments) *apperror.AppError {
	err := s.repo.GetAllAppointments(apps)
	if err != nil {
		s.logger.Error("Could not get all appointments", zap.Error(err))
		return apperror.
			New(apperror.InternalServerError).
			Describe("Could not get all appointments")
	}

	return nil
}

func (s *serviceImpl) GetAppointmentsById(apps *models.Appointments, appId string) *apperror.AppError {
	if !utils.IsValidUUID(appId) {
		return apperror.
			New(apperror.InvalidAppointmentId).
			Describe("Invalid appointment id")
	}

	err := s.repo.GetAppointmentsById(apps, appId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return apperror.
			New(apperror.AppointmentNotFound).
			Describe("Could not find the specified appointment")
	} else if err != nil {
		s.logger.Error("Could not get appointment by id", zap.Error(err))
		return apperror.
			New(apperror.InternalServerError).
			Describe("Could not get appointment by id")
	}

	return nil
}

func (s *serviceImpl) GetAppointmentsByOwnerId(apps *[]models.Appointments, userId string) *apperror.AppError {
	err := s.repo.GetAppointmentsByOwnerId(apps, userId)
	if err != nil {
		s.logger.Error("Could not get appointments by owner id", zap.Error(err))
		return apperror.
			New(apperror.InternalServerError).
			Describe("Could not get appointments by owner id")
	}

	return nil
}

func (s *serviceImpl) CreateAppointments(apps *[]models.Appointments) *apperror.AppError {
	err := s.repo.CreateAppointments(apps)
	if err != nil {
		s.logger.Error("Could not create appointments", zap.Error(err))
		return apperror.
			New(apperror.InternalServerError).
			Describe("Could not create appointments")
	}

	return nil
}

func (s *serviceImpl) UpdateAppointments(app *models.Appointments, appId string) *apperror.AppError {
	err := s.repo.UpdateAppointments(app, appId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return apperror.
			New(apperror.AppointmentNotFound).
			Describe("Could not update the specified appointment")
	} else if err != nil {
		s.logger.Error("Could not update appointments", zap.Error(err))
		return apperror.
			New(apperror.InternalServerError).
			Describe("Could not update appointments")
	}

	return nil
}

func (s *serviceImpl) DeleteAppointments(appId string) *apperror.AppError {
	err := s.repo.DeleteAppointments(appId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return apperror.
			New(apperror.AppointmentNotFound).
			Describe("Could not delete the specified appointment")
	} else if err != nil {
		s.logger.Error("Could not delete appointments", zap.Error(err))
		return apperror.
			New(apperror.InternalServerError).
			Describe("Could not delete appointments")
	}

	return nil
}
