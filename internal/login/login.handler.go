// internal/login/handler.go
package login

import (
	"net/http"
	"time"

	"github.com/brain-flowing-company/pprp-backend/config"
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/gofiber/fiber/v2"
)

type Handler interface {
	Login(c *fiber.Ctx) error
}

type handlerImpl struct {
	service Service
	cfg     *config.Config
}

func NewHandler(service Service, cfg *config.Config) Handler {
	return &handlerImpl{
		service,
		cfg,
	}
}

func (h *handlerImpl) Login(c *fiber.Ctx) error {
	// Parse login request from the request body
	var loginRequest models.LoginRequest
	if err := c.BodyParser(&loginRequest); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	// Authenticate user
	token, err := h.service.AuthenticateUser(loginRequest.Email, loginRequest.Password)
	if err != nil {
		return c.Status(err.Code).JSON(err)
	}

	// Set JWT token as a cookie
	c.Cookie(&fiber.Cookie{
		Name:    "session",
		Value:   token,
		Expires: time.Now().Add(time.Duration(h.cfg.JWTMaxAge) * time.Second),
	})

	// Return a success response
	return c.SendStatus(http.StatusOK)
}
