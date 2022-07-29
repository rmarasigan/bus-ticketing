package busroute

import (
	"context"
	"errors"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/rmarasigan/bus-ticketing/pkg/api"
	"github.com/rmarasigan/bus-ticketing/pkg/cw/kvp"
	cw "github.com/rmarasigan/bus-ticketing/pkg/cw/logger"
	"github.com/rmarasigan/bus-ticketing/pkg/models"
	"github.com/rmarasigan/bus-ticketing/pkg/service"
)

var (
	svc dynamodbiface.DynamoDBAPI
)

// Get is the Bus Route API request GET method that will process the incoming request and returns an API Gateway response.
//
// Endpoint:
//  https://{api-id}.execute.api.{region}.amazonaws.com/{stage}/bus/route/
//
// Query Parameters accepted:
//  to: indicating the destination of the bus
//  from: indicating the starting point of a bus
//  departure: expected arrival time on the destination and is in 24-hour format
//  arrival: expected arrival time on the destination and is in 24-hour format
func Get(ctx context.Context, request *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	// Initialize Dynamodb Session
	service.DynamodbSession()
	svc = service.DynamoDBClient

	tablename := os.Getenv("BUS_ROUTE_TABLE")
	queryToRoute := request.QueryStringParameters["to"]
	queryFromRoute := request.QueryStringParameters["from"]
	queryDeparture := request.QueryStringParameters["departure"]
	queryArrival := request.QueryStringParameters["arrival"]

	if queryToRoute != "" || queryFromRoute != "" {
		return BusRoute(tablename, queryFromRoute, queryToRoute)
	}

	if queryDeparture != "" && queryArrival != "" {
		return BusRouteSchedule(tablename, queryDeparture, queryArrival)
	}

	return ListBusRoute(tablename)
}

// BusRouteInfo returns an API Gateway response of a specific bus route information.
func BusRouteInfo(tablename string, id string) (*models.BusRoute, error) {
	busRoute := new(models.BusRoute)

	// Construct the key condition builder with value.
	// WHERE id = id_value
	key := expression.Key("id").Equal(expression.Value(id))

	// Using the key condition to create a DynamoDB expression.
	expr, err := expression.NewBuilder().WithKeyCondition(key).Build()
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "DynamoDBExpression", Message: "Failed to build dynamodb expression."})
		return nil, err
	}

	// Use the built expression to populate the DynamoDB Query's API input parameters.
	input := &dynamodb.QueryInput{
		TableName:                 aws.String(tablename),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
	}

	// Returns one or more items.
	result, err := svc.Query(input)
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "DynamoDBQuery", Message: "Failed to query input."}, kvp.Attribute{Key: "tablename", Value: tablename})
		return nil, err
	}

	// Check if there are items returned.
	if len(result.Items) == 0 {
		err := errors.New("bus route information not found")

		cw.Error(err, &cw.Logs{Code: "DynamoDBAPI", Message: "No data found"})
		return nil, err
	}

	// Unmarshal a map into actual bus route.
	err = service.DynamoDBAttributeResponse(busRoute, result.Items[0])
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "DynamoDBAttributeResponse", Message: "Failed to unmarshal bus route response."})
		return nil, err
	}

	return busRoute, nil
}

// ValidateBusRoute
func ValidateBusRoute(tablename string, route *models.BusRoute) (bool, error) {
	// Construct the filter builder with a name and value.
	// WHERE bus_unit = bus_unit_value AND from_route = from_route_value AND route_to = route_to_value
	toRoute := expression.Name("to_route").Equal(expression.Value(route.ToRoute))
	fromRoute := expression.Name("from_route").Equal(expression.Value(route.FromRoute))
	filter := expression.Name("bus_unit").Equal(expression.Value(route.BusUnit)).And(fromRoute).And(toRoute)

	// Using the filter expression to create a DynamoDB expression.
	expr, err := expression.NewBuilder().WithFilter(filter).Build()
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "DynamoDBExpression", Message: "Failed to build a dynamodb expression."})
		return false, err
	}

	// Use the built expression to populate the DynamoDB Scan API input parameters.
	input := &dynamodb.ScanInput{
		TableName:                 aws.String(tablename),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
	}

	// Returns one or more items.
	result, err := svc.Scan(input)
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "DynamoDBScan", Message: "Failed to scan input."})
		return false, err
	}

	return (len(result.Items) > 0), nil
}

