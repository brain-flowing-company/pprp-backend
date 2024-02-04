package register

import (
	"fmt"
	"net/mail"

	"github.com/brain-flowing-company/pprp-backend/apperror"
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/brain-flowing-company/pprp-backend/utils"
)

type Service interface {
	CreateUser(*models.Users) error
}

type serviceImpl struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &serviceImpl{
		repo,
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
		return apperror.InternalServerError
	}

	user.Password = hashedPassword

	if err := s.repo.CreateUser(user); err != nil {
		fmt.Println(err)
		return apperror.InternalServerError
	}
	return nil
}
