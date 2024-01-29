package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws/ratelimit"
	"github.com/aws/aws-sdk-go-v2/aws/retry"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/spf13/cobra"
	"github.com/tilotech/tilores-cli/internal/pkg/step"
)

const (
	maxBatchWriteSize   = 25
	eraseTableRoutines  = 25
	eraseBucketRoutines = 30
)

// eraseCmd represents the erase command
var eraseCmd = &cobra.Command{
	Use:   "erase",
	Short: "Erases all data submitted to " + applicationName + " in your AWS account.",
	Long: `Erases all data submitted to ` + applicationName + ` in your AWS account.

Warning: This process generates costs, for large amounts of data we recommend redeploying (destroy then deploy) instead.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		fmt.Println(colorYellow, "Warning: This process generates costs, for large amounts of data we recommend to redeploying (destroy then deploy) instead", colorReset)
		fmt.Println(colorRed, "Are you sure you want to erase ALL DATA!!!", colorReset)
		fmt.Println("only \"yes\" will proceed")
		answer := ""
		_, _ = fmt.Scanln(&answer)
		if answer == "yes" {
			return
		}
		cobra.CheckErr(fmt.Errorf("no confirmation, exiting"))
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		cfg, err := config.LoadDefaultConfig(ctx, func(o *config.LoadOptions) error {
			o.Region = region
			o.SharedConfigProfile = profile
			o.Retryer = func() aws.Retryer {
				standard := retry.NewStandard(func(so *retry.StandardOptions) {
					so.RateLimiter = ratelimit.NewTokenRateLimit(5000)
				})
				r := retry.AddWithMaxAttempts(standard, 10)
				r = retry.AddWithErrorCodes(r,
					(*types.ProvisionedThroughputExceededException)(nil).ErrorCode(),
					string(types.BatchStatementErrorCodeEnumThrottlingError),
				)
				return r
			}
			return nil
		})
		cobra.CheckErr(err)
		ddbClient := dynamodb.NewFromConfig(cfg)
		s3Client := s3.NewFromConfig(cfg)
		tables, buckets, err := erasableResources()
		cobra.CheckErr(err)
		fmt.Println("Detected the following resources to erase:")
		fmt.Println(strings.Join(append(tables, buckets...), "\n"))
		start := time.Now()
		fmt.Println("Started...")
		err = eraseAll(ctx, ddbClient, s3Client, tables, buckets)
		fmt.Printf("Erase took %v\n", time.Since(start))
		cobra.CheckErr(err)
	},
	SilenceUsage: true,
}

func init() {
	rootCmd.AddCommand(eraseCmd)

	eraseCmd.Flags().StringVar(&region, "region", "", "The deployments AWS region.")
	_ = eraseCmd.MarkFlagRequired("region")

	eraseCmd.Flags().StringVar(&profile, "profile", "", "The AWS credentials profile.")

	eraseCmd.Flags().StringVar(&workspace, "workspace", "default", "The deployments workspace/environment e.g. dev, prod.")
}

func eraseAll(ctx context.Context, ddbClient *dynamodb.Client, s3Client *s3.Client, tables, buckets []string) error {
	errCh := make(chan error)
	tableReqCh := make(chan *dynamodb.BatchWriteItemInput)
	bucketReqCh := make(chan *s3.DeleteObjectsInput)
	go createEraseTableRequests(ctx, ddbClient, tables, tableReqCh, errCh)
	go createEraseBucketRequests(ctx, s3Client, buckets, bucketReqCh, errCh)
	wg := sync.WaitGroup{}
	for i := 0; i < eraseTableRoutines; i++ {
		wg.Add(1)
		go processEraseTableRequests(ctx, ddbClient, &wg, tableReqCh, errCh)
	}
	for i := 0; i < eraseBucketRoutines; i++ {
		wg.Add(1)
		go processEraseBucketRequests(ctx, s3Client, &wg, bucketReqCh, errCh)
	}
	go func() {
		wg.Wait()
		close(errCh)
	}()
	err := collectErrors(errCh)
	if err != nil {
		return err
	}
	return nil
}

func erasableResources() (tables []string, buckets []string, err error) {
	steps := []step.Step{
		step.TerraformRequire,
		step.Chdir("deployment/tilores"),
		step.TerraformInitFast,
		step.TerraformNewWorkspace(workspace),
		step.Chdir("../.."),
	}
	err = step.Execute(steps)
	if err != nil {
		return nil, nil, err
	}
	cmd := exec.Command("terraform", "-chdir=deployment/tilores", "show", "-json")
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "TF_WORKSPACE="+workspace)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, nil, err
	}
	currentState := struct {
		Values struct {
			RootModule struct {
				ChildModules []struct {
					Resources []struct {
						Address string
						Values  struct {
							Name   string
							Bucket string
						}
					}
				} `json:"child_modules"`
			} `json:"root_module"`
		}
	}{}
	err = json.Unmarshal(out, &currentState)
	if err != nil {
		return nil, nil, err
	}

	for _, module := range currentState.Values.RootModule.ChildModules {
		for _, resource := range module.Resources {
			if strings.HasPrefix(resource.Address, "module.tilores.aws_dynamodb_table.") {
				tables = append(tables, resource.Values.Name)
			}
			if strings.HasPrefix(resource.Address, "module.tilores.aws_s3_bucket.") {
				buckets = append(buckets, resource.Values.Bucket)
			}
		}
	}
	return tables, buckets, nil
}

func collectErrors(errCh <-chan error) error {
	errorsCount := 0
	var lastError error
	for err := range errCh {
		errorsCount++
		lastError = err
	}
	if errorsCount != 0 {
		return fmt.Errorf("%v error(s) occurred, last error was: %v", errorsCount, lastError)
	}
	return nil
}

func createEraseTableRequests(ctx context.Context, ddbClient *dynamodb.Client, tables []string, reqCh chan<- *dynamodb.BatchWriteItemInput, errCh chan<- error) {
	defer close(reqCh)
	for _, table := range tables {
		describeTableInput := &dynamodb.DescribeTableInput{
			TableName: aws.String(table),
		}

		describeTableOutput, err := ddbClient.DescribeTable(ctx, describeTableInput)
		if err != nil {
			errCh <- err
			continue
		}

		keyAttributes := make([]string, len(describeTableOutput.Table.KeySchema))
		for i, key := range describeTableOutput.Table.KeySchema {
			keyAttributes[i] = aws.ToString(key.AttributeName)
		}

		key := strings.Join(keyAttributes, ",")

		paginator := dynamodb.NewScanPaginator(ddbClient, &dynamodb.ScanInput{
			TableName:            aws.String(table),
			ProjectionExpression: &key,
		})

		for paginator.HasMorePages() {
			scanOutput, err := paginator.NextPage(ctx)
			if err != nil {
				errCh <- err
				break
			}
			scannedItems := scanOutput.Items

			writeRequests := make([]types.WriteRequest, len(scannedItems))

			for i := range scannedItems {
				writeRequest := types.WriteRequest{
					DeleteRequest: &types.DeleteRequest{
						Key: scannedItems[i],
					},
				}
				writeRequests[i] = writeRequest
			}
			chunks := int(math.Ceil(float64(len(writeRequests)) / float64(maxBatchWriteSize)))
			for i := 0; i < chunks; i++ {
				chunk := writeRequests[i*maxBatchWriteSize : int(math.Min(float64(i*maxBatchWriteSize+maxBatchWriteSize), float64(len(writeRequests))))]
				reqCh <- &dynamodb.BatchWriteItemInput{
					RequestItems: map[string][]types.WriteRequest{
						table: chunk,
					},
				}
			}
		}
	}
}

func processEraseTableRequests(ctx context.Context, ddbClient *dynamodb.Client, wg *sync.WaitGroup, reqCh chan *dynamodb.BatchWriteItemInput, errCh chan<- error) {
	defer wg.Done()
	for batchWriteItemInput := range reqCh {
		processEraseTableRequest(ctx, ddbClient, batchWriteItemInput, errCh)
	}
}

func processEraseTableRequest(ctx context.Context, ddbClient *dynamodb.Client, req *dynamodb.BatchWriteItemInput, errCh chan<- error) {
	batchWriteItemOutput, err := ddbClient.BatchWriteItem(ctx, req)
	if err != nil {
		errCh <- err
		return
	}
	// recall with unprocessed items
	if len(batchWriteItemOutput.UnprocessedItems) > 0 {
		processEraseTableRequest(ctx, ddbClient, &dynamodb.BatchWriteItemInput{
			RequestItems: batchWriteItemOutput.UnprocessedItems,
		}, errCh)
	}
}

func createEraseBucketRequests(ctx context.Context, s3Client *s3.Client, buckets []string, reqCh chan<- *s3.DeleteObjectsInput, errCh chan<- error) {
	defer close(reqCh)
	for _, bucket := range buckets {
		paginator := s3.NewListObjectsV2Paginator(s3Client, &s3.ListObjectsV2Input{
			Bucket: aws.String(bucket),
		})
		for paginator.HasMorePages() {
			listOutput, err := paginator.NextPage(ctx)
			if err != nil {
				errCh <- err
				break
			}
			listedObjects := listOutput.Contents
			if len(listedObjects) == 0 {
				continue
			}

			objectIdentifiers := make([]s3types.ObjectIdentifier, len(listedObjects))
			for i, object := range listedObjects {
				objectIdentifiers[i].Key = object.Key
			}
			reqCh <- &s3.DeleteObjectsInput{
				Bucket: aws.String(bucket),
				Delete: &s3types.Delete{
					Objects: objectIdentifiers,
					Quiet:   aws.Bool(true),
				},
			}
		}
	}
}

func processEraseBucketRequests(ctx context.Context, s3Client *s3.Client, wg *sync.WaitGroup, reqCh <-chan *s3.DeleteObjectsInput, errCh chan<- error) {
	defer wg.Done()
	for deleteObjectsInput := range reqCh {
		deleteObjectsOutput, err := s3Client.DeleteObjects(ctx, deleteObjectsInput)
		if err != nil {
			errCh <- err
			return
		}
		for _, s3err := range deleteObjectsOutput.Errors {
			errCh <- fmt.Errorf("object %v, code %v, message %v", s3err.Key, s3err.Code, s3err.Message)
		}
	}
}
