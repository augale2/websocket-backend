package main

import (
	"context"
	"log"
	"net"

	"websocket-backend/services/document-service/internal/db"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	pb "websocket-backend/services/document-service/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// server implements the DocumentServiceServer interface.
type server struct {
	pb.UnimplementedDocumentServiceServer
	dbClient *db.DynamoDBClient
}

// CreateDocument implements the CreateDocument RPC.
func (s *server) CreateDocument(ctx context.Context, req *pb.CreateDocumentRequest) (*pb.CreateDocumentResponse, error) {
	docID, err := s.dbClient.CreateDocument(ctx, req.GetTitle(), req.GetContent())

	if err != nil {
		return nil, err
	}
	log.Printf("Created document with ID: %s", docID)

	return &pb.CreateDocumentResponse{
		DocumentId: docID,
		Message:    "Document created successfully",
	}, nil
}

// UpdateDocument implements the UpdateDocument RPC.
func (s *server) UpdateDocument(ctx context.Context, req *pb.UpdateDocumentRequest) (*pb.UpdateDocumentResponse, error) {
	err := s.dbClient.UpdateDocument(ctx, req.GetDocumentId(), req.GetTitle(), req.GetContent())
	if err != nil {
		return nil, err
	}
	log.Printf("Updated document with ID: %s", req.GetDocumentId())
	return &pb.UpdateDocumentResponse{
		Message: "Document updated successfully",
	}, nil
}

// GetDocument implements the GetDocument RPC.
func (s *server) GetDocument(ctx context.Context, req *pb.GetDocumentRequest) (*pb.GetDocumentResponse, error) {
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
	s := grpc.NewServer()
	pb.RegisterDocumentServiceServer(s, &server{dbClient: dbClient})
	reflection.Register(s)
	log.Println("Document Service gRPC server is running on port 50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
