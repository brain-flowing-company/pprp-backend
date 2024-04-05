package ratings

import (
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
		reviewQuery := `INSERT INTO reviews (review_id, property_id, dweller_user_id, rating, review) VALUES (?,?,?,?,?)`
		if err := tx.Exec(reviewQuery, reviews.ReviewId, reviews.PropertyId, reviews.DwellerUserId, reviews.Rating, reviews.Review).Error; err != nil {
			return err
		}
		return nil
	})
}
