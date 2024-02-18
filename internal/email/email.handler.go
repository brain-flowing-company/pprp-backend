package email

import (
	"github.com/brain-flowing-company/pprp-backend/apperror"
	"github.com/brain-flowing-company/pprp-backend/utils"
	"github.com/gofiber/fiber/v2"
)

type Handler interface {
	SendEmail(c *fiber.Ctx) error
}

type handlerImpl struct {
	service Service
}

func NewHandler(service Service) Handler {
	return &handlerImpl{
		service,
	}
}

type SendEmailBody struct {
	Email string `json:"email"`
}

func (h *handlerImpl) SendEmail(c *fiber.Ctx) error {
	body := SendEmailBody{}

	bodyErr := c.BodyParser(&body)
	if bodyErr != nil {
		return utils.ResponseError(c, apperror.
			New(apperror.InvalidBody).
			Describe("Invalid request body"))
	}

	appErr := h.service.SendEmail(body.Email)
	if appErr != nil {
		return utils.ResponseError(c, appErr)
	}

	return utils.ResponseMessage(c, 200, "Email sent successfully")
}
