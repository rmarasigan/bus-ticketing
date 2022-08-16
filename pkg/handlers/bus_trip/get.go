package bustrip

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

// Get is the Bus Trip API request GET method that will process the request.
//
// Query Parameters:
//  unit: bus unit ID and for the identification of specific bus trip
//
// Get all bus trip
//  https://{api-id}.execute.api.{region}.amazonaws.com/{stage}/bus/trip
//
// Sample response:
//  [
// 	{
// 		"id": "602929062182",
// 		"bus_unit": "WXZ-ABCD123",
// 		"bus_route": "GRGTLNT17301930602929",
// 		"seats_left": 25,
// 		"date_created": "1660533116"
//    }
//  ]
//
// Get specific bus trip with specific unit
//  https://{api-id}.execute.api.{region}.amazonaws.com/{stage}/bus/trip?unit={value}
//
// Sample response:
//  [
// 	{
// 		"id": "602929062182",
// 		"bus_unit": "WXZ-ABCD123",
// 		"bus_route": "GRGTLNT17301930602929",
// 		"seats_left": 10,
// 		"date_created": "1660621891"
//    }
//  ]
func Get(ctx context.Context, request *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	// Initialize Dynamodb Session
	service.DynamodbSession()
	svc = service.DynamoDBClient

	tablename := os.Getenv("BUS_TRIP_TABLE")
	queryUnit := request.QueryStringParameters["unit"]

	if tablename == "" {
		err := errors.New("dynamodb table on env is not implemented")

		cw.Error(err, &cw.Logs{Code: "DynamoDBConfig", Message: "BusTicketing_BusTripTable not set on env."})
		return api.StatusUnhandledRequest(err)
	}

	if queryUnit != "" {
		return FilterBusTrip(tablename, queryUnit)
	}

	return ListBusTrip(tablename)
}

// BusTripInfo returns an API Gateway response of a specific bus trip information.
func BusTripInfo(tablename string, id string) (*models.BusTrip, error) {
	busTrip := new(models.BusTrip)

	// Construct the key condition builder with value.
	// WHERE id = id_value
	key := expression.Key("id").Equal(expression.Value(id))

	// Using the key condition to create a DynamoDB expression.
	expr, err := expression.NewBuilder().WithKeyCondition(key).Build()
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "DynamoDBExpression", Message: "Failed to build dynamodb expression."})
		return nil, err
	}

	// Use the build expression to populate the DynamoDB Query's API input parameters.
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

	// Check if there are items returned
	if len(result.Items) == 0 {
		cw.Info(&cw.Logs{Code: "DynamoDBAPI", Message: "No data found."})
		return busTrip, nil
	}

	// Unmarshal a map into actual bus trip.
	err = service.DynamoDBAttributeResponse(busTrip, result.Items[0])
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "DynamoDBAttribtueResponse", Message: "Failed to unmarshal bus trip response."})
		return nil, err
	}

	return busTrip, nil
}

// FilterBusTripRoute returns a BusTrip object that is filtered using bus route
// ID and validates if the trip is within the day.
func FilterBusTripRoute(tablename string, route string) (*models.BusTrip, error) {
	trip := new(models.BusTrip)
	busTrips := new([]models.BusTrip)

	// Construct the filter builder with a name and value.
	// WHERE bus_route = bus_route_value
	filter := expression.Name("bus_route").Equal(expression.Value(route))

	expr, err := expression.NewBuilder().WithFilter(filter).Build()
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "DynamoDBExpression", Message: "Failed to create dynamodb expression."})
		return nil, err
	}

	// Use the filter expression to populate the DynamoDB Scan API input parameters.
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
		return nil, err
	}

	// Checks if there are items returned.
	if len(result.Items) == 0 {
		cw.Info(&cw.Logs{Code: "DynamoDBAPI", Message: "There's no entry of bus trip data."})
		return nil, nil
	}

	// Unmarshal a map into actual bus trip.
	err = service.DynamoDBAttributesResponse(busTrips, result.Items)
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "DynamoDBAttribtueResponse", Message: "Failed to unmarshal bus trip response."})
		return nil, err
	}

	for _, busTrip := range *busTrips {
		tripWithinRange, err := busTrip.IsWithinTheDay()
		if err != nil {
			cw.Error(err, &cw.Logs{Code: "FilterBusTripRoute", Message: "Failed to check if the trip schedule is within the day."})
			return nil, err
		}

		// Checks if the bus trip requested is within the day
		// to return an object.
		if tripWithinRange {
			trip = &busTrip
			break

		} else {
			trip = nil
		}
	}

	return trip, nil
}

// FilterBusTrip returns a list of items of bus trip that is filtered by "unit" parameter.
//
// Endpoint:
//  https://{api-id}.execute.api.{region}.amazonaws.com/{stage}/bus/trip?unit={value}
//
// Query Parameters accepted:
//  unit: bus unit ID and for the identification of specific bus trip
func FilterBusTrip(tablename string, unit string) (*events.APIGatewayProxyResponse, error) {
	busTrips := new([]models.BusTrip)

	// Construct the filter builder with a name and value.
	// WHERE bus_unit = bus_unit_value
	filter := expression.Name("bus_unit").Equal(expression.Value(unit))

	expr, err := expression.NewBuilder().WithFilter(filter).Build()
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "DynamoDBExpression", Message: "Failed to create dynamodb expression."})
		return api.StatusBadRequest(err)
	}

	// Use the filter expression to populate the DynamoDB Scan API input parameters.
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
		cw.Info(&cw.Logs{Code: "DynamoDBAPI", Message: "There's no entry of bus trip data."})
		return api.StatusOK("no bus trip data")
	}

	// Unmarshal a map into actual bus route object.
	err = service.DynamoDBAttributesResponse(busTrips, result.Items)
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "DynamoDBAttributesResponse", Message: "Failed to unmarshal result to busTrips"})
		return api.StatusBadRequest(err)
	}

	return api.StatusOK(busTrips)
}

// ListBusTrip returns an API Gateway response of all the items list of bus trip.
//
// Endpoint:
//  https://{api-id}.execute.api.{region}.amazonaws.com/{stage}/bus/trip/
func ListBusTrip(tablename string) (*events.APIGatewayProxyResponse, error) {
	busTrips := new([]models.BusTrip)

	// Use the built expression to populate the DynamoDB Scan API input parameters.
	input := &dynamodb.ScanInput{TableName: aws.String(tablename)}

	// Returns one or more items.
	result, err := svc.Scan(input)
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "DynamoDBScan", Message: "Failed to scan input."}, kvp.Attribute{Key: "tablename", Value: tablename})
		return api.StatusBadRequest(err)
	}

	// Checks if there are items returned.
	if len(result.Items) == 0 {
		cw.Info(&cw.Logs{Code: "DynamoDBAPI", Message: "There's no entry of bus trip data."})
		return api.StatusOK("no bus trip data")
	}

	// Unmarshal a map into actual bus route object.
	err = service.DynamoDBAttributesResponse(busTrips, result.Items)
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "DynamoDBAttributesResponse", Message: "Failed to unmarshal result to busTrips"})
		return api.StatusBadRequest(err)
	}

	return api.StatusOK(busTrips)
}
