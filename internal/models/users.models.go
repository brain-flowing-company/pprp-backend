package models

import (
	"github.com/brain-flowing-company/pprp-backend/internal/enums"
	"github.com/google/uuid"
)

type Users struct {
	UserId              uuid.UUID             `json:"user_id"                      gorm:"default:uuid_generate_v4()"`
	RegisteredType      enums.RegisteredTypes `json:"registered_type"              example:"EMAIL"`
	Email               string                `json:"email"                        form:"email"                        gorm:"unique" example:"email@email.com"`
	Password            string                `json:"password"                     form:"password"                     gorm:"default:null" example:"password1234"`
	FirstName           string                `json:"first_name"                   form:"first_name"                   example:"John"`
	LastName            string                `json:"last_name"                    form:"last_name"                    example:"Doe"`
	PhoneNumber         string                `json:"phone_number"                 form:"phone_number"                 gorm:"unique" example:"0812345678"`
	ProfileImageUrl     string                `json:"profile_image_url"            form:"profile_image_url"            gorm:"default:null" example:"https://image_url.com/abcd"`
	BankName            enums.BankNames       `json:"bank_name"                    form:"bank_name"                    gorm:"default:null" example:"KBANK"`
	BankAccountNumber   string                `json:"bank_account_number"          form:"bank_account_number"          gorm:"default:null" example:"1234567890"`
	CitizenId           string                `json:"citizen_id"                   form:"citizen_id"                   gorm:"default:null; unique" example:"1234567890123"`
	CitizenCardImageUrl string                `json:"citizen_card_image_url"       form:"citizen_card_image_url"       gorm:"default:null" example:"https://image_url.com/abcd"`
	IsVerified          bool                  `json:"is_verified"                  form:"is_verified"                  gorm:"default:null" example:"false"`
	CommonModels
}

type RegisteringUsers struct {
	UserId          uuid.UUID             `form:"-" swaggerignore:"true"`
	RegisteredType  enums.RegisteredTypes `form:"-" swaggerignore:"true"`
	Email           string                `form:"email"        example:"email@email.com"`
	Password        string                `form:"password"     example:"password1234"`
	FirstName       string                `form:"first_name"   example:"John"`
	LastName        string                `form:"last_name"    example:"Doe"`
	PhoneNumber     string                `form:"phone_number" example:"0812345678"`
	ProfileImageUrl string                `form:"-" swaggerignore:"true"`
	CommonModels    `swaggerignore:"true"`
}

func (r RegisteringUsers) TableName() string {
	return "users"
}

type UpdatingUserPersonalInfo struct {
	UserId          uuid.UUID `form:"-"            swaggerignore:"true"`
	FirstName       string    `form:"first_name"   example:"John"`
	LastName        string    `form:"last_name"    example:"Doe"`
	PhoneNumber     string    `form:"phone_number" example:"0812345678"`
	ProfileImageUrl string    `form:"-"            swaggerignore:"true"`
	CommonModels    `swaggerignore:"true"`
}

func (r UpdatingUserPersonalInfo) TableName() string {
	return "users"
}

func (u Users) TableName() string {
	return "users"
}

type LoginRequests struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
