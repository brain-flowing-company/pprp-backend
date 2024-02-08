package middleware

import (
	"fmt"

	"github.com/brain-flowing-company/pprp-backend/apperror"
	"github.com/brain-flowing-company/pprp-backend/config"
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/brain-flowing-company/pprp-backend/utils"
	"github.com/gofiber/fiber/v2"
)

func NewAuthMiddlware(cfg *config.Config) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		cookie := new(models.Cookie)

		err := c.CookieParser(cookie)
		if err != nil {
			fmt.Println(err)
			return utils.ResponseError(c, apperror.Unauthorized)
		}

		claim, err := utils.ParseToken(cookie.Session, cfg.JWTSecret)
		if err != nil {
			fmt.Println(err)
			return utils.ResponseError(c, apperror.Unauthorized)
		}

		c.Locals("email", claim.Session.Email)

		return c.Next()
	}
}
