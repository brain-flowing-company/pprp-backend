package ratings

import (
	"github.com/brain-flowing-company/pprp-backend/apperror"
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/brain-flowing-company/pprp-backend/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Handler interface {
	CreateRating(c *fiber.Ctx) error
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
	reviews := models.Reviews{
		ReviewId: uuid.New(),
	}
	if err := c.BodyParser(&reviews); err != nil {
		return utils.ResponseError(c, apperror.New(apperror.BadRequest).Describe("Failed to parse body"))
	}
	if utils.IsValidRating(string(reviews.Rating)) {
		return utils.ResponseError(c, apperror.New(apperror.BadRequest).Describe("Invalid rating"))
	}
	if err := h.service.CreateRating(&reviews); err != nil {
		return utils.ResponseError(c, err)
	}
	return utils.ResponseStatus(c, fiber.StatusCreated)
}
