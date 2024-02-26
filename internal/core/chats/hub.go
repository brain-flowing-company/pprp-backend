package chats

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
)

type Hub struct {
	sync.Mutex
	clients map[uuid.UUID]*WebsocketClients
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[uuid.UUID]*WebsocketClients),
	}
}

func (h *Hub) GetUser(userId uuid.UUID) *WebsocketClients {
	return h.clients[userId]
}

func (h *Hub) IsUserOnline(userId uuid.UUID) bool {
	_, online := h.clients[userId]
	return online
}

func (h *Hub) IsUserInChat(sendUserId uuid.UUID, recvUserId uuid.UUID) bool {
	sendUser, sendOnline := h.clients[sendUserId]
	recvUser, recvOnline := h.clients[recvUserId]

	if !sendOnline || !recvOnline {
		return false
	}

	return sendUser.RecvUserId != nil && recvUser.RecvUserId != nil &&
		*sendUser.RecvUserId == recvUserId && *recvUser.RecvUserId == sendUserId
}

func (h *Hub) Register(client *WebsocketClients) {
	h.Lock()
	fmt.Println("Registering", client.UserId)
	_, ok := h.clients[client.UserId]
	if !ok {
		h.clients[client.UserId] = client
	}
	h.Unlock()
}

func (h *Hub) Unregister(client *WebsocketClients) {
	h.Lock()
	fmt.Println("Unregistering", client.UserId)
	_, ok := h.clients[client.UserId]
	if ok {
		delete(h.clients, client.UserId)
		close(client.OutBoundMessages)
	}
	h.Unlock()
}
