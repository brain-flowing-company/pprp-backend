package users

import (
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/gofiber/fiber/v2"
)

type Handler interface {
	GetAllUsers(c *fiber.Ctx) error
	GetUserById(c *fiber.Ctx) error
	Register(c *fiber.Ctx) error
	UpdateUser(c *fiber.Ctx) error
	DeleteUser(c *fiber.Ctx) error
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
	users := []models.Users{}

	err := h.service.GetAllUsers(&users)
	if err != nil {
		return c.Status(err.Code).JSON(err)
	}

	return c.JSON(users)
}

func (h *handlerImpl) GetUserById(c *fiber.Ctx) error {
	userId := c.Params("userId")
	user := models.Users{}

	err := h.service.GetUserById(&user, userId)
	if err != nil {
		return c.Status(err.Code).JSON(err)
	}

	return c.JSON(user)
}

func (h *handlerImpl) Register(c *fiber.Ctx) error {
	user := models.Users{}

	bodyErr := c.BodyParser(&user)
	if bodyErr != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"message": "Invalid body",
		})
	}

	err := h.service.Register(&user)
	if err != nil {
		return c.Status(err.Code).JSON(fiber.Map{
			"message": err.Name,
		})
	}

	return c.JSON(user) // TODO: don't return user
}

func (h *handlerImpl) UpdateUser(c *fiber.Ctx) error {
	userId := c.Params("userId")
	user := models.Users{}

	bodyErr := c.BodyParser(&user)
	if bodyErr != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"message": "Invalid body",
		})
	}

	err := h.service.UpdateUser(&user, userId)
	if err != nil {
		return c.Status(err.Code).JSON(fiber.Map{
			"message": err.Name,
		})
	}

	return c.JSON(user)
}

func (h *handlerImpl) DeleteUser(c *fiber.Ctx) error {
	userId := c.Params("userId")

	err := h.service.DeleteUser(userId)
	if err != nil {
		return c.Status(err.Code).JSON(fiber.Map{
			"message": err.Name,
		})
	}

	return nil
}
