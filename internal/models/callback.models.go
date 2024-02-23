package models

import "github.com/brain-flowing-company/pprp-backend/internal/enums"

type Callbacks struct {
	Email string `query:"email"`
	Code  string `query:"code"`
	State string `query:"state"`
}

type CallbackResponses struct {
	Email          string                `json:"email"`
	RegisteredType enums.RegisteredTypes `json:"registered_type"`
	SessionType    enums.SessionType     `json:"session_type,omitempty"`
}
