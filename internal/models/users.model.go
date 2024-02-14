package models

import (
	"time"

	"github.com/brain-flowing-company/pprp-backend/internal/consts"
	"github.com/google/uuid"
)

type Users struct {
	UserId                    uuid.UUID             `gorm:"default:uuid_generate_v4()"`
	RegisteredType            consts.RegisteredType `example:"EMAIL"`
	Email                     string                `json:"email"        form:"email" gorm:"unique"  example:"email@email.com"`
	Password                  string                `json:"password"     form:"password" gorm:"default:null"  example:"password1234"`
	FirstName                 string                `json:"first_name"   form:"first_name" example:"John"`
	LastName                  string                `json:"last_name"    form:"last_name" example:"Doe"`
	PhoneNumber               string                `json:"phone_number" form:"phone_number" gorm:"unique"  example:"0812345678"`
	ProfileImageUrl           string                `json:"profile_image_url" form:"profile_image_url" gorm:"default:null"  example:"https://image_url.com/abcd"`
	CreditCardCardholderName  string                `json:"credit_cardholder_name" form:"credit_cardholder_name" gorm:"default:null"  example:"JOHN DOE"`
	CreditCardNumber          string                `json:"credit_card_number" form:"credit_card_number" gorm:"default:null"  example:"1234567890123456"`
	CreditCardExpirationMonth string                `json:"credit_card_expiration_month" form:"credit_card_expiration_month" gorm:"default:null"  example:"12"`
	CreditCardExpirationYear  string                `json:"credit_card_expiration_year" form:"credit_card_expiration_year" gorm:"default:null"  example:"2023"`
	CreditCardCVV             string                `json:"credit_card_cvv" form:"credit_card_cvv" gorm:"default:null"  example:"123"`
	BankName                  consts.BankName       `json:"bank_name" form:"bank_name" gorm:"default:null"  example:"KBANK"`
	BankAccountNumber         string                `json:"bank_account_number" form:"bank_account_number" gorm:"default:null"  example:"1234567890"`
	CitizenId                 string                `json:"citizen_id" form:"citizen_id" gorm:"default:null; unique"  example:"1234567890123"`
	CitizenCardImageUrl       string                `json:"citizen_card_image_url" form:"citizen_card_image_url" gorm:"default:null"  example:"https://image_url.com/abcd"`
	IsVerified                bool                  `json:"is_verified" form:"is_verified" gorm:"default:null"  example:"false"`
	CreatedAt                 *time.Time            `gorm:"autoCreateTime"`
	UpdatedAt                 *time.Time            `gorm:"autoUpdateTime"`
	DeletedAt                 *time.Time            `gorm:"default:null"`
}

func (u Users) TableName() string {
	return "users"
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
