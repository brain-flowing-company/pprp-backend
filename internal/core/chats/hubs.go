package chats

import (
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/brain-flowing-company/pprp-backend/apperror"
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/google/uuid"
)

type Hub struct {
	sync.Mutex
	clients map[uuid.UUID]*WebsocketClients
	service Service
}

func NewHub(service Service) *Hub {
	return &Hub{
		clients: make(map[uuid.UUID]*WebsocketClients),
		service: service,
	}
}

func (h *Hub) GetUser(userId uuid.UUID) *WebsocketClients {
	return h.clients[userId]
}

func (h *Hub) SendNotificationMessage(attatchment interface{}, content string, senderId uuid.UUID, receiverId uuid.UUID) *apperror.AppError {
	var readAt *time.Time
	now := time.Now()
	if h.IsReceiverInChat(senderId, receiverId) {
		readAt = &now
	}

	msg := &models.Messages{
		MessageId:   uuid.New(),
		SenderId:    senderId,
		ReceiverId:  receiverId,
		ChatId:      senderId,
		Author:      true,
		ReadAt:      readAt,
		Content:     content,
		SentAt:      now,
		Attatchment: models.MessageAttatchments{},
	}

	switch attch := attatchment.(type) {
	case *models.CreatingAppointments:
		msg.Attatchment.AppointmentId = &attch.AppointmentId
	case *models.CreatingAgreements:
		msg.Attatchment.AgreementId = &attch.AgreementId
	default:
		return apperror.New(apperror.BadRequest).Describe(fmt.Sprintf("notification message does not support %v", reflect.TypeOf(attatchment)))
	}

	err := h.service.SaveMessages(msg)
	if err != nil {
		return err
	}

	if h.IsUserOnline(senderId) {
		msg.ChatId = receiverId
		msg.Author = true
		h.GetUser(senderId).SendOutBoundMessage(msg.ToOutBound())
	}

	if h.IsUserOnline(receiverId) {
		msg.ChatId = senderId
		msg.Author = false
		h.GetUser(receiverId).SendOutBoundMessage(msg.ToOutBound())
	}

	return nil
}

func (h *Hub) IsUserOnline(userId uuid.UUID) bool {
	_, online := h.clients[userId]
	return online
}

func (h *Hub) IsReceiverInChat(sendUserId uuid.UUID, recvUserId uuid.UUID) bool {
	recvUser, recvOnline := h.clients[recvUserId]

	if !recvOnline {
		return false
	}

	return recvUser.RecvUserId != nil && *recvUser.RecvUserId == sendUserId
}

func (h *Hub) IsUserBothInChat(sendUserId uuid.UUID, recvUserId uuid.UUID) bool {
	return h.IsReceiverInChat(sendUserId, recvUserId) &&
		h.IsReceiverInChat(recvUserId, sendUserId)
}

func (h *Hub) Register(client *WebsocketClients) {
	h.Lock()
	_, ok := h.clients[client.UserId]
	if !ok {
		h.clients[client.UserId] = client
	}
	h.Unlock()
}

func (h *Hub) Unregister(client *WebsocketClients) {
	h.Lock()
	_, ok := h.clients[client.UserId]
	if ok {
		delete(h.clients, client.UserId)
		client.Close()
	}
	h.Unlock()
}
