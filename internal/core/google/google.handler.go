package google

import (
	"net/http"

	"github.com/brain-flowing-company/pprp-backend/config"
	"github.com/brain-flowing-company/pprp-backend/internal/utils"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type Handler interface {
	GoogleLogin(*fiber.Ctx) error
}

type handlerImpl struct {
	service Service
	logger  *zap.Logger
	cfg     *config.Config
}

func NewHandler(logger *zap.Logger, cfg *config.Config, service Service) Handler {
	return &handlerImpl{
		service,
		logger,
		cfg,
	}
}

// @router      /api/v1/oauth/google [get]
// @summary     Login with Google
// @description Redirect to this endpoint to login with Google OAuth2. When logging in is completed, the redirection to /register in client will occur.
// @tags        auth
// @success     307
// @failure     500 {model} models.ErrorResponses
func (h *handlerImpl) GoogleLogin(c *fiber.Ctx) error {
	url, err := h.service.GoogleLogin()
	if err != nil {
		return utils.ResponseError(c, err)
	}

	return c.Redirect(url, http.StatusTemporaryRedirect)
}