// BusRoute returns a list of items of bus route that is filtered by "from" and "to" route parameter.
//
// Endpoint:
//  https://{api-id}.execute.api.{region}.amazonaws.com/{stage}/bus/route?from={value}&to={value}
//
// Query Parameters accepted:
//  to: indicating the destination of the bus
//  from: indicating the starting point of a bus
func BusRoute(tablename string, from string, to string) (*events.APIGatewayProxyResponse, error) {
	busRoutes := new([]models.BusRoute)

	// Construct the filter builder with a name and value.
	// WHERE from_route = from_route_value OR route_to = route_to_value
	filter := expression.Name("from_route").Equal(expression.Value(from)).Or(expression.Name("to_route").Equal(expression.Value(to)))

	// Using the filter expression to create a DynamoDB expression.
	expr, err := expression.NewBuilder().WithFilter(filter).Build()
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "DynamoDBExpression", Message: "Failed to build a dynamodb expression."})
		return api.StatusBadRequest(err)
	}

	// Use the built expression to populate the DynamoDB Scan API input parameters.
	input := &dynamodb.ScanInput{
		TableName:                 aws.String(tablename),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
	}

	// Returns one or more items.
	result, err := svc.Scan(input)
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "DynamoDBScan", Message: "Failed to can input."}, kvp.Attribute{Key: "tablename", Value: tablename})
		return api.StatusBadRequest(err)
	}

	// Checks if there are items returned.
	if len(result.Items) == 0 {
		cw.Info(&cw.Logs{Code: "DynamoDBAPI", Message: "There's no bus route data found"})
		return api.StatusOK("no bus routes data found")
	}

	// Unmarshal a map into actual bus routes object.
	err = service.DynamoDBAttributesResponse(busRoutes, result.Items)
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "DynamoDBAttributesResponse", Message: "Failed to unmarshal result to busRoutes"})
		return api.StatusBadRequest(err)
	}

	return api.StatusOK(busRoutes)
}

// BusRouteSchedule returns a list of items of bus route that is filtered by "departure" and arrival "time".
// parameter.
//
// Endpoint:
//  https://{api-id}.execute.api.{region}.amazonaws.com/{stage}/bus/route?departure={value}&arrival={value}
//
// Query Parameters:
//  departure: expected arrival time on the destination and is in 24-hour format (e.g. 15:00)
//  arrival: expected arrival time on the destination and is in 24-hour format (e.g. 18:30)
func BusRouteSchedule(tablename string, departure string, arrival string) (*events.APIGatewayProxyResponse, error) {
	busRoutes := new([]models.BusRoute)

	// Construct the filter builder with a name and value.
	// WHERE departure_time >= departure_time_value AND arrival_time <= arrival_time_value
	filter := expression.Name("departure_time").GreaterThanEqual(expression.Value(departure)).And(expression.Name("arrival_time").LessThanEqual(expression.Value(arrival)))

	// Using the filter expression to create a DynamoDB expression.
	expr, err := expression.NewBuilder().WithFilter(filter).Build()
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "DynamoDBExpression", Message: "Failed to build dynamodb expression."})
		return api.StatusBadRequest(err)
	}

	// Use the built expression to populate the DynamoDB Scan API input parameters.
	input := &dynamodb.ScanInput{
		TableName:                 aws.String(tablename),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
	}

	// Returns one or more items.
	result, err := svc.Scan(input)
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "DynamoDBScan", Message: "Failed to scan input."}, kvp.Attribute{Key: "tablename", Value: tablename})
		return api.StatusBadRequest(err)
	}

	// Checks if there are items returned.
	if len(result.Items) == 0 {
		cw.Info(&cw.Logs{Code: "DynamoDBAPI", Message: "There's no bus route data found"})
		return api.StatusOK("no bus routes data found")
	}

	// Unmarshal a map into actual bus route object
	err = service.DynamoDBAttributesResponse(busRoutes, result.Items)
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "DynamoDBAttributesResponse", Message: "Failed to unmarshal result to busRoutes"})
		return api.StatusBadRequest(err)
	}

	return api.StatusOK(busRoutes)
}

// ListBusRoute returns an API Gateway response of all the items list of bus route.
//
// Endpoint:
//  https://{api-id}.execute.api.{region}.amazonaws.com/{stage}/bus/route/
func ListBusRoute(tablename string) (*events.APIGatewayProxyResponse, error) {
	busRoutes := new([]models.BusRoute)
	input := &dynamodb.ScanInput{TableName: aws.String(tablename)}

	result, err := svc.Scan(input)
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "DynamoDBScan", Message: "Failed to scan input."}, kvp.Attribute{Key: "tablename", Value: tablename})
		return api.StatusBadRequest(err)
	}

	if len(result.Items) == 0 {
		cw.Info(&cw.Logs{Code: "DynamoDBAPI", Message: "There's no entry of bus route data"})
		return api.StatusOK("no bus routes data")
	}

	err = service.DynamoDBAttributesResponse(busRoutes, result.Items)
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "DynamoDBAttributesResponse", Message: "Failed to unmarshal result to busRoutes"})
		return api.StatusBadRequest(err)
	}

	return api.StatusOK(busRoutes)
}
