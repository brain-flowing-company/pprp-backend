package ratings

import (
	"github.com/brain-flowing-company/pprp-backend/apperror"
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"gorm.io/gorm"
)

type Repository interface {
	CreateRating(*models.Reviews) error
}

type repositoryImpl struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repositoryImpl{
		db,
	}
}

func (r *repositoryImpl) CreateRating(reviews *models.Reviews) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var userCount int64
		if err := tx.Model(&models.Users{}).Where("user_id = ?", reviews.DwellerUserId).Count(&userCount).Error; err != nil {
			return err
		}
		if userCount == 0 {
			return apperror.
				New(apperror.UserNotFound).
				Describe("FK constraint in user table")
		}
		var propertyCount int64
		if err := tx.Model(&models.Properties{}).Where("property_id = ?", reviews.PropertyId).Count(&propertyCount).Error; err != nil {
			return err
		}
		if propertyCount == 0 {
			return apperror.
				New(apperror.PropertyNotFound).
				Describe("FK constraint in property table")
		}
		reviewQuery := `INSERT INTO review (review_id, property_id, dweller_user_id, rating, review) VALUES (?,?,?,?,?)`
		if err := tx.Exec(reviewQuery, reviews.ReviewId, reviews.PropertyId, reviews.DwellerUserId, reviews.Rating, reviews.Review).Error; err != nil {
			return err
		}
		return nil
	})
}
