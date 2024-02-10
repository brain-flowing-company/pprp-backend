package models

type Session struct {
	Email          string         `json:"email"`
	RegisteredType RegisteredType `json:"registered_type"`
}
