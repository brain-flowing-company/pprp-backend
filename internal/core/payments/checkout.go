package payments

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/checkout/session"
)

func Checkout(c *fiber.Ctx) error {
	stripe.Key = "sk_test_51OmWT2BayMsgzLXzrhGhYbxvTA6QtQvBwVhU2GYCNX6GFhGgVovQSapIhDKftcwpLOvqyrruOj0Tw7HfAcfJT5sd00YBwEU9aw"
	type PaymentRequest struct {
		Name  string `json:"name"`
		Price int64  `json:"price"`
	}
	var payment PaymentRequest
	payment.Name = "T-shirt"
	payment.Price = 50000

	params := &stripe.CheckoutSessionParams{
		Mode: stripe.String(string(stripe.CheckoutSessionModePayment)),
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
		SuccessURL: stripe.String("http://localhost:4242/success"),
		CancelURL:  stripe.String("http://localhost:4242/cancel"),
	}

	s, _ := session.New(params)
	fmt.Println(s.URL)

	return c.Redirect(s.URL, fiber.StatusSeeOther)
}
