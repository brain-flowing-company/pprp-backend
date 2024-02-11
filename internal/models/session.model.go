package models

type Session struct {
	Email          string         `json:"email,omitempty"`
	RegisteredType RegisteredType `json:"registered_type"`
}
