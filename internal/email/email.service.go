package email

import (
	"net/smtp"

	"github.com/brain-flowing-company/pprp-backend/apperror"
	"github.com/brain-flowing-company/pprp-backend/config"
	"github.com/brain-flowing-company/pprp-backend/utils"
	"go.uber.org/zap"
)

type Service interface {
	SendEmail(string) *apperror.AppError
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

func (s *serviceImpl) SendEmail(email string) *apperror.AppError {

	if !utils.IsValidEmail(email) {
		return apperror.
			New(apperror.InvalidEmail).
			Describe("Invalid email")
	}

	var countEmail int64
	findEmailErr := s.repo.CountEmail(&countEmail, email)
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

	smtpHost := s.cfg.SmtpHost
	smtpPort := s.cfg.SmtpPort
	from := s.cfg.Email
	password := s.cfg.EmailPassword
	to := []string{email}

	message := []byte("To: " + email + "\r\n" +
		"Subject: Email Verification from suechaokhai.com\r\n" +
		"\r\n" +
		"Welcome to Sue Chao Khai by Brain-Flowing Company :)")

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		s.logger.Error("Could not send email", zap.Error(err))
		return apperror.
			New(apperror.InternalServerError).
			Describe("Could not send email. Please try again later")
	}

	return nil
}
