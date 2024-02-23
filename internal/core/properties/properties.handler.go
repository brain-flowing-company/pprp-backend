package properties

import (
	"net/http"

	"github.com/brain-flowing-company/pprp-backend/apperror"
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/brain-flowing-company/pprp-backend/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Handler interface {
	GetPropertyById(c *fiber.Ctx) error
	GetAllProperties(c *fiber.Ctx) error
	CreateProperty(c *fiber.Ctx) error
	UpdatePropertyById(c *fiber.Ctx) error
	DeletePropertyById(c *fiber.Ctx) error
	SeachProperties(c *fiber.Ctx) error
	GetOrSearchProperties(c *fiber.Ctx) error
	AddFavoriteProperty(c *fiber.Ctx) error
	RemoveFavoriteProperty(c *fiber.Ctx) error
	GetMyFavoriteProperties(c *fiber.Ctx) error
	GetTop10Properties(c *fiber.Ctx) error
}

type handlerImpl struct {
	service Service
}

func NewHandler(service Service) Handler {
	return &handlerImpl{
		service,
	}
}

// @router      /api/v1/property/:propertyId [get]
// @summary     Get property by propertyId
// @description Get property by its id
// @tags        property
// @produce     json
// @param	    propertyId path string true "Property id"
// @success     200	{object} models.Properties
// @failure     400 {object} models.ErrorResponses "Invalid property id"
// @failure     404 {object} models.ErrorResponses "Property id not found"
// @failure     500 {object} models.ErrorResponses
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
// @failure     500 {object} models.ErrorResponses
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

// @router      /api/v1/property [post]
// @summary     Create a property
// @description Create a property with the provided details
// @tags        property
// @produce     json
// @param       body body models.Properties true "Property details"
// @success     200	{object} []models.Properties
// @failure     400 {object} models.ErrorResponses "Invalid request body"
// @failure	    403 {object} models.ErrorResponses "Unauthorized"
// @failure     404 {object} models.ErrorResponses "Property id not found"
// @failure     500 {object} models.ErrorResponses "Could not create property"
func (h *handlerImpl) CreateProperty(c *fiber.Ctx) error {
	property := models.Properties{}
	if err := c.BodyParser(&property); err != nil {
		return utils.ResponseError(c, apperror.
			New(apperror.InvalidBody).
			Describe("Invalid request body"))
	}

	userId := c.Locals("session").(models.Sessions).UserId
	property.OwnerId = userId

	err := h.service.CreateProperty(&property)
	if err != nil {
		return utils.ResponseError(c, err)
	}

	return c.JSON(property)
}

// @router      /api/v1/property/:propertyId [put]
// @summary     Update a property
// @description Update a property, owned by the current user, by its id with the provided details
// @tags        property
// @produce     json
// @param	    propertyId path string true "Property id"
// @param       body body models.Properties true "Property details"
// @success     200	{object} []models.Properties
// @failure     400 {object} models.ErrorResponses "Invalid request body"
// @failure	    403 {object} models.ErrorResponses "Unauthorized"
// @failure     404 {object} models.ErrorResponses "Property id not found"
// @failure     500 {object} models.ErrorResponses "Could not update property"
func (h *handlerImpl) UpdatePropertyById(c *fiber.Ctx) error {
	propertyId := c.Params("propertyId")
	property := models.Properties{}
	if err := c.BodyParser(&property); err != nil {
		return utils.ResponseError(c, apperror.
			New(apperror.InvalidBody).
			Describe("Invalid request body"))
	}

	userId := c.Locals("session").(models.Sessions).UserId
	property.OwnerId = userId
	property.PropertyId, _ = uuid.Parse(propertyId)

	err := h.service.UpdatePropertyById(&property, propertyId)
	if err != nil {
		return utils.ResponseError(c, err)
	}

	return c.JSON(property)
}

// @router      /api/v1/property/:propertyId [delete]
// @summary     Delete a property
// @description Delete a property, owned by the current user, by its id
// @tags        property
// @produce     json
// @param	    propertyId path string true "Property id"
// @success     200	{object} models.MessageResponses "Property deleted"
// @failure     400 {object} models.ErrorResponses "Invalid request body"
// @failure	    403 {object} models.ErrorResponses "Unauthorized"
// @failure     404 {object} models.ErrorResponses "Property id not found"
// @failure     500 {object} models.ErrorResponses "Could not delete property"
func (h *handlerImpl) DeletePropertyById(c *fiber.Ctx) error {
	propertyId := c.Params("propertyId")

	err := h.service.DeletePropertyById(propertyId)
	if err != nil {
		return utils.ResponseError(c, err)
	}

	return utils.ResponseMessage(c, http.StatusOK, "Property deleted")
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

// @router      /api/v1/property/:propertyId [post]
// @summary     Add property to favorites
// @description Add property to the current user favorites
// @tags        property
// @produce     json
// @param       propertyId path string true "Property id"
// @success     200	{object} models.MessageResponses "Property added to favorites"
// @failure	    403 {object} models.ErrorResponses "Unauthorized"
// @failure     404 {object} models.ErrorResponses "Property id not found"
// @failure     500 {object} models.ErrorResponses "Could not add favorite property"
func (h *handlerImpl) AddFavoriteProperty(c *fiber.Ctx) error {
	propertyId := c.Params("propertyId")
	userId := c.Locals("session").(models.Sessions).UserId

	err := h.service.AddFavoriteProperty(propertyId, userId)
	if err != nil {
		return utils.ResponseError(c, err)
	}

	return utils.ResponseMessage(c, http.StatusOK, "Property added to favorites")
}

// @router      /api/v1/property/:propertyId [delete]
// @summary     Remove property to favorites
// @description Remove property to the current user favorites
// @tags        property
// @produce     json
// @param       propertyId path string true "Property id"
// @success     200	{object} models.MessageResponses "Property removed from favorites"
// @failure	    403 {object} models.ErrorResponses "Unauthorized"
// @failure     404 {object} models.ErrorResponses "Property id not found"
// @failure     500 {object} models.ErrorResponses "Could not remove favorite property"
func (h *handlerImpl) RemoveFavoriteProperty(c *fiber.Ctx) error {
	propertyId := c.Params("propertyId")
	userId := c.Locals("session").(models.Sessions).UserId

	err := h.service.RemoveFavoriteProperty(propertyId, userId)
	if err != nil {
		return utils.ResponseError(c, err)
	}

	return utils.ResponseMessage(c, http.StatusOK, "Property removed from favorites")
}

// @router      /api/v1/user/me/favorites [get]
// @summary     Get my favorite properties
// @description Get all properties that the current user has added to favorites
// @tags        property
// @produce     json
// @success     200	{object} []models.Properties
// @failure	    403 {object} models.ErrorResponses "Unauthorized"
// @failure     500 {object} models.ErrorResponses "Could not get favorite properties"
func (h *handlerImpl) GetMyFavoriteProperties(c *fiber.Ctx) error {
	userId := c.Locals("session").(models.Sessions).UserId.String()

	properties := []models.Properties{}
	err := h.service.GetFavoritePropertiesByUserId(&properties, userId)
	if err != nil {
		return utils.ResponseError(c, err)
	}

	return c.JSON(properties)
}

// @router      /api/v1/properties/top10 [get]
// @summary     Get top 10 properties
// @description Get top 10 properties with the most favorites, sorted by the number of favorites then by the newest properties
// @tags        property
// @produce     json
// @success     200	{object} []models.Properties
// @failure     500 {object} models.ErrorResponses "Could not get top 10 properties"
func (h *handlerImpl) GetTop10Properties(c *fiber.Ctx) error {
	properties := []models.Properties{}
	err := h.service.GetTop10Properties(&properties)
	if err != nil {
		return utils.ResponseError(c, err)
	}

	return c.JSON(properties)
}
