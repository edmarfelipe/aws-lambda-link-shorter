package storage_test

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/edmarfelipe/aws-lambda/storage"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestLinkStorage(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	cfg := setup(context.Background(), t)
	linkStorage := storage.NewLinkStorage(cfg)

	err := linkStorage.CreateTable(context.Background())
	assert.NoError(t, err)

	link := storage.Link{
		Original: "https://www.google.com",
		Title:    "Google",
		Hash:     "123",
	}

	err = linkStorage.Create(context.Background(), link)
	assert.NoError(t, err)

	result, err := linkStorage.GetLinkByHash(context.Background(), "123")
	assert.NoError(t, err)
	assert.Equal(t, link, result)
}

func setup(ctx context.Context, t *testing.T) aws.Config {
	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "amazon/dynamodb-local:latest",
			ExposedPorts: []string{"8000/tcp"},
			Cmd:          []string{"-jar", "DynamoDBLocal.jar", "-sharedDb"},
			WaitingFor:   wait.NewHostPortStrategy("8000"),
		},
		Started: true,
	})
	if err != nil {
		t.Error(err)
	}

	t.Cleanup(func() {
		if err := c.Terminate(context.Background()); err != nil {
			t.Fatalf("failed to terminate container: %s", err.Error())
		}
	})

	url, err := c.Endpoint(ctx, "http")
	if err != nil {
		t.Error(err)
	}

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion("us-east-1"),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{URL: url}, nil
			})),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     "dummy",
				SecretAccessKey: "dummy",
				SessionToken:    "dummy",
			},
		}),
	)
	if err != nil {
		t.Error(err)
	}
	return cfg
}
