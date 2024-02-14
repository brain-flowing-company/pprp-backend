package models

import "github.com/brain-flowing-company/pprp-backend/internal/consts"

type Session struct {
	Email          string                `json:"email,omitempty"           example:"admim@email.com"`
	RegisteredType consts.RegisteredType `json:"registered_type,omitempty" example:"EMAIL / GOOGLE"`
	SessionType    SessionType           `json:"session_type,omitempty"    example:"LOGIN / REGISTER"`
}

type SessionType string

const (
	SessionRegister SessionType = "REGISTER"
	SessionLogin    SessionType = "LOGIN"
)
