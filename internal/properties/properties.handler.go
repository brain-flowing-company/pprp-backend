package properties

import (
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/brain-flowing-company/pprp-backend/utils"
	"github.com/gofiber/fiber/v2"
)

type Handler interface {
	GetPropertyById(c *fiber.Ctx) error
	GetAllProperties(c *fiber.Ctx) error
	SeachProperties(c *fiber.Ctx) error
	GetOrSearchProperties(c *fiber.Ctx) error
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
// @success     200	{object} models.Properties
// @failure     400 {object} models.ErrorResponse "Invalid property id"
// @failure     404 {object} models.ErrorResponse "Property id not found"
// @failure     500 {object} models.ErrorResponse
func (h *handlerImpl) GetPropertyById(c *fiber.Ctx) error {
	propertyId := c.Params("propertyId")

	property := models.Properties{}
	err := h.service.GetPropertyById(&property, propertyId)
	if err != nil {
		return utils.ResponseError(c, err)
	}

	return c.JSON(property)
}

// @router      /api/v1/properties [get]
// @summary     Get all properties or search properties
// @description If a query parameter is provided, search properties by project name or description. Otherwise, get all properties.
// @tags        property
// @produce     json
// @param       query query string true "Search query"
// @success     200	{object} []models.Properties
// @failure     500 {object} models.ErrorResponse
func (h *handlerImpl) GetOrSearchProperties(c *fiber.Ctx) error {
	query := c.Query("query")
	if query != "" {
		return h.SeachProperties(c)
	} else {
		return h.GetAllProperties(c)
	}
}

func (h *handlerImpl) GetAllProperties(c *fiber.Ctx) error {
	properties := []models.Properties{}
	err := h.service.GetAllProperties(&properties)
	if err != nil {
		return utils.ResponseError(c, err)
	}

	return c.JSON(properties)
}

func (h *handlerImpl) SeachProperties(c *fiber.Ctx) error {
	query := c.Query("query")

	properties := []models.Properties{}
	err := h.service.SearchProperties(&properties, query)
	if err != nil {
		return utils.ResponseError(c, err)
	}

	return c.JSON(properties)
}
