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
	// Checkout(c)
	payment := models.Payments{
		PaymentId: uuid.New(),
		UserId:    uuid.MustParse("123e4567-e89b-12d3-a456-426614174001"),
	}
	if err := c.BodyParser(&payment); err != nil {
		return utils.ResponseError(c, apperror.New(apperror.InvalidBody).Describe("Invalid payment body"))
	}

	// userId := c.Locals("session").(models.Sessions).UserId
	// fmt.Println(userId)
	// payment.UserId = userId
	fmt.Println("payment = ", payment)
	if err := h.service.CreatePayment(&payment); err != nil {
		return utils.ResponseError(c, err)
	}

	return utils.ResponseMessage(c, http.StatusOK, "Payment created successfully")

}
