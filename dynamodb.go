package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/google/uuid"
)

type DynamoDBRepository struct {
	client    *dynamodb.Client
	tableName string
}

// NewDynamoDBRepository crea una nueva instancia del repository
func NewDynamoDBRepository(tableName, region string) (*DynamoDBRepository, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("error loading AWS config: %w", err)
	}

	client := dynamodb.NewFromConfig(cfg)

	return &DynamoDBRepository{
		client:    client,
		tableName: tableName,
	}, nil
}

// CreateOrder guarda un nuevo pedido en DynamoDB
func (r *DynamoDBRepository) CreateOrder(orderName, userName string) (*Order, error) {
	order := Order{
		OrderID:   uuid.New().String(),
		OrderName: orderName,
		UserName:  userName,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}

	item, err := attributevalue.MarshalMap(order)
	if err != nil {
		return nil, fmt.Errorf("error marshaling order: %w", err)
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	}

	_, err = r.client.PutItem(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("error putting item in DynamoDB: %w", err)
	}

	log.Printf("Order created successfully: %s", order.OrderID)
	return &order, nil
}

