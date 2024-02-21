package users

import (
	"errors"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/brain-flowing-company/pprp-backend/apperror"
	"github.com/brain-flowing-company/pprp-backend/config"
	"github.com/brain-flowing-company/pprp-backend/internal/enums"
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/brain-flowing-company/pprp-backend/storage"
	"github.com/brain-flowing-company/pprp-backend/utils"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Service interface {
	GetAllUsers(*[]models.Users) *apperror.AppError
	GetUserById(*models.Users, string) *apperror.AppError
	Register(*models.RegisteringUser, *multipart.FileHeader) *apperror.AppError
	UpdateUser(*models.UpdatingUserPersonalInfo, *multipart.FileHeader) *apperror.AppError
	DeleteUser(string) *apperror.AppError
	GetUserByEmail(*models.Users, string) *apperror.AppError
}

type serviceImpl struct {
	repo    Repository
	logger  *zap.Logger
	storage storage.Storage
	cfg     *config.Config
}

func NewService(logger *zap.Logger, cfg *config.Config, repo Repository, storage storage.Storage) Service {
	return &serviceImpl{
		repo,
		logger,
		storage,
		cfg,
	}
}

func (s *serviceImpl) GetAllUsers(users *[]models.Users) *apperror.AppError {
	err := s.repo.GetAllUsers(users)

	if err != nil {
		s.logger.Error("Could not get all users", zap.Error(err))
		return apperror.
			New(apperror.InternalServerError).
			Describe("Could not get all users. Please try again later.")
	}

	return nil
}

func (s *serviceImpl) GetUserById(user *models.Users, userId string) *apperror.AppError {
	if !utils.IsValidUUID(userId) {
		return apperror.
			New(apperror.InvalidUserId).
			Describe("Invalid user id")
	}

	err := s.repo.GetUserById(user, userId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return apperror.
			New(apperror.UserNotFound).
			Describe("Could not find the specified user")
	} else if err != nil {
		s.logger.Error("Could not get user by id", zap.String("id", userId), zap.Error(err))
		return apperror.
			New(apperror.InternalServerError).
			Describe("Could not get user information. Please try again later.")
	}

	return nil
}

func (s *serviceImpl) Register(user *models.RegisteringUser, profileImage *multipart.FileHeader) *apperror.AppError {
	var countEmail int64
	if s.repo.CountEmail(&countEmail, user.Email) != nil {
		return apperror.
			New(apperror.InternalServerError).
			Describe("Could not get all emails")
	}

	if countEmail > 0 {
		return apperror.
			New(apperror.EmailAlreadyExists).
			Describe("Email already exists")
	}

	var countPhoneNumber int64
	if s.repo.CountPhoneNumber(&countPhoneNumber, user.PhoneNumber) != nil {
		return apperror.
			New(apperror.InternalServerError).
			Describe("Could not get all phone numbers")
	}

	if countPhoneNumber > 0 {
		return apperror.
			New(apperror.PhoneNumberAlreadyExists).
			Describe("Phone number already exists")
	}

	if user.RegisteredType == enums.EMAIL {
		if !utils.IsValidEmail(user.Email) {
			return apperror.
				New(apperror.InvalidEmail).
				Describe("Invalid email format")
		}

		if !utils.IsValidPassword(user.Password) {
			return apperror.
				New(apperror.InvalidPassword).
				Describe("Password should longer than 8 characters and contain alphabet and numeric characters")
		}

		hashedPassword, hashErr := utils.HashPassword(user.Password)
		if hashErr != nil {
			s.logger.Error("Could not create user", zap.Error(hashErr))
			return apperror.
				New(apperror.InternalServerError).
				Describe("Could not create user. Please try again later")
		}

		user.Password = string(hashedPassword)
	}

	url, apperr := s.uploadProfileImage(user.UserId, profileImage)
	if apperr != nil {
		return apperr
	}

	user.ProfileImageUrl = url

	err := s.repo.CreateUser(user)
	if err != nil {
		s.logger.Error("Could not create user", zap.Error(err))
		return apperror.
			New(apperror.InternalServerError).
			Describe("Could not create user. Please try again later")
	}

	return nil
}

func (s *serviceImpl) UpdateUser(user *models.UpdatingUserPersonalInfo, profileImage *multipart.FileHeader) *apperror.AppError {
	url, apperr := s.uploadProfileImage(user.UserId, profileImage)
	if apperr != nil {
		return apperr
	}
	user.ProfileImageUrl = url

	err := s.repo.UpdateUser(user, user.UserId.String())
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return apperror.
			New(apperror.UserNotFound).
			Describe("Could not find the specified user")
	} else if err != nil {
		s.logger.Error("Could not update user info", zap.String("id", user.UserId.String()), zap.Error(err))
		return apperror.
			New(apperror.InternalServerError).
			Describe("Could not update user information. Please try again later")
	}

	return nil
}

