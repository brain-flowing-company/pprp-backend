package google

import (
	"context"
	"encoding/json"
	"time"

	"github.com/brain-flowing-company/pprp-backend/apperror"
	"github.com/brain-flowing-company/pprp-backend/config"
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/brain-flowing-company/pprp-backend/utils"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Service interface {
	GoogleLogin() string
	ExchangeToken(context.Context, *models.GoogleExchangeToken) (string, *apperror.AppError)
}

type serviceImpl struct {
	authCfg *oauth2.Config
	logger  *zap.Logger
	cfg     *config.Config
}

func NewService(cfg *config.Config, logger *zap.Logger) Service {
	return &serviceImpl{
		&oauth2.Config{
			ClientID:     cfg.GoogleClientId,
			ClientSecret: cfg.GoogleSecret,
			RedirectURL:  cfg.GoogleRedirect,
			Endpoint:     google.Endpoint,
			Scopes:       cfg.GoogleScopes,
		},
		logger,
		cfg,
	}
}

func (s *serviceImpl) GoogleLogin() string {
	return s.authCfg.AuthCodeURL("state")
}

func (s *serviceImpl) ExchangeToken(c context.Context, excToken *models.GoogleExchangeToken) (string, *apperror.AppError) {
	oauthToken, err := s.authCfg.Exchange(c, excToken.Code)
	if err != nil {
		s.logger.Error("Could not exchange token from google", zap.Error(err))
		return "", apperror.
			New(apperror.ServiceUnavailable).
			Describe("Google OAuth failed")
	}

	client := s.authCfg.Client(c, oauthToken)

	res, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		s.logger.Error("Could not get userinfo", zap.Error(err))
		return "", apperror.
			New(apperror.ServiceUnavailable).
			Describe("Google OAuth failed")
	}

	googleInfo := models.GoogleUserInfo{}

	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(&googleInfo)
	if err != nil {
		s.logger.Error("Could not decode json body", zap.Error(err))
		return "", apperror.
			New(apperror.InternalServerError).
			Describe("Google OAuth failed")
	}

	session := models.Session{
		Email:          googleInfo.Email,
		RegisteredType: models.GOOGLE,
	}

	token, err := utils.CreateJwtToken(session, time.Duration(s.cfg.SessionExpire*int(time.Second)), s.cfg.JWTSecret)
	if err != nil {
		s.logger.Error("Could not create JWT token", zap.Error(err))
		return "", apperror.
			New(apperror.InternalServerError).
			Describe("Google OAuth failed")
	}

	return token, nil
}
