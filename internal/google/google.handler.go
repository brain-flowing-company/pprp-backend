package google

import (
	"net/http"

	"github.com/brain-flowing-company/pprp-backend/apperror"
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/brain-flowing-company/pprp-backend/utils"
	"github.com/gofiber/fiber/v2"
)

type Handler interface {
	GoogleLogin(*fiber.Ctx) error
	ExchangeToken(c *fiber.Ctx) error
}

type handlerImpl struct {
	service Service
}

func NewHandler(service Service) Handler {
	return &handlerImpl{
		service,
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
		return utils.ResponseError(c, apperror.InternalServerError)
	}

	token, apperr := h.service.ExchangeToken(c.Context(), &excToken)
	if apperr != nil {
		return utils.ResponseError(c, apperr)
	}

	return c.SendString(token)
}
