package authclient

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	authpb "websocket-backend/services/auth-service/proto"
)

// CreateAuthClient establishes a gRPC connection to the Auth Service using insecure credentials.
// In production, you should use proper transport credentials (e.g., TLS).
func CreateAuthClient(authServiceAddr string) (authpb.AuthServiceClient, *grpc.ClientConn, error) {
	conn, err := grpc.NewClient(authServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}
	client := authpb.NewAuthServiceClient(conn)
	return client, conn, nil
}
