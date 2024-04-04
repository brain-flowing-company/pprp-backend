package payments

import (
	"github.com/brain-flowing-company/pprp-backend/apperror"
	"github.com/brain-flowing-company/pprp-backend/config"
	"github.com/brain-flowing-company/pprp-backend/internal/enums"
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/brain-flowing-company/pprp-backend/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Handler interface {
	CreatePayment(c *fiber.Ctx) error
	GetPaymentByUserId(c *fiber.Ctx) error
	GetHistoryPaymentByUserId(c *fiber.Ctx) error
}

type handlerImpl struct {
	service Service
	cfg     *config.Config
}

func NewHandler(cfg *config.Config, service Service) Handler {
	return &handlerImpl{
		service,
		cfg,
	}
}

// @router /api/v1/checkout [get]
// @summary     Create payment
// @description Create payment
// @tags        payments
// @produce     json
// @param       agreement_id query string true "Agreement ID"
// @param       name query string true "Name"
// @param       price query string true "Price"
// @param       payment_method query string true "Payment Method"
// @success     200	{object}	models.Payments
// @failure     400 {object}	models.ErrorResponses
// @failure     401 {object}	models.ErrorResponses
// @failure     500 {object}	models.ErrorResponses
func (h *handlerImpl) CreatePayment(c *fiber.Ctx) error {
	session, ok := c.Locals("session").(models.Sessions)
	if !ok {
		return utils.ResponseError(c, apperror.New(apperror.Unauthorized).Describe("Unauthorized"))
	}
	payment := models.Payments{
		PaymentId: uuid.New(),
		UserId:    session.UserId,
		IsSuccess: false,
	}
	err := c.QueryParser(&payment)
	if err != nil {
		return utils.ResponseError(c, apperror.New(apperror.BadRequest).Describe("Failed to parse query"))
	}
	// Parse agreement_id
	agreementID := c.Query("agreement_id")
	if agreementID != "" {
		parsedAgreementID, err := uuid.Parse(agreementID)
		if err != nil {
			return utils.ResponseError(c, apperror.New(apperror.BadRequest).Describe("Invalid agreement_id format"))
		}
		payment.AgreementId = parsedAgreementID
	}
	// Parse payment_method
	paymentMethodStr := c.Query("payment_method")
	if paymentMethodStr != "" {
		// Use the enums package to parse the payment method
		paymentMethod := enums.PaymentMethods(paymentMethodStr)
		switch paymentMethod {
		case enums.CREDIT_CARD:
			payment.PaymentMethod = enums.CREDIT_CARD
		case enums.PROMPTPAY:
			payment.PaymentMethod = enums.PROMPTPAY
		default:
			return utils.ResponseError(c, apperror.New(apperror.BadRequest).Describe("Invalid payment_method"))
		}
	}
	// Check if the required fields are empty
	if payment.Price <= 0 {
		return utils.ResponseError(c, apperror.New(apperror.InvalidBody).Describe("Price is required and must be greater than 0"))
	}
	if payment.Name == "" {
		return utils.ResponseError(c, apperror.New(apperror.InvalidBody).Describe("Name is required"))
	}
	if payment.AgreementId == uuid.Nil {
		return utils.ResponseError(c, apperror.New(apperror.InvalidBody).Describe("Agreement id is required"))
	}
	if payment.PaymentMethod == "" {
		return utils.ResponseError(c, apperror.New(apperror.InvalidBody).Describe("Payment method is required"))
	}

	if err := h.service.CreatePayment(&payment); err != nil {
		return utils.ResponseError(c, err)
	}
	surl, err2 := CheckoutV2(c, payment.Name, payment.Price, string(payment.PaymentMethod), h.cfg)
	if err2 != nil {
		return utils.ResponseError(c, apperror.New(apperror.InternalServerError).Describe("Failed to create payment"))
	}
	return c.JSON(fiber.Map{
		"payment_id": payment.PaymentId,
		"success":    true,
		"message":    "Payment created successfully",
		"url":        surl,
	})

}

// @router /api/v1/payments [get]
// @summary     Get payment by user id
// @description Get payment by user id
// @tags        payments
// @produce     json
// @success     200	{object}	models.MyPaymentsResponse
// @failure     400 {object}	models.ErrorResponses
// @failure     401 {object}	models.ErrorResponses
// @failure     500 {object}	models.ErrorResponses
func (h *handlerImpl) GetPaymentByUserId(c *fiber.Ctx) error {
	userId := c.Locals("session").(models.Sessions).UserId
	payments := models.MyPaymentsResponse{}
	if err := h.service.GetPaymentByUserId(&payments, userId); err != nil {
		return utils.ResponseError(c, err)
	}
	err := c.JSON(payments)
	if err != nil {
		return utils.ResponseError(c, apperror.New(apperror.InternalServerError).Describe("Failed to get payment by user id"))
	}
	return nil
}

// @router /api/v1/payments/history [get]
// @summary     Get history payment by user id
// @description Get history payment by user id
// @tags        payments
// @produce     json
// @success     200	{object}	models.HistoryResponse
// @failure     400 {object}	models.ErrorResponses
// @failure     401 {object}	models.ErrorResponses
// @failure     500 {object}	models.ErrorResponses
func (h *handlerImpl) GetHistoryPaymentByUserId(c *fiber.Ctx) error {
	userId := c.Locals("session").(models.Sessions).UserId
	var HistoryResponse []models.HistoryResponse
	if err := h.service.GetHistoryPaymentByUserId(&HistoryResponse, userId); err != nil {
		return utils.ResponseError(c, err)
	}
	err := c.JSON(HistoryResponse)
	if err != nil {
		return utils.ResponseError(c, apperror.New(apperror.InternalServerError).Describe("Failed to get payment by user id"))
	}
	return nil
}
