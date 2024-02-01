package users

import (
	"fmt"

	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/gofiber/fiber/v2"
)

type Handler interface {
	CreateUser(c *fiber.Ctx) error
	GetAllUsers(c *fiber.Ctx) error
}

type handlerImpl struct {
	service Service
}

func NewHandler(service Service) Handler {
	return &handlerImpl{
		service,
	}
}

func (h *handlerImpl) GetAllUsers(c *fiber.Ctx) error {
	users := models.Users{}
	err := h.service.GetAllUsers(&users)
	if err != nil {
		return c.Status(err.Code).JSON(err)
	}

	return c.JSON(users)
}

func (h *handlerImpl) CreateUser(c *fiber.Ctx) error {

	user := models.Users{}

	bodyErr := c.BodyParser(&user)

	if bodyErr != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"message": "Invalid body",
		})
	}

	fmt.Println(user)

	err := h.service.CreateUser(&user)
	if err != nil {
		return c.Status(err.Code).JSON(fiber.Map{
			"message": err.Name,
		})
	}

	return c.JSON(user) // TODO: don't return user
}
