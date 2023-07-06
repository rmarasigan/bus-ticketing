package query

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/rmarasigan/bus-ticketing/api/schema"
	"github.com/rmarasigan/bus-ticketing/internal/app/env"
	awswrapper "github.com/rmarasigan/bus-ticketing/internal/aws_wrapper"
	"github.com/rmarasigan/bus-ticketing/internal/trail"
)

// getUserAccount returns the user account information.
func getUserAccount(ctx context.Context, tablename, id, username string) (schema.User, error) {
	var user schema.User

	// Create a composite key expression
	key := expression.KeyAnd(expression.Key("username").Equal(expression.Value(username)), expression.Key("id").Equal(expression.Value(id)))

	// Create a names list representing the list of item attribute names
	// to be returned.
	var namesList = []expression.NameBuilder{
		expression.Name("user_type"),
		expression.Name("first_name"),
		expression.Name("last_name"),
		expression.Name("username"),
		expression.Name("email"),
		expression.Name("mobile_number"),
		expression.Name("address"),
	}

	// SELECT id, user_type, first_name, last_name, username, address, email, mobile_number
	projection := expression.NamesList(expression.Name("id"), namesList...)

	// Build an expression to retrieve the item from the DynamoDB
	expr, err := expression.NewBuilder().WithKeyCondition(key).WithProjection(projection).Build()
	if err != nil {
		return user, err
	}

	// Build the query input parameter
	params := &dynamodb.QueryInput{
		TableName:                 aws.String(env.USERS_TABLE),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ProjectionExpression:      expr.Projection(),
		KeyConditionExpression:    expr.KeyCondition(),
	}

	result, err := awswrapper.DynamoDBQuery(ctx, params)
	if err != nil {
		return user, err
	}

	// Unmarshal a map into actual user struct which the front-end can
	// understand as a JSON.
	if result.Count > 0 {
		err := awswrapper.DynamoDBUnmarshalMap(&user, result.Items[0])
		if err != nil {
			return user, err
		}
	}

	return user, nil
}

// GetUserAccountRecords checks if the DynamoDB Table is configured on the environment, and
// returns either the specific user account or a list of user account records.
func GetUserAccountRecords(ctx context.Context, id, username string) ([]schema.User, error) {
	var (
		users     []schema.User
		tablename = env.USERS_TABLE
	)

	// Check if the DynamoDB Table is configured
	if tablename == "" {
		trail.Error("dynamodb USERS_TABLE is not configured on the environment")
		err := errors.New("dynamodb USERS_TABLE environment variable is not set")

		return users, err
	}

	// ********** Fetching a specific user account record ********** //
	if id != "" && username != "" {
		user, err := getUserAccount(ctx, tablename, id, username)
		if err != nil {
			return users, err
		}

		if user == (schema.User{}) {
			return users, nil
		}

		users = append(users, user)
		return users, nil
	}

	// **************** List of user account records **************** //
	// Create a names list representing the list of item attribute names
	// to be returned.
	var namesList = []expression.NameBuilder{
		expression.Name("user_type"),
		expression.Name("first_name"),
		expression.Name("last_name"),
		expression.Name("username"),
		expression.Name("email"),
		expression.Name("mobile_number"),
		expression.Name("address"),
	}

	// SELECT id, user_type, first_name, last_name, username, address, email, mobile_number
	projection := expression.NamesList(expression.Name("id"), namesList...)

	// Build an expression to retrieve the item from the DynamoDB
	expr, err := expression.NewBuilder().WithProjection(projection).Build()
	if err != nil {
		return users, err
	}

	// Use the build expression to populate the DynamoDB Scan API
	var params = &dynamodb.ScanInput{
		TableName:                 aws.String(tablename),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ProjectionExpression:      expr.Projection(),
	}

	result, err := awswrapper.DynamoDBScan(ctx, params)
	if err != nil {
		return nil, err
	}

	if result.Count > 0 {
		// Unmarshal a map into actual user struct which the front-end can
		// understand as a JSON.
		err = awswrapper.DynamoDBUnmarshalListOfMaps(&users, result.Items)
		if err != nil {
			return nil, err
		}
	}

	return users, nil
}

// GetUserAccountById checks if the DynamoDB Table is configured on the environment, and
// fetch the user account by id and returns the user account information.
func GetUserAccountById(ctx context.Context, id string) (schema.User, error) {
	var (
		user      schema.User
		tablename = env.USERS_TABLE
	)

	// Check if the DynamoDB Table is configured
	if tablename == "" {
		trail.Error("dynamodb USERS_TABLE is not configured on the environment")
		err := errors.New("dynamodb USERS_TABLE environment variable is not set")

		return user, err
	}

	// Construct the filter builder with a name that contains a specified value.
	// WHERE id = id_value
	filter := expression.Name("id").Equal(expression.Value(id))

	result, err := FilterItems(ctx, tablename, filter)
	if err != nil {
		return user, err
	}

	if result.Count > 0 {
		// Unmarshal a map into actual user struct which the front-end can
		// understand as a JSON.
		err = awswrapper.DynamoDBUnmarshalMap(&user, result.Items[0])
		if err != nil {
			return user, err
		}
	}

	return user, nil
}

// CreateUserAccount checks if the DynamoDB Table is configured on the environment, and
// creates a new user account.
func CreateUserAccount(ctx context.Context, data interface{}) error {
	var tablename = env.USERS_TABLE

	// Check if the DynamoDB Table is configured
	if tablename == "" {
		trail.Error("dynamodb USERS_TABLE is not configured on the environment")
		err := errors.New("dynamodb USERS_TABLE environment variable is not set")

		return err
	}

	// Save the User information into the DynamoDB Table
	err := InsertItem(ctx, tablename, data)
	if err != nil {
		trail.Error("failed to insert a new user")
		return err
	}

	return nil
}

// UpdateUserAcccount checks if the DynamoDB Table is configured on the environment and
// updates the user account’s information or record.
func UpdateUserAcccount(ctx context.Context, key map[string]types.AttributeValue, update expression.UpdateBuilder) (schema.User, error) {
	var (
		user      schema.User
		tablename = env.USERS_TABLE
	)

	// Check if the DynamoDB Table is configured
	if tablename == "" {
		trail.Error("dynamodb USERS_TABLE is not configured on the environment")
		err := errors.New("dynamodb USERS_TABLE environment variable is not set")

		return user, err
	}

	result, err := UpdateItem(ctx, tablename, key, update)
	if err != nil {
		trail.Error("failed to update the user account")
		return user, err
	}

	// Unmarshal a map into actual user struct which the front-end can
	// understand as a JSON.
	err = awswrapper.DynamoDBUnmarshalMap(&user, result.Attributes)
	if err != nil {
		return user, err
	}

	return user, nil
}
