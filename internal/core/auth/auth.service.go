// internal/login/service.go
package auth

import (
	"context"
	"time"

	"github.com/brain-flowing-company/pprp-backend/apperror"
	"github.com/brain-flowing-company/pprp-backend/config"
	"github.com/brain-flowing-company/pprp-backend/internal/core/emails"
	"github.com/brain-flowing-company/pprp-backend/internal/core/google"
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/brain-flowing-company/pprp-backend/internal/utils"
	"go.uber.org/zap"
)

type Service interface {
	AuthenticateUser(email, password string) (string, *apperror.AppError)
	Callback(ctx context.Context, callback *models.Callbacks, callbackResponse *models.CallbackResponses) *apperror.AppError
}

type serviceImpl struct {
	repo          Repository
	cfg           *config.Config
	logger        *zap.Logger
	googleService google.Service
	emailService  emails.Service
}

func NewService(logger *zap.Logger, cfg *config.Config, repo Repository, googleService google.Service, emailService emails.Service) Service {
	return &serviceImpl{
		repo,
		cfg,
		logger,
		googleService,
		emailService,
	}
}

func (s *serviceImpl) AuthenticateUser(email, password string) (string, *apperror.AppError) {
	// Retrieve user by email
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return "", apperror.
			New(apperror.UserNotFound).
			Describe("User does not exist")
	}

	// Check password
	if !utils.ComparePassword(user.Password, password) {
		return "", apperror.
			New(apperror.InvalidCredentials).
			Describe("Credentials do not match")
	}

	session := models.Sessions{
		UserId: user.UserId,
		Email:  email,
	}

	token, err := utils.CreateJwtToken(session, time.Duration(s.cfg.SessionExpire*int(time.Second)), s.cfg.JWTSecret)
	if err != nil {
		s.logger.Error("Could not create JWT token", zap.Error(err))
		return "", apperror.
			New(apperror.InternalServerError).
			Describe("Could not login. Please try again later")
	}

	return token, nil
}

func (s *serviceImpl) Callback(ctx context.Context, callback *models.Callbacks, callbackResponse *models.CallbackResponses) *apperror.AppError {
	var err *apperror.AppError

	prefix := s.cfg.EmailCodePrefix
	if callback.Code[:len(prefix)] == s.cfg.EmailCodePrefix {
		err = s.emailService.VerifyEmail(callback, callbackResponse)
	} else {
		err = s.googleService.ExchangeToken(ctx, callback, callbackResponse)
	}

	return err
}
