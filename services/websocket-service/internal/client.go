package internal

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

func NewClient(hub *Hub, conn *websocket.Conn) *Client {
	return &Client{
		hub:  hub,
		conn: conn,
		send: make(chan []byte, 256),
	}
}

func (c *Client) ReadPump() {
	defer func() {
		c.hub.unregister <- c.conn
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()

		if err != nil {
			log.Printf("Read error: %v", err)
			break
		}

		c.hub.broadcast <- message
	}
}

func (c *Client) WritePump() {
	defer c.conn.Close()

	for msg := range c.send {
		err := c.conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Printf("Write error: %v", err)
			break
		}
	}
}
