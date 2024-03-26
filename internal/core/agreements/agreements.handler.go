package agreements

import (
	"fmt"

	"github.com/brain-flowing-company/pprp-backend/apperror"
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/brain-flowing-company/pprp-backend/internal/utils"
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

// @router  /api/v1/agreements [get]
// @summary  Get all agreements
// @description  Get all agreements
// @tags agreements
// @produce json
// @success 200 {object} []models.Agreements
// @failure 500 {object} models.ErrorResponses
func (h *handlerImpl) GetAllAgreements(c *fiber.Ctx) error {
	var apps []models.Agreements
	err := h.service.GetAllAgreements(&apps)
	if err != nil {
		return apperror.New(apperror.InternalServerError).Describe("Error getting all agreements")
	}
	return c.JSON(apps)

}

// @router  /api/v1/agreement/:agreementId [get]
// @summary  Get agreement by id
// @description  Get an agreement by its id
// @tags agreements
// @produce json
// @success 200 {object} models.AgreementDetails
// @failure 400 {object} models.MessageResponses "Invalid agreement id"
// @failure 404 {object} models.MessageResponses "Agreement not found"
// @failure 500 {object} models.MessageResponses
func (h *handlerImpl) GetAgreementById(c *fiber.Ctx) error {
	agreementId := c.Params("agreementId")
	agreement := &models.AgreementDetails{}

	err := h.service.GetAgreementById(agreement, agreementId)
	if err != nil {
		return utils.ResponseError(c, err)
	}
	
	return c.JSON(agreement)
}

// @router  /api/v1/agreements/owner/:userId [get]
// @summary  Get agreements by owner id
// @description  Get all agreements by owner id
// @tags agreements
// @produce json
// @success 200 {object} []models.Agreements
// @failure 500 {object} models.ErrorResponses
func (h *handlerImpl) GetAgreementsByOwnerId(c *fiber.Ctx) error {
	id := c.Params("userId")
	var apps []models.Agreements
	err := h.service.GetAgreementsByOwnerId(&apps, id)
	if err != nil {
		return apperror.New(apperror.InternalServerError).Describe("Error getting agreements by owner id")
	}
	return c.JSON(apps)
}

// @router  /api/v1/agreements/dweller/:userId [get]
// @summary  Get agreements by dweller id
// @description  Get all agreements by dweller id
// @tags agreements
// @produce json
// @success 200 {object} []models.Agreements
// @failure 500 {object} models.ErrorResponses
func (h *handlerImpl) GetAgreementsByDwellerId(c *fiber.Ctx) error {
	id := c.Params("userId")
	var apps []models.Agreements
	err := h.service.GetAgreementsByDwellerId(&apps, id)
	if err != nil {
		return apperror.New(apperror.InternalServerError).Describe("Error getting agreements by dweller id")
	}
	return c.JSON(apps)
}

// @router  /api/v1/agreements [post]
// @summary  Create an agreement
// @description  Create an agreement by parsing the body
// @tags agreements
// @produce json
// @success 201 {object} models.MessageResponses "Agreement created successfully"
// @failure 500 {object} models.ErrorResponses
func (h *handlerImpl) CreateAgreement(c *fiber.Ctx) error {
	agreement := &models.CreatingAgreements{}
	err := c.BodyParser(agreement)
	if err != nil {
		return utils.ResponseError(c, apperror.
			New(apperror.BadRequest).
			Describe(fmt.Sprintf("Could not parse body: %v", err.Error())))
	}

	apperr := h.service.CreateAgreement(agreement)
	if apperr != nil {
		return utils.ResponseError(c, apperr)
	}

	return utils.ResponseMessage(c, fiber.StatusCreated, "Agreement created successfully")
}

// @router  /api/v1/agreement/:agreementId [delete]
// @summary  Delete an agreement
// @description  Delete an agreement by its id
// @tags agreements
// @produce json
// @success 204
// @failure 500 {object} models.ErrorResponses
func (h *handlerImpl) DeleteAgreement(c *fiber.Ctx) error {
	id := c.Params("agreementId")
	err := h.service.DeleteAgreement(id)
	if err != nil {
		return apperror.New(apperror.InternalServerError).Describe("Error deleting agreement")
	}
	return c.SendStatus(fiber.StatusNoContent)
}
