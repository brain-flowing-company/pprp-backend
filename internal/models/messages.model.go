package models

import (
	"time"

	"github.com/google/uuid"
)

type Messages struct {
	MessageId  uuid.UUID `json:"message_id"  example:"27b79b15-a56f-464a-90f7-bab515ba4c02"`
	SenderId   uuid.UUID `json:"sender_id"   example:"27b79b15-a56f-464a-90f7-bab515ba4c02"`
	ReceiverId uuid.UUID `json:"receiver_id" example:"27b79b15-a56f-464a-90f7-bab515ba4c02"`
	Content    string    `json:"content"     example:"hello, world"`
	Read       bool      `json:"read"        example:"false"`
	CreatedAt  time.Time `json:"created_at"  example:"2024-02-22T03:06:53.313735Z"`
}

type RawMessages struct {
	Content    string    `json:"content"`
	ReceiverId uuid.UUID `json:"receiver_id"`
	CreatedAt  time.Time `json:"created_at"`
	Etag       string    `json:"etag"`
}
