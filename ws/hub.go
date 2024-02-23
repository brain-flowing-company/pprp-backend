package ws

import (
	"fmt"

	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Hub struct {
	Clients     map[uuid.UUID]*WebsocketClients
	Register    chan *WebsocketClients
	Unregister  chan *WebsocketClients
	SendMessage chan *models.Messages
	db          *gorm.DB
}

func NewHub(db *gorm.DB) *Hub {
	return &Hub{
		Clients:     make(map[uuid.UUID]*WebsocketClients),
		Register:    make(chan *WebsocketClients),
		Unregister:  make(chan *WebsocketClients),
		SendMessage: make(chan *models.Messages),
		db:          db,
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
	_, ok := h.Clients[client.UserId]
	if !ok {
		h.Clients[client.UserId] = client
	}
}

func (h *Hub) unregister(client *WebsocketClients) {
	_, ok := h.Clients[client.UserId]
	if ok {
		h.Clients[client.UserId] = client
	}
}

func (h *Hub) sendMessage(msg *models.Messages) {
	fmt.Println(msg)
}
