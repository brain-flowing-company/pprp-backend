package ratings

import (
	"fmt"

	"github.com/brain-flowing-company/pprp-backend/apperror"
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	CreateRating(*models.Reviews) error
	GetRatingByPropertyId(uuid.UUID, *[]models.RatingResponse) error
	GetAllRatings(*[]models.RatingResponse) error
	GetRatingByPropertyIdSortedByRating(propertyId uuid.UUID, ratings *[]models.RatingResponse) error
	GetRatingByPropertyIdSortedByNewest(propertyId uuid.UUID, ratings *[]models.RatingResponse) error
	UpdateRatingStatus(updateStatus *models.UpdateRatingStatus, ratingId uuid.UUID) error
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
		reviewQuery := `INSERT INTO reviews (review_id, property_id, dweller_user_id, rating, review) VALUES (?,?,?,?,?)`
		if err := tx.Exec(reviewQuery, reviews.ReviewId, reviews.PropertyId, reviews.DwellerUserId, reviews.Rating, reviews.Review).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *repositoryImpl) GetRatingByPropertyId(propertyId uuid.UUID, ratings *[]models.RatingResponse) error {
	fmt.Println("propertyId", propertyId)
	var propertyCount int64
	if err := r.db.Model(&models.Properties{}).Where("property_id = ?", propertyId).Count(&propertyCount).Error; err != nil {
		return err
	}
	if propertyCount == 0 {
		return apperror.
			New(apperror.PropertyNotFound).
			Describe("Property not found")
	}
	err := r.db.Table("reviews").
		Select("reviews.*, _users.first_name ,_users.last_name").
		Joins("LEFT JOIN _users ON reviews.dweller_user_id = _users.user_id").
		Where("property_id = ?", propertyId).
		Scan(&ratings).Error
	if err != nil {
		return err
	}
	return nil

}

func (r *repositoryImpl) GetAllRatings(ratings *[]models.RatingResponse) error {
	err := r.db.Table("reviews").
		Select("reviews.*, _users.first_name ,_users.last_name").
		Joins("LEFT JOIN _users ON reviews.dweller_user_id = _users.user_id").
		Scan(&ratings).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *repositoryImpl) GetRatingByPropertyIdSortedByRating(propertyId uuid.UUID, ratings *[]models.RatingResponse) error {
	fmt.Println("propertyId", propertyId)
	var propertyCount int64
	if err := r.db.Model(&models.Properties{}).Where("property_id = ?", propertyId).Count(&propertyCount).Error; err != nil {
		return err
	}
	if propertyCount == 0 {
		return apperror.
			New(apperror.PropertyNotFound).
			Describe("Property not found")
	}
	err := r.db.Table("reviews").
		Select("reviews.*, _users.first_name ,_users.last_name").
		Joins("LEFT JOIN _users ON reviews.dweller_user_id = _users.user_id").
		Where("property_id = ?", propertyId).
		Order("rating desc").
		Scan(&ratings).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *repositoryImpl) GetRatingByPropertyIdSortedByNewest(propertyId uuid.UUID, ratings *[]models.RatingResponse) error {
	fmt.Println("propertyId", propertyId)
	var propertyCount int64
	if err := r.db.Model(&models.Properties{}).Where("property_id = ?", propertyId).Count(&propertyCount).Error; err != nil {
		return err
	}
	if propertyCount == 0 {
		return apperror.
			New(apperror.PropertyNotFound).
			Describe("Property not found")
	}
	err := r.db.Table("reviews").
		Select("reviews.*, _users.first_name ,_users.last_name").
		Joins("LEFT JOIN _users ON reviews.dweller_user_id = _users.user_id").
		Where("property_id = ?", propertyId).
		Order("created_at desc").
		Scan(&ratings).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *repositoryImpl) UpdateRatingStatus(updateStatus *models.UpdateRatingStatus, ratingId uuid.UUID) error {
	if err := r.db.Model(&models.Reviews{}).First(&models.Reviews{}, "review_id = ?", ratingId).Error; err != nil {
		return err
	}
	return r.db.Model(&models.Reviews{}).Where("review_id = ?", ratingId).Updates(updateStatus).Error
}
