package chats

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/brain-flowing-company/pprp-backend/apperror"
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/brain-flowing-company/pprp-backend/internal/utils"
	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
)

type WebsocketClients struct {
	conn    *websocket.Conn
	hub     *Hub
	UserId  uuid.UUID
	Message chan *models.Messages
}

func NewClient(conn *websocket.Conn, hub *Hub, userId uuid.UUID) *WebsocketClients {
	return &WebsocketClients{
		conn:    conn,
		hub:     hub,
		UserId:  userId,
		Message: make(chan *models.Messages),
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

		if c.UserId == raw.ReceiverId {
			errCh <- errors.New("could not send message to yourself")
			continue
		}

		msg := &models.Messages{
			MessageId:  uuid.New(),
			SenderId:   c.UserId,
			Read:       false,
			ReceiverId: raw.ReceiverId,
			Content:    raw.Content,
			CreatedAt:  raw.CreatedAt,
			Tag:        raw.Tag,
		}

		c.hub.SendMessage <- msg
	}
}
