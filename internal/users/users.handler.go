package users

import (
	"fmt"
	"net/http"

	"github.com/brain-flowing-company/pprp-backend/apperror"
	"github.com/brain-flowing-company/pprp-backend/internal/consts"
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/brain-flowing-company/pprp-backend/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Handler interface {
	GetAllUsers(c *fiber.Ctx) error
	GetUserById(c *fiber.Ctx) error
	GetCurrentUser(c *fiber.Ctx) error
	Register(c *fiber.Ctx) error
	UpdateUser(c *fiber.Ctx) error
	DeleteUser(c *fiber.Ctx) error
	GetRegisteredType(c *fiber.Ctx) error
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
// @failure     500 {object} models.ErrorResponse
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
// @failure     400 {object} models.ErrorResponse "Invalid user id"
// @failure     404 {object} models.ErrorResponse "User not found"
// @failure     500 {object} models.ErrorResponse
func (h *handlerImpl) GetUserById(c *fiber.Ctx) error {
	userId := c.Params("userId")
	user := models.Users{}

	err := h.service.GetUserById(&user, userId)
	if err != nil {
		return utils.ResponseError(c, err)
	}

	return c.JSON(user)
}

// @router      /api/v1/register [post]
// @summary     Register *use cookies*
// @description Create a user by prasing the body
// @tags        users
// @produce     json
// @success     200	{object} models.Users
// @failure     400 {object} models.ErrorResponse "Invalid user info"
// @failure     500 {object} models.ErrorResponse
func (h *handlerImpl) Register(c *fiber.Ctx) error {
	user := &models.Users{
		UserId: uuid.New(),
	}

	err := c.BodyParser(user)
	if err != nil {
		fmt.Println(err.Error())
		return utils.ResponseError(c, apperror.
			New(apperror.BadRequest).
			Describe("Could not parse form data"))
	}

	if session, ok := c.Locals("session").(models.Session); !ok {
		user.RegisteredType = consts.EMAIL
	} else {
		user.RegisteredType = session.RegisteredType
	}

	profileImage, _ := c.FormFile("profile_image")

	apperr := h.service.Register(user, profileImage)
	if apperr != nil {
		return utils.ResponseError(c, apperr)
	}

	return utils.ResponseMessage(c, http.StatusCreated, "User created")
}

// @router      /api/v1/user/:userId [put]
// @summary     Update user by id
// @description Update a user with the given id by parsing the body
// @tags        users
// @produce     json
// @success     200	{object} models.Users
// @failure     400 {object} models.ErrorResponse "Invalid user info"
// @failure     404 {object} models.ErrorResponse "User not found"
// @failure     500 {object} models.ErrorResponse
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

	return utils.ResponseMessage(c, http.StatusOK, "User updated")
}

// @router      /api/v1/user/:userId [delete]
// @summary     Delete user by id
// @description Delete a user by its id
// @tags        users
// @produce     json
// @success     200
// @failure     400 {object} models.ErrorResponse "Invalid user id"
// @failure     404 {object} models.ErrorResponse "User not found"
// @failure     500 {object} models.ErrorResponse
func (h *handlerImpl) DeleteUser(c *fiber.Ctx) error {
	userId := c.Params("userId")

	err := h.service.DeleteUser(userId)
	if err != nil {
		return utils.ResponseError(c, err)
	}

	return utils.ResponseMessage(c, http.StatusOK, "User deleted")
}

// @router      /api/v1/user/me [get]
// @summary     Get current user info *use cookies*
// @description Get current user info
// @tags        users
// @produce     json
// @success     200 {object} models.Users
// @failure     500 {object} models.ErrorResponse
func (h *handlerImpl) GetCurrentUser(c *fiber.Ctx) error {
	session := c.Locals("session").(models.Session)
	user := models.Users{}
	err := h.service.GetUserByEmail(&user, session.Email)
	if err != nil {
		return utils.ResponseError(c, err)
	}
	return c.JSON(user)
}

// @router      /api/v1/user/me/registered [get]
// @summary     Get user registered type *use cookies*
// @description Get user registered type
// @tags        users
// @produce     json
// @success     200 {object} models.Session
func (h *handlerImpl) GetRegisteredType(c *fiber.Ctx) error {
	session, ok := c.Locals("session").(models.Session)
	if !ok {
		session = models.Session{}
	}

	return c.JSON(session)
}
