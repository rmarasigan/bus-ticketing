package query

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/rmarasigan/bus-ticketing/api/schema"
	"github.com/rmarasigan/bus-ticketing/internal/app/env"
	awswrapper "github.com/rmarasigan/bus-ticketing/internal/aws_wrapper"
	"github.com/rmarasigan/bus-ticketing/internal/trail"
)

// GetBusRoute checks if the DynamoDB Table is configured on the environment, and
// fetch and returns the bus unit route information.
func GetBusRoute(ctx context.Context, id, busId string) (schema.BusRoute, error) {
	var (
		route     schema.BusRoute
		tablename = env.BUS_ROUTE_TABLE
	)

	// Check if the DynamoDB Table is configured
	if tablename == "" {
		trail.Error("dynamodb BUS_ROUTE_TABLE is not configured on the environment")
		err := errors.New("dynamodb BUS_ROUTE_TABLE environment variable is not set")

		return route, err
	}

	// Create a composite key expression
	key := expression.KeyAnd(expression.Key("id").Equal(expression.Value(id)),
		expression.Key("bus_id").Equal(expression.Value(busId)))

	// Create a names list representing the properties of the bus route
	// that is going to be returned
	var namesList = []expression.NameBuilder{
		expression.Name("bus_id"),
		expression.Name("bus_unit_id"),
		expression.Name("currency_code"),
		expression.Name("rate"),
		expression.Name("available"),
		expression.Name("departure_time"),
		expression.Name("arrival_time"),
		expression.Name("from_route"),
		expression.Name("to_route"),
	}

	// SELECT id, bus_id, bus_unit_id, currency_code, rate, available,
	// departure_time, arrival_time, from_route, to_route
	projection := expression.NamesList(expression.Name("id"), namesList...)

	// Build an expression to retrieve the item from the DynamoDB
	expr, err := expression.NewBuilder().WithKeyCondition(key).WithProjection(projection).Build()
	if err != nil {
		return route, err
	}

	// Build the query params
	params := &dynamodb.QueryInput{
		TableName:                 aws.String(tablename),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ProjectionExpression:      expr.Projection(),
		KeyConditionExpression:    expr.KeyCondition(),
	}

	result, err := awswrapper.DynamoDBQuery(ctx, params)
	if err != nil {
		return route, err
	}

	// Unmarshal a map into actual Bus Route struct which front-end can
	// understand as a JSON
	if result.Count > 0 {
		err := awswrapper.DynamoDBUnmarshalMap(&route, result.Items[0])
		if err != nil {
			return route, err
		}
	}

	return route, nil
}

// CreateBusRoute checks if the DynamoDB Table is configured on the environment, and
// creates a new bus route record.
func CreateBusRoute(ctx context.Context, data interface{}) error {
	var tablename = env.BUS_ROUTE_TABLE

	// Check if the DynamoDB Table is configured
	if tablename == "" {
		trail.Error("dynamodb BUS_ROUTE_TABLE is not configured on the environment")
		err := errors.New("dynamodb BUS_ROUTE_TABLE environment variable is not set")

		return err
	}

	// Save the Bus Route information into the DynamoDB Table
	err := InsertItem(ctx, tablename, data)
	if err != nil {
		trail.Error("failed to insert a new bus route")
		return err
	}

	return nil
}

// FilterBusRoute checks if the DynamoDB Table is configured on the environment, fetches
// and returns a list of bus routes information.
func FilterBusRoute(ctx context.Context, route schema.BusRouteFilter) ([]schema.BusRoute, error) {
	var (
		routes     []schema.BusRoute
		tablename  = env.BUS_ROUTE_TABLE
		filter     expression.ConditionBuilder
		filterList []expression.ConditionBuilder
	)

	// Check if the DynamoDB Table is configured
	if tablename == "" {
		trail.Error("dynamodb BUS_ROUTE_TABLE is not configured on the environment")
		err := errors.New("dynamodb BUS_ROUTE_TABLE environment is not set")

		return routes, err
	}

	trail.Debug("FilterBusRoute:  %v", route)

	// Construct the filter builder with a name that contains a specified value.
	if route.Available != nil {
		filterList = append(filterList, expression.Name("available").Equal(expression.Value(route.Available)))
	}

	if route.Departure != "" {
		filterList = append(filterList, expression.Name("departure_time").Equal(expression.Value(route.Departure)))
	}

	if route.Arrival != "" {
		filterList = append(filterList, expression.Name("arrival_time").Equal(expression.Value(route.Arrival)))
	}

	if route.FromRoute != "" {
		filterList = append(filterList, expression.Name("from_route").Equal(expression.Value(route.FromRoute)))
	}

	if route.ToRoute != "" {
		filterList = append(filterList, expression.Name("to_route").Equal(expression.Value(route.ToRoute)))
	}

	filter = expression.Name("bus_id").Equal(expression.Value(route.BusID)).And(expression.Name("bus_unit_id").Equal(expression.Value(route.BusUnitID)), filterList...)
	result, err := FilterItems(ctx, tablename, filter)
	if err != nil {
		return routes, err
	}

	if result.Count > 0 {
		// Unmarshal a map into actual bus route struct which the front-end can
		// understand as a JSON.
		err = awswrapper.DynamoDBUnmarshalListOfMaps(&routes, result.Items)
		if err != nil {
			return routes, err
		}
	}

	return routes, nil
}
