package ratings

import (
	"fmt"

	"github.com/brain-flowing-company/pprp-backend/apperror"
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/brain-flowing-company/pprp-backend/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Handler interface {
	CreateRating(c *fiber.Ctx) error
	GetRatingByPropertyId(c *fiber.Ctx) error
	GetAllRatings(c *fiber.Ctx) error
	GetRatingByPropertyIdSortedByRating(c *fiber.Ctx) error
	GetRatingByPropertyIdSortedByNewest(c *fiber.Ctx) error
	UpdateRatingStatus(c *fiber.Ctx) error
	DeleteRating(c *fiber.Ctx) error
}

type handlerImpl struct {
	service Service
}

func NewHandler(service Service) Handler {
	return &handlerImpl{
		service,
	}
}

// @router /api/v1/ratings [post]
// @summary Create rating
// @description Create rating
// @tags ratings
// @produce json
// @param rating body int true "rating"
// @param review body string true "review"
// @param property_id body string true "property_id"
// @success 200 {object} models.RatingResponse
// @failure 400 {object} models.ErrorResponses
// @failure 401 {object} models.ErrorResponses
// @failure 500 {object} models.ErrorResponses
func (h *handlerImpl) CreateRating(c *fiber.Ctx) error {
	session, ok := c.Locals("session").(models.Sessions)
	if !ok {
		return utils.ResponseError(c, apperror.New(apperror.Unauthorized).Describe("Unauthorized"))
	}
	reviews := models.Reviews{
		ReviewId:      uuid.New(),
		DwellerUserId: session.UserId,
	}
	if err := c.BodyParser(&reviews); err != nil {
		return utils.ResponseError(c, apperror.New(apperror.BadRequest).Describe("Failed to parse body"))
	}
	if !utils.IsValidRating(reviews.Rating) {
		return utils.ResponseError(c, apperror.New(apperror.BadRequest).Describe("Invalid rating"))
	}
	if err := h.service.CreateRating(&reviews); err != nil {
		return utils.ResponseError(c, err)
	}
	return c.JSON(fiber.Map{
		"success":         true,
		"message":         "Rating created successfully",
		"rating":          reviews.Rating,
		"review":          reviews.Review,
		"property_id":     reviews.PropertyId,
		"dweller_user_id": reviews.DwellerUserId,
	})
}

// @router /api/v1/ratings/:propertyId [get]
// @summary Get rating by property id
// @description Get rating by property id
// @tags ratings
// @produce json
// @param propertyId path string true "propertyId"
// @success 200 {object} models.RatingResponse
// @failure 400 {object} models.ErrorResponses
// @failure 401 {object} models.ErrorResponses
// @failure 500 {object} models.ErrorResponses
func (h *handlerImpl) GetRatingByPropertyId(c *fiber.Ctx) error {
	propertyId := c.Params("propertyId")
	fmt.Println("propertyId", propertyId)
	parsedPropertyID, err := uuid.Parse(propertyId)
	if err != nil {
		return utils.ResponseError(c, apperror.New(apperror.BadRequest).Describe("Invalid property ID"))
	}
	var ratings []models.RatingResponse
	if err := h.service.GetRatingByPropertyId(parsedPropertyID, &ratings); err != nil {
		return utils.ResponseError(c, err)
	}
	return c.JSON(ratings)
}

// @router /api/v1/ratings [get]
// @summary Get all ratings
// @description Get all ratings
// @tags ratings
// @produce json
// @success 200 {object} models.RatingResponse
// @failure 400 {object} models.ErrorResponses
// @failure 401 {object} models.ErrorResponses
// @failure 500 {object} models.ErrorResponses
func (h *handlerImpl) GetAllRatings(c *fiber.Ctx) error {
	var ratings []models.RatingResponse
	if err := h.service.GetAllRatings(&ratings); err != nil {
		return utils.ResponseError(c, err)
	}
	return c.JSON(ratings)
}

// @router /api/v1/ratings/sorted/:propertyId [get]
// @summary Get rating by property id sorted by rating
// @description Get rating by property id sorted by rating
// @tags ratings
// @produce json
// @param propertyId path string true "propertyId"
// @success 200 {object} models.RatingResponse
// @failure 400 {object} models.ErrorResponses
// @failure 401 {object} models.ErrorResponses
// @failure 500 {object} models.ErrorResponses
func (h *handlerImpl) GetRatingByPropertyIdSortedByRating(c *fiber.Ctx) error {
	propertyId := c.Params("propertyId")
	fmt.Println("propertyId", propertyId)
	parsedPropertyID, err := uuid.Parse(propertyId)
	if err != nil {
		return utils.ResponseError(c, apperror.New(apperror.BadRequest).Describe("Invalid property ID"))
	}
	var ratings []models.RatingResponse
	if err := h.service.GetRatingByPropertyIdSortedByRating(parsedPropertyID, &ratings); err != nil {
		return utils.ResponseError(c, err)
	}
	return c.JSON(ratings)
}

// @router /api/v1/ratings/newest/:propertyId [get]
// @summary Get rating by property id sorted by newest
// @description Get rating by property id sorted by newest
// @tags ratings
// @produce json
// @param propertyId path string true "propertyId"
// @success 200 {object} models.RatingResponse
// @failure 400 {object} models.ErrorResponses
// @failure 401 {object} models.ErrorResponses
// @failure 500 {object} models.ErrorResponses
func (h *handlerImpl) GetRatingByPropertyIdSortedByNewest(c *fiber.Ctx) error {
	propertyId := c.Params("propertyId")
	fmt.Println("propertyId", propertyId)
	parsedPropertyID, err := uuid.Parse(propertyId)
	if err != nil {
		return utils.ResponseError(c, apperror.New(apperror.BadRequest).Describe("Invalid property ID"))
	}
	var ratings []models.RatingResponse
	if err := h.service.GetRatingByPropertyIdSortedByNewest(parsedPropertyID, &ratings); err != nil {
		return utils.ResponseError(c, err)
	}
	return c.JSON(ratings)
}

// apiv1.Patch("/ratings/:ratingId", mw.AuthMiddlewareWrapper(ratingsHandler.UpdateRatingStatus))

// @router /api/v1/ratings/:ratingId [patch]
// @summary Update rating status
// @description Update rating status
// @tags ratings
// @produce json
// @param ratingId path string true "ratingId"
// @param review body string true "review"
// @param rating body int true "rating"
// @success 200 {object} models.MessageResponses
// @failure 400 {object} models.ErrorResponses
// @failure 401 {object} models.ErrorResponses
// @failure 500 {object} models.ErrorResponses
func (h *handlerImpl) UpdateRatingStatus(c *fiber.Ctx) error {
	updatingRating := models.UpdateRatingStatus{}
	err := c.BodyParser(&updatingRating)
	if err != nil {
		return utils.ResponseError(c, apperror.New(apperror.BadRequest).Describe("Failed to parse body"))
	}
	if !utils.IsValidRating(updatingRating.Rating) {
		return utils.ResponseError(c, apperror.New(apperror.BadRequest).Describe("Invalid rating"))
	}
	ratingId := c.Params("ratingId")
	ratingIdParsed, err := uuid.Parse(ratingId)
	if err != nil {
		return utils.ResponseError(c, apperror.New(apperror.BadRequest).Describe("Invalid rating ID"))
	}
	apperr := h.service.UpdateRatingStatus(&updatingRating, ratingIdParsed)
	if apperr != nil {
		return utils.ResponseError(c, apperr)
	}
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Rating status updated successfully",
	})

}

// @router /api/v1/ratings/:ratingId [delete]
// @summary Delete rating
// @description Delete rating
// @tags ratings
// @produce json
// @param ratingId path string true "ratingId"
// @success 200 {object} models.MessageResponses
// @failure 400 {object} models.ErrorResponses
// @failure 401 {object} models.ErrorResponses
// @failure 500 {object} models.ErrorResponses
func (h *handlerImpl) DeleteRating(c *fiber.Ctx) error {
	ratingId := c.Params("ratingId")
	err := h.service.DeleteRating(ratingId)
	if err != nil {
		return utils.ResponseError(c, err)
	}
	return c.JSON(fiber.Map{
		"message": "Rating deleted",
	})

}
