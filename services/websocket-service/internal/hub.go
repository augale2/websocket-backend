package internal

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Hub struct {
	clients map[*websocket.Conn]bool

	broadcast chan []byte

	register chan *websocket.Conn

	unregister chan *websocket.Conn

	mu sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*websocket.Conn]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
	}
}
func (h *Hub) RegisterConnection(conn *websocket.Conn) {
	h.register <- conn
}

func (h *Hub) Broadcast(message []byte) {
	h.broadcast <- message
}

func (h *Hub) Run() {
	for {
		select {
		case conn := <-h.register:
			h.mu.Lock()
			h.clients[conn] = true
			h.mu.Unlock()
			log.Println("Client registered")

		case conn := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[conn]; ok {
				delete(h.clients, conn)
				conn.Close()
				log.Println("Client unregistered")
			}
			h.mu.Unlock()
		case message := <-h.broadcast:
			h.mu.Lock()

			for conn := range h.clients {
				if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
					log.Printf("Error broadcasting: %v", err)
					conn.Close()
					delete(h.clients, conn)
				}
			}

			h.mu.Unlock()

		}
	}
}