func (s *serviceImpl) DeleteUser(userId string) *apperror.AppError {
	if !utils.IsValidUUID(userId) {
		return apperror.
			New(apperror.InvalidUserId).
			Describe("Invalid user id")
	}

	err := s.repo.DeleteUser(userId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return apperror.
			New(apperror.UserNotFound).
			Describe("Could not find specified user")
	} else if err != nil {
		s.logger.Error("Could not delete user", zap.String("id", userId), zap.Error(err))
		return apperror.
			New(apperror.InternalServerError).
			Describe("Could not delete user. Please try again later")
	}

	return nil
}

func (s *serviceImpl) GetUserByEmail(user *models.Users, email string) *apperror.AppError {
	// Actaully, this shouldn't trigger unless data in database is somehow fucked
	if !utils.IsValidEmail(email) {
		s.logger.Error("Invalid email format", zap.String("email", email))
		return apperror.
			New(apperror.InternalServerError).
			Describe("Invalid email format. Try re-logging in")
	}

	// Same here
	err := s.repo.GetUserByEmail(user, email)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.Error("Could not find specified user", zap.String("email", email), zap.Error(err))
		return apperror.
			New(apperror.InternalServerError).
			Describe("Could not find specified user. Try re-logging in")
	} else if err != nil {
		s.logger.Error("Could not get current user info", zap.String("email", email), zap.Error(err))
		return apperror.
			New(apperror.InternalServerError).
			Describe("Could not get user information. Please try again later.")
	}

	return nil
}

func (s *serviceImpl) uploadProfileImage(userId uuid.UUID, profileImage *multipart.FileHeader) (string, *apperror.AppError) {
	if profileImage == nil {
		return "", nil
	}

	file, err := profileImage.Open()
	if err != nil {
		return "", apperror.
			New(apperror.InternalServerError).
			Describe("Could not upload profile image")
	}

	ext := filepath.Ext(profileImage.Filename)

	ip := utils.NewImageProcessor()

	switch strings.ToLower(ext) {
	case ".jpg":
		fallthrough

	case ".jpeg":
		err = ip.LoadJPEG(file)

	case ".png":
		err = ip.LoadPNG(file)

	default:
		return "", apperror.
			New(apperror.InvalidProfileImageExtension).
			Describe(fmt.Sprintf("App does not support %v extension", ext))
	}

	if err != nil {
		s.logger.Error("Could not load image", zap.Error(err))
		return "", apperror.
			New(apperror.InternalServerError).
			Describe("Could not process image")
	}

	err = ip.Resize(1024)
	if err != nil {
		s.logger.Error("Could not resize image", zap.Error(err))
		return "", apperror.
			New(apperror.InternalServerError).
			Describe("Could not process image")
	}

	err = ip.SquareCropped()
	if err != nil {
		s.logger.Error("Could not sqaure crop image", zap.Error(err))
		return "", apperror.
			New(apperror.InternalServerError).
			Describe("Could not process image")
	}

	processedFile, err := ip.Save()
	if err != nil {
		s.logger.Error("Could not create new image", zap.Error(err))
		return "", apperror.
			New(apperror.InternalServerError).
			Describe("Could not process image")
	}

	url, err := s.storage.Upload(fmt.Sprintf("profiles/%v.jpeg", userId.String()), processedFile)
	if err != nil {
		return "", apperror.
			New(apperror.InternalServerError).
			Describe("Could not upload profile image")
	}

	return url, nil
}
