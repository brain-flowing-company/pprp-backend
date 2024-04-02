package models

import (
	"time"

	"github.com/brain-flowing-company/pprp-backend/internal/enums"
	"github.com/google/uuid"
)

type OutBoundPayload interface {
	ToOutBound() *OutBoundMessages
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

type Messages struct {
	MessageId   uuid.UUID           `json:"message_id"    example:"27b79b15-a56f-464a-90f7-bab515ba4c02"`
	ChatId      uuid.UUID           `json:"chat_id"       example:"27b79b15-a56f-464a-90f7-bab515ba4c02" gorm:"-"`
	SenderId    uuid.UUID           `json:"-"             example:"27b79b15-a56f-464a-90f7-bab515ba4c02"`
	ReceiverId  uuid.UUID           `json:"-"             example:"27b79b15-a56f-464a-90f7-bab515ba4c02"`
	Content     string              `json:"content"       example:"hello, world"`
	ReadAt      *time.Time          `json:"read_at"       example:"2024-02-22T03:06:53.313735Z"`
	SentAt      time.Time           `json:"sent_at"       example:"2024-02-22T03:06:53.313735Z"`
	Author      bool                `json:"author"        example:"true"                                 gorm:"-"`
	Tag         string              `json:"-"             gorm:"-"`
	Attatchment MessageAttatchments `json:"attatchment"   gorm:"embedded"`
}

type MessageAttatchments struct {
	MessageId     uuid.UUID  `json:"-"`
	PropertyId    *uuid.UUID `json:"property_id,omitempty"    gorm:"->"`
	AppointmentId *uuid.UUID `json:"appointment_id,omitempty" gorm:"->"`
	AgreementId   *uuid.UUID `json:"agreement_id,omitempty"   gorm:"->"`
}

func (e *Messages) ToOutBound() *OutBoundMessages {
	tmp := *e
	return &OutBoundMessages{
		Event:   enums.OUTBOUND_MSG,
		Tag:     e.Tag,
		Payload: tmp,
	}
}

type ReadEvents struct {
	ChatId uuid.UUID `json:"chat_id"`
	ReadAt time.Time `json:"read_at"`
}

func (e *ReadEvents) ToOutBound() *OutBoundMessages {
	tmp := *e
	return &OutBoundMessages{
		Event:   enums.OUTBOUND_READ,
		Payload: tmp,
	}
}

type ChatPreviews struct {
	UserId          uuid.UUID `json:"user_id"           example:"123e4567-e89b-12d3-a456-426614174000"`
	ProfileImageUrl string    `json:"profile_image_url" example:"www.image.com/profile"`
	FirstName       string    `json:"first_name"        example:"John"`
	LastName        string    `json:"last_name"         example:"Doe"`
	UnreadMessages  int64     `json:"unread_messages"   example:"9"`
	Content         string    `json:"content"           example:"hello, world"`
}

type OKResponses struct{}

func (e *OKResponses) ToOutBound() *OutBoundMessages {
	return &OutBoundMessages{
		Event: enums.OUTBOUND_OK,
	}
}
