package main

import (
	"log"
	"net"
	"net/http"

	"websocket-backend/services/websocket-service/internal"

	"github.com/gorilla/websocket"

	pb "websocket-backend/services/websocket-service/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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

	go func() {
		http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
			serveWs(hub, w, r)
		})
		log.Println("WebSocket HTTP server is running on port 8080")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("Failed to listen on port 50053: %v", err)
	}
	grpcServer := grpc.NewServer()
	wsGRPCServer := &internal.WebsocketGRPCServer{Hub: hub}
	pb.RegisterWebsocketServiceServer(grpcServer, wsGRPCServer)
	reflection.Register(grpcServer)
	log.Println("WebSocket gRPC server is running on port 50053")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
