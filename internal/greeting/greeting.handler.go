package greeting

import (
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

// @router      /greeting [get]
// @summary     Greeting
// @description hello, world endpoint
// @produce     json
// @success     200	{object}	dto.GreetingResponse
func (h *handlerImpl) Greeting(c *fiber.Ctx) error {
	res := h.service.Greeting()

	return c.JSON(res)
}
