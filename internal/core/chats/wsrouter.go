package chats

import (
	"encoding/json"
	"fmt"

	"github.com/brain-flowing-company/pprp-backend/apperror"
	"github.com/brain-flowing-company/pprp-backend/internal/enums"
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/brain-flowing-company/pprp-backend/internal/utils"
	"github.com/gofiber/contrib/websocket"
)

type handlerFunc func(*models.InBoundMessages) *apperror.AppError

type WebsocketRouter struct {
	conn             *websocket.Conn
	handlers         map[enums.MessageInboundEvents]handlerFunc
	outBoundMessages chan *models.OutBoundMessages
}

func NewWebsocketRouter(conn *websocket.Conn) *WebsocketRouter {
	return &WebsocketRouter{
		conn:             conn,
		handlers:         make(map[enums.MessageInboundEvents]handlerFunc),
		outBoundMessages: make(chan *models.OutBoundMessages),
	}
}

func (r *WebsocketRouter) On(e enums.MessageInboundEvents, h handlerFunc) {
	r.handlers[e] = h
}

func (r *WebsocketRouter) Send(msg *models.OutBoundMessages) {
	r.outBoundMessages <- msg
}

func (r *WebsocketRouter) Listen() {
	errch := make(chan *apperror.AppError)
	term := make(chan bool)

	go r.handleWrite()
	go r.handleRead(term, errch)

	for {
		select {
		case <-term:
			return

		case err := <-errch:
			utils.WebsocketError(r.conn, err)
		}
	}
}

func (r *WebsocketRouter) Close() {
	close(r.outBoundMessages)
}

func (r *WebsocketRouter) handleWrite() {
	for {
		msg, isAlive := <-r.outBoundMessages
		if !isAlive {
			return
		}

		r.conn.WriteJSON(msg)
	}
}

func (r *WebsocketRouter) handleRead(term chan bool, errch chan *apperror.AppError) {
	for {
		_, data, err := r.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Println(err)
			}
			term <- true
			break
		}

		var msg models.InBoundMessages
		err = json.Unmarshal(data, &msg)
		if err != nil {
			errch <- apperror.
				New(apperror.BadRequest).
				Describe("could not parse json")
			continue
		}

		handler, ok := r.handlers[msg.Event]
		if !ok {
			continue
		}

		apperr := handler(&msg)
		if apperr != nil {
			errch <- apperr
		}
	}
}
