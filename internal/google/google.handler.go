package google

import (
	"net/http"
	"time"

	"github.com/brain-flowing-company/pprp-backend/apperror"
	"github.com/brain-flowing-company/pprp-backend/config"
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/brain-flowing-company/pprp-backend/utils"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type Handler interface {
	GoogleLogin(*fiber.Ctx) error
	ExchangeToken(c *fiber.Ctx) error
}

type handlerImpl struct {
	service Service
	logger  *zap.Logger
	cfg     *config.Config
}

func NewHandler(service Service, logger *zap.Logger, cfg *config.Config) Handler {
	return &handlerImpl{
		service,
		logger,
		cfg,
	}
}

func (h *handlerImpl) GoogleLogin(c *fiber.Ctx) error {
	url := h.service.GoogleLogin()

	return c.Redirect(url, http.StatusTemporaryRedirect)
}

func (h *handlerImpl) ExchangeToken(c *fiber.Ctx) error {
	excToken := models.GoogleExchangeToken{}

	err := c.QueryParser(&excToken)
	if err != nil {
		h.logger.Error("Could not parse query", zap.Error(err))
		return utils.ResponseError(c, apperror.InternalServerError)
	}

	token, apperr := h.service.ExchangeToken(c.Context(), &excToken)
	if apperr != nil {
		return utils.ResponseError(c, apperr)
	}

	c.Cookie(&fiber.Cookie{
		Name:    "session",
		Value:   token,
		Expires: time.Now().Add(time.Duration(h.cfg.SessionExpire) * time.Second),
	})

	return c.Redirect(h.cfg.LoginRedirect, http.StatusPermanentRedirect)
}
