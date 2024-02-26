package models

import (
	"time"

	"github.com/brain-flowing-company/pprp-backend/internal/enums"
	"github.com/google/uuid"
)

type OutBoundPayload interface {
	ToOutBound(tag string) *OutBoundMessages
}

type Messages struct {
	MessageId  uuid.UUID  `json:"message_id"    example:"27b79b15-a56f-464a-90f7-bab515ba4c02"`
	SenderId   uuid.UUID  `json:"sender_id"     example:"27b79b15-a56f-464a-90f7-bab515ba4c02"`
	ReceiverId *uuid.UUID `json:"receiver_id"   example:"27b79b15-a56f-464a-90f7-bab515ba4c02"`
	Content    string     `json:"content"       example:"hello, world"`
	ReadAt     *time.Time `json:"read_at"       example:"2024-02-22T03:06:53.313735Z"`
	SentAt     time.Time  `json:"sent_at"       example:"2024-02-22T03:06:53.313735Z"`
}

func (e *Messages) ToOutBound(tag string) *OutBoundMessages {
	return &OutBoundMessages{
		Event:   enums.OUTBOUND_MSG,
		Tag:     tag,
		Payload: e,
	}
}

type ReadEvents struct {
	SenderId   uuid.UUID `json:"sender_id"`
	ReceiverId uuid.UUID `json:"receiver_id"`
	ReadAt     time.Time `json:"read_at"`
}

func (e *ReadEvents) ToOutBound(tag string) *OutBoundMessages {
	return &OutBoundMessages{
		Event:   enums.OUTBOUND_READ,
		Tag:     tag,
		Payload: e,
	}
}

type OutBoundMessages struct {
	Event   enums.MessageOutboundEvents `json:"event"`
	Tag     string                      `json:"tag,omitempty"`
	Payload interface{}                 `json:"payload"`
}

type InBoundMessages struct {
	Event   enums.MessageInboundEvents `json:"event"`
	Content string                     `json:"content"`
	SentAt  time.Time                  `json:"sent_at"`
	Tag     string                     `json:"tag"`
}
