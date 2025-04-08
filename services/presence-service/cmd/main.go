package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	// Import Presence Service proto.
	pb "websocket-backend/services/presence-service/proto"
	// Import Auth Service proto.
	authpb "websocket-backend/services/auth-service/proto"
	// Import the Auth client helper from your Auth Service public package.
	auth "websocket-backend/services/auth-service/pkg/authclient"
)

// presenceServer implements the PresenceServiceServer interface.
type presenceServer struct {
	pb.UnimplementedPresenceServiceServer
	// In a real implementation, you might have a Presence tracker here.
}

// validateToken uses the Auth Service to validate the JWT token.
func validateToken(ctx context.Context, token string) error {
	// Create an Auth client. Adjust the address if needed.
	authClient, conn, err := auth.CreateAuthClient("localhost:50052")
	if err != nil {
		return fmt.Errorf("could not connect to auth service: %v", err)
	}
	defer conn.Close()

	validateCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	resp, err := authClient.ValidateToken(validateCtx, &authpb.ValidateTokenRequest{
		Token: token,
	})
	if err != nil {
		return fmt.Errorf("token validation error: %v", err)
	}
	if !resp.Valid {
		return fmt.Errorf("invalid token: %s", resp.Message)
	}
	return nil
}

// UpdatePresence validates the token and updates the presence for a user.
func (s *presenceServer) UpdatePresence(ctx context.Context, req *pb.UpdatePresenceRequest) (*pb.UpdatePresenceResponse, error) {
	if err := validateToken(ctx, req.GetToken()); err != nil {
		return nil, err
	}
	// For demo, simply log that presence is updated.
	log.Printf("Updated presence for user: %s", req.GetUserId())
	return &pb.UpdatePresenceResponse{
		Message: "Presence updated successfully",
	}, nil
}

// GetOnlineUsers validates the token and returns a list of online user IDs.
func (s *presenceServer) GetOnlineUsers(ctx context.Context, req *pb.GetOnlineUsersRequest) (*pb.GetOnlineUsersResponse, error) {
	if err := validateToken(ctx, req.GetToken()); err != nil {
		return nil, err
	}
	// For demo purposes, return a static list. In a real app, query your presence tracker.
	onlineUsers := []string{"user1", "user2", "user3"}
	log.Printf("Returning online users")
	return &pb.GetOnlineUsersResponse{
		UserIds: onlineUsers,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterPresenceServiceServer(grpcServer, &presenceServer{})
	reflection.Register(grpcServer)
	log.Println("Presence Service gRPC server is running on port 50053")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
