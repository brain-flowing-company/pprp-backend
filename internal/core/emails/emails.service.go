package emails

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"time"

	"github.com/brain-flowing-company/pprp-backend/apperror"
	"github.com/brain-flowing-company/pprp-backend/config"
	"github.com/brain-flowing-company/pprp-backend/internal/enums"
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/brain-flowing-company/pprp-backend/internal/utils"
	"go.uber.org/zap"
)

type Service interface {
	SendVerificationEmail([]string) *apperror.AppError
	VerifyEmail(*models.Callbacks, *models.CallbackResponses) *apperror.AppError
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

func (s *serviceImpl) SendVerificationEmail(emails []string) *apperror.AppError {
	userEmail := emails[0]
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

	code := s.cfg.EmailCodePrefix + utils.RandomString(16)

	emailVerificationCodeExpire := s.cfg.AuthVerificationExpire
	expiredAt := time.Now().Add(time.Duration(emailVerificationCodeExpire) * time.Second)

	verificationData := models.EmailVerificationCodes{
		Email:     userEmail,
		Code:      code,
		ExpiredAt: expiredAt,
	}

	if s.repo.CreateEmailVerificationCode(&verificationData) != nil {
		s.logger.Error("Could not create email verification data", zap.Error(findEmailErr))
		return apperror.
			New(apperror.InternalServerError).
			Describe("Could not send email. Please try again later")
	}

	link := fmt.Sprintf("%v?email=%v&code=%v", s.cfg.AuthRedirect, userEmail, code)
	subject := "Email Verification from suechaokhai.com"
	emailStructure := models.VerificationEmails{
		// VerificationLink: "https://www.youtube.com/@oreo10baht",
		// VerificationLink: "http://localhost:8000/email/verify?email=" + userEmail + "&code=" + code,
		VerificationLink: link,
	}

	return s.sendEmail(emails, subject, emailStructure)
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
	emailTempErr := t.Execute(&body, emailStructure)
	if emailTempErr != nil {
		s.logger.Error("Could not execute email template", zap.Error(emailTempErr))
		return apperror.
			New(apperror.InternalServerError).
			Describe("Could not send email. Please try again later")
	}

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

func (s *serviceImpl) VerifyEmail(verificationReq *models.Callbacks, callbackResponse *models.CallbackResponses) *apperror.AppError {
	userEmail := verificationReq.Email
	userCode := verificationReq.Code

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

	verificationData := models.EmailVerificationCodes{}

	var countData int64
	countDataErr := s.repo.CountEmailVerificationCode(&countData, userEmail)
	if countDataErr != nil {
		s.logger.Error("Could not count email verification data", zap.Error(countDataErr))
		return apperror.
			New(apperror.InternalServerError).
			Describe("Could not verify email. Please try again later")
	} else if countData == 0 {
		return apperror.
			New(apperror.EmailVerificationDataNotFound).
			Describe("Could not verify email. Please try again later")
	}

	getDataErr := s.repo.GetEmailVerificationCodeByEmail(&verificationData, userEmail)
	if getDataErr != nil {
		s.logger.Error("Could not get email verification data", zap.Error(getDataErr))
		return apperror.
			New(apperror.InternalServerError).
			Describe("Could not verify email. Please try again later")
	}

	if verificationData.ExpiredAt.Before(time.Now()) {
		err := s.repo.DeleteEmailVerificationCode(userEmail)
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

	deleteErr := s.repo.DeleteEmailVerificationCode(userEmail)
	if deleteErr != nil {
		s.logger.Error("Could not delete email verification data", zap.Error(deleteErr))
		return apperror.
			New(apperror.InternalServerError).
			Describe("Server Error. Please try again later")
	}

	*callbackResponse = models.CallbackResponses{
		Email:          userEmail,
		RegisteredType: enums.EMAIL,
		SessionType:    enums.SessionRegister,
	}

	return nil
}
