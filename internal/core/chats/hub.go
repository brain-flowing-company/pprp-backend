package chats

import (
	"fmt"
	"sync"
	"time"

	"github.com/brain-flowing-company/pprp-backend/apperror"
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/brain-flowing-company/pprp-backend/internal/utils"
	"github.com/google/uuid"
)

type Hub struct {
	Clients     map[uuid.UUID]*WebsocketClients
	Register    chan *WebsocketClients
	Unregister  chan *WebsocketClients
	SendMessage chan *models.Messages
	ChatRepo    Repository
}

func NewHub(chatRepo Repository) *Hub {
	return &Hub{
		Clients:     make(map[uuid.UUID]*WebsocketClients),
		Register:    make(chan *WebsocketClients),
		Unregister:  make(chan *WebsocketClients),
		SendMessage: make(chan *models.Messages),
		ChatRepo:    chatRepo,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.register(client)

		case client := <-h.Unregister:
			h.unregister(client)

		case msg := <-h.SendMessage:
			h.sendMessage(msg)
		}
	}
}

func (h *Hub) register(client *WebsocketClients) {
	fmt.Println("Registering", client.UserId)
	_, ok := h.Clients[client.UserId]
	if !ok {
		h.Clients[client.UserId] = client
	}
}

func (h *Hub) unregister(client *WebsocketClients) {
	fmt.Println("Unregistering", client.UserId)
	_, ok := h.Clients[client.UserId]
	if ok {
		delete(h.Clients, client.UserId)
		close(client.Message)
	}
}

func (h *Hub) sendMessage(msg *models.Messages) {
	sendUser, sendUserOnline := h.Clients[msg.SenderId]

	var recvUser *WebsocketClients
	recvUserOnline := false
	if msg.ReceiverId != nil {
		recvUser, recvUserOnline = h.Clients[*msg.ReceiverId]
	}

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		err := h.ChatRepo.CreateMessages(msg)
		if sendUserOnline && err != nil {
			fmt.Println(err)
			utils.WebsocketError(sendUser.conn, apperror.
				New(apperror.InternalServerError).
				Describe("Could not send message"))
		}
	}()

	go func() {
		defer wg.Done()
		if sendUserOnline {
			sendUser.Message <- msg
		}

		if recvUserOnline && recvUser.RecvUserId != nil && *recvUser.RecvUserId == sendUser.UserId {
			now := time.Now()
			msg.ReadAt = &now
			msg.Tag = ""
			recvUser.Message <- msg
		}
	}()

	wg.Wait()
}
