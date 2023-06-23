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

// GetBusUnit checks if the DynamoDB Table is configured on the environment, and
// fetch and returns the bus line unit information.
func GetBusUnit(ctx context.Context, code, busId string) (schema.BusUnit, error) {
	var (
		unit      schema.BusUnit
		tablename = env.BUS_UNIT
	)

	// Check if the DynamoDB Table is configured
	if tablename == "" {
		trail.Error("dynamodb BUS_UNIT_TABLE is not configured on the environment")
		err := errors.New("dynamodb BUS_UNIT_TABLE environment variable is not set")

		return unit, err
	}

	// Create a composite key expression
	key := expression.KeyAnd(expression.Key("code").Equal(expression.Value(code)), expression.Key("bus_id").Equal(expression.Value(busId)))

	// Create a names list representing the list of item attribute names
	// to be returned.
	var namesList = []expression.NameBuilder{
		expression.Name("bus_id"),
		expression.Name("active"),
		expression.Name("min_capacity"),
		expression.Name("max_capacity"),
	}

	// SELECT code, bus_id, active, min_capacity, max_capacity
	projection := expression.NamesList(expression.Name("code"), namesList...)

	// Build an expression to retrieve the item from the DynamoDB
	expr, err := expression.NewBuilder().WithKeyCondition(key).WithProjection(projection).Build()
	if err != nil {
		return unit, err
	}

	// Build the query params parameter
	params := &dynamodb.QueryInput{
		TableName:                 aws.String(tablename),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ProjectionExpression:      expr.Projection(),
		KeyConditionExpression:    expr.KeyCondition(),
	}

	result, err := awswrapper.DynamoDBQuery(ctx, params)
	if err != nil {
		return unit, err
	}

	// Unmarshal a map into actual use which front-end can understand as a JSON
	if result.Count > 0 {
		err := awswrapper.DynamoDBUnmarshalMap(&unit, result.Items[0])
		if err != nil {
			return unit, err
		}
	}

	return unit, nil
}

// CreateBusUnit checks if the DynamoDB Table is configured on the environment, and
// creates a new bus unit record/information.
func CreateBusUnit(ctx context.Context, data interface{}) error {
	var tablename = env.BUS_UNIT

	// Check if the DynamoDB Table is configured
	if tablename == "" {
		trail.Error("dynamodb BUS_UNIT_TABLE is not configured on the environment")
		err := errors.New("dynamodb BUS_UNIT_TABLE environment variable is not set")

		return err
	}

	// Save the Bus Unit information into the DynamoDB Table
	err := InsertItem(ctx, tablename, data)
	if err != nil {
		trail.Error("failed to insert bus unit")
		return err
	}

	return nil
}

// UpdateBusUnit checks if the DynamoDB Table is configured on the environment and
// updates the bus unit's information or record.
func UpdateBusUnit(ctx context.Context, key map[string]types.AttributeValue, update expression.UpdateBuilder) (schema.BusUnit, error) {
	var (
		unit      schema.BusUnit
		tablename = env.BUS_UNIT
	)

	// Check if the DynamoDB Table is configured
	if tablename == "" {
		trail.Error("dynamodb BUS_UNI_TABLE is not configured on the environment")
		err := errors.New("dynamodb BUS_UNIT_TABLE environment variable is not set")

		return unit, err
	}

	result, err := UpdateItem(ctx, tablename, key, update)
	if err != nil {
		trail.Error("failed to update the bus unit record")
		return unit, err
	}

	// Unmarshal a map into actual bus struct which the front-end can
	// understand as a JSON
	err = awswrapper.DynamoDBUnmarshalMap(&unit, result.Attributes)
	if err != nil {
		return unit, err
	}

	return unit, nil
}

// FilterBusUnit checks if the DynamoDB Table is configured on the environment,
// fetches and returns a list of bus unit information.
func FilterBusUnit(ctx context.Context, code, busId string, active *bool) ([]schema.BusUnit, error) {
	var (
		unitList  []schema.BusUnit
		tablename = env.BUS_UNIT
	)

	// Check if the DynamoDB Table is configured
	if tablename == "" {
		trail.Error("dynamodb BUS_UNIT_TABLE is not configured on the environment")
		err := errors.New("dynamodb BUS_UNIT_TABLE environment variable is not set")

		return unitList, err
	}

	// Construct the filter builder with a name that contains a specified value.
	var (
		filter          expression.ConditionBuilder
		busIdExpression = expression.Name("bus_id").Equal(expression.Value(busId))
	)

	switch {
	case code != "":
		filter = busIdExpression.And(expression.Name("code").Equal(expression.Value(code)))

	case active != nil:
		filter = busIdExpression.And(expression.Name("active").Equal(expression.Value(active)))

	case code != "" && active != nil:
		filter = busIdExpression.And(expression.Name("active").Equal(expression.Value(active)), expression.Name("code").Equal(expression.Value(code)))

	default:
		filter = busIdExpression
	}

	result, err := FilterItems(ctx, tablename, filter)
	if err != nil {
		return unitList, err
	}

	if result.Count > 0 {
		// Unmarshal a map into actual bus unit struct which the front-end can
		// understand as a JSON.
		err = awswrapper.DynamoDBUnmarshalListOfMaps(&unitList, result.Items)
		if err != nil {
			return unitList, err
		}
	}

	return unitList, nil
}
