package property

import (
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"gorm.io/gorm"
)

type Repository interface {
	GetPropertyById(*models.Property, string) error
}

type repositoryImpl struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repositoryImpl{
		db,
	}
}

func (repo *repositoryImpl) GetPropertyById(result *models.Property, id string) error {
	return repo.db.Model(&models.Property{}).
		Preload("PropertyImages").
		Preload("SellingProperty").
		Preload("RentingProperty").
		First(result, "property_id = ?", id).Error
}
