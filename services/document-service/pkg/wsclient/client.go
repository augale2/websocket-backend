// websocket-backend/services/document-service/pkg/wsclient/client.go
package wsclient

import (
	ws "websocket-backend/services/websocket-service/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// CreateWSClient establishes a gRPC connection to the websocket service using insecure credentials.
// In production, replace insecure.NewCredentials() with proper transport credentials.
func CreateWSClient(address string) (ws.WebsocketServiceClient, *grpc.ClientConn, error) {
	// Directly connect without a context.
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}
	client := ws.NewWebsocketServiceClient(conn)
	return client, conn, nil
}
