package models

import (
	"github.com/brain-flowing-company/pprp-backend/internal/enums"
	"github.com/google/uuid"
)

type CreditCards struct {
	CreditCardId   uuid.UUID       `json:"credit_card_id"               gorm:"default:uuid_generate_v4()"`
	CardNickname   string          `json:"card_nickname"                example:"My Card"`
	CardholderName string          `json:"cardholder_name"              example:"JOHN DOE"`
	CardNumber     string          `json:"card_number"                  example:"1234567890123456"`
	ExpireMonth    string          `json:"expire_month"             example:"12"`
	ExpireYear     string          `json:"expire_year"              example:"2023"`
	CVV            string          `json:"cvv"                          example:"123"`
	CardColor      enums.CardColor `json:"card_color"                   example:"LIGHT BLUE"`
	CommonModels
}
