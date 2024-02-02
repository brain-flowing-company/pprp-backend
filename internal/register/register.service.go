package register

import (
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
	if err := user.HashPassword(); err != nil {
		return err
	}
	if err := s.repo.CreateUser(user); err != nil {
		return err
	}
	return nil
}
