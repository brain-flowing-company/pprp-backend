package models

import (
	"time"

	"github.com/google/uuid"
)

type GoogleOAuthStates struct {
	Code      uuid.UUID
	ExpiredAt time.Time
}

func (g GoogleOAuthStates) TableName() string {
	return "google_oauth_states"
}
