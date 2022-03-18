package cmd

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
)

func TestMain(m *testing.M) {
	// set fake default region for tests
	os.Setenv("AWS_REGION", "eu-west-1")

	flag.Parse()
	if testing.Short() {
		code := m.Run()
		os.Exit(code)
	}

	pool, resource := startLocalStack()
	err := waitForServicesAndInit(pool)
	if err != nil {
		stopLocalStack(pool, resource)
		log.Fatalf("Could not init services: %s", err)
	}
	code := m.Run()
	stopLocalStack(pool, resource)
	os.Exit(code)
}

func startLocalStack() (*dockertest.Pool, *dockertest.Resource) {
	host := os.Getenv("LOCALSTACK_HOSTNAME_EXTERNAL")
	if host != "" {
		return nil, nil
	}

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "localstack/localstack",
		Tag:        "0.13",
		PortBindings: map[docker.Port][]docker.PortBinding{
			"4566/tcp": {{HostPort: "4566"}}, // Port for all services
		},
		Env: []string{
			"SERVICES=dynamodb",
			"ENABLE_CONFIG_UPDATES=1",
		},
	})
	if err != nil {
		log.Fatalf("Cloud not start resource: %s", err)
	}
	return pool, resource
}

func waitForServicesAndInit(pool *dockertest.Pool) error {
	var host string
	var retry func(op func() error) error
	if pool == nil {
		host = os.Getenv("LOCALSTACK_HOSTNAME_EXTERNAL")
		retry = func(op func() error) error {
			i := 0
			for {
				err := op()
				if err == nil || i > 30 {
					return err
				}
				i++
				time.Sleep(1 * time.Second)
			}
		}
	} else {
		host = "localhost"
		retry = pool.Retry
	}

	os.Setenv("LOCALSTACK_HOST", host)

	// wait for DynamoDB to become available
	err := retry(func() error {
		_, err := createTestTable("pulse")
		return err
	})

	if err != nil {
		return err
	}
	return nil
}

func stopLocalStack(pool *dockertest.Pool, resource *dockertest.Resource) {
	if pool == nil || resource == nil {
		return
	}
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
}

func newTestConfig(ctx context.Context) (aws.Config, error) {
	host := os.Getenv("LOCALSTACK_HOST")
	endpointResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:           fmt.Sprintf("http://%v:4566", host),
			SigningRegion: region,
		}, nil
	})
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("not", "required", "")),
		config.WithEndpointResolverWithOptions(endpointResolver),
	)
	if err != nil {
		return aws.Config{}, err
	}
	return cfg, nil
}

func createTestTable(tableName string) (*dynamodb.CreateTableOutput, error) {
	ctx := context.Background()
	cfg, err := newTestConfig(ctx)
	if err != nil {
		return nil, err
	}
	ddbClient := dynamodb.NewFromConfig(cfg)

	attributeDefinitions := []types.AttributeDefinition{
		{
			AttributeName: aws.String("id"),
			AttributeType: types.ScalarAttributeTypeS,
		},
	}

	keySchema := []types.KeySchemaElement{
		{
			AttributeName: aws.String("id"),
			KeyType:       types.KeyTypeHash,
		},
	}

	return ddbClient.CreateTable(ctx, &dynamodb.CreateTableInput{
		TableName:            aws.String(tableName),
		AttributeDefinitions: attributeDefinitions,
		KeySchema:            keySchema,
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(1000),
			WriteCapacityUnits: aws.Int64(10000),
		},
	})
}
