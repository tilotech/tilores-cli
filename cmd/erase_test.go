package cmd

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEraseTable(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}
	table := "records"
	_, err := createTestTable(table)
	require.NoError(t, err)

	ctx := context.Background()
	cfg, err := newTestConfig(ctx)
	require.NoError(t, err)
	ddbClient := dynamodb.NewFromConfig(cfg)

	start := time.Now()
	err = fillTable(ctx, ddbClient, table, 500)
	require.NoError(t, err)
	fmt.Println("fill table took", time.Since(start))

	recordsTable, err := ddbClient.DescribeTable(ctx, &dynamodb.DescribeTableInput{TableName: aws.String(table)})
	require.NoError(t, err)
	require.Equal(t, aws.Int64(500), recordsTable.Table.ItemCount)

	err = eraseAll(ctx, ddbClient, nil, []string{table}, nil)
	assert.NoError(t, err)
	fmt.Printf("Finished, erase took %v\n", time.Since(start))

	recordsTable, err = ddbClient.DescribeTable(ctx, &dynamodb.DescribeTableInput{TableName: aws.String(table)})
	require.NoError(t, err)
	assert.Equal(t, aws.Int64(0), recordsTable.Table.ItemCount)
}

func fillTable(ctx context.Context, ddbClient *dynamodb.Client, table string, itemCount int) error {
	for i := 0; i < itemCount; i++ {
		_, err := ddbClient.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(table),
			Item: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{Value: strconv.Itoa(i)},
			},
		})
		if err != nil {
			return err
		}
	}
	return nil
}
