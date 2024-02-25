package chats

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/brain-flowing-company/pprp-backend/apperror"
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/brain-flowing-company/pprp-backend/internal/utils"
	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
)

type WebsocketClients struct {
	conn       *websocket.Conn
	hub        *Hub
	UserId     uuid.UUID
	RecvUserId *uuid.UUID
	Message    chan *models.Messages
}

func NewClient(conn *websocket.Conn, hub *Hub, userId uuid.UUID) *WebsocketClients {
	return &WebsocketClients{
		conn:       conn,
		hub:        hub,
		UserId:     userId,
		RecvUserId: nil,
		Message:    make(chan *models.Messages),
	}
}

func (c *WebsocketClients) Listen() {
	errCh := make(chan error)
	term := make(chan bool)

	go c.writerHandler()
	go c.readHandler(term, errCh)

	for {
		select {
		case <-term:
			return

		case err := <-errCh:
			utils.WebsocketError(c.conn, apperror.
				New(apperror.InternalServerError).
				Describe(err.Error()))
		}
	}

}

func (c *WebsocketClients) writerHandler() {
	for {
		msg, isAlive := <-c.Message
		if !isAlive {
			return
		}

		c.conn.WriteJSON(msg)
	}
}

func (c *WebsocketClients) readHandler(term chan bool, errCh chan error) {
	for {
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Println(err)
			}
			term <- true
			break
		}

		var raw models.RawMessages
		err = json.Unmarshal(data, &raw)
		if err != nil {
			errCh <- err
			continue
		}

		msg := &models.Messages{
			MessageId:  uuid.New(),
			SenderId:   c.UserId,
			ReceiverId: *c.RecvUserId,
			ReadAt:     nil,
			Content:    raw.Content,
			SentAt:     raw.SentAt,
			Tag:        raw.Tag,
		}

		send := false
		for _, tag := range strings.Split(msg.Tag, ";") {
			key, val := utils.SplitByFirstString(tag, "=")
			switch strings.ToLower(key) {
			case "tag":
				msg.Tag = tag
				send = true

			case "join":
				uuid, err := uuid.Parse(val)
				if err != nil {
					errCh <- errors.New("invalid receiver uuid")
					continue
				}

				if c.UserId == uuid {
					errCh <- errors.New("could not send message to yourself")
					continue
				}

				c.RecvUserId = &uuid

			case "leave":
				c.RecvUserId = nil
			}
		}

		if send {
			c.hub.SendMessage <- msg
		}
	}
}
