package middleware

import (
	"net/http"

	"github.com/brain-flowing-company/pprp-backend/config"
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/brain-flowing-company/pprp-backend/utils"
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

func (m *Middleware) AuthMiddlewareWrapper(next func(*fiber.Ctx) error) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		_, ok := c.Locals("session").(models.Session)
		if !ok {
			return utils.ResponseMessage(c, http.StatusUnauthorized, "Unauthorized")
		}

		return next(c)
	}
}

func (m *Middleware) SessionMiddleware(c *fiber.Ctx) error {
	cookie := new(models.Cookie)

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
