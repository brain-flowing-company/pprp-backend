package models

import (
	"time"

	"github.com/google/uuid"
)

type ChatsResponses struct {
	UnreadCount int64     `json:"unread_count"   example:"9"`
	SenderId    uuid.UUID `json:"sender_id"      example:"123e4567-e89b-12d3-a456-426614174000"`
	CreatedAt   time.Time `json:"created_at"     example:"2024-02-22T03:06:53.313735Z"`
	Content     string    `json:"latest_message" example:"hello, world"`
}
