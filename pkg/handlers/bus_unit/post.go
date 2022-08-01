package busunit

import (
	"context"
	"errors"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/rmarasigan/bus-ticketing/pkg/api"
	"github.com/rmarasigan/bus-ticketing/pkg/cw/kvp"
	cw "github.com/rmarasigan/bus-ticketing/pkg/cw/logger"
	"github.com/rmarasigan/bus-ticketing/pkg/handlers/bus"
	"github.com/rmarasigan/bus-ticketing/pkg/models"
	"github.com/rmarasigan/bus-ticketing/pkg/service"
	"github.com/rmarasigan/bus-ticketing/pkg/validate"
)

// Post is the Bus Unit API request POST method that will process the "create" and "update" type requests.
// To process the request, the request query parameter "type" is required and the value must be either
// "create", or "update", also the request query parameter "bus". If none of the said "type" parameter
// values is satisfied it will return an API Gateway response of an HTTP StatusNotImplemented.
//
// Query Parameter:
//  bus: it is the bus id and a required parameter.
//  type: the value must be either "create" or "update" and a required parameter.
//
// Endpoint:
//  https://{api-id}.execute.api.{region}.amazonaws.com/{stage}/bus/unit/?type={value}&bus={value}
func Post(ctx context.Context, request *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	service.DynamodbSession()
	svc = service.DynamoDBClient

	busTable := os.Getenv("BUS_TABLE")
	tablename := os.Getenv("BUS_UNIT_TABLE")

	queryBus := request.QueryStringParameters["bus"]
	queryType := request.QueryStringParameters["type"]

	if tablename == "" {
		err := errors.New("dynamodb table on env is not implemented")

		cw.Error(err, &cw.Logs{Code: "DynamoDBConfig", Message: "BusTicketing_BusUnitTable not set on env."})
		return api.StatusUnhandledRequest(err)
	}

	if busTable == "" {
		err := errors.New("bus table on env is not implemented")

		cw.Error(err, &cw.Logs{Code: "DynamoDBConfig", Message: "BusTicketing_BusTable not set on env."})
		return api.StatusUnhandledRequest(err)
	}

	if queryType == "" {
		err := errors.New("method.request.querystring.type is not implemented")

		cw.Error(err, &cw.Logs{Code: "APIParameter", Message: "Query stirng type is not implemented."})
		return api.StatusBadRequest(err)
	}

	switch queryType {
	case "create":
		if queryBus == "" {
			err := errors.New("method.request.querystring.bus is not implemented")

			cw.Error(err, &cw.Logs{Code: "APIParameter", Message: "Query stirng bus is not implemented."})
			return api.StatusBadRequest(err)
		}

		_, err := bus.BusInformation(busTable, queryBus)
		if err != nil {
			err = errors.New("bus id not found")

			cw.Error(err, &cw.Logs{Code: "BusUnitPostAPI", Message: "Bus ID is not found."})
			return api.StatusBadRequest(err)
		}

		return CreateBusUnit(tablename, queryBus, []byte(request.Body))

	case "update":
		return UpdateBusUnit(tablename, []byte(request.Body))

	default:
		return api.StatusUnhandledRequest(errors.New("request not implemented"))
	}
}

// CreateBusUnit creates a new item to the DynamoDB table. Bus ID is a required field to
// connect the Bus Unit information to the parent Bus. After saving the Bus Unit information
// it will return an API Gateway response.
//
// Endpoint:
//  https://{api-id}.execute.api.{region}.amazonaws.com/{stage}/bus/unit?type=create&bus={value}
//
// Payload Parameters:
//  active: whether the bus is on trip
//  code: unique identification of a bus unit
//  capacity: the number of passengers of a bus unit
//
// Payload Request:
//  {
// 	"capacity": 30,
// 	"active": true,
// 	"code": "XYZ123",
//  }
func CreateBusUnit(tablename string, busID string, body []byte) (*events.APIGatewayProxyResponse, error) {
	unit := new(models.BusUnit)

	// Checks if the request payload body is set.
	if len(body) == 0 {
		err := errors.New("payload is not set")

		cw.Error(err, &cw.Logs{Code: "CreateBusUnit", Message: "Request cannot be processed as payload is not set."})
		return api.StatusBadRequest(err)
	}

	err := api.ParseJSON(body, unit)
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "ParseJSONError", Message: "Failed to parse data to bus unit."})
		return api.StatusBadRequest(err)
	}

	// Validate bus unit information before saving
	isValid := validate.CreateBusUnit(unit)
	if isValid != "" {
		err := errors.New(isValid)

		cw.Error(err, &cw.Logs{Code: "CreateBusUnit", Message: "Validate creation of bus unit."})
		return api.StatusBadRequest(err)
	}

	unit.Bus = busID
	unit.SetValues()

	// Checks whether the bus unit code exist or not.
	busUnitExist, err := ValidateBusUnitCode(tablename, unit.Code)
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "ValidateBusUnitCode", Message: "Failed to validate bus unit code."})
		return api.StatusBadRequest(err)
	}

	if busUnitExist {
		err := errors.New("bus unit code already exist")

		cw.Error(err, &cw.Logs{Code: "ValidateBusUnitCode", Message: "The bus unit code parameter passed already exist."})
		return api.StatusBadRequest(err)
	}

	// Creating the data that you want to send to DynamoDB
	value, err := dynamodbattribute.MarshalMap(unit)
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "MarshalMap", Message: "Failed to marshal bus unit"})
		return api.StatusBadRequest(err)
	}

	input := &dynamodb.PutItemInput{
		Item:      value,
		TableName: aws.String(tablename),
	}

	// Creates a new item or replaces an old item with a new item.
	_, err = svc.PutItem(input)
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "DynamoDBPutItem", Message: "Failed to add item to the table"}, kvp.Attribute{Key: "tablename", Value: tablename})
		return api.StatusBadRequest(err)
	}

	return api.StatusOK(unit)
}

