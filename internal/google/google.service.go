package google

import (
	"context"
	"encoding/json"

	"github.com/brain-flowing-company/pprp-backend/apperror"
	"github.com/brain-flowing-company/pprp-backend/config"
	"github.com/brain-flowing-company/pprp-backend/internal/models"
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
	}
}

func (s *serviceImpl) GoogleLogin() string {
	return s.authCfg.AuthCodeURL("state")
}

func (s *serviceImpl) ExchangeToken(c context.Context, excToken *models.GoogleExchangeToken) (string, *apperror.AppError) {
	token, err := s.authCfg.Exchange(c, excToken.Code)
	if err != nil {
		s.logger.Error("Could not exchange token from google", zap.Error(err))
		return "", apperror.ServiceUnavailable
	}

	client := s.authCfg.Client(c, token)

	res, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		s.logger.Error("Could not get userinfo", zap.Error(err))
		return "", apperror.ServiceUnavailable
	}

	userInfo := models.GoogleUserInfo{}

	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(&userInfo)
	if err != nil {
		s.logger.Error("Could not decode json body", zap.Error(err))
		return "", apperror.InternalServerError
	}

	return userInfo.Email, nil
}
