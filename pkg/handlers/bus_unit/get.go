package busunit

import (
	"context"
	"errors"
	"os"
	"strconv"

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

// Get is the Bus Unit API request GET method that will process the request.
//
// Query Parameters:
// 	capacity: the number of passenger of a bus unit.
// 	active: whether the bus is on trip. accepts true or false value.
//
// Get all bus unit
//  https://{api-id}.execute.api.{region}.amazonaws.com/{stage}/bus/unit/
//
// Sample response:
//  [
//     {
//         "id": "WXZ-ABCD123",
//         "bus": "WXZ-587390",
//         "code": "abcd123",
//         "active": true,
//         "capacity": 40,
//         "date_created": "1658889765"
//     },
//     {
//         "id": "BCDFGH-EFG456",
//         "bus": "BCDFGH-587390",
//         "code": "efg456",
//         "active": false,
//         "capacity": 35,
//         "date_created": "1658889595"
//     }
//  ]
//
// Get active bus unit
//  https://{api-id}.execute.api.{region}.amazonaws.com/{stage}/bus/unit/?active={value}
//
// Sample response:
//  [
//     {
//         "id": "WXZ-ABCD123",
//         "bus": "WXZ-587390",
//         "code": "abcd123",
//         "active": true,
//         "capacity": 40,
//         "date_created": "1658889765"
//     }
//  ]
//
// Get list of bus unit with specific capacity
//  https://{api-id}.execute.api.{region}.amazonaws.com/{stage}/bus/unit/?capacity={value}
//
// Sample response:
//  [
//     {
//         "id": "WXZ-ABCD123",
//         "bus": "WXZ-587390",
//         "code": "abcd123",
//         "active": true,
//         "capacity": 40,
//         "date_created": "1658889765"
//     }
//  ]
func Get(ctx context.Context, request *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	// Creates DynamoDB Session
	service.DynamodbSession()
	svc = service.DynamoDBClient

	tablename := os.Getenv("BUS_UNIT_TABLE")

	queryActive := request.QueryStringParameters["active"]
	queryCapacity := request.QueryStringParameters["capacity"]

	if queryActive != "" {
		active, err := strconv.ParseBool(queryActive)
		if err != nil {
			cw.Error(err, &cw.Logs{Code: "StrconvParseBool", Message: "Failed to convert active string to bool"})
			return api.StatusBadRequest(err)
		}

		return FilterActiveBusUnit(tablename, active)
	}

	if queryCapacity != "" {
		capacity, err := strconv.Atoi(queryCapacity)
		if err != nil {
			cw.Error(err, &cw.Logs{Code: "StrconvAtoi", Message: "Failed to convert capacity string to int"})
			return api.StatusBadRequest(err)
		}

		return FilterCapacityBusUnit(tablename, capacity)
	}

	return ListBusUnit(tablename)
}

// BusUnitInformation fetches the information about the bus unit.
func BusUnitInformation(tablename string, id string) (*models.BusUnit, error) {
	busUnit := new(models.BusUnit)
	svc = service.DynamoDBClient

	// Construct the key condition builder with value.
	// WHERE id = id_value
	key := expression.Key("id").Equal(expression.Value(id))

	// Using the key condition and projection to create a DynamoDB expression.
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

	if len(result.Items) == 0 {
		cw.Info(&cw.Logs{Code: "DynamoDBAPI", Message: "No data found"})
		return nil, errors.New("bus unit information not found")
	}

	// Unmarshal a map into actual user which front-end can uderstand as a JSON
	err = service.DynamoDBAttributeResponse(busUnit, result.Items[0])
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "DynamoDBAttributeResponse", Message: "Failed to unmarshal bus unit response"})
		return nil, err
	}

	return busUnit, nil
}

// ValidateBusUnitCode returns a boolean and error value to check whether the bus unit code already exists or not.
func ValidateBusUnitCode(tablename string, code string) (bool, error) {
	// Construct the filter builder with a name and value.
	// WHERE code = code_value
	filter := expression.Name("code").Equal(expression.Value(code))

	// Using the filter expression to create a DynamoDB expression.
	expr, err := expression.NewBuilder().WithFilter(filter).Build()
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "DynamoDBExpression", Message: "Failed to build dynamodb expression."})
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

