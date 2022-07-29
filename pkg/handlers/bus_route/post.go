package busroute

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/rmarasigan/bus-ticketing/pkg/api"
	"github.com/rmarasigan/bus-ticketing/pkg/common"
	"github.com/rmarasigan/bus-ticketing/pkg/cw/kvp"
	cw "github.com/rmarasigan/bus-ticketing/pkg/cw/logger"
	busunit "github.com/rmarasigan/bus-ticketing/pkg/handlers/bus_unit"
	"github.com/rmarasigan/bus-ticketing/pkg/models"
	"github.com/rmarasigan/bus-ticketing/pkg/service"
	"github.com/rmarasigan/bus-ticketing/pkg/validate"
)

// Post is the Bus Route API request POST method that will process the "create" and "update" type request of the specific
// bus unit. To process the request, the request query parameter "type" is required and the value must be either "create"
// or "update". If none of the said "type" parameter values is satisfied it will return an API Gateway response of an
// HTTP StatusNotImplemented.
//
// Endpoint:
//  https://{api-id}.execute.api.{region}.amazonaws.com/{stage}/bus/route?type={value}&unit={value}
//
// Query Parameters:
//  unit: it is the bus unit id and is a required parameter.
//  type: the value must be either "create", or "update" and is a required parameter.
func Post(ctx context.Context, request *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	// Initialize DynamoDB Session
	service.DynamodbSession()
	svc = service.DynamoDBClient

	tablename := os.Getenv("BUS_ROUTE_TABLE")
	queryType := request.QueryStringParameters["type"]

	if queryType == "" {
		return api.StatusUnhandledRequest(errors.New("method.querystring.parameter.type is not implemented"))
	}

	switch queryType {
	case "create":
		queryBusUnit := request.QueryStringParameters["unit"]
		if queryBusUnit == "" {
			return api.StatusUnhandledRequest(errors.New("method.querystring.parameter.unit is not implemented"))
		}

		return CreateBusRoute(tablename, queryBusUnit, []byte(request.Body))

	case "update":
		return UpdateBusRoute(tablename, []byte(request.Body))

	default:
		return api.StatusUnhandledRequest(errors.New("request not implemented"))
	}
}

// CreateBusRoute creates a new item entry to the the DynamoDB table. Payload parameters will be validated.
// As it passed the validation, it will then save the Bus Route information and returns an API Gateway response.
//
// Country currency codes: https://www.iban.com/currency-codes
//
// Endpoint:
//  https://{api-id}.execute.api.{region}.amazonaws.com/{stage}/bus/route?type=create&unit={value}
//
// Payload Parameters:
//  rate: fare charged to the passenger
//  currency_code: medium of exchange for goods and services
//  departure_time: expected departure time on the starting point and is in 24-hour format
//  arrival_time: expected arrival time on the destination and is in 24-hour format
//  from_route: indicating the starting point of a bus
//  to_route: indicating the destination of the bus
//  available: defines if the bus is available for that route
//
// Payload Request:
//  {
//     "rate": 19.20
//     "currency_code": "USD"
//     "departure_time": "06:30",
//     "arrival_time": "08:15",
//     "from_route": "Boston",
//     "to_route": "New York"
//  }
func CreateBusRoute(tablename string, unit string, body []byte) (*events.APIGatewayProxyResponse, error) {
	busRoute := new(models.BusRoute)
	busUnitTable := os.Getenv("BUS_UNIT_TABLE")

	// Checks if the request payload body is set.
	if len(body) == 0 || body == nil {
		err := errors.New("payload is not set")

		cw.Error(err, &cw.Logs{Code: "CreateBusRoute", Message: "Request cannot be processed as payload is not set."})
		return api.StatusBadRequest(err)
	}

	err := api.ParseJSON(body, busRoute)
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "ParseJSONError", Message: "Failed to parse bus route json."})
		return api.StatusBadRequest(err)
	}

	// Validate if rate is in proper format
	rate := fmt.Sprint(busRoute.Rate)
	isValidRate, err := common.IsNumberOnly(rate)
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "IsNumberOnly", Message: "Failed to validate rate if IsNumberOnly."})
		return api.StatusBadRequest(err)
	}

	if !isValidRate {
		err := errors.New("bus route rate is not a number")

		cw.Error(err, &cw.Logs{Code: "IsNumberOnly", Message: "Bus route rate is not in a number format."}, kvp.Attribute{Key: "rate", Value: rate})
		return api.StatusBadRequest(err)
	}

	// Get the bus unit information and check if it exist
	busUnit, err := busunit.BusUnitInformation(busUnitTable, unit)
	if err != nil {
		err = errors.New("bus unit does not exist")

		cw.Error(err, &cw.Logs{Code: "BusUnitInformation", Message: "The bus unit does not exist."}, kvp.Attribute{Key: "tablename", Value: busUnitTable})
		return api.StatusBadRequest(err)
	}

	// Validate if the required fields are not empty.
	isValid := validate.CreateBusRoute(busRoute)
	if isValid != "" {
		err := errors.New(isValid)

		cw.Error(err, &cw.Logs{Code: "CreateBusRoute", Message: "Validate creation of bus route."})
		return api.StatusBadRequest(err)
	}

	// Setting values
	busRoute.SetValues()
	busRoute.BusUnit = unit
	busRoute.Bus = busUnit.Bus

	// Validate if the record already exist
	busRouteExist, err := ValidateBusRoute(tablename, busRoute)
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "ValidateBusRoute", Message: "Failed to validate creation of bus route."})
		return api.StatusBadRequest(err)
	}

	if busRouteExist {
		err := errors.New("bus unit route already exist")

		cw.Error(err, &cw.Logs{Code: "ValidateBusRoute", Message: "Bus unit route entry already exist."})
		return api.StatusBadRequest(err)
	}

	// Converting the record to dynamodb.AttributeValue type.
	value, err := dynamodbattribute.MarshalMap(busRoute)
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "DynamoDBMarshalMap", Message: "Failed to marshal bus route map."})
		return api.StatusBadRequest(err)
	}

	// Creating the data that you want to send to DynamoDB
	input := &dynamodb.PutItemInput{
		Item:      value, // Map of attribute name-value pairs, one for each attribute
		TableName: aws.String(tablename),
	}

	// Creates a new item or replaces an old item with a new item.
	_, err = svc.PutItem(input)
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "DynamoDBPutItem", Message: "Failed to add item to the table."}, kvp.Attribute{Key: "tablename", Value: tablename})
		return api.StatusBadRequest(err)
	}

	return api.StatusOK(busRoute)
}

