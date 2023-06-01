package main

import (
	"context"
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/rmarasigan/bus-ticketing/api"
	"github.com/rmarasigan/bus-ticketing/api/schema"
	"github.com/rmarasigan/bus-ticketing/internal/app/query"
	"github.com/rmarasigan/bus-ticketing/internal/app/validate"
	"github.com/rmarasigan/bus-ticketing/internal/utility"
)

func main() {
	lambda.Start(handler)
}

// It receives the Amazon API Gateway event record data as input, validates the
// request body, checks the validated request body if the user credentials are
// valid or not, updates the user account’s last login, and responds with a 200
// OK HTTP Status.
//
// Endpoint:
//  https://{api_id}.execute-api.{region}.amazonaws.com/prod/user/login
//
// Sample API Payload:
// 	{
// 		"username": "j.doe",
// 		"password": "j.doe1234"
// 	}
//
// Sample API Response:
// 	{
// 		"id": "CSTMR-855048",
// 		"user_type": "CUSTOMER",
// 		"first_name": "John",
// 		"last_name": "Doe",
// 		"username": "j.doe",
// 		"address": "South Calorina",
// 		"email": "j.doe@outlook.com",
// 		"mobile_number": "11223344556"
// 	}
func handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var user = new(schema.User)

	// Umarshal the received JSON-encoded data
	err := utility.ParseJSON([]byte(request.Body), user)
	if err != nil {
		user.Error(err, "JSONError", "failed to unmarshal the JSON-encoded data",
			utility.KVP{Key: "request", Value: request.Body})

		return api.StatusInternalServerError()
	}

	// Validate if the required fields are not empty
	err = validate.LogInFields(*user)
	if err != nil {
		user.Error(err, "LogInFields", "failed to login the user account")
		return api.StatusBadRequest(err)
	}

	// Checks whether the user credentials are valid or not
	existing, account, err := validate.UserAccountExists(ctx, user.Username, user.Password)
	if err != nil {
		user.Error(err, "UserAccountExists", "failed to validate user account credentials")
		return api.StatusInternalServerError()
	}

	// If the user account does not exist, return a 400 BadRequest HTTP Status
	if !existing {
		err := errors.New("the username or password you entered is incorrect")
		account.Error(err, "UserAccountExists", "incorrect credentials")

		return api.StatusBadRequest(err)
	}

	// Create a composite key that has both the partition/primary key
	// and the sort key of the item.
	var compositKey = map[string]types.AttributeValue{
		"id":       &types.AttributeValueMemberS{Value: account.ID},
		"username": &types.AttributeValueMemberS{Value: account.Username},
	}

	// Construct the update builder
	update := expression.Set(expression.Name("last_login"), expression.Value(account.LastLogIn()))

	// Update the User’s Last Login into the DynamoDB Table
	_, err = query.UpdateUserAcccount(ctx, compositKey, update)
	if err != nil {
		account.Error(err, "DynamoDBError", "failed to update the user last login")
		return api.StatusInternalServerError()
	}

	return api.StatusOK(account)
}
