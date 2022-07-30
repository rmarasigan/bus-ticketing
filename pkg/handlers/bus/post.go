package bus

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
	"github.com/rmarasigan/bus-ticketing/pkg/models"
	"github.com/rmarasigan/bus-ticketing/pkg/service"
	"github.com/rmarasigan/bus-ticketing/pkg/validate"
)

// Post is the Bus API request POST method that will process the "create" and "update" type requests.
// To process the request, the request query parameter "type" is required, and the value must be either
// "create" or "update". If none of the said type parameter values is satisfied, it will return an API
// Gateway response of an HTTP StatusNotImplemented.
//
// Endpoint:
//  https://{api-id}.execute.api.{region}.amazonaws.com/{stage}/bus?type={value}
func Post(ctx context.Context, request *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	// Intialize DynamoDB Session
	service.DynamodbSession()
	svc = service.DynamoDBClient

	tablename := os.Getenv("BUS_TABLE")
	queryType := request.QueryStringParameters["type"]

	if tablename == "" {
		err := errors.New("dynamodb table on env is not implemented")

		cw.Error(err, &cw.Logs{Code: "DynamoDBConfig", Message: "BusTicketing_BusTable not set on env."})
		return api.StatusUnhandledRequest(err)
	}

	if queryType == "" {
		err := errors.New("method.request.querystring.type is not set")

		cw.Error(err, &cw.Logs{Code: "APIParameter", Message: "Query string type is not implemented."})
		return api.StatusUnhandledRequest(err)
	}

	switch queryType {
	case "create":
		return CreateBus(tablename, []byte(request.Body))

	case "update":
		return UpdateBus(tablename, []byte(request.Body))

	default:
		return api.StatusUnhandledRequest(errors.New("request not implemented"))
	}
}

// CreateBus creates a new item entry of bus company information and validates if the required fields
// are not empty before saving it to the DynamoDB table. After saving the Bus information it will
// return an API Gateway response of an HTTP StatusOK.
//
// Endpoint:
//  https://{api-id}.execute.api.{region}.amazonaws.com/{stage}/bus?type=create
//
// Payload Parameters:
//  owner: bus company owner
//  company: name of the company and serves as your sort key
//  address: bus company complete address
//  email: bus company email
//  mobile_number: bus company mobile number
//
// Payload Request:
//  {
//    "owner": "Thando Oyibo Emmett",
//    "company": "Rail Bus Way",
//    "address": "1986 Bogisich Junctions, Hamillhaven, Kansas",
//    "email": "thando.emmet@outlook.com",
//    "mobile_number": "+1-335-908-1432"
//  }
func CreateBus(tablename string, body []byte) (*events.APIGatewayProxyResponse, error) {
	bus := new(models.Bus)

	// Checks if the request payload body is set.
	if len(body) == 0 {
		err := errors.New("payload is not set")

		cw.Error(err, &cw.Logs{Code: "CreateBus", Message: "Request cannot be processed as payload is not set."})
		return api.StatusBadRequest(err)
	}

	err := api.ParseJSON(body, bus)
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "ParseJSONError", Message: "Failed to parse bus json."})
		return api.StatusBadRequest(err)
	}

	// Validate if the required fields are not empty.
	isValid := validate.CreateBus(bus)
	if isValid != "" {
		err := errors.New(isValid)

		cw.Error(err, &cw.Logs{Code: "CreateBus", Message: "Validate creation of bus."})
		return api.StatusBadRequest(err)
	}

	bus.SetValues()

	// Checks if the bus company already exist.
	companyExist, err := ValidateBusCompany(tablename, bus.Company, bus.Address)
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "ValidateBusCompany", Message: "Failed to validate bus company name."})
		return api.StatusBadRequest(err)
	}

	// Returns error message of "company name already exist".
	if companyExist {
		err := errors.New("company already exist")

		cw.Error(err, &cw.Logs{Code: "ValidateBusCompany", Message: "The company parameters passed already exist."})
		return api.StatusBadRequest(err)
	}

	// Converting the record to dynamodb.AttributeValue type.
	value, err := dynamodbattribute.MarshalMap(bus)
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "DynamoDBMarshalMap", Message: " Failed to marshal bus data."})
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

	return api.StatusOK(bus)
}

// UpdateBus validates the field before saving the item to the DynamoDB table and updates and returns
// the updated bus information. After updating the bus information, it will return an API Gateway
// response of an HTTP StatusOK.
//
// Endpoint:
//  https://{api-id}.execute.api.{region}.amazonaws.com/{stage}/bus?type=update
//
// Payload Parameter accepts:
//  id: bus ID as the primary key and is a required field
//  owner: bus company owner
//  address: bus company complete address
//  email: bus company email
//  mobile_number: bus company mobile number
//
// Payload Request:
//  {
//    "id": "RLBSW-589710"
//    "owner": "Thando Oyibo Emmett",
//    "address": "1986 Bogisich Junctions, Hamillhaven, Kansas",
//    "email": "thando.emmet@outlook.com",
//    "mobile_number": "+1-335-908-1432"
//  }
func UpdateBus(tablename string, body []byte) (*events.APIGatewayProxyResponse, error) {
	bus := new(models.Bus)

	err := api.ParseJSON(body, bus)
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "ParseJSONError", Message: "Failed to parse bus json."})
		return api.StatusBadRequest(err)
	}

	// Check if the bus ID is implemented or not.
	// Cannot update item without the primary key.
	if bus.ID == "" {
		err := errors.New("bus id is not set")

		cw.Error(err, &cw.Logs{Code: "UpdateBusInformation", Message: "Bus ID is not implemented."})
		return api.StatusUnhandledRequest(err)
	}

	// Get bus information using the ID
	busInfo, err := BusInformation(tablename, bus.ID)
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "BusInformation", Message: "Failed to get bus information."}, kvp.Attribute{Key: "tablename", Value: tablename})
		return api.StatusBadRequest(err)
	}

	// Validate bus information
	bus.ValidateUpdate(busInfo)
	if bus.Company != "" && bus.Company != busInfo.Company {
		err := errors.New("cannot update company name")

		cw.Error(err, &cw.Logs{Code: "ValidateBusUpdate", Message: "Cannot update company name, composite primary key."})
		return api.StatusBadRequest(err)
	}

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
		cw.Error(err, &cw.Logs{Code: "DynamoDBExpression", Message: "Failed to build an expression."})
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
		cw.Error(err, &cw.Logs{Code: "DynamoDBUpdateItem", Message: "Failed to update item."}, kvp.Attribute{Key: "tablename", Value: tablename})
		return api.StatusBadRequest(err)
	}

	// Returns a bus response in JSON format.
	err = service.DynamoDBAttributeResponse(bus, result.Attributes)
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "DynamoDBAttributeResponse", Message: "Failed to unmarshal bus response."})
		return api.StatusBadRequest(err)
	}

	return api.StatusOK(bus)
}
