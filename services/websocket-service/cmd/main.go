package main

import (
	"log"
	"net/http"

	"websocket-backend/services/websocket-service/internal"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // In production, validate the origin!
	},
}

func serveWs(hub *internal.Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Upgrade error: %v", err)
		return
	}
	// Register connection using exported method.
	hub.RegisterConnection(conn)
	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				break
			}
			// Use the exported Broadcast method.
			hub.Broadcast(message)
		}
	}()
}

func main() {
	hub := internal.NewHub()
	go hub.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	log.Println("WebSocket Service is running on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
