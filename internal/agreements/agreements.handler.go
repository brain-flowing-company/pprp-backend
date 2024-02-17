package agreements

import (
	"github.com/brain-flowing-company/pprp-backend/apperror"
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/gofiber/fiber/v2"
)

type Handler interface {
	GetAllAgreements(c *fiber.Ctx) error
	GetAgreementById(c *fiber.Ctx) error
	GetAgreementsByOwnerId(c *fiber.Ctx) error
	GetAgreementsByDwellerId(c *fiber.Ctx) error
	CreateAgreement(c *fiber.Ctx) error
	DeleteAgreement(c *fiber.Ctx) error
}
type handlerImpl struct {
	service Service
}

func NewHandler(service Service) Handler {
	return &handlerImpl{
		service,
	}
}

func (h *handlerImpl) GetAllAgreements(c *fiber.Ctx) error {
	var apps []models.Agreement
	err := h.service.GetAllAgreements(&apps)
	if err != nil {
		return apperror.New(apperror.InternalServerError).Describe("Error getting all agreements")
	}
	return c.JSON(apps)

}

func (h *handlerImpl) GetAgreementById(c *fiber.Ctx) error {
	id := c.Params("agreementId")
	var app models.Agreement
	err := h.service.GetAgreementById(&app, id)
	if err != nil {
		return apperror.New(apperror.InternalServerError).Describe("Error getting agreement by id")
	}
	return c.JSON(app)
}

func (h *handlerImpl) GetAgreementsByOwnerId(c *fiber.Ctx) error {
	id := c.Params("userId")
	var apps []models.Agreement
	err := h.service.GetAgreementsByOwnerId(&apps, id)
	if err != nil {
		return apperror.New(apperror.InternalServerError).Describe("Error getting agreements by owner id")
	}
	return c.JSON(apps)
}

func (h *handlerImpl) GetAgreementsByDwellerId(c *fiber.Ctx) error {
	id := c.Params("userId")
	var apps []models.Agreement
	err := h.service.GetAgreementsByDwellerId(&apps, id)
	if err != nil {
		return apperror.New(apperror.InternalServerError).Describe("Error getting agreements by dweller id")
	}
	return c.JSON(apps)
}

func (h *handlerImpl) CreateAgreement(c *fiber.Ctx) error {
	var creatingAgreement models.CreatingAgreement
	if err := c.BodyParser(&creatingAgreement); err != nil {
		return apperror.New(apperror.InvalidBody).Describe("Invalid body")
	}
	err := h.service.CreateAgreement(&creatingAgreement)
	if err != nil {
		return apperror.New(apperror.InternalServerError).Describe("Error creating agreement")
	}
	return c.SendStatus(fiber.StatusCreated)
}

func (h *handlerImpl) DeleteAgreement(c *fiber.Ctx) error {
	id := c.Params("agreementId")
	err := h.service.DeleteAgreement(id)
	if err != nil {
		return apperror.New(apperror.InternalServerError).Describe("Error deleting agreement")
	}
	return c.SendStatus(fiber.StatusNoContent)
}
