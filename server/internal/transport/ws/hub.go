package ws

import (
	"sync"

	"github.com/gorilla/websocket"
	realtime_model "github.com/nougght/monitoring-system/server/internal/model/realtime"
)

// Hub maintains the set of active clients and broadcasts messages to them.
type Hub struct {
	clients map[*Client]struct{}
	mu      sync.RWMutex

	// Inbound messages from the clients
	broadcast chan *realtime_model.Message

	sub       chan string
	unsub     chan string
	subFunc   func(subject string) // subjects that 	subFunc   func(subject string)
	unsubFunc func(subjec string)
}

func NewHub(subFunc func(subject string), unsubFunc func(subject string)) *Hub {
	return &Hub{
		broadcast: make(chan *realtime_model.Message),
		clients:   make(map[*Client]struct{}),
		subFunc:   subFunc,
		unsubFunc: unsubFunc,
		sub:       make(chan string),
		unsub:     make(chan string),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case s := <-h.sub:
			h.subFunc(s)
		case s := <-h.unsub:
			h.unsubFunc(s)
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.toSend <- message:
				default:
				}
			}
		}
	}
}

func (h *Hub) RegisterClient(conn *websocket.Conn) {
	client := NewClient(conn, func(subject string) {
		h.sub <- subject
	})
	h.clients[client] = struct{}{}
	client.Run()
}

func (h *Hub) SendMessage(subject string, message *realtime_model.Message) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.clients {
		if client.IsSub(subject) {
			select {
			case client.toSend <- message:
			default:
			}
		}
	}
}
