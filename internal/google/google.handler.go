package google

import (
	"net/http"

	"github.com/brain-flowing-company/pprp-backend/apperror"
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
}

func NewHandler(service Service, logger *zap.Logger) Handler {
	return &handlerImpl{
		service,
		logger,
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

	return c.SendString(token)
}
