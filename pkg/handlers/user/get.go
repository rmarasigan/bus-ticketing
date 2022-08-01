package user

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

// Get is the Users API request GET method that will process the request.
func Get(ctx context.Context, request *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	// Creates Dynamodb Session
	service.DynamodbSession()

	tablename := os.Getenv("USERS_TABLE")
	queryUserID := request.QueryStringParameters["id"]

	if tablename == "" {
		err := errors.New("dynamodb table on env is not implemented")

		cw.Error(err, &cw.Logs{Code: "DynamoDBConfig", Message: "BusTicketing_UsersTable not set on env."})
		return api.StatusUnhandledRequest(err)
	}

	if queryUserID == "" {
		err := errors.New("method.request.querystring.id is not implemented")

		cw.Error(err, &cw.Logs{Code: "APIParameter", Message: "Query string id is not implemented."})
		return api.StatusBadRequest(err)
	}

	response, err := UserInformation(tablename, queryUserID, service.DynamoDBClient)
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "UserInformation", Message: "Failed to get user information"})
		return api.StatusBadRequest(err)
	}

	return api.StatusOK(response)
}

// ValidateUsername returns boolean and error value to check whether the username already exist or not.
func ValidateUsername(tablename string, username string, svc dynamodbiface.DynamoDBAPI) (bool, error) {
	// Create the names list projection of names to project.
	// SELECT id, type, first_name, last_name, username, address, email, mobile_number
	projection := expression.NamesList(expression.Name("id"), expression.Name("type"),
		expression.Name("first_name"), expression.Name("last_name"),
		expression.Name("username"), expression.Name("address"),
		expression.Name("email"), expression.Name("mobile_number"))

	// Construct the filter builder with a name and value.
	// WHERE username = username_value
	filter := expression.Name("username").Equal(expression.Value(username))

	// Using the filter and projections create a DynamoDB expression.
	expr, err := expression.NewBuilder().WithFilter(filter).WithProjection(projection).Build()
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "DynamoDBExpression", Message: "Failed to build an expression"})
		return false, err
	}

	// Use the built expression to populate the DynamoDB Scan API input parameters.
	input := &dynamodb.ScanInput{
		TableName:                 aws.String(tablename),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
	}

	// Returns one or more items and item attributes.
	result, err := svc.Scan(input)
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "DynamoDBScan", Message: "Failed to scan input"}, kvp.Attribute{Key: "tablename", Value: tablename})
		return false, err
	}

	return (len(result.Items) > 0), nil
}

// UserInformation returns a user account information.
func UserInformation(tablename string, id string, svc dynamodbiface.DynamoDBAPI) (*models.UserResponse, error) {
	response := new(models.UserResponse)

	// Construct the Key condition builder
	// WHERE id = id_value
	key := expression.Key("id").Equal(expression.Value(id))

	// Create the names list projection of names to project.
	// SELECT id, type, first_name, last_name, username, address, email, mobile_number
	projection := expression.NamesList(expression.Name("id"), expression.Name("type"),
		expression.Name("first_name"), expression.Name("last_name"),
		expression.Name("username"), expression.Name("address"),
		expression.Name("email"), expression.Name("mobile_number"))

	// Using the key condition and projection together as a DynamoDB expression.
	expr, err := expression.NewBuilder().WithKeyCondition(key).WithProjection(projection).Build()
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "DynamoDBExpression", Message: "Failed to build an expression"})
		return nil, err
	}

	// Use the built expression to populate the DynamoDB Query's API input parameters.
	input := &dynamodb.QueryInput{
		TableName:                 aws.String(tablename),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ProjectionExpression:      expr.Projection(),
		KeyConditionExpression:    expr.KeyCondition(),
	}

	// Returns one or more items and item attributes.
	result, err := svc.Query(input)
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "DynamoDBQuery", Message: "Failed to query input"}, kvp.Attribute{Key: "tablename", Value: tablename})
		return nil, err
	}

	// Checks if there are items returned
	if len(result.Items) == 0 {
		cw.Info(&cw.Logs{Code: "DynamoDBAPI", Message: "No data found"})
		return nil, errors.New("user does not exist")
	}

	// Unmarshal a map into actual user which front-end can uderstand as a JSON
	err = service.DynamoDBAttributeResponse(response, result.Items[0])
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "DynamoDBAttributeResponse", Message: "Failed to unmarshal user record"})
		return nil, err
	}

	return response, nil
}
