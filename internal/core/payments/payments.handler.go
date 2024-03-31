package payments

import (
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

func (h *handlerImpl) CreatePayment(c *fiber.Ctx) error {
	payment := models.Payments{
		PaymentId: uuid.New(),
		UserId:    c.Locals("session").(models.Sessions).UserId,
		IsSuccess: false,
	}
	if err := c.BodyParser(&payment); err != nil {
		return utils.ResponseError(c, apperror.New(apperror.InvalidBody).Describe("Invalid payment body"))
	}
	if payment.Price == 0 {
		return utils.ResponseError(c, apperror.New(apperror.InvalidBody).Describe("Price is required"))
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
	err := CheckoutV2(c, payment.Name, payment.Price)
	if err != nil {
		return utils.ResponseError(c, apperror.New(apperror.InternalServerError).Describe("Failed to create payment"))
	}

	return utils.ResponseMessage(c, http.StatusOK, "Payment created successfully")

}

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
