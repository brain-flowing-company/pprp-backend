package greeting

import (
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/gofiber/fiber/v2"
)

type Handler interface {
	Greeting(c *fiber.Ctx) error
	UserGreeting(c *fiber.Ctx) error
}

type handlerImpl struct {
	service Service
}

func NewHandler(service Service) Handler {
	return &handlerImpl{
		service,
	}
}

// @router      /api/v1/greeting [get]
// @summary     Greeting
// @description says hello, world
// @tags        greeting
// @produce     json
// @success     200	{object}	models.Greeting
func (h *handlerImpl) Greeting(c *fiber.Ctx) error {
	msg := models.Greeting{}
	h.service.Greeting(&msg)

	return c.JSON(msg)
}

// @router      /api/v1/user/greeting [get]
// @summary     Greeting with auth required
// @description says hello to current user
// @tags        greeting
// @produce     json
// @success     200	{object}	models.Greeting
// @failure     401 {object}	apperror.AppError
func (h *handlerImpl) UserGreeting(c *fiber.Ctx) error {
	email := (c.Locals("email").(string))

	msg := models.Greeting{}
	h.service.UserGreeting(&msg, email)

	return c.JSON(msg)
}
