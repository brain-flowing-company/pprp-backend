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

// @router      /api/v1/register [post]
// @summary     Register user
// @description Register user
// @tags        register
// @produce     json
// @param       user body User true "User"
// @success     200	{object} User
// @failure     400 {object} AppError
// @failure     404 {object} AppError
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
