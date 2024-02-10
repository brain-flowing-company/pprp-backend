package users

import (
	"fmt"

	"github.com/brain-flowing-company/pprp-backend/apperror"
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/brain-flowing-company/pprp-backend/utils"
	"github.com/gofiber/fiber/v2"
)

type Handler interface {
	GetAllUsers(c *fiber.Ctx) error
	GetUserById(c *fiber.Ctx) error
	GetCurrentUserFromLocalStorage(c *fiber.Ctx) error
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
// @failure     500 {object} apperror.AppError
func (h *handlerImpl) GetAllUsers(c *fiber.Ctx) error {
	users := []models.Users{}

	err := h.service.GetAllUsers(&users)
	if err != nil {
		return utils.ResponseError(c, err)
	}

	return c.JSON(users)
}

// @router      /api/v1/user/:userId [get]
// @summary     Get user by id
// @description Get a user by its id
// @tags        users
// @produce     json
// @success     200	{object} models.Users
// @failure     400 {object} apperror.AppError "Invalid user id"
// @failure     404 {object} apperror.AppError "User not found"
// @failure     500 {object} apperror.AppError
func (h *handlerImpl) GetUserById(c *fiber.Ctx) error {
	userId := c.Params("userId")
	user := models.Users{}

	err := h.service.GetUserById(&user, userId)
	if err != nil {
		return utils.ResponseError(c, err)
	}

	return c.JSON(user)
}

// @router      /api/v1/user/register [post]
// @summary     Register
// @description Create a user by prasing the body
// @tags        users
// @produce     json
// @success     200	{object} models.Users
// @failure     400 {object} apperror.AppError "Invalid user info"
// @failure     401 {object} apperror.AppError "Invalid password"
// @failure     500 {object} apperror.AppError
func (h *handlerImpl) Register(c *fiber.Ctx) error {
	user := models.Users{}

	bodyErr := c.BodyParser(&user)
	if bodyErr != nil {
		return utils.ResponseError(c, apperror.InvalidBody)
	}

	err := h.service.Register(&user)
	if err != nil {
		return utils.ResponseError(c, err)
	}

	return c.JSON(user) // TODO: don't return user
}

// @router      /api/v1/user/:userId [put]
// @summary     Update user by id
// @description Update a user with the given id by parsing the body
// @tags        users
// @produce     json
// @success     200	{object} models.Users
// @failure     400 {object} apperror.AppError "Invalid user info"
// @failure     404 {object} apperror.AppError "User not found"
// @failure     500 {object} apperror.AppError
func (h *handlerImpl) UpdateUser(c *fiber.Ctx) error {
	userId := c.Params("userId")
	user := models.Users{}

	bodyErr := c.BodyParser(&user)
	if bodyErr != nil {
		return utils.ResponseError(c, apperror.InvalidBody)
	}

	err := h.service.UpdateUser(&user, userId)
	if err != nil {
		return utils.ResponseError(c, err)
	}

	return c.JSON(user)
}

// @router      /api/v1/user/:userId [delete]
// @summary     Delete user by id
// @description Delete a user by its id
// @tags        users
// @produce     json
// @success     200
// @failure     400 {object} apperror.AppError "Invalid user id"
// @failure     404 {object} apperror.AppError "User not found"
// @failure     500 {object} apperror.AppError
func (h *handlerImpl) DeleteUser(c *fiber.Ctx) error {
	userId := c.Params("userId")

	err := h.service.DeleteUser(userId)
	if err != nil {
		return utils.ResponseError(c, err)
	}

	return c.JSON(fiber.Map{
		"message": "User deleted successfully",
	})
}

// @router      /api/v1/user/me [get]
// @summary     Get current user info
// @description Get current user info
// @tags        users
// @produce     json
// @success     200 {object} models.Users
// @failure     401 {object} apperror.AppError
// @failure     500 {object} apperror.AppError
func (h *handlerImpl) GetCurrentUserFromLocalStorage(c *fiber.Ctx) error {
	fmt.Println("HERE")
	email := c.Locals("email").(string)
	user := models.Users{}
	err := h.service.GetUserByEmail(&user, email)
	if err != nil {
		return c.Status(err.Code).JSON(fiber.Map{
			"message": err.Name,
		})
	}
	return c.JSON(user)
}
