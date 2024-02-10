package property

import (
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/brain-flowing-company/pprp-backend/utils"
	"github.com/gofiber/fiber/v2"
)

type Handler interface {
	GetPropertyById(c *fiber.Ctx) error
	GetAllProperties(c *fiber.Ctx) error
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
// @failure     400 {object} model.ErrorResponse "Invalid property id"
// @failure     404 {object} model.ErrorResponse "Property id not found"
// @failure     500 {object} model.ErrorResponse
func (h *handlerImpl) GetPropertyById(c *fiber.Ctx) error {
	propertyId := c.Params("propertyId")

	property := models.Property{}
	err := h.service.GetPropertyById(&property, propertyId)
	if err != nil {
		return utils.ResponseError(c, err)
	}

	return c.JSON(property)
}

// @router      /api/v1/properties [get]
// @summary     Get all properties
// @description Get all properties
// @tags        property
// @produce     json
// @success     200	{object} []models.Property
// @failure     500 {object} model.ErrorResponse
func (h *handlerImpl) GetAllProperties(c *fiber.Ctx) error {
	properties := []models.Property{}
	err := h.service.GetAllProperties(&properties)
	if err != nil {
		return utils.ResponseError(c, err)
	}

	return c.JSON(properties)
}
