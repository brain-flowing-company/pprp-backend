package register

import (
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/gofiber/fiber/v2"
)

type Handler interface {
	CreateUser(*fiber.Ctx) error
}

type handlerImpl struct {
	service Service
}

func NewHandler(service Service) Handler {
	return &handlerImpl{
		service,
	}
}

func (h *handlerImpl) CreateUser(c *fiber.Ctx) error {
	user := models.Users{}
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	if err := h.service.CreateUser(&user); err != nil {
		return c.Status(500).JSON(err)
	}

	return c.JSON(user)
}
