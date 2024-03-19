package properties

import (
	"database/sql"
	"fmt"

	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/brain-flowing-company/pprp-backend/internal/utils"
	"gorm.io/gorm"
)

type Repository interface {
	GetAllProperties(*models.AllPropertiesResponses, string, string, *utils.PaginatedQuery, *utils.SortedQuery) error
	GetPropertyById(*models.Properties, string) error
	GetPropertyByOwnerId(*models.MyPropertiesResponses, string, *utils.PaginatedQuery) error
	CreateProperty(*models.PropertyInfos) error
	UpdatePropertyById(*models.PropertyInfos, string) error
	DeletePropertyById(string) error
	CountProperty(*int64, string) error
	CountPropertyImages(*int64, string) error
	AddFavoriteProperty(*models.FavoriteProperties) error
	RemoveFavoriteProperty(string, string) error
	GetFavoritePropertiesByUserId(*models.MyFavoritePropertiesResponses, string, *utils.PaginatedQuery) error
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

func (repo *repositoryImpl) GetAllProperties(properties *models.AllPropertiesResponses, query string, userId string, paginated *utils.PaginatedQuery, sorted *utils.SortedQuery) error {
	return repo.db.Transaction(func(tx *gorm.DB) error {
		if err := repo.db.Model(&models.Properties{}).
			Raw(`
				SELECT COUNT(*) AS total
				FROM (
					SELECT properties.*
					FROM properties
					WHERE LOWER(property_name) LIKE @query OR LOWER(property_description) LIKE @query
				) AS props
				LEFT JOIN favorite_properties ON (
					favorite_properties.property_id = props.property_id AND
					favorite_properties.user_id = @user_id
				)
				`,
				sql.Named("user_id", userId),
				sql.Named("query", "%"+query+"%")).
			First(&properties.Total).Error; err != nil {
			return err
		}

		if err := repo.db.Model(&models.Properties{}).
			Raw(`
				SELECT
					props.*,
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
				)`+sorted.SortedSQL()+" "+paginated.PaginatedSQL(),
				sql.Named("user_id", userId),
				sql.Named("query", "%"+query+"%")).
			Scan(&properties.Properties).Error; err != nil {
			return err
		}

		for i, property := range properties.Properties {
			if err := repo.db.Model(&models.PropertyImages{}).
				Raw(`
					SELECT image_url 
					FROM property_images 
					WHERE property_id = @property_id`,
					sql.Named("property_id", property.PropertyId)).
				Pluck("image_url", &properties.Properties[i].PropertyImages).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (repo *repositoryImpl) GetPropertyById(property *models.Properties, propertyId string) error {
	return repo.db.Transaction(func(tx *gorm.DB) error {
		if err := repo.db.Model(&models.Properties{}).
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
			Scan(property).Error; err != nil {
			return err
		}

		if err := repo.db.Model(&models.PropertyImages{}).
			Raw(`
				SELECT image_url
				FROM property_images
				WHERE property_id = @property_id
				`, sql.Named("property_id", property.PropertyId)).
			Pluck("image_url", &property.PropertyImages).Error; err != nil {
			return err
		}

		return nil
	})

}

func (repo *repositoryImpl) GetPropertyByOwnerId(properties *models.MyPropertiesResponses, ownerId string, paginated *utils.PaginatedQuery) error {
	return repo.db.Transaction(func(tx *gorm.DB) error {
		if err := repo.db.Model(&models.Properties{}).
			Raw(`
				SELECT COUNT(*) AS total
				FROM (
					SELECT properties.*
					FROM properties
					WHERE properties.owner_id = @user_id
				) AS props
				LEFT JOIN favorite_properties ON (
					favorite_properties.property_id = props.property_id AND
					favorite_properties.user_id = @user_id
				)
			`, sql.Named("user_id", ownerId)).
			First(&properties.Total).Error; err != nil {
			return err
		}

		if err := repo.db.Model(&models.Properties{}).
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
				LIMIT @limit OFFSET @offset
				`,
				sql.Named("owner_id", ownerId),
				sql.Named("limit", paginated.Limit),
				sql.Named("offset", paginated.Offset)).
			Scan(&properties.Properties).Error; err != nil {
			return err
		}

		for i, property := range properties.Properties {
			if err := repo.db.Model(&models.PropertyImages{}).
				Raw(`
					SELECT image_url 
					FROM property_images 
					WHERE property_id = @property_id`,
					sql.Named("property_id", property.PropertyId)).
				Pluck("image_url", &properties.Properties[i].PropertyImages).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (repo *repositoryImpl) CreateProperty(property *models.PropertyInfos) error {
	return repo.db.Transaction(func(tx *gorm.DB) error {
		propertyQuery := `INSERT INTO properties (property_id, owner_id, property_name, property_description, property_type, address, alley, street, sub_district, district, province, country, postal_code, bedrooms, bathrooms, furnishing, floor, floor_size, floor_size_unit, unit_number)
            VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`
		if err := tx.Exec(propertyQuery, property.PropertyId,
			property.OwnerId, property.PropertyName, property.PropertyDescription, property.PropertyType,
			property.Address, property.Alley, property.Street, property.SubDistrict, property.District,
			property.Province, property.Country, property.PostalCode, property.Bedrooms, property.Bathrooms,
			property.Furnishing, property.Floor, property.FloorSize, property.FloorSizeUnit, property.UnitNumber,
		).Error; err != nil {
			return err
		}

		if len(property.ImageUrls) != 0 {
			imageQuery := `INSERT INTO property_images (property_id, image_url) VALUES (?, ?);`
			for _, imageUrl := range property.ImageUrls {
				if err := tx.Exec(imageQuery, property.PropertyId, imageUrl).Error; err != nil {
					return err
				}
			}
		}

		if property.Price != 0 {
			sellingQuery := `INSERT INTO selling_properties (property_id, price, is_sold) VALUES (?, ?, ?);`
			if err := tx.Exec(sellingQuery, property.PropertyId, property.Price, property.IsSold).Error; err != nil {
				return err
			}
		}

		if property.PricePerMonth != 0 {
			rentingQuery := `INSERT INTO renting_properties (property_id, price_per_month, is_occupied) VALUES (?, ?, ?);`
			if err := tx.Exec(rentingQuery, property.PropertyId, property.PricePerMonth, property.IsOccupied).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (repo *repositoryImpl) UpdatePropertyById(property *models.PropertyInfos, propertyId string) error {
	fmt.Println(property)
	return repo.db.Transaction(func(tx *gorm.DB) error {
		var existingProperty models.Properties
		if err := tx.Model(&models.Properties{}).First(&existingProperty, "property_id = ?", propertyId).Error; err != nil {
			return err
		}

		propertyQuery := `UPDATE properties SET property_name = ?, property_description = ?, property_type = ?, address = ?, alley = ?, street = ?, sub_district = ?, district = ?, province = ?, country = ?, postal_code = ?, bedrooms = ?, bathrooms = ?, furnishing = ?, floor = ?, floor_size = ?, floor_size_unit = ?, unit_number = ?, updated_at = CURRENT_TIMESTAMP WHERE property_id = ?`
		if err := tx.Exec(propertyQuery,
			property.PropertyName, property.PropertyDescription, property.PropertyType, property.Address,
			property.Alley, property.Street, property.SubDistrict, property.District, property.Province,
			property.Country, property.PostalCode, property.Bedrooms, property.Bathrooms, property.Furnishing,
			property.Floor, property.FloorSize, property.FloorSizeUnit, property.UnitNumber, propertyId,
		).Error; err != nil {
			return err
		}

		if err := tx.Where("property_id = ?", propertyId).Delete(&models.PropertyImages{}).Error; err != nil {
			return err
		} else if len(property.ImageUrls) != 0 {
			createImageQuery := `INSERT INTO property_images (property_id, image_url) VALUES (?, ?);`
			updateImageQuery := `UPDATE property_images SET deleted_at = NULL WHERE property_id = ? AND image_url = ?;`
			for _, imageUrl := range property.ImageUrls {
				if err := tx.Model(&models.PropertyImages{}).First(&models.PropertyImages{}, "property_id = ? AND image_url = ?", propertyId, imageUrl).Error; err == gorm.ErrRecordNotFound {
					if err := tx.Exec(createImageQuery, propertyId, imageUrl).Error; err != nil {
						return err
					}
				} else if err == nil {
					if err := tx.Exec(updateImageQuery, propertyId, imageUrl).Error; err != nil {
						return err
					}
				} else {
					return err
				}
			}
		}

		if property.Price != existingProperty.SellingProperty.Price || property.IsSold != existingProperty.SellingProperty.IsSold {
			sellingQuery := `UPDATE selling_properties SET price = ?, is_sold = ?, updated_at = CURRENT_TIMESTAMP WHERE property_id = ?;`
			if err := tx.Exec(sellingQuery, property.Price, property.IsSold, propertyId).Error; err != nil {
				return err
			}
		}

		if property.PricePerMonth != existingProperty.RentingProperty.PricePerMonth || property.IsOccupied != existingProperty.RentingProperty.IsOccupied {
			rentingQuery := `UPDATE renting_properties SET price_per_month = ?, is_occupied = ?, updated_at = CURRENT_TIMESTAMP WHERE property_id = ?;`
			if err := tx.Exec(rentingQuery, property.PricePerMonth, property.IsOccupied, propertyId).Error; err != nil {
				return err
			}
		}

		return nil
	})
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

func (repo *repositoryImpl) CountPropertyImages(countPropertyImages *int64, propertyId string) error {
	return repo.db.Model(&models.PropertyImages{}).Where("property_id = ?", propertyId).Count(countPropertyImages).Error
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

func (repo *repositoryImpl) GetFavoritePropertiesByUserId(properties *models.MyFavoritePropertiesResponses, userId string, paginated *utils.PaginatedQuery) error {
	return repo.db.Transaction(func(tx *gorm.DB) error {
		if err := repo.db.Model(&models.Properties{}).
			Raw(`
				SELECT COUNT(*) AS total
				FROM favorite_properties
				LEFT JOIN properties
				ON favorite_properties.property_id = properties.property_id
				WHERE favorite_properties.user_id = @user_id
			`, sql.Named("user_id", userId)).
			First(&properties.Total).Error; err != nil {
			return err
		}

		if err := repo.db.Model(&models.Properties{}).
			Raw(`
				SELECT
					props.*,
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
				LIMIT @limit OFFSET @offset
				`,
				sql.Named("user_id", userId),
				sql.Named("limit", paginated.Limit),
				sql.Named("offset", paginated.Offset)).
			Scan(&properties.Properties).Error; err != nil {
			return err
		}

		for i, property := range properties.Properties {
			if err := repo.db.Model(&models.PropertyImages{}).
				Raw(`
						SELECT image_url 
						FROM property_images 
						WHERE property_id = @property_id`,
					sql.Named("property_id", property.PropertyId)).
				Pluck("image_url", &properties.Properties[i].PropertyImages).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (repo *repositoryImpl) GetTop10Properties(properties *[]models.Properties, userId string) error {
	return repo.db.Transaction(func(tx *gorm.DB) error {
		if err := repo.db.Model(&models.Properties{}).
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
						COALESCE(count_property_favorite.favorites, 0) AS favorite_count,
						properties.created_at
					FROM properties
					LEFT JOIN (
						SELECT property_id,
							COUNT(property_id) as favorites
						FROM favorite_properties
						GROUP BY property_id
					) AS count_property_favorite ON count_property_favorite.property_id = properties.property_id
					ORDER BY favorite_count DESC, properties.created_at DESC, properties.property_id DESC
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
			Scan(properties).Error; err != nil {
			return err
		}

		for i, property := range *properties {
			if err := repo.db.Model(&models.PropertyImages{}).
				Raw(`
					SELECT image_url
					FROM property_images
					WHERE property_id = @property_id
					`, sql.Named("property_id", property.PropertyId)).
				Pluck("image_url", &(*properties)[i].PropertyImages).Error; err != nil {
				return err
			}
		}

		return nil
	})

}
