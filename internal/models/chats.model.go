package models

import (
	"github.com/brain-flowing-company/pprp-backend/internal/enums"
	"github.com/google/uuid"
)

type ChatsResponses struct {
	UnreadMessages int64     `json:"unread_messages" example:"9"`
	UserId         uuid.UUID `json:"user_id"         example:"123e4567-e89b-12d3-a456-426614174000"`
	Content        string    `json:"content"         example:"hello, world"`
}

func (e *ChatsResponses) ToOutBound() *OutBoundMessages {
	return &OutBoundMessages{
		Event:   enums.OUTBOUND_CHATS,
		Payload: e,
	}
}
