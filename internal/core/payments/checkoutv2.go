package payments

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/checkout/session"
)

func CheckoutV2(c *fiber.Ctx, name string, price float64) error {
	stripe.Key = os.Getenv("PAYMENT_SECRET_KEY")
	type PaymentRequest struct {
		Name  string `json:"name"`
		Price int64  `json:"price"`
	}
	var payment PaymentRequest
	payment.Name = name
	payment.Price = int64(price) * 100

	params := &stripe.CheckoutSessionParams{
		Mode: stripe.String(string(stripe.CheckoutSessionModePayment)),
		PaymentMethodTypes: []*string{
			stripe.String(string(stripe.PaymentMethodTypeCard)),
			stripe.String(string(stripe.PaymentMethodTypePromptPay)),
		},
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			&stripe.CheckoutSessionLineItemParams{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String("thb"),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String(payment.Name),
					},
					UnitAmount: stripe.Int64(payment.Price),
				},
				Quantity: stripe.Int64(1),
			},
		},
		SuccessURL: stripe.String(os.Getenv("SUCCESS_URL")),
		CancelURL:  stripe.String(os.Getenv("CANCEL_URL")),
	}

	s, _ := session.New(params)
	fmt.Println(s.URL)

	return c.Redirect(s.URL, fiber.StatusSeeOther)
}
