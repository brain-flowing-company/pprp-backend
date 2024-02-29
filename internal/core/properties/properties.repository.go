package properties

import (
	"fmt"

	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"gorm.io/gorm"
)

type Repository interface {
	GetAllProperties(*[]models.Properties) error
	GetPropertyById(*models.Properties, string) error
	GetPropertyByOwnerId(*[]models.Properties, string) error
	CreateProperty(*models.Properties) error
	UpdatePropertyById(*models.Properties, string) error
	DeletePropertyById(string) error
	CountProperty(*int64, string) error
	SearchProperties(*[]models.Properties, string) error
	AddFavoriteProperty(*models.FavoriteProperties) error
	RemoveFavoriteProperty(string, string) error
	GetFavoritePropertiesByUserId(*[]models.Properties, string) error
	GetTop10Properties(*[]models.Properties) error
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

func (repo *repositoryImpl) GetPropertyByOwnerId(result *[]models.Properties, ownerId string) error {
	return repo.db.Model(&models.Properties{}).
		Preload("PropertyImages").
		Preload("SellingProperty").
		Preload("RentingProperty").
		Where("owner_id = ?", ownerId).
		Find(result).Error
}

func (repo *repositoryImpl) CreateProperty(property *models.Properties) error {
	return repo.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(property).Error; err != nil {
			fmt.Println("Property")
			return err
		}

		if len(property.PropertyImages) != 0 {
			if err := tx.Create(&property.PropertyImages).Error; err != nil {
				fmt.Println("Image")
				return err
			}
		}

		if err := tx.Create(&property.SellingProperty).Error; err != nil {
			fmt.Println("Selling")
			return err
		}

		if err := tx.Create(&property.RentingProperty).Error; err != nil {
			fmt.Println("Renting")
			return err
		}

		return nil
	})
}

func (repo *repositoryImpl) UpdatePropertyById(property *models.Properties, propertyId string) error {
	err := repo.db.First(&models.Properties{}, "property_id = ?", propertyId).Error
	if err != nil {
		return err
	}

	return repo.db.Model(&models.Properties{}).Where("property_id = ?", propertyId).Updates(property).Error
}

func (repo *repositoryImpl) DeletePropertyById(propertyId string) error {
	err := repo.db.First(&models.Properties{}, "property_id = ?", propertyId).Error
	if err != nil {
		return err
	}

	return repo.db.Where("property_id = ?", propertyId).Delete(&models.Properties{}).Error
}

func (repo *repositoryImpl) CountProperty(countProperty *int64, propertyId string) error {
	return repo.db.Model(&models.Properties{}).Where("property_id = ?", propertyId).Count(countProperty).Error
}

func (repo *repositoryImpl) SearchProperties(result *[]models.Properties, query string) error {
	return repo.db.Model(&models.Properties{}).
		Preload("PropertyImages").
		Preload("SellingProperty").
		Preload("RentingProperty").
		Where("LOWER(property_name) LIKE ? OR LOWER(property_description) LIKE ?", "%"+query+"%", "%"+query+"%").
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

func (repo *repositoryImpl) GetTop10Properties(properties *[]models.Properties) error {
	countPropertyFavorite := repo.db.Model(&models.FavoriteProperties{}).
		Select("property_id, COUNT(property_id) as favorites").
		Group("property_id")

	return repo.db.Model(&models.Properties{}).
		Preload("PropertyImages").
		Preload("SellingProperty").
		Preload("RentingProperty").
		Select("properties.*, COALESCE(count_property_favorite.favorites, 0) AS favorite_count").
		Joins("LEFT JOIN (?) AS count_property_favorite ON count_property_favorite.property_id = properties.property_id", countPropertyFavorite).
		Order("favorite_count DESC").
		Order("created_at DESC").
		Order("property_id DESC").
		Limit(10).
		Find(properties).Error
}
