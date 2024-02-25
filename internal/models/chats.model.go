package models

import (
	"time"

	"github.com/google/uuid"
)

type ChatsResponses struct {
	UnreadMessages int64     `json:"unread_messages" example:"9"`
	UserId         uuid.UUID `json:"user_id"         example:"123e4567-e89b-12d3-a456-426614174000"`
	SentAt         time.Time `json:"sent_at"         example:"2024-02-22T03:06:53.313735Z"`
	Content        string    `json:"content"         example:"hello, world"`
}
