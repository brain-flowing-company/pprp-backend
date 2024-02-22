package properties

import (
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"gorm.io/gorm"
)

type Repository interface {
	GetPropertyById(*models.Properties, string) error
	GetAllProperties(*[]models.Properties) error
	SearchProperties(*[]models.Properties, string) error
}

type repositoryImpl struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repositoryImpl{
		db,
	}
}

func (repo *repositoryImpl) GetPropertyById(result *models.Properties, id string) error {
	return repo.db.Model(&models.Properties{}).
		Preload("PropertyImages").
		Preload("SellingProperty").
		Preload("RentingProperty").
		First(result, "property_id = ?", id).Error
}

func (repo *repositoryImpl) GetAllProperties(result *[]models.Properties) error {
	return repo.db.Model(&models.Properties{}).
		Preload("PropertyImages").
		Preload("SellingProperty").
		Preload("RentingProperty").
		Find(result).Error
}

func (repo *repositoryImpl) SearchProperties(result *[]models.Properties, query string) error {
	return repo.db.Model(&models.Properties{}).
		Preload("PropertyImages").
		Preload("SellingProperty").
		Preload("RentingProperty").
		Where("LOWER(project_name) LIKE ? OR LOWER(description) LIKE ?", "%"+query+"%", "%"+query+"%").
		Find(result).Error
}
