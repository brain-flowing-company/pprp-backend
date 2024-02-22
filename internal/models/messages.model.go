package models

import (
	"time"

	"github.com/google/uuid"
)

type Messages struct {
	MessageId  uuid.UUID `json:"message_id"`
	SenderId   uuid.UUID `json:"sender_id"`
	ReceiverId uuid.UUID `json:"receiver_id"`
	Content    string    `json:"content"`
	Read       bool      `json:"read"`
	CreatedAt  time.Time `json:"created_at"`
}
