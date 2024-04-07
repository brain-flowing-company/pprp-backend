package middleware

import (
	"net/http"

	"github.com/brain-flowing-company/pprp-backend/config"
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/brain-flowing-company/pprp-backend/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type Middleware struct {
	cfg *config.Config
}

func NewMiddleware(cfg *config.Config) Middleware {
	return Middleware{
		cfg,
	}
}

func (m *Middleware) SessionMiddleware(c *fiber.Ctx) error {
	cookie := new(models.Cookies)

	err := c.CookieParser(cookie)
	if err != nil {
		return c.Next()
	}

	claim, err := utils.ParseToken(cookie.Session, m.cfg.JWTSecret)
	if err == nil {
		c.Locals("session", claim.Session)
	}

	return c.Next()
}

func (m *Middleware) WithAuthentication(next func(*fiber.Ctx) error) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		_, ok := c.Locals("session").(models.Sessions)
		if !ok {
			return utils.ResponseMessage(c, http.StatusUnauthorized, "Unauthorized")
		}

		return next(c)
	}
}

func (m *Middleware) WithOwnerAccess(next func(*fiber.Ctx) error) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		session, ok := c.Locals("session").(models.Sessions)
		if !ok || !session.IsOwner {
			return utils.ResponseMessage(c, http.StatusForbidden, "Require owner access")
		}

		return next(c)
	}
}