// UpdateBusRoute updates and validates the field before saving the item to the DynamoDB table.
// After updating the bus route information, it returns an API Gateway response. ID parameter
// cannot be updated.
//
// Country currency codes: https://www.iban.com/currency-codes
//
// Endpoint
//  https://{api-id}.execute.api.{region}.amazonaws.com/{stage}/bus/route?type=update
//
// Payload Parameters:
//  id: unique bus route ID as the primary key and is a required field
//  rate: fare charged to the passenger
//  currency_code: medium of exchange for goods and services
//  departure_time: expected departure time on the starting point and is in 24-hour format
//  arrival_time: expected arrival time on the destination and is in 24-hour format
//  from_route: indicating the starting point of a bus
//  route_to: indicating the destination of the bus
//  available: defines if the bus is available for that route
//
// Payload Request:
//  {
//     "id": "BSTNSNWRK06300815590636",
//     "rate": 19,
//     "currency_code": "USD",
//     "departure_time": "06:30",
//     "arrival_time": "08:15",
//     "from_route": "Boston",
//     "to_route": "New York"
//  }
func UpdateBusRoute(tablename string, body []byte) (*events.APIGatewayProxyResponse, error) {
	busRoute := new(models.BusRoute)

	// Checks if the request payload body is set.
	if len(body) == 0 {
		err := errors.New("payload is not set")

		cw.Error(err, &cw.Logs{Code: "UpdateBusRoute", Message: "Request cannot be processed as payload."})
		return api.StatusBadRequest(err)
	}

	err := api.ParseJSON(body, busRoute)
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "ParseJSONError", Message: "Failed to parse data to bus route."})
		return api.StatusBadRequest(err)
	}

	// Validate if rate is in proper format.
	rate := fmt.Sprint(busRoute.Rate)
	isValidRate, err := common.IsNumberOnly(rate)
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "IsNumberOnly", Message: "Failed to validate rate if IsNumberOnly."})
		return api.StatusBadRequest(err)
	}

	if !isValidRate {
		err := errors.New("bus route rate is not a number")

		cw.Error(err, &cw.Logs{Code: "IsNumberOnly", Message: "Bus route rate is not in a number format."}, kvp.Attribute{Key: "rate", Value: rate})
		return api.StatusBadRequest(err)
	}

	// Get the bus route information.
	routeInfo, err := BusRouteInfo(tablename, busRoute.ID)
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "BusRouteInfo", Message: "Failed to get bus route info."}, kvp.Attribute{Key: "tablename", Value: tablename})
		return api.StatusBadRequest(err)
	}

	// Validate bus route update information before updating.
	busRoute.ValidateUpdate(routeInfo)
	if busRoute.ID != "" && busRoute.ID != routeInfo.ID {
		err := errors.New("cannot update bus route id")

		cw.Error(err, &cw.Logs{Code: "ValidateUpdate", Message: "Cannot update bus route id, composite primary key."})
		return api.StatusBadRequest(err)
	}

	compositePrimaryKey := map[string]*dynamodb.AttributeValue{
		"id":  {S: aws.String(busRoute.ID)},
		"bus": {S: aws.String(routeInfo.Bus)}}

	// Construct the update builder.
	update := expression.Set(expression.Name("rate"), expression.Value(busRoute.Rate)).
		Set(expression.Name("currency_code"), expression.Value(busRoute.Currency)).
		Set(expression.Name("departure_time"), expression.Value(busRoute.DepartureTime)).
		Set(expression.Name("arrival_time"), expression.Value(busRoute.ArrivalTime)).
		Set(expression.Name("from_route"), expression.Value(busRoute.FromRoute)).
		Set(expression.Name("to_route"), expression.Value(busRoute.ToRoute))

	// Using the update to create a DynamoDB Expression.
	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "DynamoDBExpression", Message: "Failed to build an expression."})
		return api.StatusBadRequest(err)
	}

	// Use the built expression to populate the DynamoDB Update Item API input parameters.
	input := &dynamodb.UpdateItemInput{
		TableName:                 aws.String(tablename),
		Key:                       compositePrimaryKey,
		ReturnValues:              aws.String("ALL_NEW"), // Returns all of the attribute of the item (after the UpdateItem operation)
		UpdateExpression:          expr.Update(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	}

	// Update an item in a table.
	result, err := svc.UpdateItem(input)
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "DynamoDBUpdateItem", Message: "Failed to update item."}, kvp.Attribute{Key: "tablename", Value: tablename})
		return api.StatusBadRequest(err)
	}

	// Returns a bus route response.
	err = service.DynamoDBAttributeResponse(busRoute, result.Attributes)
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "DynamoDBAttributeResponse", Message: "Failed to unmarshal bus route response."})
		return api.StatusBadRequest(err)
	}

	return api.StatusOK(busRoute)
}
