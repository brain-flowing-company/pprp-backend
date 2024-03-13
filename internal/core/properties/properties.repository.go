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
	GetFavoritePropertiesByUserId(*[]models.Properties, string) error
	GetTop10Properties(*[]models.Properties, string) error
}

type repositoryImpl struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repositoryImpl{
		db,
	}
}

func (repo *repositoryImpl) GetAllProperties(properties *[]models.Properties, userId string) error {
	return repo.db.Model(&models.Properties{}).
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
		Scan(properties).Error
}

func (repo *repositoryImpl) GetPropertyById(property *models.Properties, propertyId string) error {
	return repo.db.Model(&models.Properties{}).
		Raw(`
		SELECT properties.*,
			selling_properties.price, 
			selling_properties.is_sold,
			renting_properties.price_per_month,
			renting_properties.is_occupied
		FROM properties
		LEFT JOIN selling_properties ON properties.property_id = selling_properties.property_id
		LEFT JOIN renting_properties ON properties.property_id = renting_properties.property_id
		WHERE properties.property_id = @property_id
		`, sql.Named("property_id", propertyId)).
		Scan(property).Error
}

func (repo *repositoryImpl) GetPropertyByOwnerId(properties *[]models.Properties, ownerId string) error {
	return repo.db.Model(&models.Properties{}).
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
				WHERE properties.owner_id = @owner_id
			) AS props
			LEFT JOIN favorite_properties ON (
				favorite_properties.property_id = props.property_id AND
				favorite_properties.user_id = @owner_id
			)
		`, sql.Named("owner_id", ownerId)).
		Scan(properties).Error
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
			return err
		}
		property_id := propertyTemp.PropertyId

		if len(property.PropertyImages) != 0 {
			imageQuery := `INSERT INTO property_images (property_id, image_url) VALUES (?, ?);`
			for _, image := range property.PropertyImages {
				if err := tx.Exec(imageQuery, property_id, image.ImageUrl).Error; err != nil {
					return err
				}
			}
		}

		if property.SellingProperty.Price != 0 {
			sellingQuery := `INSERT INTO selling_properties (property_id, price, is_sold) VALUES (?, ?, ?);`
			if err := tx.Exec(sellingQuery, property_id, property.SellingProperty.Price, property.SellingProperty.IsSold).Error; err != nil {
				return err
			}
		}

		if property.RentingProperty.PricePerMonth != 0 {
			rentingQuery := `INSERT INTO renting_properties (property_id, price_per_month, is_occupied) VALUES (?, ?, ?);`
			if err := tx.Exec(rentingQuery, property_id, property.RentingProperty.PricePerMonth, property.RentingProperty.IsOccupied).Error; err != nil {
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

func (repo *repositoryImpl) SearchProperties(properties *[]models.Properties, query string, userId string) error {
	return repo.db.Model(&models.Properties{}).
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
		Scan(properties).Error
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
		Raw(`
		SELECT props.*,
			TRUE AS is_favorite
			FROM favorite_properties
			LEFT JOIN (
				SELECT properties.*,
				selling_properties.price,
				selling_properties.is_sold,
				renting_properties.price_per_month,
				renting_properties.is_occupied
				FROM properties
				LEFT JOIN selling_properties ON properties.property_id = selling_properties.property_id
				LEFT JOIN renting_properties ON properties.property_id = renting_properties.property_id
			) AS props ON favorite_properties.property_id = props.property_id
			WHERE favorite_properties.user_id = @user_id
		`, sql.Named("user_id", userId)).
		Scan(properties).Error
}

func (repo *repositoryImpl) GetTop10Properties(properties *[]models.Properties, userId string) error {
	return repo.db.Model(&models.Properties{}).
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
				FROM (
					SELECT properties.property_id,
						COALESCE(count_property_favorite.favorites, 0) AS favorite_count
					FROM properties
					LEFT JOIN (
						SELECT property_id,
							COUNT(property_id) as favorites
						FROM favorite_properties
						GROUP BY property_id
					) AS count_property_favorite ON count_property_favorite.property_id = properties.property_id
					ORDER BY favorite_count DESC, created_at DESC, property_id DESC
					LIMIT 10
				) AS top10
				LEFT JOIN properties ON top10.property_id = properties.property_id
				LEFT JOIN selling_properties ON top10.property_id = selling_properties.property_id
				LEFT JOIN renting_properties ON top10.property_id = renting_properties.property_id
			) AS props
			LEFT JOIN favorite_properties ON (
				favorite_properties.property_id = props.property_id AND
				favorite_properties.user_id = @user_id
			)
		`, sql.Named("user_id", userId)).
		Scan(properties).Error
}
