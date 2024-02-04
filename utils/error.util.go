package utils

import (
	"github.com/brain-flowing-company/pprp-backend/apperror"
	"github.com/gofiber/fiber/v2"
)

func ResponseError(c *fiber.Ctx, err *apperror.AppError) error {
	return c.Status(err.Code).JSON(err)
}
