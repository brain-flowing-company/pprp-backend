package payments

import (
	"fmt"

	"github.com/brain-flowing-company/pprp-backend/config"
	"github.com/gofiber/fiber/v2"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/checkout/session"
)

func CheckoutV2(c *fiber.Ctx, name string, price float64, paymentMethod string, cfg *config.Config) (string, error) {
	stripe.Key = cfg.STRIPE_SECRET_KEY
	type PaymentRequest struct {
		Name  string `json:"name"`
		Price int64  `json:"price"`
	}
	var payment PaymentRequest
	payment.Name = name
	payment.Price = int64(price) * 100
	var method string
	if paymentMethod == "PROMPTPAY" {
		method = string(stripe.PaymentMethodTypePromptPay)
	}
	if paymentMethod == "CREDIT_CARD" {
		method = string(stripe.PaymentMethodTypeCard)
	}

	params := &stripe.CheckoutSessionParams{
		Mode: stripe.String(string(stripe.CheckoutSessionModePayment)),
		PaymentMethodTypes: []*string{
			stripe.String(method),
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
		SuccessURL: stripe.String(cfg.FRONTEND_URL + "/success"),
		CancelURL:  stripe.String(cfg.FRONTEND_URL + "/cancel"),
	}

	s, _ := session.New(params)
	fmt.Println(s.URL)

	return s.URL, nil
}
