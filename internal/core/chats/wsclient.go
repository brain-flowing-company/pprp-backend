package chats

import (
	"fmt"
	"time"

	"github.com/brain-flowing-company/pprp-backend/apperror"
	"github.com/brain-flowing-company/pprp-backend/internal/enums"
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
)

type WebsocketClients struct {
	router     *WebsocketRouter
	hub        *Hub
	Service    Service
	UserId     uuid.UUID
	RecvUserId *uuid.UUID
}

func NewClient(conn *websocket.Conn, hub *Hub, service Service, userId uuid.UUID) *WebsocketClients {
	return &WebsocketClients{
		router:     NewWebsocketRouter(conn),
		hub:        hub,
		Service:    service,
		UserId:     userId,
		RecvUserId: nil,
	}
}

func (client *WebsocketClients) SendMessage(msg *models.OutBoundMessages) {
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
	var readAt *time.Time
	now := time.Now()
	if client.hub.IsUserInChat(client.UserId, *client.RecvUserId) {
		readAt = &now
	}

	msg := &models.Messages{
		MessageId:  uuid.New(),
		SenderId:   client.UserId,
		ReceiverId: client.RecvUserId,
		ReadAt:     readAt,
		Content:    inbound.Content,
		SentAt:     inbound.SentAt,
	}

	err := client.Service.SaveMessages(msg)
	if err != nil {
		return err
	}

	if client.hub.IsUserOnline(client.UserId) {
		client.SendMessage(msg.ToOutBound().SetTag(inbound.Tag))
	}

	if client.hub.IsUserInChat(client.UserId, *client.RecvUserId) {
		client.hub.GetUser(*client.RecvUserId).SendMessage(msg.ToOutBound())
	} else if client.hub.IsUserOnline(*client.RecvUserId) {
		chatResponse := models.ChatPreviews{
			Content:        inbound.Content,
			UnreadMessages: 1,
			UserId:         client.UserId,
		}
		client.hub.GetUser(*client.RecvUserId).SendMessage(chatResponse.ToOutBound())
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

	fmt.Println("Joining", uuid)
	client.RecvUserId = &uuid

	apperr := client.Service.ReadMessages(uuid, client.UserId)
	if apperr != nil {
		return apperr
	}

	if client.hub.IsUserInChat(client.UserId, uuid) {
		read := models.ReadEvents{
			SenderId:   uuid,
			ReceiverId: client.UserId,
			ReadAt:     time.Now(),
		}
		client.hub.GetUser(uuid).SendMessage(read.ToOutBound())
	}

	chatResponse := models.ChatPreviews{
		UnreadMessages: 0,
		UserId:         uuid,
	}
	client.SendMessage(chatResponse.ToOutBound())

	return nil
}

func (client *WebsocketClients) inBoundLeftHandler(inbound *models.InBoundMessages) *apperror.AppError {
	fmt.Println("Leaving", client.RecvUserId)
	client.RecvUserId = nil

	return nil
}
