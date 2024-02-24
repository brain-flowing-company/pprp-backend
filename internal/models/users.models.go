package models

import (
	"time"

	"github.com/brain-flowing-company/pprp-backend/internal/enums"
	"github.com/google/uuid"
)

type Users struct {
	UserId          uuid.UUID             `json:"user_id"                      gorm:"default:uuid_generate_v4()"`
	RegisteredType  enums.RegisteredTypes `json:"registered_type"              example:"EMAIL"`
	Email           string                `json:"email"                        form:"email"                        gorm:"unique" example:"email@email.com"`
	Password        string                `json:"password"                     form:"password"                     gorm:"default:null" example:"password1234"`
	FirstName       string                `json:"first_name"                   form:"first_name"                   example:"John"`
	LastName        string                `json:"last_name"                    form:"last_name"                    example:"Doe"`
	PhoneNumber     string                `json:"phone_number"                 form:"phone_number"                 gorm:"unique" example:"0812345678"`
	ProfileImageUrl string                `json:"profile_image_url"            form:"profile_image_url"            gorm:"default:null" example:"https://image_url.com/abcd"`
	IsVerified      bool                  `json:"is_verified"                  gorm:"default:null" example:"false"`
	CommonModels
}

func (u Users) TableName() string {
	return "users"
}

type UserFinancialInformations struct {
	UserId            uuid.UUID       `form:"-"          swaggerignore:"true"    json:"user_id"`
	CreditCards       []CreditCards   `gorm:"foreignKey:UserId; references:UserId;"          swaggerignore:"true"    json:"credit_cards"`
	BankName          enums.BankNames `json:"bank_name"                    form:"bank_name"                    gorm:"default:null" example:"KBANK"`
	BankAccountNumber string          `json:"bank_account_number"          form:"bank_account_number"          gorm:"default:null" example:"1234567890"`
	CommonModels
}

func (uf UserFinancialInformations) TableName() string {
	return "user_financial_informations"
}

type CreditCards struct {
	UserId         uuid.UUID        `json:"-" swaggerignore:"true"`
	TagNumber      int64            `json:"tag_number" example:"1234"`
	CardNickname   string           `json:"card_nickname" example:"John's Card"`
	CardholderName string           `json:"cardholder_name" example:"John Doe"`
	CardNumber     string           `json:"card_number" example:"1234567890123456"`
	ExpireMonth    string           `json:"expire_month" example:"12"`
	ExpireYear     string           `json:"expire_year" example:"2023"`
	CVV            string           `json:"cvv" example:"123"`
	CardColor      enums.CardColors `json:"card_color" example:"BLUE"`
}

func (cc CreditCards) TableName() string {
	return "credit_cards"
}

type UserVerifications struct {
	UserId              uuid.UUID `form:"-"          swaggerignore:"true"    json:"user_id"`
	CitizenId           string    `form:"citizen_id" example:"1100111111111" json:"citizen_id"`
	CitizenCardImageUrl string    `form:"-"          swaggerignore:"true"    json:"citizen_card_image_url"`
	VerifiedAt          time.Time `form:"-"          swaggerignore:"true"    json:"verified_at"`
}

func (uv UserVerifications) TableName() string {
	return "user_verifications"
}

type RegisteringUsers struct {
	UserId          uuid.UUID             `form:"-" swaggerignore:"true"`
	RegisteredType  enums.RegisteredTypes `form:"registered_type" exmaple:"EMAIL / GOOGLE"`
	Email           string                `form:"email"           example:"email@email.com"`
	Password        string                `form:"password"        example:"password1234"`
	FirstName       string                `form:"first_name"      example:"John"`
	LastName        string                `form:"last_name"       example:"Doe"`
	PhoneNumber     string                `form:"phone_number"    example:"0812345678"`
	ProfileImageUrl string                `form:"-" swaggerignore:"true"`
	CommonModels    `swaggerignore:"true"`
}

func (r RegisteringUsers) TableName() string {
	return "users"
}

type UpdatingUserPersonalInfos struct {
	UserId          uuid.UUID `form:"-"            swaggerignore:"true"`
	FirstName       string    `form:"first_name"   example:"John"`
	LastName        string    `form:"last_name"    example:"Doe"`
	PhoneNumber     string    `form:"phone_number" example:"0812345678"`
	ProfileImageUrl string    `form:"-"            swaggerignore:"true"`
	CommonModels    `swaggerignore:"true"`
}

func (r UpdatingUserPersonalInfos) TableName() string {
	return "users"
}

type LoginRequests struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
