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

type DynamoDBClient struct {
	Client *dynamodb.Client
	Table  string
}

func NewDynamoDBClient(ctx context.Context, tableName string) (*DynamoDBClient, error) {
	cfg, err := config.LoadDefaultConfig((ctx))

	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %v", err)

	}

	client := dynamodb.NewFromConfig(cfg)

	return &DynamoDBClient{
		Client: client,
		Table:  tableName,
	}, nil
}

func (db *DynamoDBClient) CreateDocument(ctx context.Context, title, content string) (string, error) {

	documentID := uuid.New().String()
	_, err := db.Client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(db.Table),
		Item: map[string]types.AttributeValue{
			"DocumentId": &types.AttributeValueMemberS{Value: documentID},
			"Title":      &types.AttributeValueMemberS{Value: title},
			"Content":    &types.AttributeValueMemberS{Value: content},
		},
	})

	if err != nil {
		return "", fmt.Errorf("failed to put item: %v", err)
	}

	return documentID, nil
}

func (db *DynamoDBClient) UpdateDocument(ctx context.Context, documentID, title, content string) error {

	updateExpr := "SET #c = :c"
	exprAttrNames := map[string]string{
		"#c": "Content",
	}
	exprAttrValues := map[string]types.AttributeValue{
		":c": &types.AttributeValueMemberS{Value: content},
	}

	// If title is provided, include it in the update expression
	if title != "" {
		updateExpr += ", #t = :t"
		exprAttrNames["#t"] = "Title"
		exprAttrValues[":t"] = &types.AttributeValueMemberS{Value: title}
	}

	_, err := db.Client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(db.Table),
		Key: map[string]types.AttributeValue{
			"DocumentId": &types.AttributeValueMemberS{Value: documentID},
		},
		UpdateExpression:          aws.String(updateExpr),
		ExpressionAttributeNames:  exprAttrNames,
		ExpressionAttributeValues: exprAttrValues,
	})

	if err != nil {
		return fmt.Errorf("failed to update item: %v", err)
	}
	return nil
}

func (db *DynamoDBClient) GetDocument(ctx context.Context, documentID string) (map[string]types.AttributeValue, error) {
	output, err := db.Client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(db.Table),
		Key: map[string]types.AttributeValue{
			"DocumentId": &types.AttributeValueMemberS{Value: documentID},
		},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get item: %v", err)

	}

	if output.Item == nil {
		return nil, fmt.Errorf("document not found")
	}

	return output.Item, nil
}
