package greetings

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
	msg := models.Greetings{}
	h.service.Greeting(&msg)

	return c.JSON(msg)
}

// @router      /api/v1/user/greeting [get]
// @summary     Greeting with auth required *use cookies*
// @description says hello to current user
// @tags        greeting
// @produce     json
// @success     200	{object}	models.Greeting
// @failure     401 {object}	models.ErrorResponse
func (h *handlerImpl) UserGreeting(c *fiber.Ctx) error {
	session := (c.Locals("session").(models.Sessions))

	msg := models.Greetings{}
	h.service.UserGreeting(&msg, session.Email)

	return c.JSON(msg)
}
