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
	var namesList = []expression.NameBuilder{
		expression.Name("name"),
		expression.Name("owner"),
		expression.Name("email"),
		expression.Name("address"),
		expression.Name("company"),
		expression.Name("mobile_number"),
	}

	// SELECT id, name, owner, email, address, company, mobile_number
	projection := expression.NamesList(expression.Name("id"), namesList...)

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

	// Save the Bus Line information into the DynamoDB Table
	err := InsertItem(ctx, tablename, data)
	if err != nil {
		trail.Error("failed to insert a new bus line")
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

	result, err := UpdateItem(ctx, tablename, key, update)
	if err != nil {
		trail.Error("failed to update the bus line record")
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

// FilterBusLine checks if the DynamoDB Table is configured on the environment,
// fetches and returns a list of bus line information.
func FilterBusLine(ctx context.Context, name, company string) ([]schema.Bus, error) {
	var (
		busList   []schema.Bus
		tablename = env.BUS_TABLE
	)

	// Check if the DynamoDB Table is configured
	if tablename == "" {
		trail.Error("dynamodb BUS_TABLE is not configured on the environment")
		err := errors.New("dynamodb BUS_TABLE environment variable is not set")

		return busList, err
	}

	// Construct the filter builder with a name that contains a specified value.
	// WHERE name LIKE %name_value% OR company LIKE or %company_value%
	var filter expression.ConditionBuilder
	if name != "" && company != "" {
		filter = expression.Name("name").Contains(name).Or(expression.Name("company").Contains(company))
	} else {
		if name != "" {
			filter = expression.Name("name").Contains(name)
		}

		if company != "" {
			filter = expression.Name("company").Contains(company)
		}
	}

	result, err := FilterItems(ctx, tablename, filter)
	if err != nil {
		return busList, err
	}

	if result.Count > 0 {
		// Unmarshal a map into actual bus struct which the front-end can
		// understand as a JSON.
		err = awswrapper.DynamoDBUnmarshalListOfMaps(&busList, result.Items)
		if err != nil {
			return busList, err
		}
	}

	return busList, nil
}
