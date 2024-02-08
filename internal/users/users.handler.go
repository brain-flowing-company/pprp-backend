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

// @router      /api/v1/users [get]
// @summary     Get all users
// @description Get all users
// @tags        users
// @produce     json
// @success     200	{object} []models.Users
// @failure     400 {object} apperror.AppError
// @failure     404 {object} apperror.AppError
func (h *handlerImpl) GetAllUsers(c *fiber.Ctx) error {
	users := []models.Users{}

	err := h.service.GetAllUsers(&users)
	if err != nil {
		return c.Status(err.Code).JSON(err)
	}

	return c.JSON(users)
}

// @router      /api/v1/users/:userId [get]
// @summary     Get user by id
// @description Get a user by its id
// @tags        users
// @produce     json
// @success     200	{object} models.Users
// @failure     400 {object} apperror.AppError
// @failure     404 {object} apperror.AppError
func (h *handlerImpl) GetUserById(c *fiber.Ctx) error {
	userId := c.Params("userId")
	user := models.Users{}

	err := h.service.GetUserById(&user, userId)
	if err != nil {
		return c.Status(err.Code).JSON(err)
	}

	return c.JSON(user)
}

// @router      /api/v1/users [post]
// @summary     Register
// @description Create a user by prasing the body
// @tags        users
// @produce     json
// @success     200	{object} models.Users
// @failure     400 {object} apperror.AppError
// @failure     404 {object} apperror.AppError
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

// @router      /api/v1/users/:userId [put]
// @summary     Update user by id
// @description Update a user with the given id by parsing the body
// @tags        users
// @produce     json
// @success     200	{object} models.Users
// @failure     400 {object} apperror.AppError
// @failure     404 {object} apperror.AppError
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

// @router      /api/v1/users/:userId [delete]
// @summary     Delete user by id
// @description Delete a user by its id
// @tags        users
// @produce     json
// @success     200
// @failure     400 {object} apperror.AppError
// @failure     404 {object} apperror.AppError
func (h *handlerImpl) DeleteUser(c *fiber.Ctx) error {
	userId := c.Params("userId")

	err := h.service.DeleteUser(userId)
	if err != nil {
		return c.Status(err.Code).JSON(fiber.Map{
			"message": err.Name,
		})
	}

	return c.JSON(fiber.Map{
		"message": "User deleted successfully",
	})
}
