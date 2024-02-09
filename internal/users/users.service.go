package users

import (
	"errors"

	"github.com/brain-flowing-company/pprp-backend/apperror"
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/brain-flowing-company/pprp-backend/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Service interface {
	GetAllUsers(*[]models.Users) *apperror.AppError
	GetUserById(*models.Users, string) *apperror.AppError
	Register(*models.Users) *apperror.AppError
	UpdateUser(*models.Users, string) *apperror.AppError
	DeleteUser(string) *apperror.AppError
	GetUserByEmail(*models.Users, string) *apperror.AppError
}

type serviceImpl struct {
	repo   Repository
	logger *zap.Logger
}

func NewService(repo Repository, logger *zap.Logger) Service {
	return &serviceImpl{
		repo,
		logger,
	}
}

func (service *serviceImpl) GetAllUsers(users *[]models.Users) *apperror.AppError {
	err := service.repo.GetAllUsers(users)

	if err != nil {
		service.logger.Error("Could not get all users", zap.Error(err))
		return apperror.InternalServerError
	}

	return nil
}

func (s *serviceImpl) GetUserById(user *models.Users, userId string) *apperror.AppError {
	if !utils.IsValidUUID(userId) {
		return apperror.InvalidUserId
	}

	err := s.repo.GetUserById(user, userId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return apperror.UserNotFound
	} else if err != nil {
		s.logger.Error("Could not get user by id", zap.String("id", userId), zap.Error(err))
		return apperror.InternalServerError
	}

	return nil
}

func (s *serviceImpl) Register(user *models.Users) *apperror.AppError {
	if !utils.IsValidEmail(user.Email) {
		return apperror.InvalidEmail
	}

	if s.repo.GetUserByEmail(&models.Users{}, user.Email) == nil {
		return apperror.EmailAlreadyExists
	}

	if user.Password != "" {
		if !utils.IsValidPassword(user.Password) {
			return apperror.InvalidPassword
		}

		hashedPassword, hashErr := utils.HashPassword(user.Password)
		if hashErr != nil {
			return apperror.PasswordCannotBeHashed
		}
		user.Password = string(hashedPassword)
	}

	err := s.repo.CreateUser(user)
	if err != nil {
		s.logger.Error("Could not create user", zap.Error(err))
		return apperror.InternalServerError
	}

	return nil
}

func (s *serviceImpl) UpdateUser(user *models.Users, userId string) *apperror.AppError {
	if !utils.IsValidUUID(userId) {
		return apperror.InvalidUserId
	}

	err := s.repo.UpdateUser(user, userId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return apperror.UserNotFound
	} else if err != nil {
		s.logger.Error("Could not update user info", zap.String("id", userId), zap.Error(err))
		return apperror.InternalServerError
	}

	return nil
}

func (s *serviceImpl) DeleteUser(userId string) *apperror.AppError {
	if !utils.IsValidUUID(userId) {
		return apperror.InvalidUserId
	}

	err := s.repo.DeleteUser(userId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return apperror.UserNotFound
	} else if err != nil {
		s.logger.Error("Could not delete user", zap.String("id", userId), zap.Error(err))
		return apperror.InternalServerError
	}

	return nil
}

func (s *serviceImpl) GetUserByEmail(user *models.Users, email string) *apperror.AppError {
	if !utils.IsValidEmail(email) {
		return apperror.InvalidEmail
	}

	err := s.repo.GetUserByEmail(user, email)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return apperror.UserNotFound
	} else if err != nil {
		s.logger.Error("Could not get current user info", zap.String("email", email), zap.Error(err))
		return apperror.InternalServerError
	}

	return nil
}
