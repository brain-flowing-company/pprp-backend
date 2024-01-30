package greeting

import (
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/gofiber/fiber/v2"
)

type Handler interface {
	Greeting(c *fiber.Ctx) error
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
// @description hello, world endpoint
// @tags        greeting
// @produce     json
// @success     200	{object}	models.Greeting
func (h *handlerImpl) Greeting(c *fiber.Ctx) error {
	msg := models.Greeting{}
	h.service.Greeting(&msg)

	return c.JSON(msg)
}
