package google

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/brain-flowing-company/pprp-backend/apperror"
	"github.com/brain-flowing-company/pprp-backend/config"
	"github.com/brain-flowing-company/pprp-backend/internal/enums"
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/brain-flowing-company/pprp-backend/internal/utils"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gorm.io/gorm"
)

type Service interface {
	GoogleLogin() (string, *apperror.AppError)
	ExchangeToken(context.Context, *models.Callbacks, *models.CallbackResponses) *apperror.AppError
}

type serviceImpl struct {
	authCfg *oauth2.Config
	logger  *zap.Logger
	cfg     *config.Config
	repo    Repository
}

func NewService(logger *zap.Logger, cfg *config.Config, repo Repository) Service {
	return &serviceImpl{
		&oauth2.Config{
			ClientID:     cfg.GoogleClientId,
			ClientSecret: cfg.GoogleSecret,
			RedirectURL:  cfg.AuthRedirect,
			Endpoint:     google.Endpoint,
			Scopes:       cfg.GoogleScopes,
		},
		logger,
		cfg,
		repo,
	}
}

func (s *serviceImpl) GoogleLogin() (string, *apperror.AppError) {
	state := &models.GoogleOAuthStates{
		Code:      uuid.New(),
		ExpiredAt: time.Now().Add(time.Duration(s.cfg.AuthVerificationExpire)),
	}

	err := s.repo.CreateState(state)
	if err != nil {
		return "", apperror.
			New(apperror.InternalServerError).
			Describe("Could not authenticate with google. Please try again later.")
	}

	return s.authCfg.AuthCodeURL(state.Code.String()), nil
}

func (s *serviceImpl) ExchangeToken(c context.Context, callback *models.Callbacks, callbackResponse *models.CallbackResponses) *apperror.AppError {
	err := s.repo.GetState(&models.GoogleOAuthStates{}, callback.State)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return apperror.
			New(apperror.BadRequest).
			Describe("State does not match")
	} else if err != nil {
		s.logger.Error("Could not get state", zap.Error(err))
		return apperror.
			New(apperror.InternalServerError).
			Describe("Google OAuth failed")
	}

	deleteErr := s.repo.DeleteState(callback.State)
	if deleteErr != nil {
		s.logger.Error("Could not delete state", zap.Error(deleteErr))
		return apperror.
			New(apperror.InternalServerError).
			Describe("Google OAuth failed")
	}

	oauthToken, err := s.authCfg.Exchange(c, callback.Code)
	if err != nil {
		s.logger.Error("Could not exchange token from google", zap.Error(err))
		return apperror.
			New(apperror.ServiceUnavailable).
			Describe("Google OAuth failed")
	}

	client := s.authCfg.Client(c, oauthToken)

	res, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		s.logger.Error("Could not get userinfo", zap.Error(err))
		return apperror.
			New(apperror.ServiceUnavailable).
			Describe("Google OAuth failed")
	}

	googleInfo := struct{ Email string }{}
	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(&googleInfo)
	if err != nil {
		s.logger.Error("Could not decode json body", zap.Error(err))
		return apperror.
			New(apperror.InternalServerError).
			Describe("Google OAuth failed")
	}

	registered := true
	user := models.Users{}
	err = s.repo.GetUserByEmail(&user, googleInfo.Email)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		registered = false
	} else if err != nil {
		s.logger.Error("Could not get user", zap.Error(err))
		return apperror.
			New(apperror.InternalServerError).
			Describe("Google OAuth failed")
	}

	*callbackResponse = models.CallbackResponses{
		Email:          googleInfo.Email,
		RegisteredType: enums.GOOGLE,
		SessionType:    enums.SessionRegister,
	}

	if registered {
		session := models.Sessions{
			UserId:  user.UserId,
			Email:   googleInfo.Email,
			IsOwner: user.IsVerified,
		}

		token, err := utils.CreateJwtToken(session, time.Duration(s.cfg.SessionExpire*int(time.Second)), s.cfg.JWTSecret)
		if err != nil {
			s.logger.Error("Could not create JWT token", zap.Error(err))
			return apperror.
				New(apperror.InternalServerError).
				Describe("Could not login. Please try again later")
		}

		callbackResponse.Token = token
		callbackResponse.SessionType = enums.SessionLogin
	}

	return nil
}
