package query

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/rmarasigan/bus-ticketing/api/schema"
	"github.com/rmarasigan/bus-ticketing/internal/app/env"
	awswrapper "github.com/rmarasigan/bus-ticketing/internal/aws_wrapper"
	"github.com/rmarasigan/bus-ticketing/internal/trail"
)

// GetBusLine checks if the DynamoDB Table is configured on the environment, and
// fetch and returns the bus line information.
func GetBusLine(ctx context.Context, id, name string) (schema.Bus, error) {
	var (
		bus       schema.Bus
		tablename = env.BUS_TABLE
	)

	// Check if the DynamoDB Table is configured
	if tablename == "" {
		trail.Error("dynamodb BUS_TABLE is not configured on the environment")
		err := errors.New("dynamodb BUS_TABLE environment variable is not set")

		return bus, err
	}

	// Create a primary key expression
	key := expression.Key("name").Equal(expression.Value(name))

	// Create a names list representing the list of item attribute names
	// to be returned.
	var namesListt = []expression.NameBuilder{
		expression.Name("name"),
		expression.Name("owner"),
		expression.Name("email"),
		expression.Name("address"),
		expression.Name("company"),
		expression.Name("mobile_number"),
	}

	// SELECT id, name, owner, email, address, company, mobile_number
	projection := expression.NamesList(expression.Name("id"), namesListt...)

	// Construct the filter builder with a name and value.
	// WHERE id = id_value
	filter := expression.Name("id").Equal(expression.Value(id))

	// Build an expression to retrieve the item from the DynamoDB
	expr, err := expression.NewBuilder().WithKeyCondition(key).WithFilter(filter).WithProjection(projection).Build()
	if err != nil {
		return bus, err
	}

	// Build the query params parameter
	params := &dynamodb.QueryInput{
		TableName:                 aws.String(tablename),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		KeyConditionExpression:    expr.KeyCondition(),
	}

	result, err := awswrapper.DynamoDBQuery(ctx, params)
	if err != nil {
		return bus, err
	}

	// Unmarshal a map into actual use which front-end can understand as a JSON
	if result.Count > 0 {
		err := awswrapper.DynamoDBUnmarshalMap(&bus, result.Items[0])
		if err != nil {
			return bus, err
		}
	}

	return bus, nil
}

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

// UpdateBusLine checks if the DynamoDB Table is configured on the environment and
// updates the bus line's information or record.
func UpdateBusLine(ctx context.Context, key map[string]types.AttributeValue, update expression.UpdateBuilder) (schema.Bus, error) {
	var (
		bus       schema.Bus
		tablename = env.BUS_TABLE
	)

	// Check if the DynamoDB Table is configured
	if tablename == "" {
		trail.Error("dynamodb BUS_TABLE is not configured on the environment")
		err := errors.New("dynamodb BUS_TABLE environment variable is not set")

		return bus, err
	}

	// Using the update expression to create a DynamoDB Expression
	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		trail.Error("failed to build a DynamoDB Expression")
		return bus, err
	}

	// Use the built expression to populate the DynamoDB Update Item API
	var params = &dynamodb.UpdateItemInput{
		Key:                       key,
		TableName:                 aws.String(tablename),
		UpdateExpression:          expr.Update(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ReturnValues:              types.ReturnValueAllNew,
	}

	result, err := awswrapper.DynamoDBUpdateItem(ctx, params)
	if err != nil {
		trail.Error("fail to update the bus line information")
		return bus, err
	}

	// Unmarshal a map into actual bus struct which the front-end can
	// understand as a JSON.
	err = awswrapper.DynamoDBUnmarshalMap(&bus, result.Attributes)
	if err != nil {
		return bus, err
	}

	return bus, nil
}
