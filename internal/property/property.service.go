package property

import (
	"errors"
	"fmt"

	"github.com/brain-flowing-company/pprp-backend/apperror"
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"gorm.io/gorm"
)

type Service interface {
	GetPropertyById(*models.Property, string) *apperror.AppError
}

type serviceImpl struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &serviceImpl{
		repo,
	}
}

func (s *serviceImpl) GetPropertyById(property *models.Property, propertyId string) *apperror.AppError {
	err := s.repo.GetPropertyById(property, propertyId)
	fmt.Println(err)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return apperror.InvalidPropertyId
	} else if err != nil {
		return apperror.InternalServerError
	}

	return nil
}
