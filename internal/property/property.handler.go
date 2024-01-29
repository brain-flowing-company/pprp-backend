package property

import (
	"github.com/brain-flowing-company/pprp-backend/internal/dto"
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

// @router      /api/property/:id [get]
// @summary     Get property by id
// @description Get property by id
// @produce     json
// @success     200	{object}	dto.GreetingResponse
func (h *handlerImpl) GetPropertyById(c *fiber.Ctx) error {
	propertyId := c.Params("propertyId")

	property := models.Property{}
	err := h.service.GetPropertyById(&property, propertyId)
	if err != nil {
		return c.Status(err.Code).JSON(err)
	}

	res := dto.GetPropertyByIdResponse{
		Property: property,
	}

	return c.JSON(res)
}
