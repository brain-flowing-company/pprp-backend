package agreements

import (
	"fmt"
	"net/http"

	"github.com/brain-flowing-company/pprp-backend/apperror"
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/brain-flowing-company/pprp-backend/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type Handler interface {
	GetAllAgreements(c *fiber.Ctx) error
	GetAgreementById(c *fiber.Ctx) error
	GetMyAgreements(c *fiber.Ctx) error
	CreateAgreement(c *fiber.Ctx) error
	DeleteAgreement(c *fiber.Ctx) error
	UpdateAgreementStatus(c *fiber.Ctx) error
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
// @summary  Get all agreements  *use cookies*
// @description  Get all agreements
// @tags agreements
// @produce json
// @success 200 {object} []models.Agreements
// @failure 500 {object} models.ErrorResponses
func (h *handlerImpl) GetAllAgreements(c *fiber.Ctx) error {
	var apps []models.Agreements

	err := h.service.GetAllAgreements(&apps)
	if err != nil {
		return utils.ResponseError(c, err)
	}
	
	return c.JSON(apps)

}

// @router  /api/v1/agreements/:agreementId [get]
// @summary  Get agreement by id *use cookies*
// @description  Get an agreement by its id
// @tags agreements
// @produce json
// @param agreementId path string true "Agreement ID"
// @success 200 {object} models.AgreementDetails
// @failure 400 {object} models.ErrorResponses "Invalid agreement id"
// @failure 404 {object} models.ErrorResponses "Could not find the specified agreement"
// @failure 500 {object} models.ErrorResponses "Could not get agreement by id"
func (h *handlerImpl) GetAgreementById(c *fiber.Ctx) error {
	agreementId := c.Params("agreementId")
	agreement := &models.AgreementDetails{}

	err := h.service.GetAgreementById(agreement, agreementId)
	if err != nil {
		return utils.ResponseError(c, err)
	}
	
	return c.JSON(agreement)
}

// @router      /api/v1/user/me/agreements [get]
// @summary     Get my agreements *use cookies*
// @description Get all agreements related to the user
// @tags        agreements
// @produce     json
// @success     200	{object} models.MyAgreementResponses
// @failure     500 {object} models.ErrorResponses "Could not get my agreements"
func (h *handlerImpl) GetMyAgreements(c *fiber.Ctx) error {
	userId := c.Locals("session").(models.Sessions).UserId.String()

	var agreements models.MyAgreementResponses
	err := h.service.GetAgreementByUserId(&agreements, userId)
	if err != nil {
		return utils.ResponseError(c, err)
	}

	return c.JSON(agreements)
}

// @router  /api/v1/agreements [post]
// @summary  Create an agreement *use cookies*
// @description  Create an agreement by parsing the body
// @tags agreements
// @produce json
// @param body body models.CreatingAgreements true "Agreement to create"
// @success 201 {object} models.MessageResponses "Agreement created successfully"
// @failure 500 {object} models.ErrorResponses "Could not create agreement"
func (h *handlerImpl) CreateAgreement(c *fiber.Ctx) error {
	agreement := &models.CreatingAgreements{
		OwnerUserId: c.Locals("session").(models.Sessions).UserId,
	}
	
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

// @router  /api/v1/agreements/:agreementId [delete]
// @summary  Delete an agreement by id *use cookies*
// @description  Delete an agreement by its id
// @tags agreements
// @produce json
// @param agreementId path string true "Agreement ID"
// @success 200 {object} models.MessageResponses "Agreement deleted"
// @failure 500 {object} models.ErrorResponses "Could not delete agreement"
func (h *handlerImpl) DeleteAgreement(c *fiber.Ctx) error {
	agreementId := c.Params("agreementId")

	err := h.service.DeleteAgreement(agreementId)
	if err != nil {
		return utils.ResponseError(c, err)
	}
	
	return utils.ResponseMessage(c, fiber.StatusOK, "Agreement deleted")
}

// @router      /api/v1/agreements/:agreementId [patch]
// @summary     Update an agreement status by id *use cookies*
// @description Update an agreement status by id with **status** and **cancelled_message**(optional)
// @tags        agreements
// @produce     json
// @param       agreementId path string true "Agreement ID"
// @param       body body models.UpdatingAgreementStatus true "Agreement status and cancelled message(optional)"
// @success     200	{object} models.MessageResponses "Agreement state updated"
// @failure     400 {object} models.ErrorResponses "Invalid agreement id"
// @failure     500 {object} models.ErrorResponses "Could not update agreement status"
func (h *handlerImpl) UpdateAgreementStatus(c *fiber.Ctx) error {
	updatingAgreement := models.UpdatingAgreementStatus{}
	err := c.BodyParser(&updatingAgreement)
	if err != nil {
		return utils.ResponseError(c, apperror.
			New(apperror.BadRequest).
			Describe(fmt.Sprintf("Could not parse body: %v", err.Error())))
	}

	agreementId := c.Params("agreementId")

	apperr := h.service.UpdateAgreementStatus(&updatingAgreement, agreementId)
	if apperr != nil {
		return utils.ResponseError(c, apperr)
	}

	return utils.ResponseMessage(c, http.StatusOK, "Agreement state updated")
}