// FilterActiveBusUnit returns an API Gateway response of filtered list of bus unit which is active or not.
func FilterActiveBusUnit(tablename string, active bool) (*events.APIGatewayProxyResponse, error) {
	busUnit := new([]models.BusUnit)

	// Construct the filter builder with a name and value.
	// WHERE active = active_value
	filter := expression.Name("active").Equal(expression.Value(aws.Bool(active)))

	// Using the filter expression to create a DynamoDB expression.
	expr, err := expression.NewBuilder().WithFilter(filter).Build()
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "DynamoDBExpression", Message: "Failed to build dynamodb expression."})
		return api.StatusBadRequest(err)
	}

	// Use the built expression to populate the DynamoDB Scan API input parameters.
	input := &dynamodb.ScanInput{
		TableName:                 aws.String(tablename),
		FilterExpression:          expr.Filter(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	}

	// Returns one or more items.
	result, err := svc.Scan(input)
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "DynamoDBScan", Message: "Failed to scan input."}, kvp.Attribute{Key: "tablename", Value: tablename})
		return api.StatusBadRequest(err)
	}

	// Checks if there are items returned.
	if len(result.Items) == 0 {
		cw.Info(&cw.Logs{Code: "DynamoDBAPI", Message: "No data found."})
		return api.StatusOK("no active bus unit")
	}

	// Unmarshal a map into actual bus unit.
	err = service.DynamoDBAttributesResponse(busUnit, result.Items)
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "DynamoDBAttributesResponse", Message: "Failed to unmarshal bus unit response"})
		return api.StatusBadRequest(err)
	}

	return api.StatusOK(busUnit)
}

// FilterCapacityBusUnit returns an API Gateway response of filtered list of bus unit which has the specific
// number of passenger (capacity).
func FilterCapacityBusUnit(tablename string, capacity int) (*events.APIGatewayProxyResponse, error) {
	busUnit := new([]models.BusUnit)

	// Construct the filter builder with a name and value.
	// WHERE capacity = capacity_value
	filter := expression.Name("capacity").Equal(expression.Value(aws.Int(capacity)))

	// Using the filter expression to create a DynamoDB expression.
	expr, err := expression.NewBuilder().WithFilter(filter).Build()
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "DynamoDBExpression", Message: "Failed to build dynamodb expression."})
		return api.StatusBadRequest(err)
	}

	// Use the built expression to populate the DynamoDB Scan API input parameters.
	input := &dynamodb.ScanInput{
		TableName:                 aws.String(tablename),
		FilterExpression:          expr.Filter(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	}

	// Returns one or more items.
	result, err := svc.Scan(input)
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "DynamoDBScan", Message: "Failed to scan input."}, kvp.Attribute{Key: "tablename", Value: tablename})
		return api.StatusBadRequest(err)
	}

	// Checks if there are items returned.
	if len(result.Items) == 0 {
		cw.Info(&cw.Logs{Code: "DynamoDBAPI", Message: "No data found."})
		return api.StatusOK("no active bus unit")
	}

	// Unmarshal a map into actual bus unit.
	err = service.DynamoDBAttributesResponse(busUnit, result.Items)
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "DynamoDBAttributesResponse", Message: "Failed to unmarshal bus unit response"})
		return api.StatusBadRequest(err)
	}

	return api.StatusOK(busUnit)
}

// ListBusUnit returns an API Gateway response of all the items list of bus unit.
func ListBusUnit(tablename string) (*events.APIGatewayProxyResponse, error) {
	busUnits := new([]models.BusUnit)
	input := &dynamodb.ScanInput{TableName: aws.String(tablename)}

	// Returns one or more items.
	result, err := svc.Scan(input)
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "DynamoDBScan", Message: "Failed to scan input."}, kvp.Attribute{Key: "tablename", Value: tablename})
		return api.StatusBadRequest(err)
	}

	// Checks if there are items returned.
	if len(result.Items) == 0 {
		cw.Info(&cw.Logs{Code: "DynamoDBAPI", Message: "No data found."})
		return api.StatusOK("no bus unit data available.")
	}

	// Unmarshal a map into actual bus unit.
	err = service.DynamoDBAttributesResponse(busUnits, result.Items)
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "DynamoDBAttributesResponse", Message: "Failed to unmarshal bus unit response."})
		return api.StatusBadRequest(err)
	}

	return api.StatusOK(busUnits)
}
