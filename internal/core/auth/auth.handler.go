// internal/login/handler.go
package auth

import (
	"github.com/brain-flowing-company/pprp-backend/apperror"
	"github.com/brain-flowing-company/pprp-backend/config"
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/brain-flowing-company/pprp-backend/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type Handler interface {
	Login(*fiber.Ctx) error
	Logout(*fiber.Ctx) error
	Callback(c *fiber.Ctx) error
}

type handlerImpl struct {
	service Service
	cfg     *config.Config
}

func NewHandler(cfg *config.Config, service Service) Handler {
	return &handlerImpl{
		service,
		cfg,
	}
}

// @router      /api/v1/login [post]
// @summary     Login with email
// @description Login with email and password
// @tags        auth
// @produce     json
// @success     200	{object} models.Properties
// @failure     400 {object} models.ErrorResponses "Empty or invalid credentials"
// @failure     401 {object} models.ErrorResponses "Password mismatch"
// @failure     404 {object} models.ErrorResponses "User not found"
// @failure     500 {object} models.ErrorResponses
func (h *handlerImpl) Login(c *fiber.Ctx) error {
	// Parse login request from the request body
	var loginRequest models.LoginRequests
	if err := c.BodyParser(&loginRequest); err != nil {
		return utils.ResponseError(c, apperror.
			New(apperror.BadRequest).
			Describe("Empty credential"))
	}

	// Authenticate user
	token, err := h.service.AuthenticateUser(loginRequest.Email, loginRequest.Password)
	if err != nil {
		return utils.ResponseError(c, err)
	}

	// Set JWT token as a cookie
	c.Cookie(utils.CreateSessionCookie(token, h.cfg.SessionExpire))

	// Return a success response
	return c.JSON(fiber.Map{
		"success": true,
	})
}

// @router      /api/v1/logout [post]
// @summary     Logout
// @description Logout
// @tags        auth
// @success     200
func (h *handlerImpl) Logout(c *fiber.Ctx) error {
	c.Cookie(utils.CreateSessionCookie("", 0))

	return c.JSON(fiber.Map{
		"success": true,
	})
}

func (h *handlerImpl) Callback(c *fiber.Ctx) error {
	callback := models.Callbacks{}
	err := c.QueryParser(&callback)
	if err != nil {
		return utils.ResponseError(c, apperror.
			New(apperror.InvalidCallbackRequest).
			Describe("Invalid callback request"))
	}

	var callbackResponse models.CallbackResponses
	apperr := h.service.Callback(c.Context(), &callback, &callbackResponse)
	if apperr != nil {
		return utils.ResponseError(c, apperr)
	}

	return c.JSON(callbackResponse)
}
