package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {

		return true
	},
}

func handleConnections(w http.ResponseWriter, r *http.Request) {

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade Error: ", err)
		return
	}
	defer ws.Close()

	for {
		_, msg, err := ws.ReadMessage()

		if err != nil {
			log.Println("Read Error: ", err)
			break
		}
		log.Printf("Received: %s", msg)

	}
}

func main() {
	http.HandleFunc("/ws", handleConnections)

	log.Println("Websocket server starting on :8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
