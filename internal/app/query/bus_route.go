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

// getBusRoute returns the bus unit route information.
func getBusRoute(ctx context.Context, tablename, id, busId string) (schema.BusRoute, error) {
	var route schema.BusRoute

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
		expression.Name("active"),
		expression.Name("departure_time"),
		expression.Name("arrival_time"),
		expression.Name("from_route"),
		expression.Name("to_route"),
	}

	// SELECT id, bus_id, bus_unit_id, currency_code, rate, active,
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

// GetBusRouteRecords checks if the DynamoDB Table is configured on the environment, and
// returns either the specific bus line route or a list of bus line route records.
func GetBusRouteRecords(ctx context.Context, id, busId string) ([]schema.BusRoute, error) {
	var (
		routes    []schema.BusRoute
		tablename = env.BUS_ROUTE_TABLE
	)

	// Check if the DynamoDB Table is configured
	if tablename == "" {
		trail.Error("dynamodb BUS_ROUTE_TABLE is not configured on the environment")
		err := errors.New("dynamodb BUS_ROUTE_TABLE environment variable is not set")

		return routes, err
	}

	// ********** Fetching a specific bus route record ********** //
	if id != "" && busId != "" {
		route, err := getBusRoute(ctx, tablename, id, busId)
		if err != nil {
			return routes, err
		}

		if route == (schema.BusRoute{}) {
			return routes, nil
		}

		routes = append(routes, route)
		return routes, nil
	}

	// **************** List of bus route records **************** //
	// Create a names list representing the properties of the bus route
	// that is going to be returned
	var namesList = []expression.NameBuilder{
		expression.Name("bus_id"),
		expression.Name("bus_unit_id"),
		expression.Name("currency_code"),
		expression.Name("rate"),
		expression.Name("active"),
		expression.Name("departure_time"),
		expression.Name("arrival_time"),
		expression.Name("from_route"),
		expression.Name("to_route"),
	}

	// SELECT id, bus_id, bus_unit_id, currency_code, rate, active,
	// departure_time, arrival_time, from_route, to_route
	projection := expression.NamesList(expression.Name("id"), namesList...)

	// Build an expression to retrieve the item from the DynamoDB
	expr, err := expression.NewBuilder().WithProjection(projection).Build()
	if err != nil {
		return routes, err
	}

	// Use the build expression to populate the DynamoDB Scan API
	var params = &dynamodb.ScanInput{
		TableName:                 aws.String(tablename),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ProjectionExpression:      expr.Projection(),
	}

	result, err := awswrapper.DynamoDBScan(ctx, params)
	if err != nil {
		return nil, err
	}

	if result.Count > 0 {
		// Unmarshal a map into actual bus route struct which the front-end can
		// understand as a JSON.
		err = awswrapper.DynamoDBUnmarshalListOfMaps(&routes, result.Items)
		if err != nil {
			return nil, err
		}
	}

	return routes, nil
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

// UpdateBusRoute checks if the DynamoDB Table is configured on the environment and
// updates the bus routes information or record.
func UpdateBusRoute(ctx context.Context, key map[string]types.AttributeValue, update expression.UpdateBuilder) (schema.BusRoute, error) {
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

	result, err := UpdateItem(ctx, tablename, key, update)
	if err != nil {
		trail.Error("failed to update the bus route record")
		return route, err
	}

	// Unmarshal a map into actual bus route struct which the front-end
	// can understan as a JSON
	err = awswrapper.DynamoDBUnmarshalMap(&route, result.Attributes)
	if err != nil {
		return route, err
	}

	return route, nil
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

	// Construct the filter builder with a name that contains a specified value.
	if route.Active != nil {
		filterList = append(filterList, expression.Name("active").Equal(expression.Value(route.Active)))
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

	if route.BusUnitID != "" {
		filter = expression.Name("bus_id").Equal(expression.Value(route.BusID)).And(expression.Name("bus_unit_id").Equal(expression.Value(route.BusUnitID)), filterList...)

	} else {
		if len(filterList) > 0 {
			filter = expression.And(expression.Name("bus_id").Equal(expression.Value(route.BusID)), filterList[0], filterList[1:]...)

		} else {
			filter = expression.Name("bus_id").Equal(expression.Value(route.BusID))
		}
	}

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
