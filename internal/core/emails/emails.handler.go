package emails

import (
	"net/http"

	"github.com/brain-flowing-company/pprp-backend/apperror"
	"github.com/brain-flowing-company/pprp-backend/config"
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/brain-flowing-company/pprp-backend/internal/utils"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type Handler interface {
	SendVerificationEmail(c *fiber.Ctx) error
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

// @router      /api/v1/email [post]
// @summary     Send verification email
// @description Send a verification email to the user
// @tags        emails
// @produce     json
// @param       body body models.SendingEmailRequests true "User email address"
// @success     200 {object} models.MessageResponses
// @failure     500 {object} models.ErrorResponses
func (h *handlerImpl) SendVerificationEmail(c *fiber.Ctx) error {
	body := models.SendingEmailRequests{}

	bodyErr := c.BodyParser(&body)
	if bodyErr != nil {
		h.logger.Error("Could not parse body", zap.Error(bodyErr))
		return utils.ResponseError(c, apperror.
			New(apperror.InvalidBody).
			Describe("Invalid request body"))
	}

	appErr := h.service.SendVerificationEmail(body.Emails)
	if appErr != nil {
		return utils.ResponseError(c, appErr)
	}

	return utils.ResponseMessage(c, http.StatusOK, "Email sent successfully")
}
