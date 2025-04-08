package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "websocket-backend/services/document-service/proto"

	authpb "websocket-backend/services/auth-service/proto"

	auth "websocket-backend/services/auth-service/pkg/authclient"

	"websocket-backend/services/document-service/internal/db"

	wsclient "websocket-backend/services/document-service/pkg/wsclient"

	ws "websocket-backend/services/websocket-service/proto"
)

type server struct {
	pb.UnimplementedDocumentServiceServer
	dbClient *db.DynamoDBClient
}

func publishDocumentEvent(ctx context.Context, documentID, eventType, title, content string) error {
	// Connect to the websocket service on port 50053.
	wsClient, conn, err := wsclient.CreateWSClient("localhost:50053")
	if err != nil {
		return fmt.Errorf("failed to create ws client: %v", err)
	}
	defer conn.Close()

	// Create an event with the current timestamp.
	event := &ws.DocumentEvent{
		DocumentId: documentID,
		EventType:  eventType,
		Title:      title,
		Content:    content,
		Timestamp:  time.Now().Unix(),
	}

	// Call the PublishDocumentEvent method with a timeout.
	callCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	res, err := wsClient.PublishDocumentEvent(callCtx, event)
	if err != nil {
		return fmt.Errorf("PublishDocumentEvent call failed: %v", err)
	}
	if !res.Success {
		return fmt.Errorf("websocket service error: %s", res.Message)
	}
	log.Printf("Event published: %s for document %s", eventType, documentID)
	return nil
}

func validateUserToken(ctx context.Context, token string) error {

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
func (s *server) CreateDocument(ctx context.Context, req *pb.CreateDocumentRequest) (*pb.CreateDocumentResponse, error) {
	if err := validateUserToken(ctx, req.GetToken()); err != nil {
		return nil, err
	}

	docID, err := s.dbClient.CreateDocument(ctx, req.GetTitle(), req.GetContent())
	if err != nil {
		return nil, err
	}
	log.Printf("Created document with ID: %s", docID)

	// Publish a "created" event.
	if err := publishDocumentEvent(ctx, docID, "created", req.GetTitle(), req.GetContent()); err != nil {
		log.Printf("Error publishing 'document_created' event: %v", err)
		// Decide whether to fail the request or simply log the error.
	}

	return &pb.CreateDocumentResponse{
		DocumentId: docID,
		Message:    "Document created successfully",
	}, nil
}

func (s *server) UpdateDocument(ctx context.Context, req *pb.UpdateDocumentRequest) (*pb.UpdateDocumentResponse, error) {
	if err := validateUserToken(ctx, req.GetToken()); err != nil {
		return nil, err
	}
	err := s.dbClient.UpdateDocument(ctx, req.GetDocumentId(), req.GetTitle(), req.GetContent())
	if err != nil {
		return nil, err
	}
	log.Printf("Updated document with ID: %s", req.GetDocumentId())

	// Publish an "updated" event.
	if err := publishDocumentEvent(ctx, req.GetDocumentId(), "updated", req.GetTitle(), req.GetContent()); err != nil {
		log.Printf("Error publishing 'document_updated' event: %v", err)
	}

	return &pb.UpdateDocumentResponse{
		Message: "Document updated successfully",
	}, nil
}

func (s *server) GetDocument(ctx context.Context, req *pb.GetDocumentRequest) (*pb.GetDocumentResponse, error) {

	if err := validateUserToken(ctx, req.GetToken()); err != nil {
		return nil, err
	}
	item, err := s.dbClient.GetDocument(ctx, req.GetDocumentId())
	if err != nil {
		return nil, err
	}
	log.Printf("DynamoDB item: %+v", item)
	title := item["Title"].(*types.AttributeValueMemberS).Value
	content := item["Content"].(*types.AttributeValueMemberS).Value
	log.Printf("Fetched document with ID: %s", req.GetDocumentId())
	return &pb.GetDocumentResponse{
		DocumentId: req.GetDocumentId(),
		Title:      title,
		Content:    content,
	}, nil
}

func main() {
	ctx := context.Background()
	dbClient, err := db.NewDynamoDBClient(ctx, "Documents")
	if err != nil {
		log.Fatalf("Failed to create DynamoDB client: %v", err)
	}
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterDocumentServiceServer(grpcServer, &server{dbClient: dbClient})
	reflection.Register(grpcServer)
	log.Println("Document Service gRPC server is running on port 50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
