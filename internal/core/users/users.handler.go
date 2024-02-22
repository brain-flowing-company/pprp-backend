package users

import (
	"fmt"
	"net/http"

	"github.com/brain-flowing-company/pprp-backend/apperror"
	"github.com/brain-flowing-company/pprp-backend/internal/enums"
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/brain-flowing-company/pprp-backend/internal/utils"
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
// @failure     500 {object} models.ErrorResponses
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
// @failure     400 {object} models.ErrorResponses "Invalid user id"
// @failure     404 {object} models.ErrorResponses "User not found"
// @failure     500 {object} models.ErrorResponses
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
// @description Create user with formData **\***upload profile image in formData with field `profile_image`. Available formats are .png / .jpg / .jpeg
// @tags        users
// @produce     json
// @param       formData formData models.RegisteringUsers true "User information"
// @success     200	{object} models.Users
// @failure     400 {object} models.ErrorResponses "Invalid user info"
// @failure     500 {object} models.ErrorResponses
func (h *handlerImpl) Register(c *fiber.Ctx) error {
	user := &models.RegisteringUsers{
		UserId: uuid.New(),
	}

	err := c.BodyParser(user)
	if err != nil {
		return utils.ResponseError(c, apperror.
			New(apperror.BadRequest).
			Describe(fmt.Sprintf("Could not parse form data: %v", err.Error())))
	}

	if session, ok := c.Locals("session").(models.Sessions); !ok {
		user.RegisteredType = enums.EMAIL
	} else {
		user.RegisteredType = session.RegisteredType
	}

	fmt.Println(user)

	profileImage, _ := c.FormFile("profile_image")
	apperr := h.service.Register(user, profileImage)
	if apperr != nil {
		return utils.ResponseError(c, apperr)
	}

	return utils.ResponseMessage(c, http.StatusCreated, "User created")
}

// @router      /api/v1/user/ [put]
// @summary     Update current user personal information *use cookies*
// @description Update specifying userId with formData **\***upload profile image in formData with field `profile_image`. Available formats are .png / .jpg / .jpeg
// @tags        users
// @produce     json
// @param       formData formData models.UpdatingUserPersonalInfo true "User information"
// @success     200	{object} models.Users
// @failure     400 {object} models.ErrorResponses "Invalid user info"
// @failure     404 {object} models.ErrorResponses "User not found"
// @failure     500 {object} models.ErrorResponses
func (h *handlerImpl) UpdateUser(c *fiber.Ctx) error {
	session, ok := c.Locals("session").(models.Sessions)
	if !ok {
		session = models.Sessions{}
	}

	user := models.UpdatingUserPersonalInfo{UserId: session.UserId}
	bodyErr := c.BodyParser(&user)
	if bodyErr != nil {
		return utils.ResponseError(c, apperror.InvalidBody)
	}

	profileImage, _ := c.FormFile("profile_image")

	apperr := h.service.UpdateUser(&user, profileImage)
	if apperr != nil {
		return utils.ResponseError(c, apperr)
	}

	return utils.ResponseMessage(c, http.StatusOK, "User updated")
}

// @router      /api/v1/user/:userId [delete]
// @summary     Delete user by id  *use cookies*
// @description Delete a user by its id
// @tags        users
// @produce     json
// @success     200
// @failure     400 {object} models.ErrorResponses "Invalid user id"
// @failure     404 {object} models.ErrorResponses "User not found"
// @failure     500 {object} models.ErrorResponses
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
// @failure     500 {object} models.ErrorResponses
func (h *handlerImpl) GetCurrentUser(c *fiber.Ctx) error {
	session := c.Locals("session").(models.Sessions)
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
// @success     200 {object} models.Sessions
func (h *handlerImpl) GetRegisteredType(c *fiber.Ctx) error {
	session, ok := c.Locals("session").(models.Sessions)
	if !ok {
		session = models.Sessions{}
	}

	return c.JSON(session)
}
