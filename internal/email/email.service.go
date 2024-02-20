package email

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"

	"github.com/brain-flowing-company/pprp-backend/apperror"
	"github.com/brain-flowing-company/pprp-backend/config"
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/brain-flowing-company/pprp-backend/utils"
	"go.uber.org/zap"
)

type Service interface {
	SendVerificationEmail(string) *apperror.AppError
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

	to := []string{userEmail}
	subject := "Email Verification from suechaokhai.com"
	emailStructure := models.VerificationEmail{
		VerificationLink: "https://www.youtube.com/@oreo10baht",
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
