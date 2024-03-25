package users

import (
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	GetAllUsers(*[]models.Users) error
	GetUserById(*models.Users, string) error
	GetUserFinancialInforamtionById(*models.UserFinancialInformations, string) error
	GetUserByEmail(*models.Users, string) error
	CreateUser(*models.RegisteringUsers) error
	UpdateUserById(*models.UpdatingUserPersonalInfos, string) error
	UpdateUserFinancialInformationById(*models.UserFinancialInformations, string) error
	DeleteUser(string) error
	CountEmail(*int64, string) error
	CountPhoneNumber(*int64, uuid.UUID, string) error
	CreateUserVerification(*models.UserVerifications) error
	CountUserVerification(cnt *int64, userId uuid.UUID) error
}

type repositoryImpl struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repositoryImpl{
		db,
	}
}

func (repo *repositoryImpl) GetAllUsers(users *[]models.Users) error {
	return repo.db.Find(users).Error
}

func (repo *repositoryImpl) GetUserById(user *models.Users, userId string) error {
	return repo.db.First(user, "user_id = ?", userId).Error
}

func (repo *repositoryImpl) GetUserFinancialInforamtionById(userFinancialInformation *models.UserFinancialInformations, userId string) error {
	return repo.db.Model(&models.UserFinancialInformations{}).
		Preload("CreditCards").
		First(userFinancialInformation, "user_id = ?", userId).Error
}

func (repo *repositoryImpl) GetUserByEmail(user *models.Users, email string) error {
	return repo.db.First(user, "email = ?", email).Error
}

func (repo *repositoryImpl) CreateUser(user *models.RegisteringUsers) error {
	return repo.db.Create(user).Error
}

func (repo *repositoryImpl) UpdateUserById(user *models.UpdatingUserPersonalInfos, userId string) error {
	err := repo.db.First(&models.Users{}, "user_id = ?", userId).Error
	if err != nil {
		return err
	}

	return repo.db.Where("user_id = ?", userId).Updates(user).Error
}

func (repo *repositoryImpl) UpdateUserFinancialInformationById(userFinancialInformation *models.UserFinancialInformations, userId string) error {
	err := repo.db.First(&models.UserFinancialInformations{}, "user_id = ?", userId).Error
	if err != nil {
		return err
	}

	return repo.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ?", userId).Delete(&models.CreditCards{}).Error; err != nil {
			return err
		} else if len(userFinancialInformation.CreditCards) != 0 {
			if err := tx.Create(&userFinancialInformation.CreditCards).Error; err != nil {
				return err
			}
		}

		return tx.Model(&models.UserFinancialInformations{}).Where("user_id = ?", userId).Updates(userFinancialInformation).Error
	})
}

func (repo *repositoryImpl) DeleteUser(userId string) error {
	err := repo.db.First(&models.Users{}, "user_id = ?", userId).Error
	if err != nil {
		return err
	}

	return repo.db.Where("user_id = ?", userId).Delete(&models.Users{}).Error
}

func (repo *repositoryImpl) CountEmail(count *int64, email string) error {
	return repo.db.Model(&models.Users{}).Where("email = ?", email).Count(count).Error
}

func (repo *repositoryImpl) CountPhoneNumber(count *int64, userId uuid.UUID, phoneNumber string) error {
	return repo.db.Model(&models.Users{}).
		Where("phone_number = ? AND user_id != ?", phoneNumber, userId).
		Count(count).Error
}

func (repo *repositoryImpl) CountUserVerification(cnt *int64, userId uuid.UUID) error {
	return repo.db.Model(&models.UserVerifications{}).
		Where("user_id = ?", userId).
		Count(cnt).Error
}

func (repo *repositoryImpl) CreateUserVerification(user *models.UserVerifications) error {
	return repo.db.Model(&models.UserVerifications{}).Create(user).Error
}
