package users

import (
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/gofiber/fiber/v2"
)

type Handler interface {
	CreateUser(c *fiber.Ctx) error
	// GetAllUsers(c *fiber.Ctx) error
}

type handlerImpl struct {
	service Service
}

func NewHandler(service Service) Handler {
	return &handlerImpl{
		service,
	}
}

// func (h *handlerImpl) GetAllUsers(c *fiber.Ctx) error {
// 	return c.Status(200).JSON(fiber.Map{
// 		"message": "All users have been retrieved",
// 	})
// }

func (h *handlerImpl) CreateUser(c *fiber.Ctx) error {

	// type body struct {
	// 	FirstName                 string `json:"first_name"`
	// 	LastName                  string `json:"last_name"`
	// 	Email                     string `json:"email"`
	// 	PhoneNumber               string `json:"phone_number"`
	// 	ProfileImage              string `json:"profile_image"`
	// 	CreditCardholderName      string `json:"credit_cardholder_name"`
	// 	CreditCardNumber          string `json:"credit_card_number"`
	// 	CreditCardExpirationMonth string `json:"credit_card_expiration_month"`
	// 	CreditCardExpirationYear  string `json:"credit_card_expiration_year"`
	// 	CreditCardCVV             string `json:"credit_card_cvv"`
	// 	BankName                  string `json:"bank_name"`
	// 	BankAccountNumber         string `json:"bank_account_number"`
	// 	IsVerified                bool   `json:"is_verified"`
	// }

	user := models.Users{}

	bodyErr := c.BodyParser(&user)

	if bodyErr != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"message": "Invalid body",
		})
	}

	err := h.service.CreateUser(&user)
	if err != nil {
		return c.Status(err.Code).JSON(fiber.Map{
			"message": "Failed to create user",
		})
	}

	return c.JSON(user) // TODO: don't return user
}
