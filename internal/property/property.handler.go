package property

import (
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/gofiber/fiber/v2"
)

type Handler interface {
	GetPropertyById(c *fiber.Ctx) error
}

type handlerImpl struct {
	service Service
}

func NewHandler(service Service) Handler {
	return &handlerImpl{
		service,
	}
}

// @router      /api/v1/property/:id [get]
// @summary     Get property by id
// @description Get property by its id
// @tags        property
// @produce     json
// @success     200	{object} models.Property
// @failure     400 {object} apperror.AppError
// @failure     404 {object} apperror.AppError
func (h *handlerImpl) GetPropertyById(c *fiber.Ctx) error {
	propertyId := c.Params("propertyId")

	property := models.Property{}
	err := h.service.GetPropertyById(&property, propertyId)
	if err != nil {
		return c.Status(err.Code).JSON(err)
	}

	return c.JSON(property)
}
