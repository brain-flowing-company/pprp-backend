package payments

import (
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"gorm.io/gorm"
)

type Repository interface {
	CreatePayment(*models.Payments) error
}

type repositoryImpl struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repositoryImpl{
		db,
	}
}

func (r *repositoryImpl) CreatePayment(payment *models.Payments) error {
	result := r.db.Create(payment)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