// UpdateBusUnit updates and validates the field before saving the item to the DynamoDB table.
// After updating the bus unit information, it returns an API Gateway response. Bus ID and Bus
// Unit Code parameters cannot be updated.
//
// Endpoint:
//  https://{api-id}.execute.api.{region}.amazonaws.com/{stage}/bus/unit?type=update
//
// Payload Parameter accepts:
//  active: whether the bus is on trip
//  id: unique bus unit ID as the primary key and is required field
//  capacity: the number of passengers of a bus unit
//
// Payload Request:
//  {
// 	"capacity": 30,
// 	"active": true,
// 	"code": "XYZ123",
//  }
func UpdateBusUnit(tablename string, body []byte) (*events.APIGatewayProxyResponse, error) {
	unit := new(models.BusUnit)

	// Checks if the request payload body is set.
	if len(body) == 0 || body == nil {
		err := errors.New("payload is not set")

		cw.Error(err, &cw.Logs{Code: "UpdateBusUnit", Message: "Request cannot be processed as payload is not set"})
		return api.StatusBadRequest(err)
	}

	err := api.ParseJSON(body, unit)
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "ParseJSONError", Message: "Failed to parse data to bus unit."})
		return api.StatusBadRequest(err)
	}

	// Get the bus unit information
	unitInfo, err := BusUnitInformation(tablename, unit.ID)
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "BusUnitInformation", Message: "Failed to get bus unit information."}, kvp.Attribute{Key: "tablename", Value: tablename})
		return api.StatusBadRequest(err)
	}

	// Validate bus unit update information before updating
	unit.ValidateUpdate(unitInfo)
	if unit.Bus != "" && unit.Bus != unitInfo.Bus {
		err := errors.New("cannot update bus id")

		cw.Error(err, &cw.Logs{Code: "ValidateBusUnitUpdate", Message: "Cannot update bus id, composite primary key."})
		return api.StatusBadRequest(err)
	}

	if unit.Code != "" && unit.Code != unitInfo.Code {
		err := errors.New("cannot update bus unit code")

		cw.Error(err, &cw.Logs{Code: "ValidateBusUnitUpdate", Message: "Cannot update bus unit code"})
		return api.StatusBadRequest(err)
	}

	compositePrimaryKey := map[string]*dynamodb.AttributeValue{
		"id":  {S: aws.String(unit.ID)},
		"bus": {S: aws.String(unit.Bus)}}

	// Construct the update builder
	// SET active = active_value, SET capacity = capacity_value
	update := expression.Set(expression.Name("active"), expression.Value(unit.Active)).
		Set(expression.Name("capacity"), expression.Value(unit.Capacity))

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
		ReturnValues:              aws.String("ALL_NEW"), // Returns all of the atrribute of the item (after the UpdateItem operation)
		UpdateExpression:          expr.Update(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	}

	// Update an item in a table
	result, err := svc.UpdateItem(input)
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "DynamoDBUpdateItem", Message: "Failed eto update item."}, kvp.Attribute{Key: "tablename", Value: tablename})
		return api.StatusBadRequest(err)
	}

	// Returns a bus unit response in JSON format.
	err = service.DynamoDBAttributeResponse(unit, result.Attributes)
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "DynamoDBAttributeResponse", Message: "Failed to unmarshal bus unit response."})
		return api.StatusBadRequest(err)
	}

	return api.StatusOK(unit)
}
