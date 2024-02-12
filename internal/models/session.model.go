package models

type Session struct {
	Email          string         `json:"email,omitempty"           example:"admim@email.com"`
	RegisteredType RegisteredType `json:"registered_type,omitempty" example:"EMAIL / GOOGLE"`
	SessionType    SessionType    `json:"session_type,omitempty"    example:"LOGIN / REGISTER"`
}

type SessionType string

const (
	SessionRegister SessionType = "REGISTER"
	SessionLogin    SessionType = "LOGIN"
)
