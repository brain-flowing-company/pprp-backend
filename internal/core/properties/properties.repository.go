package properties

import (
	"database/sql"
	"fmt"

	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"gorm.io/gorm"
)

type Repository interface {
	GetAllProperties(*[]models.Properties, string) error
	GetPropertyById(*models.Properties, string) error
	GetPropertyByOwnerId(*[]models.Properties, string) error
	CreateProperty(*models.Properties) error
	UpdatePropertyById(*models.Properties, string) error
	DeletePropertyById(string) error
	CountProperty(*int64, string) error
	SearchProperties(*[]models.Properties, string, string) error
	AddFavoriteProperty(*models.FavoriteProperties) error
	RemoveFavoriteProperty(string, string) error
	GetFavoritePropertiesByUserId(*models.MyFavoritePropertiesResponses, string) error
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

func (repo *repositoryImpl) GetAllProperties(results *[]models.Properties, userId string) error {
	return repo.db.Model(&models.Messages{}).
		Raw(`
		SELECT props.*,
			CASE
				WHEN favorite_properties.user_id IS NOT NULL THEN TRUE
				ELSE FALSE
			END AS is_favorite
			FROM (
				SELECT properties.*,
					selling_properties.price, 
					selling_properties.is_sold,
					renting_properties.price_per_month,
					renting_properties.is_occupied
				FROM properties
				LEFT JOIN selling_properties ON properties.property_id = selling_properties.property_id
				LEFT JOIN renting_properties ON properties.property_id = renting_properties.property_id
			) AS props
			LEFT JOIN favorite_properties ON (
				favorite_properties.property_id = props.property_id AND
				favorite_properties.user_id = @user_id
			)
		`, sql.Named("user_id", userId)).
		Scan(results).Error
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
	fmt.Println(property)
	return repo.db.Transaction(func(tx *gorm.DB) error {
		propertyQuery := `INSERT INTO properties (owner_id, property_name, property_description, property_type, address, alley, street, sub_district, district, province, country, postal_code, bedrooms, bathrooms, furnishing, floor, floor_size, floor_size_unit, unit_number)
            VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`
		if err := tx.Exec(propertyQuery,
			property.OwnerId, property.PropertyName, property.PropertyDescription, property.PropertyType,
			property.Address, property.Alley, property.Street, property.SubDistrict, property.District,
			property.Province, property.Country, property.PostalCode, property.Bedrooms, property.Bathrooms,
			property.Furnishing, property.Floor, property.FloorSize, property.FloorSizeUnit, property.UnitNumber,
		).Error; err != nil {
			fmt.Println("Property")
			return err
		}

		propertyTemp := models.Properties{}
		if err := tx.Model(&models.Properties{}).Find(&propertyTemp, "property_name = ? AND owner_id = ?", property.PropertyName, property.OwnerId).Error; err != nil {
			fmt.Println("Find")
			return err
		}
		property_id := propertyTemp.PropertyId

		if len(property.PropertyImages) != 0 {
			imageQuery := `INSERT INTO property_images (property_id, image_url) VALUES (?, ?);`
			for _, image := range property.PropertyImages {
				if err := tx.Exec(imageQuery, property_id, image.ImageUrl).Error; err != nil {
					fmt.Println("Image")
					return err
				}
			}
		}

		if property.SellingProperty.Price != 0 {
			sellingQuery := `INSERT INTO selling_properties (property_id, price, is_sold) VALUES (?, ?, ?);`
			if err := tx.Exec(sellingQuery, property_id, property.SellingProperty.Price, property.SellingProperty.IsSold).Error; err != nil {
				fmt.Println("Selling")
				return err
			}
		}

		if property.RentingProperty.PricePerMonth != 0 {
			rentingQuery := `INSERT INTO renting_properties (property_id, price_per_month, is_occupied) VALUES (?, ?, ?);`
			if err := tx.Exec(rentingQuery, property_id, property.RentingProperty.PricePerMonth, property.RentingProperty.IsOccupied).Error; err != nil {
				fmt.Println("Renting")
				return err
			}
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

func (repo *repositoryImpl) SearchProperties(results *[]models.Properties, query string, userId string) error {
	return repo.db.Model(&models.Messages{}).
		Raw(`
		SELECT props.*,
			CASE
				WHEN favorite_properties.user_id IS NOT NULL THEN TRUE
				ELSE FALSE
			END AS is_favorite
			FROM (
				SELECT properties.*,
					selling_properties.price, 
					selling_properties.is_sold,
					renting_properties.price_per_month,
					renting_properties.is_occupied
				FROM properties
				LEFT JOIN selling_properties ON properties.property_id = selling_properties.property_id
				LEFT JOIN renting_properties ON properties.property_id = renting_properties.property_id
				WHERE LOWER(property_name) LIKE @query OR LOWER(property_description) LIKE @query 
			) AS props
			LEFT JOIN favorite_properties ON (
				favorite_properties.property_id = props.property_id AND
				favorite_properties.user_id = @user_id
			)
		`, sql.Named("user_id", userId), sql.Named("query", "%"+query+"%")).
		Scan(results).Error
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

func (repo *repositoryImpl) GetFavoritePropertiesByUserId(properties *models.MyFavoritePropertiesResponses, userId string) error {
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
