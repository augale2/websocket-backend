package db

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

type UserDBClient struct {
	Client *dynamodb.Client
	Table  string
}

func NewUserDBClient(ctx context.Context, tableName string) (*UserDBClient, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %v", err)
	}
	client := dynamodb.NewFromConfig(cfg)
	return &UserDBClient{
		Client: client,
		Table:  tableName,
	}, nil
}

func (db *UserDBClient) CreateUser(ctx context.Context, username, hashedPassword string) (string, error) {

	userID := uuid.New().String()
	_, err := db.Client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(db.Table),
		Item: map[string]types.AttributeValue{
			"UserId":   &types.AttributeValueMemberS{Value: userID},
			"Username": &types.AttributeValueMemberS{Value: username},
			"Password": &types.AttributeValueMemberS{Value: hashedPassword},
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed to create user: %v", err)
	}
	return userID, nil
}

func (db *UserDBClient) GetUserByUsername(ctx context.Context, username string) (map[string]types.AttributeValue, error) {
	output, err := db.Client.Scan(ctx, &dynamodb.ScanInput{
		TableName:        aws.String(db.Table),
		FilterExpression: aws.String("Username = :u"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":u": &types.AttributeValueMemberS{Value: username},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to scan table: %v", err)
	}
	if len(output.Items) == 0 {
		return nil, fmt.Errorf("user not found")
	}
	return output.Items[0], nil
}
