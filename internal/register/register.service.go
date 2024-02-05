package register

import (
	"net/mail"

	"github.com/brain-flowing-company/pprp-backend/apperror"
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/brain-flowing-company/pprp-backend/utils"
	"go.uber.org/zap"
)

type Service interface {
	CreateUser(*models.Users) error
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

func (s *serviceImpl) CreateUser(user *models.Users) error {
	// Validate email
	_, err := mail.ParseAddress(user.Email)
	if err != nil {
		return apperror.InvalidEmail
	}

	existingUser, _ := s.repo.GetUserByEmail(user.Email)
	if existingUser != nil {
		return apperror.EmailAlreadyExists
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		s.logger.Error("Could not hash password", zap.Error(err))
		return apperror.InternalServerError
	}

	user.Password = hashedPassword

	if err := s.repo.CreateUser(user); err != nil {
		s.logger.Error("Could not create user", zap.Error(err))
		return apperror.InternalServerError
	}
	return nil
}
