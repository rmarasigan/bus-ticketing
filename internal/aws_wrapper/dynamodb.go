package awswrapper

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/rmarasigan/bus-ticketing/internal/utility"
)

var (
	dynamoClient *dynamodb.Client
)

// initDynamoClient initializes the DynamoDB client from the
// prpovided configuration.
func initDynamoClient(ctx context.Context) {
	if dynamoClient != nil {
		return
	}

	// Using the SDK's default configuration, loading additional config
	// and credentials values from the environment variables, shared
	// credentials, and shared configuration files
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(AWS_REGION))
	if err != nil {
		utility.Error(err, "DynamoClientError", "failed to load the default config")
		return
	}

	// Using the cfg value to create the DynamoDB client
	dynamoClient = dynamodb.NewFromConfig(cfg)
}

// DynamoDBScan initializes the DynamoDB Client and reads every item in a table or a secondary index.
// It returns one or more items and item attributes.
func DynamoDBScan(ctx context.Context, params *dynamodb.ScanInput) (*dynamodb.ScanOutput, error) {
	// Initialize the DynamoClient
	initDynamoClient(ctx)

	output, err := dynamoClient.Scan(ctx, params)
	if err != nil {
		return nil, err
	}

	return output, nil
}

// DynamoDBQuery initializes the DynamoDB Client and finds items based on primary key values. It
// returns all items with that partition key value.
func DynamoDBQuery(ctx context.Context, params *dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
	// Initialize the DynamoClient
	initDynamoClient(ctx)

	output, err := dynamoClient.Query(ctx, params)
	if err != nil {
		return nil, err
	}

	return output, nil
}

// DynamoDBPutItem initializes the DynamoDB Client and creates a new item, or replaces an old item with a new item.
func DynamoDBPutItem(ctx context.Context, params *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	// Initialize the DynamoClient
	initDynamoClient(ctx)

	output, err := dynamoClient.PutItem(ctx, params)
	if err != nil {
		return nil, err
	}

	return output, nil
}

// DynamoDBUpdateItem initializes the DynamoDB Client and edits an existing item's attributes, or adds a new item to
// the table if it does not already exist.
func DynamoDBUpdateItem(ctx context.Context, params *dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error) {
	// Initialize the DynamoClient
	initDynamoClient(ctx)

	output, err := dynamoClient.UpdateItem(ctx, params)
	if err != nil {
		return nil, err
	}

	return output, nil
}

// DynamoDBMarshalMap marshals Go value type to a map of AttributeValues.
func DynamoDBMarshalMap(v interface{}) (map[string]types.AttributeValue, error) {
	output, err := attributevalue.MarshalMap(v)
	if err != nil {
		return nil, err
	}

	return output, nil
}

// DynamoDBUnmarshalMap returns an item with the values of result object.
func DynamoDBUnmarshalMap(v interface{}, result map[string]types.AttributeValue) error {
	err := attributevalue.UnmarshalMap(result, v)
	if err != nil {
		return err
	}
	return nil
}

// DynamoDBUnmarshalListOfMaps returns a list of items with the values of result object.
func DynamoDBUnmarshalListOfMaps(v interface{}, result []map[string]types.AttributeValue) error {
	err := attributevalue.UnmarshalListOfMaps(result, v)
	if err != nil {
		return err
	}

	return nil
}
