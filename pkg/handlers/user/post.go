package user

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

// Post is the Users API request POST method that will process the "create" or "login" type request.
// To process the request, request query "type" is required and the value must be either "create" or "login".
func Post(ctx context.Context, request *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	// Creates Dynamodb Session
	service.DynamodbSession()

	tablename := os.Getenv("USERS_TABLE")
	queryType := request.QueryStringParameters["type"]

	if tablename == "" {
		return api.StatusUnhandledRequest(errors.New("dynamodb table on env is not implemented"))
	}

	if queryType == "" {
		return api.StatusBadRequest(errors.New("method.request.querystring.type is not set"))
	}

	switch queryType {
	case "create":
		return CreateUser(tablename, []byte(request.Body), service.DynamoDBClient)

	case "login":
		return LogIn(tablename, []byte(request.Body), service.DynamoDBClient)

	default:
		return api.StatusUnhandledRequest(errors.New("request not implemented"))
	}
}

// CreateUser creates a new user account and the user type can be Admin (1) or Customer (2).
func CreateUser(tablename string, body []byte, svc dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {
	user := new(models.User)

	err := api.ParseJSON(body, user)
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "ParseJSONError", Message: "Failed to parse user json"})
		return api.StatusBadRequest(errors.New("error parsing user json"))
	}

	// Validate if the required fields are not empty.
	isValid := validate.CreateAccount(user)
	if isValid != "" {
		err := errors.New(isValid)

		cw.Error(err, &cw.Logs{Code: "CreateUserAccount", Message: "Validate creation of account"})
		return api.StatusUnhandledRequest(err)
	}

	// Checks whether the username exist or not.
	usernameExist, err := ValidateUsername(tablename, user.Username, svc)
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "ValidateUsername", Message: "Failed to validate username"})
		return api.StatusBadRequest(err)
	}

	// Returns error message of "username exist".
	if usernameExist {
		return api.StatusBadRequest(errors.New("username exist"))
	}

	user.SetValues()

	// Converting the record to dynamodb.AttributeValue type.
	value, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "DynamoDBMarshalMap", Message: "Failed to marshal user"})
		return api.StatusBadRequest(err)
	}

	input := &dynamodb.PutItemInput{
		Item:      value,
		TableName: aws.String(tablename),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "DynamoDBPutItem", Message: "Failed to add item to the table"},
			kvp.Attribute{Key: "tablename", Value: tablename})

		return api.StatusBadRequest(err)
	}

	// Set the response message.
	response := map[string]string{
		"username":   user.Username,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"email":      user.Email,
	}

	return api.StatusOK(response)
}

// LogIn authenticates the user if the account exist and confirming credentials.
func LogIn(tablename string, body []byte, svc dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {
	user := new(models.User)
	response := new(models.UserResponse)

	err := api.ParseJSON(body, user)
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "ParseJSONError", Message: "Failed to parse user json"})
		return api.StatusBadRequest(errors.New("error parsing user json"))
	}

	// Validate if the required fields are not empty.
	isValid := validate.LogIn(user)
	if isValid != "" {
		err := errors.New(isValid)

		cw.Error(err, &cw.Logs{Code: "Login", Message: "Validate login user credentials"})
		return api.StatusUnhandledRequest(err)
	}

	// Create the names list projection of names to project.
	// SELECT id, type, first_name, last_name, username, address, email, mobile_number
	projection := expression.NamesList(expression.Name("id"), expression.Name("type"),
		expression.Name("first_name"), expression.Name("last_name"),
		expression.Name("username"), expression.Name("address"),
		expression.Name("email"), expression.Name("mobile_number"))

	// Construct the filter builder with a name and value.
	// WHERE username = username_value AND password = password_value
	filter := expression.Name("username").Equal(expression.Value(user.Username)).And(expression.Name("password").Equal(expression.Value(user.Password)))

	// Using the filter and projections to create a DynamoDB expression.
	expr, err := expression.NewBuilder().WithFilter(filter).WithProjection(projection).Build()
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "DynamoDBExpression", Message: "Failed to build an expression"})
		return api.StatusBadRequest(err)
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
		return api.StatusBadRequest(err)
	}

	// Checks if there are items returned
	if len(result.Items) == 0 {
		cw.Info(&cw.Logs{Code: "DynamoDBAPI", Message: "No data found"})
		err := errors.New("the username or password you entered is incorrect")

		return api.StatusBadRequest(err)
	}

	// Unmarshal a map into actual user which front-end can uderstand as a JSON
	err = dynamodbattribute.UnmarshalMap(result.Items[0], response)
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "DynamoDBError", Message: "Failed to unmarshal user record"})
		return api.StatusBadRequest(err)
	}

	return api.StatusOK(response)
}
