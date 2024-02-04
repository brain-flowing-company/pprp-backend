// internal/login/service.go
package login

import (
	"time"

	"github.com/brain-flowing-company/pprp-backend/apperror"
	"github.com/brain-flowing-company/pprp-backend/config"
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/brain-flowing-company/pprp-backend/utils"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	AuthenticateUser(email, password string) (string, *apperror.AppError)
}

type serviceImpl struct {
	repo Repository
	cfg  *config.Config
}

func NewService(repo Repository, cfg *config.Config) Service {
	return &serviceImpl{
		repo,
		cfg,
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

	token, err := utils.CreateJwtToken(*session, time.Duration(s.cfg.JWTMaxAge), s.cfg.JWTSecret)
	if err != nil {
		return "", apperror.InternalServerError
	}

	return token, nil
}

func (s *serviceImpl) checkPassword(user *models.Users, password string) *apperror.AppError {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return apperror.InvalidCredentials
	}
	return nil
}
