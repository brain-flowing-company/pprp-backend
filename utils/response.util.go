package utils

import (
	"github.com/brain-flowing-company/pprp-backend/apperror"
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/gofiber/fiber/v2"
)

func ResponseError(c *fiber.Ctx, err interface{}) error {
	r := models.ErrorResponse{}

	switch appErr := err.(type) {
	case *apperror.AppError:
		r.Code = appErr.Code()
		r.Name = appErr.Name()
		r.Message = appErr.Error()

	case *apperror.AppErrorType:
		r.Code = appErr.Code
		r.Name = appErr.Name
	}

	return c.Status(r.Code).JSON(r)
}

func ResponseStatus(c *fiber.Ctx, status int) error {
	return c.Status(status).Send(nil)
}

func ResponseMessage(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(fiber.Map{
		"message": message,
	})
}
