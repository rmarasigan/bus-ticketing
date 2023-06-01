package validate

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/rmarasigan/bus-ticketing/api/schema"
	"github.com/rmarasigan/bus-ticketing/internal/app/env"
	awswrapper "github.com/rmarasigan/bus-ticketing/internal/aws_wrapper"
	"github.com/rmarasigan/bus-ticketing/internal/trail"
)

// CreateUserAccount validates if the required fields are empty or not.
//
// Fields that are validated:
//  user_type, first_name, last_name, username, password, address,
//  email, mobile_number
func CreateUserAccount(user *schema.User) error {
	var fields []string

	if user.UserType == "" {
		fields = append(fields, "user_type")
	}

	if user.FirstName == "" {
		fields = append(fields, "first_name")
	}

	if user.LastName == "" {
		fields = append(fields, "last_name")
	}

	if user.Username == "" {
		fields = append(fields, "username")
	}

	if user.Password == "" {
		fields = append(fields, "password")
	}

	if user.Address == "" {
		fields = append(fields, "address")
	}

	if user.Email == "" {
		fields = append(fields, "email")
	}

	if user.MobileNumber == "" {
		fields = append(fields, "mobile_number")
	}

	if len(fields) != 0 {
		return fmt.Errorf("missing %s field(s)", strings.Join(fields, ", "))
	}

	return nil
}

// LogInFields validates if the required fields are empty or not.
//
// Fields that are validated:
//  username, password
func LogInFields(user schema.User) error {
	var fields []string

	if user.Username == "" {
		fields = append(fields, "username")
	}

	if user.Password == "" {
		fields = append(fields, "password")
	}

	if len(fields) > 0 {
		return fmt.Errorf("missing %s field(s)", strings.Join(fields, ", "))
	}

	return nil
}

// UpdateUserAccountFields validates if the field that are going to be updated
// are empty or not to set its previous value.
//
// Fields that are validated:
//  first_name, last_name, address, email, mobile_number
func UpdateUserAccountFields(user schema.User, old schema.User) schema.User {
	if user.FirstName == "" {
		user.FirstName = old.FirstName
	}

	if user.LastName == "" {
		user.LastName = old.LastName
	}

	if user.Address == "" {
		user.Address = old.Address
	}

	if user.Email == "" {
		user.Email = old.Email
	}

	if user.MobileNumber == "" {
		user.MobileNumber = old.MobileNumber
	}

	return user
}

// IsUsernameExisting checks if the DynamoDB Table is configured on the environment, and
// returns a boolean and error value to check whether the username already exist or not.
func IsUsernameExisting(ctx context.Context, username string) (bool, error) {
	var tablename = env.USERS_TABLE

	// Check if the DynamoDB Table is configured
	if tablename == "" {
		trail.Error("dynamodb USERS_TABLE is not configured on the environment")
		err := errors.New("dynamodb USERS_TABLE environment variable is not set")

		return false, err
	}

	// Create a key expression
	key := expression.Key("username").Equal(expression.Value(username))

	// Build an expression to retrieve item from the DynamoDB
	expr, err := expression.NewBuilder().WithKeyCondition(key).Build()
	if err != nil {
		return false, err
	}

	// Build the query params parameter
	params := &dynamodb.QueryInput{
		TableName:                 aws.String(env.USERS_TABLE),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
	}

	result, err := awswrapper.DynamoDBQuery(ctx, params)
	if err != nil {
		return false, err
	}

	return (result.Count > 0), nil
}

// UserAccountExists checks if the DynamoDB Table is configured on the environment, and
// returns a boolean and error value to check whether the user account credentials are correct or not.
func UserAccountExists(ctx context.Context, username, password string) (bool, schema.User, error) {
	var (
		user      schema.User
		tablename = env.USERS_TABLE
	)

	// Check if the DynamoDB Table is configured
	if tablename == "" {
		trail.Error("dynamodb USERS_TABLE is not configured on the environment")
		err := errors.New("dynamodb USERS_TABLE environment variable is not set")

		return false, user, err
	}

	// Create a key expression
	key := expression.Key("username").Equal(expression.Value(username))

	// Create a names list representing the list of item attribute names
	// to be returned.
	var namesList = []expression.NameBuilder{
		expression.Name("user_type"),
		expression.Name("first_name"),
		expression.Name("last_name"),
		expression.Name("username"),
		expression.Name("address"),
		expression.Name("email"),
		expression.Name("mobile_number"),
	}

	// SELECT id, user_type, first_name, last_name, username, address, email, mobile_number
	projection := expression.NamesList(expression.Name("id"), namesList...)

	// Construct the filter builder with a name and value.
	// WHERE password = password_value
	filter := expression.Name("password").Equal(expression.Value(password))

	// Build an expression to retrieve the item from the DynamoDB
	expr, err := expression.NewBuilder().WithKeyCondition(key).WithFilter(filter).WithProjection(projection).Build()
	if err != nil {
		return false, user, err
	}

	// Build the query params parameter
	params := &dynamodb.QueryInput{
		TableName:                 aws.String(tablename),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		KeyConditionExpression:    expr.KeyCondition(),
	}

	result, err := awswrapper.DynamoDBQuery(ctx, params)
	if err != nil {
		return false, user, err
	}

	// Unmarshal a map into actual user which front-end can uderstand as a JSON
	if result.Count > 0 {
		err := awswrapper.DynamoDBUnmarshalMap(&user, result.Items[0])
		if err != nil {
			return false, user, err
		}
	}

	return (result.Count > 0), user, nil
}
