package bus

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

// Get is the Bus API request GET method that will process the request.
func Get(ctx context.Context, request *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	// Creates DynamoDB Session
	service.DynamodbSession()
	svc = service.DynamoDBClient

	tablename := os.Getenv("BUS_TABLE")
	queryCompany := request.QueryStringParameters["company"]

	if tablename == "" {
		err := errors.New("dynamodb table on env is not implemented")

		cw.Error(err, &cw.Logs{Code: "DynamoDBConfig", Message: "BusTicketing_BusTable not set on env."})
		return api.StatusUnhandledRequest(err)
	}

	if queryCompany != "" {
		return FilterBus(tablename, queryCompany)
	}

	return ListBus(tablename)
}

// BusInformation returns a bus information.
func BusInformation(tablename string, id string) (*models.Bus, error) {
	bus := new(models.Bus)
	svc = service.DynamoDBClient

	// Construct the filter builder with a name and value.
	// WHERE id = id_value
	key := expression.Key("id").Equal(expression.Value(id))

	// Using the key condition to create a DynamoDB expression.
	expr, err := expression.NewBuilder().WithKeyCondition(key).Build()
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "DynamoDBExpression", Message: "Failed to build an expression"})
		return nil, err
	}

	// Use the built expression to populate the DynamoDB Query's API input parameters.
	input := &dynamodb.QueryInput{
		TableName:                 aws.String(tablename),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
	}

	// Returns one or more items and item attributes.
	result, err := svc.Query(input)
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "DynamoDBQuery", Message: "Failed to query input."}, kvp.Attribute{Key: "tablename", Value: tablename})
		return nil, err
	}

	// Checks if there are items returned.
	if len(result.Items) == 0 {
		cw.Info(&cw.Logs{Code: "DynamoDBAPI", Message: "No data found."})
		return nil, errors.New("bus information not found")
	}

	// Unmarshal a map into actual user which front-end can uderstand as a JSON
	err = service.DynamoDBAttributeResponse(bus, result.Items[0])
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "DynamoDBAttributesResponse", Message: "Failed to unmarshal bus response."})
		return nil, err
	}

	return bus, nil
}

func ValidateBusCompany(tablename string, company string, address string) (bool, error) {
	// Construct the condition builder with a name that contains of specified value.
	// WHERE company = company_value AND address = address_value
	companyAddress := expression.Name("address").Equal(expression.Value(address))
	filter := expression.Name("company").Equal(expression.Value(company)).And(companyAddress)

	// Using the filter condition to create a DynamoDB expression.
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
		cw.Error(err, &cw.Logs{Code: "DynamoDBScan", Message: "Failed to scan input."}, kvp.Attribute{Key: "tablename", Value: tablename})
		return false, err
	}

	return (len(result.Items) > 0), nil
}

// FilterBus returns a list of bus that is filtered by the company name.
func FilterBus(tablename string, company string) (*events.APIGatewayProxyResponse, error) {
	bus := new([]models.Bus)

	// Construct the condition builder with a name that contains of specified value.
	// WHERE company LIKE %company_value%
	filter := expression.Contains(expression.Name("company"), company)

	// Using the condition to create a DynamoDB expression.
	expr, err := expression.NewBuilder().WithFilter(filter).Build()
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "DynamoDBExpression", Message: "Failed to build an expression."})
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
		cw.Info(&cw.Logs{Code: "DynamoDBAPI", Message: "No data found."})
		return api.StatusOK("bus information not found")
	}

	// Returns a bus response in JSON formation
	err = service.DynamoDBAttributesResponse(bus, result.Items)
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "DynamoDBAttributesResponse", Message: "Failed to unmarshal bus response."})
		return api.StatusBadRequest(err)
	}

	return api.StatusOK(bus)
}

// ListBus returns a list of all bus items in the DynamoDB table.
func ListBus(tablename string) (*events.APIGatewayProxyResponse, error) {
	bus := new([]models.Bus)
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
		return api.StatusOK("no bus information data")
	}

	// Unmarshal a map into actual bus struct.
	err = service.DynamoDBAttributesResponse(bus, result.Items)
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "DynamoDBAttributesResponse", Message: "Failed to unmarshal bus response."})
		return api.StatusBadRequest(err)
	}

	return api.StatusOK(bus)
}
