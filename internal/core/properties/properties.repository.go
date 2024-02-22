package properties

import (
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"gorm.io/gorm"
)

type Repository interface {
	GetAllProperties(*[]models.Properties) error
	GetPropertyById(*models.Properties, string) error
	CreateProperty(*models.Properties) error
	UpdatePropertyById(*models.Properties) error
	DeletePropertyById(string) error
	SearchProperties(*[]models.Properties, string) error
	AddFavoriteProperty(*models.FavoriteProperties) error
	RemoveFavoriteProperty(string, string) error
	GetFavoritePropertiesByUserId(*[]models.Properties, string) error
}

type repositoryImpl struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repositoryImpl{
		db,
	}
}

func (repo *repositoryImpl) GetAllProperties(result *[]models.Properties) error {
	return repo.db.Model(&models.Properties{}).
		Preload("PropertyImages").
		Preload("SellingProperty").
		Preload("RentingProperty").
		Find(result).Error
}

func (repo *repositoryImpl) GetPropertyById(result *models.Properties, id string) error {
	return repo.db.Model(&models.Properties{}).
		Preload("PropertyImages").
		Preload("SellingProperty").
		Preload("RentingProperty").
		First(result, "property_id = ?", id).Error
}

func (repo *repositoryImpl) CreateProperty(property *models.Properties) error {
	return repo.db.Create(property).Error
}

func (repo *repositoryImpl) UpdatePropertyById(property *models.Properties) error {
	return repo.db.Save(property).Error
}

func (repo *repositoryImpl) DeletePropertyById(propertyId string) error {
	err := repo.db.First(&models.Properties{}, "property_id = ?", propertyId).Error
	if err != nil {
		return err
	}

	return repo.db.Where("property_id = ?", propertyId).Delete(&models.Properties{}).Error
}

func (repo *repositoryImpl) SearchProperties(result *[]models.Properties, query string) error {
	return repo.db.Model(&models.Properties{}).
		Preload("PropertyImages").
		Preload("SellingProperty").
		Preload("RentingProperty").
		Where("LOWER(project_name) LIKE ? OR LOWER(description) LIKE ?", "%"+query+"%", "%"+query+"%").
		Find(result).Error
}

func (repo *repositoryImpl) AddFavoriteProperty(favoriteProperty *models.FavoriteProperties) error {
	return repo.db.Create(favoriteProperty).Error
}

func (repo *repositoryImpl) RemoveFavoriteProperty(propertyId string, userId string) error {
	err := repo.db.First(&models.FavoriteProperties{}, "property_id = ? AND user_id = ?", propertyId, userId).Error
	if err != nil {
		return err
	}

	return repo.db.Where("property_id = ? AND user_id = ?", propertyId, userId).Delete(&models.FavoriteProperties{}).Error
}

func (repo *repositoryImpl) GetFavoritePropertiesByUserId(properties *[]models.Properties, userId string) error {
	return repo.db.Model(&models.Properties{}).
		Preload("PropertyImages").
		Preload("SellingProperty").
		Preload("RentingProperty").
		Joins("JOIN favorite_properties ON favorite_properties.property_id = properties.property_id").
		Where("favorite_properties.user_id = ?", userId).
		Find(properties).Error
}
