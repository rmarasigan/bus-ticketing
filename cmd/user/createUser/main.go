package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
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
// request body, saves the validated request body to the DynamoDB Table, and
// responds with a 200 OK HTTP Status.
//
// Method: POST
//
// Endpoint: https://{api_id}.execute-api.{region}.amazonaws.com/prod/user/create
//
// Sample API Payload:
//	{
//	  "user_type": "1",
//	  "first_name": "Emily",
//	  "last_name": "Davis",
//	  "username": "emilydavis",
//	  "password": "passwordabc",
//	  "address": "321 Cedar Road",
//	  "email": "emilydavis@example.com",
//	  "mobile_number": "4449876543"
//	}
func handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var user = new(schema.User)

	// Umarshal the received JSON-encoded data
	err := utility.ParseJSON([]byte(request.Body), user)
	if err != nil {
		user.Error(err, "JSONError", "failed to unmarshal the JSON-encoded data",
			utility.KVP{Key: "payload", Value: request.Body})

		return api.StatusInternalServerError(err)
	}

	// Checks whether the username exist or not
	usernameExist, err := validate.IsUsernameExisting(ctx, user.Username)
	if err != nil {
		user.Error(err, "IsUsernameExisting", "failed to validate username if it exist")
		return api.StatusInternalServerError(err)
	}

	// If the username exists, return a 400 BadRequest HTTP Status
	if usernameExist {
		err := fmt.Errorf("%s username already exist", user.Username)
		user.Error(err, "IsUsernameExisting", "already existing username")

		return api.StatusBadRequest(err)
	}

	// Set default values of user account information
	user.SetValues()

	// Inserts a new user account to the DynamoDB
	err = query.CreateUserAccount(ctx, user)
	if err != nil {
		user.Error(err, "DynamoDBError", "failed to create a new account")
		return api.StatusInternalServerError(err)
	}

	return api.StatusOKWithoutBody()
}
