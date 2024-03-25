package payments

import (
	"github.com/brain-flowing-company/pprp-backend/apperror"
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/checkout/session"
)

type Handler interface {
	CreatePayment(c *fiber.Ctx) error
}

// @router  /api/v1/payments [post]
// @summary  Create payment
// @description  Create a payment
// @tags payments
// @accept json
// @produce json
// @param payment body models.CreatingPayment true "Payment object"
// @success 201 {object} models.Payment
// @failure 400 {object} models.ErrorResponses "Invalid payment object"
// @failure 500 {object} models.ErrorResponses

func CreatePayment(c *fiber.Ctx) error {
	stripe.Key = "sk_test_51OmWT2BayMsgzLXzrhGhYbxvTA6QtQvBwVhU2GYCNX6GFhGgVovQSapIhDKftcwpLOvqyrruOj0Tw7HfAcfJT5sd00YBwEU9aw"
	var payment models.Payments
	err := c.BodyParser(&payment)
	if err != nil {
		return apperror.New(apperror.InvalidBody).Describe("Invalid payment object")
	}

	params := &stripe.CheckoutSessionParams{
		Mode: stripe.String(string(stripe.CheckoutSessionModePayment)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			&stripe.CheckoutSessionLineItemParams{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String("thb"),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String("T-shirt"),
					},
					UnitAmount: stripe.Int64(4000),
				},
				Quantity: stripe.Int64(1),
			},
		},
		SuccessURL: stripe.String("http://localhost:4242/success"),
		CancelURL:  stripe.String("http://localhost:4242/cancel"),
	}

	s, _ := session.New(params)

	if err != nil {
		return err
	}

	return c.Redirect(s.URL, fiber.StatusSeeOther)
}

// RegisterPaymentRoutes registers payment routes with the provided Fiber app instance
func RegisterPaymentRoutes(app *fiber.App, paymentHandler Handler) {
	api := app.Group("/api/v1")

	// Define routes
	api.Post("/payments", paymentHandler.CreatePayment)
}
