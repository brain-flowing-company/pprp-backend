package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Users struct {
	UserId                    uuid.UUID      `gorm:"default:uuid_generate_v4()"`
	Email                     string         `gorm:"unique" json:"email" example:"email@email.com"`
	Password                  string         `gorm:"default:null" json:"password" example:"password1234"`
	FirstName                 string         `json:"first_name" example:"John"`
	LastName                  string         `json:"last_name" example:"Doe"`
	PhoneNumber               string         `gorm:"unique" json:"phone_number" example:"0812345678"`
	ProfileImageUrl           string         `gorm:"default:null" json:"profile_image_url" example:"https://image_url.com/abcd"`
	CreditCardCardholderName  string         `gorm:"default:null" json:"credit_cardholder_name" example:"JOHN DOE"`
	CreditCardNumber          string         `gorm:"default:null" json:"credit_card_number" example:"1234567890123456"`
	CreditCardExpirationMonth string         `gorm:"default:null" json:"credit_card_expiration_month" example:"12"`
	CreditCardExpirationYear  string         `gorm:"default:null" json:"credit_card_expiration_year" example:"2023"`
	CreditCardCVV             string         `gorm:"default:null" json:"credit_card_cvv" example:"123"`
	BankName                  BankName       `gorm:"default:null" json:"bank_name" example:"KBANK"`
	BankAccountNumber         string         `gorm:"default:null" json:"bank_account_number" example:"1234567890"`
	IsVerified                bool           `gorm:"default:null" json:"is_verified" example:"false"`
	CreatedAt                 time.Time      `gorm:autoCreateTime`
	UpdatedAt                 time.Time      `gorm:autoUpdateTime`
	DeletedAt                 gorm.DeletedAt `gorm:"index"`
}

type BankName string

const (
	KBANK BankName = "KASIKORN BANK"
	BBL   BankName = "BANGKOK BANK"
	KTB   BankName = "KRUNG THAI BANK"
	BAY   BankName = "BANK OF AYUDHYA"
	CIMB  BankName = "CIMB THAI BANK"
	TTB   BankName = "TMBTHANACHART BANK"
	SCB   BankName = "SIAM COMMERCIAL BANK"
	GSB   BankName = "GOVERNMENT SAVINGS BANK"
)

func (u Users) TableName() string {
	return "users"
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
