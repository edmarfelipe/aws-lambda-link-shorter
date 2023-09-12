package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var (
	ErrLinkNotFound = errors.New("no Link found")
)

type Link struct {
	Hash     string `dynamodbav:"hash"`
	Title    string `dynamodbav:"title"`
	Original string `dynamodbav:"original"`
}

type LinkStorage interface {
	Create(ctx context.Context, link Link) error
	GetLinkByHash(ctx context.Context, hash string) (Link, error)
	CreateTable(ctx context.Context) error
}

const tableName = "links"

type linkStorage struct {
	client *dynamodb.Client
}

func NewLinkStorage(cfg aws.Config) LinkStorage {
	return &linkStorage{
		client: dynamodb.NewFromConfig(cfg),
	}
}

func (lks *linkStorage) Create(ctx context.Context, link Link) error {
	av, err := attributevalue.MarshalMap(link)
	if err != nil {
		return fmt.Errorf("failed to marshal link, %w", err)
	}

	_, err = lks.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      av,
	})
	if err != nil {
		return fmt.Errorf("failed to put link, %w", err)
	}
	return nil
}

func (lks *linkStorage) GetLinkByHash(ctx context.Context, hash string) (Link, error) {
	result, err := lks.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"hash": &types.AttributeValueMemberS{
				Value: hash,
			},
		},
	})
	if err != nil {
		return Link{}, fmt.Errorf("failed to scan Link, %w", err)
	}
	if result.Item == nil {
		return Link{}, ErrLinkNotFound
	}

	var link Link
	err = attributevalue.UnmarshalMap(result.Item, &link)
	if err != nil {
		return Link{}, fmt.Errorf("failed to unmarshal Link, %w", err)
	}
	return link, nil
}

func (lks *linkStorage) CreateTable(ctx context.Context) error {
	_, err := lks.client.CreateTable(context.Background(), &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("hash"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("hash"),
				KeyType:       types.KeyTypeHash,
			},
		},
		BillingMode: types.BillingModePayPerRequest,
		TableName:   aws.String(tableName),
	})
	return err
}
