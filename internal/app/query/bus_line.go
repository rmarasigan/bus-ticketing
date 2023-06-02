package query

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/rmarasigan/bus-ticketing/internal/app/env"
	awswrapper "github.com/rmarasigan/bus-ticketing/internal/aws_wrapper"
	"github.com/rmarasigan/bus-ticketing/internal/trail"
)

// CreateBusLine checks if the DynamoDB Table is configured on the environment, and
// creates a new bus line record.
func CreateBusLine(ctx context.Context, data interface{}) error {
	var tablename = env.BUS_TABLE

	// Check if the DynamoDB Table is configured
	if tablename == "" {
		trail.Error("dynamodb BUS_TABLE is not configured on the environment")
		err := errors.New("dynamodb BUS_TABLE environment variable is not set ")

		return err
	}

	// Marshal the user to a map of AttributeValeus
	values, err := awswrapper.DynamoDBMarshalMap(data)
	if err != nil {
		trail.Error("failed to marshal data to a map of AttributeValues")
		return err
	}

	params := &dynamodb.PutItemInput{
		Item:      values,
		TableName: aws.String(tablename),
	}

	// Save the Bus Line information into the DynamoDB Table
	_, err = awswrapper.DynamoDBPutItem(ctx, params)
	if err != nil {
		trail.Error("failed to insert a new user")
		return err
	}

	return nil
}
