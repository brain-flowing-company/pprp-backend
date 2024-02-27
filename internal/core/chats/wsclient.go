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

func (c *WebsocketClients) SendMessage(msg *models.OutBoundMessages) {
	c.router.Send(msg)
}

func (c *WebsocketClients) Listen() {
	c.router.On(enums.INBOUND_MSG, c.inBoundMsgHandler)
	c.router.On(enums.INBOUND_JOIN, c.inBoundJoinHandler)
	c.router.On(enums.INBOUND_LEFT, c.inBoundLeftHandler)
	c.router.Listen()
}

func (c *WebsocketClients) Close() {
	c.router.Close()
}

func (c *WebsocketClients) inBoundMsgHandler(inbound *models.InBoundMessages) *apperror.AppError {
	var readAt *time.Time
	now := time.Now()
	if c.hub.IsUserInChat(c.UserId, *c.RecvUserId) {
		readAt = &now
	}

	msg := &models.Messages{
		MessageId:  uuid.New(),
		SenderId:   c.UserId,
		ReceiverId: c.RecvUserId,
		ReadAt:     readAt,
		Content:    inbound.Content,
		SentAt:     inbound.SentAt,
	}

	err := c.Service.SaveMessages(msg)
	if err != nil {
		return err
	}

	if c.hub.IsUserOnline(c.UserId) {
		c.SendMessage(msg.ToOutBound().SetTag(inbound.Tag))
	}

	if c.hub.IsUserInChat(c.UserId, *c.RecvUserId) {
		c.hub.GetUser(*c.RecvUserId).SendMessage(msg.ToOutBound())
	} else if c.hub.IsUserOnline(*c.RecvUserId) {
		chatResponse := models.ChatPreviews{
			Content:        inbound.Content,
			UnreadMessages: 1,
			UserId:         c.UserId,
		}
		c.hub.GetUser(*c.RecvUserId).SendMessage(chatResponse.ToOutBound())
	}

	return nil
}

func (c *WebsocketClients) inBoundJoinHandler(inbound *models.InBoundMessages) *apperror.AppError {
	uuid, err := uuid.Parse(inbound.Content)
	if err != nil {
		return apperror.
			New(apperror.BadRequest).
			Describe("invalid receiver uuid")
	}

	if c.UserId == uuid {
		return apperror.
			New(apperror.BadRequest).
			Describe("could not send message to yourself")
	}

	fmt.Println("Joining", uuid)
	c.RecvUserId = &uuid

	apperr := c.Service.ReadMessages(uuid, c.UserId)
	if apperr != nil {
		return apperr
	}

	if c.hub.IsUserOnline(uuid) {
		read := models.ReadEvents{
			SenderId:   uuid,
			ReceiverId: c.UserId,
			ReadAt:     time.Now(),
		}
		c.hub.GetUser(uuid).SendMessage(read.ToOutBound())
	}

	chatResponse := models.ChatPreviews{
		UnreadMessages: 0,
		UserId:         uuid,
	}
	c.SendMessage(chatResponse.ToOutBound())

	return nil
}

func (c *WebsocketClients) inBoundLeftHandler(inbound *models.InBoundMessages) *apperror.AppError {
	fmt.Println("Leaving", c.RecvUserId)
	c.RecvUserId = nil

	return nil
}
