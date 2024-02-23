package models

import "github.com/brain-flowing-company/pprp-backend/internal/enums"

type Callbacks struct {
	Email string `query:"email"`
	Code  string `query:"code"`
	State string `query:"state"`
}

type CallbackResponses struct {
	Email          string                `json:"email"                  example:"johnd@email.com"`
	RegisteredType enums.RegisteredTypes `json:"registered_type"        example:"EMAIL / GOOGLE"`
	SessionType    enums.SessionType     `json:"session_type,omitempty" example:"REGISTER / LOGIN"`
	Token          string                `json:"-"                      swaggerignore:"true"`
}
