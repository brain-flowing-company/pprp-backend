package appointments

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
	GetAllAppointments(*[]models.Appointments) *apperror.AppError
	GetAppointmentsById(*models.Appointments, string) *apperror.AppError
	CreateAppointments(*models.CreatingAppointments) *apperror.AppError
	DeleteAppointments(*[]string) *apperror.AppError
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

func (s *serviceImpl) CreateAppointments(creatingApp *models.CreatingAppointments) *apperror.AppError {
	n := len(creatingApp.AppointmentDates)
	if n == 0 {
		return apperror.
			New(apperror.BadRequest).
			Describe("Appointment dates cannot be empty")
	}

	apps := make([]models.Appointments, len(creatingApp.AppointmentDates))
	for i := 0; i < n; i++ {
		apps[i] = models.Appointments{
			AppointmentId:      uuid.New(),
			PropertyId:         creatingApp.PropertyId,
			OwnerUserId:        creatingApp.OwnerUserId,
			DwellerUserId:      creatingApp.DwellerUserId,
			AppointmentDate:    creatingApp.AppointmentDates[i],
			AppointmentsStatus: models.Pending,
		}
	}

	err := s.repo.CreateAppointments(&apps)
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return apperror.
			New(apperror.DuplicateAppointment).
			Describe("Could not create appointments")
	} else if err != nil {
		s.logger.Error("Could not create appointments", zap.Error(err))
		return apperror.
			New(apperror.InternalServerError).
			Describe("Could not create appointments")
	}

	return nil
}

func (s *serviceImpl) DeleteAppointments(appIds *[]string) *apperror.AppError {
	for _, e := range *appIds {
		if !utils.IsValidUUID(e) {
			return apperror.
				New(apperror.InvalidAppointmentId).
				Describe("Invalid appointment id")
		}
	}

	err := s.repo.DeleteAppointments(appIds)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return apperror.
			New(apperror.AppointmentNotFound).
			Describe("Could not find the specified appointment")
	} else if err != nil {
		s.logger.Error("Could not delete appointments", zap.Error(err))
		return apperror.
			New(apperror.InternalServerError).
			Describe("Could not delete appointments")
	}

	return nil
}
