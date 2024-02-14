package utils

import (
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/gofiber/fiber/v2"
)

func ParseFormToUser(c *fiber.Ctx) models.Users {
	user := models.Users{
		RegisteredType:            ParseRegisteredType(c.FormValue("registered_type")),
		Email:                     c.FormValue("email"),
		Password:                  c.FormValue("password"),
		FirstName:                 c.FormValue("first_name"),
		LastName:                  c.FormValue("last_name"),
		PhoneNumber:               c.FormValue("phone_number"),
		ProfileImageUrl:           c.FormValue("profile_image_url"),
		CreditCardCardholderName:  c.FormValue("credit_cardholder_name"),
		CreditCardNumber:          c.FormValue("credit_card_number"),
		CreditCardExpirationMonth: c.FormValue("credit_card_expiration_month"),
		CreditCardExpirationYear:  c.FormValue("credit_card_expiration_year"),
		CreditCardCVV:             c.FormValue("credit_card_cvv"),
		BankName:                  ParseBankName(c.FormValue("bank_name")),
		BankAccountNumber:         c.FormValue("bank_account_number"),
		CitizenId:                 c.FormValue("citizen_id"),
		CitizenCardImageUrl:       c.FormValue("citizen_card_image_url"),
		IsVerified:                false,
	}

	if c.FormValue("is_verified") == "true" {
		user.IsVerified = true
	}

	return user
}
