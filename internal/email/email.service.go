package email

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"time"

	"github.com/brain-flowing-company/pprp-backend/apperror"
	"github.com/brain-flowing-company/pprp-backend/config"
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/brain-flowing-company/pprp-backend/utils"
	"go.uber.org/zap"
)

type Service interface {
	SendVerificationEmail(string) *apperror.AppError
	VerifyEmail(string, string) *apperror.AppError
	DeleteEmailVerificationData(string) *apperror.AppError
}

type serviceImpl struct {
	repo   Repository
	logger *zap.Logger
	cfg    *config.Config
}

func NewService(logger *zap.Logger, cfg *config.Config, repo Repository) Service {
	return &serviceImpl{
		repo,
		logger,
		cfg,
	}
}

func (s *serviceImpl) SendVerificationEmail(userEmail string) *apperror.AppError {

	if !utils.IsValidEmail(userEmail) {
		return apperror.
			New(apperror.InvalidEmail).
			Describe("Invalid email")
	}

	var countEmail int64
	findEmailErr := s.repo.CountEmail(&countEmail, userEmail)
	if findEmailErr != nil {
		s.logger.Error("Could not count email", zap.Error(findEmailErr))
		return apperror.
			New(apperror.InternalServerError).
			Describe("Could not send email. Please try again later")
	} else if countEmail > 0 {
		return apperror.
			New(apperror.EmailAlreadyExists).
			Describe("Email already exists")
	}

	code := "SCK-" + utils.RandomString(16)

	expiredAt := time.Now().Add(5 * time.Minute)

	verificationData := models.EmailVerificationData{
		Email:     userEmail,
		Code:      code,
		ExpiredAt: &expiredAt,
	}

	if s.repo.CreateEmailVerificationData(&verificationData) != nil {
		s.logger.Error("Could not create email verification data", zap.Error(findEmailErr))
		return apperror.
			New(apperror.InternalServerError).
			Describe("Could not send email. Please try again later")
	}

	to := []string{userEmail}
	subject := "Email Verification from suechaokhai.com"
	emailStructure := models.VerificationEmail{
		// VerificationLink: "https://www.youtube.com/@oreo10baht",
		VerificationLink: "http://localhost:8000/email/verify?email=" + userEmail + "&code=" + code,
	}

	return s.sendEmail(to, subject, emailStructure)
}

func (s *serviceImpl) sendEmail(to []string, subject string, emailStructure models.EmailType) *apperror.AppError {
	smtpHost := s.cfg.SmtpHost
	smtpPort := s.cfg.SmtpPort
	smtpAddr := smtpHost + ":" + smtpPort

	from := s.cfg.Email
	password := s.cfg.EmailPassword

	auth := smtp.PlainAuth("", from, password, smtpHost)

	path := emailStructure.Path()
	t, templateErr := template.ParseFiles(path)
	if templateErr != nil {
		s.logger.Error("Could not parse email template", zap.Error(templateErr))
		return apperror.
			New(apperror.InternalServerError).
			Describe("Could not send email. Please try again later")
	}

	var body bytes.Buffer
	t.Execute(&body, emailStructure)

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	message := []byte(fmt.Sprintf("Subject: %s \n%s\n\n%s", subject, mimeHeaders, body.String()))

	err := smtp.SendMail(smtpAddr, auth, from, to, message)
	if err != nil {
		s.logger.Error("Could not send email", zap.Error(err))
		return apperror.
			New(apperror.InternalServerError).
			Describe("Could not send email. Please try again later")
	}

	return nil
}

func (s *serviceImpl) VerifyEmail(userEmail string, userCode string) *apperror.AppError {
	if !utils.IsValidEmail(userEmail) {
		return apperror.
			New(apperror.InvalidEmail).
			Describe("Invalid email")
	}

	if !utils.IsValidEmailVerificationCode(userCode) {
		return apperror.
			New(apperror.InvalidEmailVerificationCode).
			Describe("Invalid verification code")
	}

	verificationData := models.EmailVerificationData{}

	getDataErr := s.repo.GetEmailVerificationDataByEmail(&verificationData, userEmail)
	if getDataErr != nil {
		s.logger.Error("Could not get email verification data", zap.Error(getDataErr))
		return apperror.
			New(apperror.InternalServerError).
			Describe("Could not verify email. Please try again later")
	}

	if verificationData.ExpiredAt.Before(time.Now()) {
		err := s.repo.DeleteEmailVerificationData(userEmail)
		if err != nil {
			s.logger.Error("Could not delete email verification data", zap.Error(err))
			return apperror.
				New(apperror.InternalServerError).
				Describe("Verification code expired")
		}
		return apperror.
			New(apperror.EmailVerificationCodeExpired).
			Describe("Verification code expired")
	}

	if userCode != verificationData.Code {
		return apperror.
			New(apperror.InvalidEmailVerificationCode).
			Describe("Invalid verification code")
	}

	err := s.repo.DeleteEmailVerificationData(userEmail)
	if err != nil {
		s.logger.Error("Could not delete email verification data", zap.Error(err))
		return apperror.
			New(apperror.InternalServerError).
			Describe("Server Error. Please try again later")
	}

	return nil
}

func (s *serviceImpl) DeleteEmailVerificationData(userEmail string) *apperror.AppError {
	if !utils.IsValidEmail(userEmail) {
		return apperror.
			New(apperror.InvalidEmail).
			Describe("Invalid email")
	}

	err := s.repo.DeleteEmailVerificationData(userEmail)
	if err != nil {
		s.logger.Error("Could not delete email verification data", zap.Error(err))
		return apperror.
			New(apperror.InternalServerError).
			Describe("Could not delete email verification data. Please try again later")
	}

	return nil
}
