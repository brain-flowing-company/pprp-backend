// internal/login/handler.go
package auth

import (
	"time"

	"github.com/brain-flowing-company/pprp-backend/apperror"
	"github.com/brain-flowing-company/pprp-backend/config"
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/brain-flowing-company/pprp-backend/utils"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type Handler interface {
	Login(*fiber.Ctx) error
	Logout(*fiber.Ctx) error
}

type handlerImpl struct {
	service Service
	cfg     *config.Config
	logger  *zap.Logger
}

func NewHandler(service Service, cfg *config.Config, logger *zap.Logger) Handler {
	return &handlerImpl{
		service,
		cfg,
		logger,
	}
}

// @router      /api/v1/login [post]
// @summary     Login with email
// @description Login with email and password
// @tags        auth
// @produce     json
// @success     200	{object} models.Property
// @failure     400 {object} models.ErrorResponse "Empty or invalid credentials"
// @failure     401 {object} models.ErrorResponse "Password mismatch"
// @failure     404 {object} models.ErrorResponse "User not found"
// @failure     500 {object} models.ErrorResponse
func (h *handlerImpl) Login(c *fiber.Ctx) error {
	// Parse login request from the request body
	var loginRequest models.LoginRequest
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
	c.Cookie(&fiber.Cookie{
		Name:     "session",
		Value:    token,
		Expires:  time.Now().Add(time.Duration(h.cfg.SessionExpire) * time.Second),
		HTTPOnly: true,
	})

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
	c.Cookie(&fiber.Cookie{
		Name:     "session",
		Expires:  time.Now(),
		HTTPOnly: true,
	})

	return c.JSON(fiber.Map{
		"success": true,
	})
}
