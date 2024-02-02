package register

import (
	"net/mail"

	"github.com/brain-flowing-company/pprp-backend/apperror"
	"github.com/brain-flowing-company/pprp-backend/internal/models"
)

type Service interface {
	CreateUser(*models.User) error
}

type serviceImpl struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &serviceImpl{
		repo,
	}
}

func (s *serviceImpl) CreateUser(user *models.User) error {
	// Validate email
	_, err := mail.ParseAddress(user.Email)
	if err != nil {
		return apperror.InvalidEmail
	}

	existingUser, _ := s.repo.GetUserByEmail(user.Email)
	if existingUser != nil {
		return apperror.EmailAlreadyExists
	}

	if err := user.HashPassword(); err != nil {
		return err
	}
	if err := s.repo.CreateUser(user); err != nil {
		return err
	}
	return nil
}
