package users

import (
	"fmt"
	"net/http"
	"time"

	"github.com/brain-flowing-company/pprp-backend/apperror"
	"github.com/brain-flowing-company/pprp-backend/config"
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/brain-flowing-company/pprp-backend/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Handler interface {
	GetAllUsers(c *fiber.Ctx) error
	GetUserById(c *fiber.Ctx) error
	GetCurrentUser(c *fiber.Ctx) error
	GetUserFinancialInformation(c *fiber.Ctx) error
	Register(c *fiber.Ctx) error
	UpdateUser(c *fiber.Ctx) error
	UpdateUserFinancialInformation(c *fiber.Ctx) error
	DeleteUser(c *fiber.Ctx) error
	GetRegisteredType(c *fiber.Ctx) error
	VerifyCitizenId(c *fiber.Ctx) error
}

type handlerImpl struct {
	logger  *zap.Logger
	cfg     *config.Config
	service Service
}

func NewHandler(logger *zap.Logger, cfg *config.Config, service Service) Handler {
	return &handlerImpl{
		logger,
		cfg,
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

// @router      /api/v1/user/me/financial-information [get]
// @summary     Get current user financial information *use cookies*
// @description Get current user financial information
// @tags        users
// @produce     json
// @success     200 {object} models.UserFinancialInformations
// @failure     400 {object} models.ErrorResponses "Invalid user id"
// @failure     403 {object} models.ErrorResponses "Unauthorized"
// @failure     500 {object} models.ErrorResponses
func (h *handlerImpl) GetUserFinancialInformation(c *fiber.Ctx) error {
	session, ok := c.Locals("session").(models.Sessions)
	if !ok {
		session = models.Sessions{}
	}
	userFinancialInformation := models.UserFinancialInformations{}

	err := h.service.GetUserFinancialInforamtionById(&userFinancialInformation, session.UserId.String())
	if err != nil {
		return utils.ResponseError(c, err)
	}

	return c.JSON(userFinancialInformation)
}

// @router      /api/v1/register [post]
// @summary     Register *use cookies*
// @description Create user with formData **\***upload profile image in formData with field `profile_image`. Available formats are .png / .jpg / .jpeg
// @tags        users
// @produce     json
// @param       formData formData models.RegisteringUsers true "User information"
// @success     200	{object} models.MessageResponses "User created"
// @failure     400 {object} models.ErrorResponses "Invalid user info"
// @failure     500 {object} models.ErrorResponses "Could not create user"
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

	profileImage, _ := c.FormFile("profile_image")
	apperr := h.service.Register(user, profileImage)
	if apperr != nil {
		return utils.ResponseError(c, apperr)
	}

	return utils.ResponseMessage(c, http.StatusCreated, "User created")
}

// @router      /api/v1/user/me/personal-information [put]
// @summary     Update current user personal information *use cookies*
// @description Update specifying userId with formData **\***upload profile image in formData with field `profile_image`. Available formats are .png / .jpg / .jpeg
// @tags        users
// @produce     json
// @param       formData formData models.UpdatingUserPersonalInfos true "User personal information"
// @success     200	{object} models.MessageResponses "User personal information updated"
// @failure     400 {object} models.ErrorResponses "Invalid user info"
// @failure     404 {object} models.ErrorResponses "User not found"
// @failure     500 {object} models.ErrorResponses "Could not update user"
func (h *handlerImpl) UpdateUser(c *fiber.Ctx) error {
	session, ok := c.Locals("session").(models.Sessions)
	if !ok {
		session = models.Sessions{}
	}

	user := models.UpdatingUserPersonalInfos{UserId: session.UserId}
	bodyErr := c.BodyParser(&user)
	if bodyErr != nil {
		return utils.ResponseError(c, apperror.InvalidBody)
	}

	profileImage, _ := c.FormFile("profile_image")

	apperr := h.service.UpdateUser(&user, profileImage)
	if apperr != nil {
		return utils.ResponseError(c, apperr)
	}

	return utils.ResponseMessage(c, http.StatusOK, "User personal information updated")
}

// @router      /api/v1/user/me/financial-information [put]
// @summary     Update the current user financial information *use cookies*
// @description Update the current user financial information with data from the body
// @tags        users
// @produce     json
// @param       body body models.UserFinancialInformations true "User financial information"
// @success     200 {object} models.MessageResponses "User financial information updated"
// @failure     400 {object} models.ErrorResponses "Invalid user financial information"
// @failure     403 {object} models.ErrorResponses "Unauthorized"
// @failure     500 {object} models.ErrorResponses "Could not update user financial information"
func (h *handlerImpl) UpdateUserFinancialInformation(c *fiber.Ctx) error {
	session, ok := c.Locals("session").(models.Sessions)
	if !ok {
		session = models.Sessions{}
	}

	userFinancialInformation := models.UserFinancialInformations{}
	bodyErr := c.BodyParser(&userFinancialInformation)
	if bodyErr != nil {
		return utils.ResponseError(c, apperror.InvalidBody)
	}

	apperr := h.service.UpdateUserFinancialInformationById(&userFinancialInformation, session.UserId.String())
	if apperr != nil {
		return utils.ResponseError(c, apperr)
	}

	return utils.ResponseMessage(c, http.StatusOK, "User financial information updated")
}

// @router      /api/v1/user/:userId [delete]
// @summary     Delete user by id  *use cookies*
// @description Delete a user by its id
// @tags        users
// @produce     json
// @success     200 {object} models.MessageResponses "User deleted"
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

// @router      /api/v1/user/me/verify [post]
// @summary     Verify user *use cookies*
// @description Verify user by citizen id and citizen id image
// @tags        users
// @produce     json
// @param       formData formData models.UserVerifications true "Verification information"
// @success     200 {object} models.MessageResponses "Verified"
// @success     500 {object} models.ErrorResponses
func (h *handlerImpl) VerifyCitizenId(c *fiber.Ctx) error {
	session, ok := c.Locals("session").(models.Sessions)
	if !ok {
		session = models.Sessions{}
	}

	user := models.UserVerifications{UserId: session.UserId}
	bodyErr := c.BodyParser(&user)
	if bodyErr != nil {
		return utils.ResponseError(c, apperror.InvalidBody)
	}

	profileImage, _ := c.FormFile("citizen_card_image")
	err := h.service.VerifyCitizenId(&user, profileImage)
	if err != nil {
		return utils.ResponseError(c, err)
	}

	session.IsOwner = true

	token, apperr := utils.CreateJwtToken(session, time.Duration(h.cfg.SessionExpire*int(time.Second)), h.cfg.JWTSecret)
	if apperr != nil {
		h.logger.Error("Could not create JWT token", zap.Error(apperr))
		return utils.ResponseError(c, apperror.
			New(apperror.InternalServerError).
			Describe("Could not login. Please try again later"))
	}

	c.Cookie(utils.CreateSessionCookie(token, h.cfg.SessionExpire))

	return utils.ResponseMessage(c, http.StatusOK, "Verified")
}
