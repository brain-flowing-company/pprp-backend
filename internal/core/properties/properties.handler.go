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
	GetMyProperties(c *fiber.Ctx) error
	CreateProperty(c *fiber.Ctx) error
	UpdatePropertyById(c *fiber.Ctx) error
	DeletePropertyById(c *fiber.Ctx) error
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

// @router      /api/v1/properties/:propertyId [get]
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
// @summary     Get or search properties
// @description Get all properties or search properties by query
// @tags        property
// @produce     json
// @param       query query string false "Search query"
// @param       limit query int false "Pagination limit per page, max 50, default 20"
// @param       page  query int false "Pagination page index as 1-based index, default 1"
// @success     200	{object} models.AllPropertiesResponses
// @failure     500 {object} models.ErrorResponses "Could not get properties"
func (h *handlerImpl) GetAllProperties(c *fiber.Ctx) error {
	query := c.Query("query")
	properties := models.AllPropertiesResponses{}

	sorted := utils.NewSortedQuery(models.Properties{})
	err := sorted.ParseQuery(c.Query("sort"))
	if err != nil {
		return utils.ResponseError(c, apperror.
			New(apperror.BadRequest).
			Describe(err.Error()))
	}

	var userId string
	if _, ok := c.Locals("session").(models.Sessions); !ok {
		userId = "00000000-0000-0000-0000-000000000000"
	} else {
		userId = c.Locals("session").(models.Sessions).UserId.String()
	}

	limit := utils.Clamp(c.QueryInt("limit", 20), 1, 50)
	page := utils.Max(c.QueryInt("page", 1), 1)

	paginated := utils.NewPaginatedQuery(page, limit)

	apperr := h.service.GetAllProperties(&properties, query, userId, paginated, sorted)
	if apperr != nil {
		return utils.ResponseError(c, apperr)
	}

	return c.JSON(properties)
}

// @router      /api/v1/user/me/properties [get]
// @summary     Get my properties
// @description Get all properties owned by the current user
// @tags        property
// @produce     json
// @param       limit query int false "Pagination limit per page, max 50, default 20"
// @param       page  query int false "Pagination page index as 1-based index, default 1"
// @success     200	{object} models.MyPropertiesResponses
// @failure	    403 {object} models.ErrorResponses "Unauthorized"
// @failure     500 {object} models.ErrorResponses
func (h *handlerImpl) GetMyProperties(c *fiber.Ctx) error {
	userId := c.Locals("session").(models.Sessions).UserId.String()

	limit := utils.Clamp(c.QueryInt("limit", 20), 1, 50)
	page := utils.Max(c.QueryInt("page", 1), 1)

	paginated := utils.NewPaginatedQuery(page, limit)

	properties := models.MyPropertiesResponses{}
	err := h.service.GetPropertyByOwnerId(&properties, userId, paginated)
	if err != nil {
		return utils.ResponseError(c, err)
	}

	return c.JSON(properties)
}

// @router      /api/v1/properties [post]
// @summary     Create a property
// @description Create a property with the provided details
// @tags        property
// @produce     json
// @param       formData formData models.PropertyInfos true "Property details"
// @success     200	{object} models.MessageResponses "Property created"
// @failure     400 {object} models.ErrorResponses "Invalid request body"
// @failure	    403 {object} models.ErrorResponses "Unauthorized"
// @failure     404 {object} models.ErrorResponses "Property id not found"
// @failure     500 {object} models.ErrorResponses "Could not create property"
func (h *handlerImpl) CreateProperty(c *fiber.Ctx) error {
	property := models.PropertyInfos{}
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

	return utils.ResponseMessage(c, http.StatusOK, "Property created")
}

// @router      /api/v1/properties/:propertyId [put]
// @summary     Update a property
// @description Update a property, owned by the current user, by its id with the provided details
// @tags        property
// @produce     json
// @param	    propertyId path string true "Property id"
// @param       formData formData models.PropertyInfos true "Property details"
// @success     200	{object} models.MessageResponses "Property updated"
// @failure     400 {object} models.ErrorResponses "Invalid request body"
// @failure	    403 {object} models.ErrorResponses "Unauthorized"
// @failure     404 {object} models.ErrorResponses "Property id not found"
// @failure     500 {object} models.ErrorResponses "Could not update property"
func (h *handlerImpl) UpdatePropertyById(c *fiber.Ctx) error {
	propertyIdString := c.Params("propertyId")
	propertyIdUuid, _ := uuid.Parse(propertyIdString)

	property := models.PropertyInfos{}
	userId := c.Locals("session").(models.Sessions).UserId

	property.OwnerId = userId
	property.PropertyId = propertyIdUuid

	if err := c.BodyParser(&property); err != nil {
		return utils.ResponseError(c, apperror.InvalidBody)
	}

	err := h.service.UpdatePropertyById(&property, propertyIdString)
	if err != nil {
		return utils.ResponseError(c, err)
	}

	return utils.ResponseMessage(c, http.StatusOK, "Property updated")
}

// @router      /api/v1/properties/:propertyId [delete]
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

// @router      /api/v1/properties/favorites/:propertyId [post]
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

// @router      /api/v1/properties/favorites/:propertyId [delete]
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
// @param       limit query int false "Pagination limit per page, max 50, default 20"
// @param       page  query int false "Pagination page index as 1-based index, default 1"
// @success     200	{object} models.MyFavoritePropertiesResponses
// @failure	    403 {object} models.ErrorResponses "Unauthorized"
// @failure     500 {object} models.ErrorResponses "Could not get favorite properties"
func (h *handlerImpl) GetMyFavoriteProperties(c *fiber.Ctx) error {
	userId := c.Locals("session").(models.Sessions).UserId.String()

	limit := utils.Clamp(c.QueryInt("limit", 20), 1, 50)
	page := utils.Max(c.QueryInt("page", 1), 1)

	paginated := utils.NewPaginatedQuery(page, limit)

	properties := models.MyFavoritePropertiesResponses{}
	err := h.service.GetFavoritePropertiesByUserId(&properties, userId, paginated)
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
	var userId string
	if _, ok := c.Locals("session").(models.Sessions); !ok {
		userId = "00000000-0000-0000-0000-000000000000"
	} else {
		userId = c.Locals("session").(models.Sessions).UserId.String()
	}

	properties := []models.Properties{}
	err := h.service.GetTop10Properties(&properties, userId)
	if err != nil {
		return utils.ResponseError(c, err)
	}

	return c.JSON(properties)
}
