package models

import (
	"github.com/google/uuid"
)

type Sessions struct {
	UserId uuid.UUID `json:"user_id,omitempty" example:"123e4567-e89b-12d3-a456-426614174000"`
	Email  string    `json:"email,omitempty"   example:"admim@email.com"`
}
