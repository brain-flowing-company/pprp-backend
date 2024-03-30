package payments

import (
	"fmt"
	"net/http"

	"github.com/brain-flowing-company/pprp-backend/apperror"
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/brain-flowing-company/pprp-backend/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Handler interface {
	CreatePayment(c *fiber.Ctx) error
	GetPaymentByUserId(c *fiber.Ctx) error
}

type handlerImpl struct {
	service Service
}

func NewHandler(service Service) Handler {
	return &handlerImpl{
		service,
	}
}

// @router      /api/v1/payments [post]
// @summary     Create Payment
// @description Create a new payment
// @tags        payments
// @produce     json
// @param       body body models.Payments true "Payment object"
// @success     200 {object} models.MessageResponses "Payment created successfully"
// @failure     400 {object} models.ErrorResponses "Invalid payment body"
// @failure     500 {object} models.ErrorResponses "Failed to create payment"
func (h *handlerImpl) CreatePayment(c *fiber.Ctx) error {
	payment := models.Payments{
		PaymentId: uuid.New(),
		UserId:    c.Locals("session").(models.Sessions).UserId,
	}

	if err := c.BodyParser(&payment); err != nil {
		return utils.ResponseError(c, apperror.New(apperror.InvalidBody).Describe("Invalid payment body"))
	}
	fmt.Println("paymentkuay = ", payment.UserId)
	if err := h.service.CreatePayment(&payment); err != nil {
		return utils.ResponseError(c, err)
	}
	err := CheckoutV2(c, payment.Name, payment.Price)
	if err != nil {
		return utils.ResponseError(c, apperror.New(apperror.InternalServerError).Describe("Failed to create payment"))
	}

	return utils.ResponseMessage(c, http.StatusOK, "Payment created successfully")

}

// @router      /api/v1/payments [get]
// @summary     Get Payments by User ID
// @description Get payments associated with the current user
// @tags        payments
// @produce     json
// @success     200 {object} models.MyPaymentsResponse "Payments retrieved successfully"
// @failure     400 {object} models.ErrorResponses "Invalid user session"
// @failure     500 {object} models.ErrorResponses "Failed to get payments by user ID"
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
