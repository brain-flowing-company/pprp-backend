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
}

type handlerImpl struct {
	service Service
}

func NewHandler(service Service) Handler {
	return &handlerImpl{
		service,
	}
}

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
	if !utils.IsValidRating(string(reviews.Rating)) {
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

func (h *handlerImpl) GetAllRatings(c *fiber.Ctx) error {
	var ratings []models.RatingResponse
	if err := h.service.GetAllRatings(&ratings); err != nil {
		return utils.ResponseError(c, err)
	}
	return c.JSON(ratings)
}

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
