// internal/login/handler.go
package login

import (
	"github.com/gofiber/fiber/v2"
)

type Handler interface {
	Login(c *fiber.Ctx) error
}

type handlerImpl struct {
	service Service
}

func NewHandler(service Service) Handler {
	return &handlerImpl{
		service,
	}
}

func (h *handlerImpl) Login(c *fiber.Ctx) error {
	// Parse login request from the request body
	var loginRequest LoginRequest
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
		Name:  "jwt",
		Value: token,
	})

	// Return a success response
	return c.SendStatus(fiber.StatusOK)
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
