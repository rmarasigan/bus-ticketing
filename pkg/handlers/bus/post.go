package bus

import (
	"context"
	"errors"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/rmarasigan/bus-ticketing/pkg/api"
	"github.com/rmarasigan/bus-ticketing/pkg/cw/kvp"
	cw "github.com/rmarasigan/bus-ticketing/pkg/cw/logger"
	"github.com/rmarasigan/bus-ticketing/pkg/models"
	"github.com/rmarasigan/bus-ticketing/pkg/service"
	"github.com/rmarasigan/bus-ticketing/pkg/validate"
)

// Post is the Bus API request POST method that will process the "create" and "update" type request.
// To process the request, request query "type" is required and the value must be either "create", or "update".
func Post(ctx context.Context, request *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	// Creates DynamoDB Session
	service.DynamodbSession()

	tablename := os.Getenv("BUS_TABLE")
	queryType := request.QueryStringParameters["type"]

	if tablename == "" {
		return api.StatusUnhandledRequest(errors.New("dynamodb table on env is not implemented"))
	}

	if queryType == "" {
		return api.StatusBadRequest(errors.New("method.request.querystring.type is not set"))
	}

	switch queryType {
	case "create":
		return CreateBus(tablename, []byte(request.Body), service.DynamoDBClient)

	case "update":
		return UpdateBus(tablename, []byte(request.Body), service.DynamoDBClient)

	default:
		return api.StatusUnhandledRequest(errors.New("request not implemented"))
	}
}

// CreateBus creates a new bus company information and validates if the required fields
// are not empty before saving to the DynamoDB table.
func CreateBus(tablename string, body []byte, svc dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {
	bus := new(models.Bus)

	err := api.ParseJSON(body, bus)
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "ParseJSONError", Message: "Failed to parse bus json"})
		return api.StatusBadRequest(err)
	}

	// Validate if the required fields are not empty.
	isValid := validate.CreateBus(bus)
	if isValid != "" {
		err := errors.New(isValid)

		cw.Error(err, &cw.Logs{Code: "CreateBus", Message: "Validate creation of bus"})
		return api.StatusUnhandledRequest(err)
	}

	bus.SetValues()

	// Converting the record to dynamodb.AttributeValue type.
	value, err := dynamodbattribute.MarshalMap(bus)
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "DynamoDBMarshalMap", Message: " Failed to marshal bus"})
		return api.StatusBadRequest(err)
	}

	// Creating the data that you want to send to dynamoDB
	input := &dynamodb.PutItemInput{
		Item:      value, // Map of attribute name-value pairs, one for each attribute
		TableName: aws.String(tablename),
	}

	// Creates a new item or replaces an old item with a new item.
	_, err = svc.PutItem(input)
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "DynamoDBPutItem", Message: "Failed to add item to the table"}, kvp.Attribute{Key: "tablename", Value: tablename})
		return api.StatusBadRequest(err)
	}

	return api.StatusOK(bus)
}

// UpdateBus updates, validates and returns the updated bus information.
func UpdateBus(tablename string, body []byte, svc dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {
	bus := new(models.Bus)

	err := api.ParseJSON(body, bus)
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "ParseJSONError", Message: "Failed to parse bus json"})
		return api.StatusBadRequest(err)
	}

	// Get bus information using the ID
	busInfo, err := BusInformation(tablename, bus.ID, svc)
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "BusInformation", Message: "Failed to get bus information"}, kvp.Attribute{Key: "tablename", Value: tablename})
		return api.StatusBadRequest(err)
	}

	// Validate bus information
	bus.ValidateUpdate(busInfo)
	compositePrimaryKey := map[string]*dynamodb.AttributeValue{
		"id":      {S: aws.String(bus.ID)},
		"company": {S: aws.String(bus.Company)}}

	// Construct the update builder
	// SET owner = owner_value, SET address = address_value,
	// SET email = email_value, SET mobile_number = mobile_number_value
	update := expression.Set(expression.Name("owner"), expression.Value(bus.Owner)).
		Set(expression.Name("address"), expression.Value(bus.Address)).
		Set(expression.Name("email"), expression.Value(bus.Email)).
		Set(expression.Name("mobile_number"), expression.Value(bus.MobileNumber))

	// Using the update to create a DynamoDB Expression.
	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "DynamoDBExpression", Message: "Failed to build an expression"})
		return api.StatusBadRequest(err)
	}

	// Use the built expression to populate the DynamoDB Update Item API input parameters.
	input := &dynamodb.UpdateItemInput{
		TableName:                 aws.String(tablename),
		Key:                       compositePrimaryKey,
		ReturnValues:              aws.String("ALL_NEW"), // Returns all of the attributes of the item (after the UpdateItem operation)
		UpdateExpression:          expr.Update(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	}

	// Update an item in a table.
	result, err := svc.UpdateItem(input)
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "DynamoDBUpdateItem", Message: "Failed to update item"}, kvp.Attribute{Key: "tablename", Value: tablename})
		return api.StatusBadRequest(err)
	}

	// Returns a bus response in JSON format.
	err = service.DynamoDBAttributeResponse(bus, result.Attributes)
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "DynamoDBAttributeResponse", Message: "Failed to unmarshal bus response"})
		return api.StatusBadRequest(err)
	}

	return api.StatusOK(bus)
}
