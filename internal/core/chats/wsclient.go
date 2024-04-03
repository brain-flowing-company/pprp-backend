package chats

import (
	"fmt"
	"time"

	"github.com/brain-flowing-company/pprp-backend/apperror"
	"github.com/brain-flowing-company/pprp-backend/internal/enums"
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type WebsocketClients struct {
	router     *WebsocketRouter
	hub        *Hub
	service    Service
	UserId     uuid.UUID
	RecvUserId *uuid.UUID
	chats      map[uuid.UUID]*models.ChatPreviews
}

func NewClient(logger *zap.Logger, conn *websocket.Conn, hub *Hub, service Service, userId uuid.UUID) (*WebsocketClients, *apperror.AppError) {
	chatPreviews := []models.ChatPreviews{}
	err := service.GetAllChats(&chatPreviews, userId, "")
	if err != nil {
		return nil, err
	}

	chats := map[uuid.UUID]*models.ChatPreviews{}
	for _, chat := range chatPreviews {
		chats[chat.UserId] = new(models.ChatPreviews)
		*chats[chat.UserId] = chat
	}

	return &WebsocketClients{
		router:     NewWebsocketRouter(logger, conn),
		hub:        hub,
		service:    service,
		UserId:     userId,
		RecvUserId: nil,
		chats:      chats,
	}, nil
}

func (client *WebsocketClients) SendOutBoundMessage(msg *models.OutBoundMessages) {
	client.router.Send(msg)
}

func (client *WebsocketClients) Listen() {
	client.router.On(enums.INBOUND_MSG, client.inBoundMsgHandler)
	client.router.On(enums.INBOUND_JOIN, client.inBoundJoinHandler)
	client.router.On(enums.INBOUND_LEFT, client.inBoundLeftHandler)
	client.router.Listen()
}

func (client *WebsocketClients) Close() {
	client.router.Close()
}

func (client *WebsocketClients) inBoundMsgHandler(inbound *models.InBoundMessages) *apperror.AppError {
	if client.RecvUserId == nil {
		return apperror.
			New(apperror.NotInChat).
			Describe("Invalid chat")
	}

	fmt.Println(inbound)

	var readAt *time.Time
	now := time.Now()
	if client.hub.IsUserBothInChat(client.UserId, *client.RecvUserId) {
		readAt = &now
	}

	msg := &models.Messages{
		MessageId:  uuid.New(),
		SenderId:   client.UserId,
		ReceiverId: *client.RecvUserId,
		ChatId:     *client.RecvUserId,
		Author:     true,
		ReadAt:     readAt,
		Content:    inbound.Content,
		SentAt:     inbound.SentAt,
	}

	err := client.service.SaveMessages(msg)
	if err != nil {
		return err
	}

	if client.hub.IsUserOnline(client.UserId) {
		msg.Tag = inbound.Tag
		msg.ChatId = *client.RecvUserId
		msg.Author = true
		client.SendOutBoundMessage(msg.ToOutBound())
	}

	if client.hub.IsUserOnline(*client.RecvUserId) {
		msg.ChatId = client.UserId
		msg.Author = false
		client.hub.GetUser(*client.RecvUserId).SendOutBoundMessage(msg.ToOutBound())
	}

	return nil
}

func (client *WebsocketClients) inBoundJoinHandler(inbound *models.InBoundMessages) *apperror.AppError {
	uuid, err := uuid.Parse(inbound.Content)
	if err != nil {
		return apperror.
			New(apperror.BadRequest).
			Describe("invalid receiver uuid")
	}

	if client.UserId == uuid {
		return apperror.
			New(apperror.BadRequest).
			Describe("could not send message to yourself")
	}

	client.RecvUserId = &uuid

	apperr := client.service.ReadMessages(uuid, client.UserId)
	if apperr != nil {
		return apperr
	}

	if client.hub.IsUserBothInChat(client.UserId, uuid) {
		read := &models.ReadEvents{
			ChatId: client.UserId,
			ReadAt: time.Now(),
		}
		client.hub.GetUser(uuid).SendOutBoundMessage(read.ToOutBound())
	}

	if (models.MessageAttatchments{}) != inbound.Attatchment {
		property := models.Properties{PropertyId: *inbound.Attatchment.PropertyId}
		apperr := client.hub.SendNotificationMessage(&property, true, "Embedded property", client.UserId, uuid)
		if apperr != nil {
			return apperr
		}
	}

	return nil
}

func (client *WebsocketClients) inBoundLeftHandler(inbound *models.InBoundMessages) *apperror.AppError {
	client.RecvUserId = nil

	return nil
}
