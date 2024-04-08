package users

import (
	"errors"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/brain-flowing-company/pprp-backend/apperror"
	"github.com/brain-flowing-company/pprp-backend/config"
	"github.com/brain-flowing-company/pprp-backend/internal/enums"
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/brain-flowing-company/pprp-backend/internal/utils"
	"github.com/brain-flowing-company/pprp-backend/storage"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Service interface {
	GetAllUsers(*[]models.Users) *apperror.AppError
	GetUserById(*models.Users, string) *apperror.AppError
	GetUserFinancialInforamtionById(*models.UserFinancialInformations, string) *apperror.AppError
	Register(*models.RegisteringUsers, *multipart.FileHeader) *apperror.AppError
	UpdateUser(*models.UpdatingUserPersonalInfos, *multipart.FileHeader) *apperror.AppError
	UpdateUserFinancialInformationById(*models.UserFinancialInformations, string) *apperror.AppError
	DeleteUser(string) *apperror.AppError
	GetUserByEmail(*models.Users, string) *apperror.AppError
	VerifyCitizenId(*models.UserVerifications, *multipart.FileHeader) *apperror.AppError
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

func (s *serviceImpl) GetUserFinancialInforamtionById(userFinancialInformation *models.UserFinancialInformations, userId string) *apperror.AppError {
	if !utils.IsValidUUID(userId) {
		return apperror.
			New(apperror.InvalidUserId).
			Describe("Invalid user id")
	}

	err := s.repo.GetUserFinancialInforamtionById(userFinancialInformation, userId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return apperror.
			New(apperror.UserNotFound).
			Describe("Could not find the specified user")
	} else if err != nil {
		s.logger.Error("Could not get user financial information by id", zap.String("id", userId), zap.Error(err))
		return apperror.
			New(apperror.InternalServerError).
			Describe("Could not get user financial information. Please try again later.")
	}

	return nil
}

func (s *serviceImpl) Register(user *models.RegisteringUsers, profileImage *multipart.FileHeader) *apperror.AppError {
	validate, validatorErr := utils.NewUserValidator()
	if validatorErr != nil {
		s.logger.Error("Could not create new user validator", zap.Error(validatorErr))
		return apperror.
			New(apperror.InternalServerError).
			Describe("Could not create user. Please try again later")
	}

	if err := validate.Struct(user); err != nil {
		s.logger.Error("Invalid user information", zap.Error(err))
		return apperror.
			New(apperror.BadRequest).
			Describe("Invalid user information")
	}

	var countEmail int64
	if s.repo.CountEmail(&countEmail, user.Email) != nil {
		return apperror.
			New(apperror.InternalServerError).
			Describe("Could not get all emails")
	} else if countEmail > 0 {
		return apperror.
			New(apperror.EmailAlreadyExists).
			Describe("Email already exists")
	}

	var countPhoneNumber int64
	if s.repo.CountPhoneNumber(&countPhoneNumber, user.UserId, user.PhoneNumber) != nil {
		return apperror.
			New(apperror.InternalServerError).
			Describe("Could not get all phone numbers")
	} else if countPhoneNumber > 0 {
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

func (s *serviceImpl) UpdateUser(user *models.UpdatingUserPersonalInfos, profileImage *multipart.FileHeader) *apperror.AppError {
	url, apperr := s.uploadProfileImage(user.UserId, profileImage)
	if apperr != nil {
		return apperr
	}
	user.ProfileImageUrl = url

	if user.PhoneNumber != "" {
		if user.PhoneNumber[0] != '0' || len(user.PhoneNumber) != 10 {
			return apperror.
				New(apperror.InvalidPhoneNumber).
				Describe("Invalid phone number format")
		}

		var count int64
		countErr := s.repo.CountPhoneNumber(&count, user.UserId, user.PhoneNumber)
		if countErr != nil {
			return apperror.
				New(apperror.InternalServerError).
				Describe("Could not count phone numbers")
		} else if count > 0 {
			return apperror.
				New(apperror.PhoneNumberAlreadyExists).
				Describe("Phone number already exists")
		}
	}

	err := s.repo.UpdateUserById(user, user.UserId.String())
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

func (s *serviceImpl) UpdateUserFinancialInformationById(userFinancialInformation *models.UserFinancialInformations, userId string) *apperror.AppError {
	if !utils.IsValidUUID(userId) {
		return apperror.
			New(apperror.InvalidUserId).
			Describe("Invalid user id")
	}

	creditCards := &userFinancialInformation.CreditCards
	for i := range *creditCards {
		(*creditCards)[i].UserId, _ = uuid.Parse(userId)
	}

	validate, validatorErr := utils.NewUserFinancialInformationValidator()
	if validatorErr != nil {
		s.logger.Error("Could not create new user financial information validator", zap.Error(validatorErr))
		return apperror.
			New(apperror.InternalServerError).
			Describe("Could not update user financial information. Please try again later")
	}

	if err := validate.Struct(userFinancialInformation); err != nil {
		s.logger.Error("Invalid user financial information", zap.Error(err))
		return apperror.
			New(apperror.BadRequest).
			Describe("Invalid user financial information")
	}

	err := s.repo.UpdateUserFinancialInformationById(userFinancialInformation, userId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return apperror.
			New(apperror.UserNotFound).
			Describe("Could not find the specified user")
	} else if err != nil {
		s.logger.Error("Could not update user financial information", zap.String("id", userId), zap.Error(err))
		return apperror.
			New(apperror.InternalServerError).
			Describe("Could not update user financial information. Please try again later")
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
	url, err := s.storage.Upload(fmt.Sprintf("profiles/%v.jpeg", userId.String()), processedFile, types.ObjectCannedACLPublicRead)

	if err != nil {
		return "", apperror.
			New(apperror.InternalServerError).
			Describe("Could not upload profile image")
	}

	return url, nil
}

func (s *serviceImpl) VerifyCitizenId(user *models.UserVerifications, profileImage *multipart.FileHeader) *apperror.AppError {
	var cnt int64
	err := s.repo.CountUserVerification(&cnt, user.UserId)
	if err != nil {
		s.logger.Error("Could not count user verification", zap.String("id", user.UserId.String()), zap.Error(err))
		return apperror.
			New(apperror.InternalServerError).
			Describe("Could not verify user. Please try again later")
	}

	if cnt > 0 {
		return apperror.
			New(apperror.UserHasVerified).
			Describe("User has already verified")
	}

	url, apperr := s.uploadCitizenImage(user.UserId, profileImage)
	if apperr != nil {
		return apperr
	}

	user.CitizenCardImageUrl = url
	user.VerifiedAt = time.Now()

	err = s.repo.CreateUserVerification(user)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return apperror.
			New(apperror.UserNotFound).
			Describe("Could not find the specified user")
	} else if err != nil {
		s.logger.Error("Could not verify user", zap.String("id", user.UserId.String()), zap.Error(err))
		return apperror.
			New(apperror.InternalServerError).
			Describe("Could not verify user. Please try again later")
	}

	return nil
}

func (s *serviceImpl) uploadCitizenImage(userId uuid.UUID, profileImage *multipart.FileHeader) (string, *apperror.AppError) {
	if profileImage == nil {
		return "", apperror.
			New(apperror.BadRequest).
			Describe("No citizen card found")
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

	processedFile, err := ip.Save()
	if err != nil {
		s.logger.Error("Could not create new image", zap.Error(err))
		return "", apperror.
			New(apperror.InternalServerError).
			Describe("Could not process image")
	}

	url, err := s.storage.Upload(fmt.Sprintf("verifications/%v.jpeg", userId.String()), processedFile, types.ObjectCannedACLPrivate)

	if err != nil {
		return "", apperror.
			New(apperror.InternalServerError).
			Describe("Could not upload profile image")
	}

	return url, nil
}
