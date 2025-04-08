// internal/grpc_server.go
package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	pb "websocket-backend/services/websocket-service/proto"
)

// WebsocketGRPCServer implements the gRPC service for document events.
type WebsocketGRPCServer struct {
	pb.UnimplementedWebsocketServiceServer
	Hub *Hub
}

// PublishDocumentEvent receives an event from the document service and forwards it to all connected websocket clients.
func (s *WebsocketGRPCServer) PublishDocumentEvent(ctx context.Context, req *pb.DocumentEvent) (*pb.EventResponse, error) {
	// Optionally, you can add any pre-processing or validation here.
	// For example, logging the received event.
	log.Printf("Received event: document_id=%s, event_type=%s", req.DocumentId, req.EventType)

	// Marshal the event into JSON. This message will be broadcasted.
	message, err := json.Marshal(req)
	if err != nil {
		errMsg := fmt.Sprintf("failed to marshal event: %v", err)
		log.Println(errMsg)
		return &pb.EventResponse{Success: false, Message: errMsg}, err
	}

	// Use the Hub to broadcast the message to all connected websocket clients.
	s.Hub.Broadcast(message)
	log.Printf("Event broadcasted: %s", message)

	return &pb.EventResponse{Success: true, Message: "Event broadcasted successfully"}, nil
}
