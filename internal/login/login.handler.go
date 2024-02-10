// internal/login/handler.go
package login

import (
	"net/http"
	"time"

	"github.com/brain-flowing-company/pprp-backend/apperror"
	"github.com/brain-flowing-company/pprp-backend/config"
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/brain-flowing-company/pprp-backend/utils"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type Handler interface {
	Login(c *fiber.Ctx) error
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
// @failure     400 {object} model.ErrorResponse "Empty or invalid credentials"
// @failure     401 {object} model.ErrorResponse "Password mismatch"
// @failure     404 {object} model.ErrorResponse "User not found"
// @failure     500 {object} model.ErrorResponse
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
		Name:    "session",
		Value:   token,
		Expires: time.Now().Add(time.Duration(h.cfg.SessionExpire) * time.Second),
	})

	// Return a success response
	return utils.ResponseStatus(c, http.StatusOK)
}
