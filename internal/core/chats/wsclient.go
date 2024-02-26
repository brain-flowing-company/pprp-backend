package chats

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/brain-flowing-company/pprp-backend/apperror"
	"github.com/brain-flowing-company/pprp-backend/internal/enums"
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/brain-flowing-company/pprp-backend/internal/utils"
	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
)

type WebsocketClients struct {
	conn             *websocket.Conn
	hub              *Hub
	Service          Service
	OutBoundMessages chan *models.OutBoundMessages
	UserId           uuid.UUID
	RecvUserId       *uuid.UUID
}

func NewClient(conn *websocket.Conn, hub *Hub, service Service, userId uuid.UUID) *WebsocketClients {
	return &WebsocketClients{
		conn:             conn,
		hub:              hub,
		Service:          service,
		OutBoundMessages: make(chan *models.OutBoundMessages),
		UserId:           userId,
		RecvUserId:       nil,
	}
}

func (c *WebsocketClients) Listen() {
	errCh := make(chan *apperror.AppError)
	term := make(chan bool)

	go c.writerHandler()
	go c.readHandler(term, errCh)

	for {
		select {
		case <-term:
			return

		case err := <-errCh:
			utils.WebsocketError(c.conn, err)
		}
	}
}

func (c *WebsocketClients) writerHandler() {
	for {
		msg, isAlive := <-c.OutBoundMessages
		if !isAlive {
			return
		}

		c.conn.WriteJSON(msg)
	}
}

func (c *WebsocketClients) readHandler(term chan bool, errCh chan *apperror.AppError) {
	for {
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Println(err)
			}
			term <- true
			break
		}

		var inbound models.InBoundMessages
		err = json.Unmarshal(data, &inbound)
		if err != nil {
			errCh <- apperror.
				New(apperror.BadRequest).
				Describe("could not parse json")
			continue
		}

		switch inbound.Event {
		case enums.INBOUND_MSG:
			err := c.inBoundMsgHandler(&inbound)
			if err != nil {
				errCh <- err
				continue
			}

		case enums.INBOUND_JOIN:
			err := c.inBoundJoinHandler(&inbound)
			if err != nil {
				errCh <- err
				continue
			}

		case enums.INBOUND_LEFT:
			c.inBoundLeftHandler()

		default:
			errCh <- apperror.
				New(apperror.BadRequest).
				Describe("invalid event type")
		}
	}
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
		c.OutBoundMessages <- msg.ToOutBound().SetTag(inbound.Tag)
	}

	if c.hub.IsUserInChat(c.UserId, *c.RecvUserId) {
		c.hub.GetUser(*c.RecvUserId).OutBoundMessages <- msg.ToOutBound()
	} else if c.hub.IsUserOnline(*c.RecvUserId) {
		chatResponse := models.ChatsResponses{
			Content:        inbound.Content,
			UnreadMessages: 1,
			UserId:         c.UserId,
		}
		c.hub.GetUser(*c.RecvUserId).OutBoundMessages <- chatResponse.ToOutBound()
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
		c.hub.GetUser(uuid).OutBoundMessages <- read.ToOutBound()
	}

	chatResponse := models.ChatsResponses{
		UnreadMessages: 0,
		UserId:         uuid,
	}
	c.OutBoundMessages <- chatResponse.ToOutBound()

	return nil
}

func (c *WebsocketClients) inBoundLeftHandler() {
	fmt.Println("Leaving", c.RecvUserId)
	c.RecvUserId = nil
}
