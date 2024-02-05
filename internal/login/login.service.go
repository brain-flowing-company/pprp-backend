// internal/login/service.go
package login

import (
	"time"

	"github.com/brain-flowing-company/pprp-backend/apperror"
	"github.com/brain-flowing-company/pprp-backend/config"
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/brain-flowing-company/pprp-backend/utils"
	"go.uber.org/zap"
)

type Service interface {
	AuthenticateUser(email, password string) (string, *apperror.AppError)
}

type serviceImpl struct {
	repo   Repository
	cfg    *config.Config
	logger *zap.Logger
}

func NewService(repo Repository, cfg *config.Config, logger *zap.Logger) Service {
	return &serviceImpl{
		repo,
		cfg,
		logger,
	}
}

func (s *serviceImpl) AuthenticateUser(email, password string) (string, *apperror.AppError) {
	// Retrieve user by email
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return "", apperror.UserNotFound
	}

	// Check password
	if !utils.ComparePassword(user.Password, password) {
		return "", apperror.InvalidCredentials
	}

	session := &models.Session{
		Email: user.Email,
	}

	token, err := utils.CreateJwtToken(*session, time.Duration(s.cfg.SessionExpire*int(time.Second)), s.cfg.JWTSecret)
	if err != nil {
		s.logger.Error("Could not create JWT token", zap.Error(err))
		return "", apperror.InternalServerError
	}

	return token, nil
}
