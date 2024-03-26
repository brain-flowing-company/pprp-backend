package appointments

import (
	"errors"

	"github.com/brain-flowing-company/pprp-backend/apperror"
	"github.com/brain-flowing-company/pprp-backend/internal/enums"
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/brain-flowing-company/pprp-backend/internal/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Service interface {
	GetAllAppointments(*[]models.AppointmentLists) *apperror.AppError
	GetAppointmentById(*models.AppointmentDetails, string) *apperror.AppError
	GetMyAppointments(*models.MyAppointmentResponse, string) *apperror.AppError
	CreateAppointment(*models.CreatingAppointments) *apperror.AppError
	DeleteAppointment(string) *apperror.AppError
	UpdateAppointmentStatus(*models.UpdatingAppointmentStatus, string) *apperror.AppError
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

func (s *serviceImpl) GetAllAppointments(appointments *[]models.AppointmentLists) *apperror.AppError {
	err := s.repo.GetAllAppointments(appointments)
	if err != nil {
		s.logger.Error("Could not get all appointments", zap.Error(err))
		return apperror.
			New(apperror.InternalServerError).
			Describe("Could not get all appointments")
	}

	return nil
}

func (s *serviceImpl) GetAppointmentById(appointment *models.AppointmentDetails, appointmentId string) *apperror.AppError {
	if !utils.IsValidUUID(appointmentId) {
		return apperror.
			New(apperror.InvalidAppointmentId).
			Describe("Invalid appointment id")
	}

	err := s.repo.GetAppointmentById(appointment, appointmentId)
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

func (s *serviceImpl) GetMyAppointments(appointments *models.MyAppointmentResponse, userId string) *apperror.AppError {
	if !utils.IsValidUUID(userId) {
		return apperror.
			New(apperror.InvalidUserId).
			Describe("Invalid user id")
	}

	err := s.repo.GetAppointmentByUserId(appointments, userId)
	if err != nil {
		s.logger.Error("Could not get my appointments", zap.Error(err))
		return apperror.
			New(apperror.InternalServerError).
			Describe("Could not get my appointments")
	}

	return nil
}

func (s *serviceImpl) GetAppointmentByOwnerId(apps []*models.Appointments, userId string) *apperror.AppError {
	if !utils.IsValidUUID(userId) {
		return apperror.
			New(apperror.InvalidUserId).
			Describe("Invalid user id")
	}

	err := s.repo.GetAppointmentByOwnerId(apps, userId)
	if err != nil {
		s.logger.Error("Could not get appointments by owner id", zap.Error(err))
		return apperror.
			New(apperror.InternalServerError).
			Describe("Could not get appointments by owner id")
	}

	return nil
}

func (s *serviceImpl) GetAppointmentByDwellerId(apps []*models.Appointments, userId string) *apperror.AppError {
	if !utils.IsValidUUID(userId) {
		return apperror.
			New(apperror.InvalidUserId).
			Describe("Invalid user id")
	}

	err := s.repo.GetAppointmentByDwellerId(apps, userId)
	if err != nil {
		s.logger.Error("Could not get appointments by dweller id", zap.Error(err))
		return apperror.
			New(apperror.InternalServerError).
			Describe("Could not get appointments by dweller id")
	}

	return nil
}

func (s *serviceImpl) CreateAppointment(creatingAppointment *models.CreatingAppointments) *apperror.AppError {
	err := s.repo.CreateAppointment(creatingAppointment)
	if err != nil {
		s.logger.Error("Could not create appointments", zap.Error(err))
		return apperror.
			New(apperror.InternalServerError).
			Describe("Could not create appointments")
	}

	return nil
}

func (s *serviceImpl) DeleteAppointment(appointmentId string) *apperror.AppError {
	if !utils.IsValidUUID(appointmentId) {
		return apperror.
			New(apperror.InvalidAppointmentId).
			Describe("Invalid appointment id")
	}

	err := s.repo.DeleteAppointment(appointmentId)
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

func (s *serviceImpl) UpdateAppointmentStatus(updatingAppointment *models.UpdatingAppointmentStatus, appointmentId string) *apperror.AppError {
	if !utils.IsValidUUID(appointmentId) {
		return apperror.
			New(apperror.InvalidAppointmentId).
			Describe("Invalid appointment id")
	}

	_, ok := enums.AppointmentStatusMap[string(updatingAppointment.Status)]
	if !ok {
		return apperror.
			New(apperror.InvalidAppointmentStatus).
			Describe("Invalid appointment status")
	}

	err := s.repo.UpdateAppointmentStatus(updatingAppointment, appointmentId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return apperror.
			New(apperror.AppointmentNotFound).
			Describe("Could not find the specified appointment")
	} else if err != nil {
		s.logger.Error("Could not update appointment statue", zap.Error(err))
		return apperror.
			New(apperror.InternalServerError).
			Describe("Could not set appointment status")
	}

	return nil
}
