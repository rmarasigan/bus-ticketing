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
// request query and body, updates the user account’s information/record and responds
// with a 200 OK HTTP Status.
//
// Method: POST
//
// Endpoint: https://{api_id}.execute-api.{region}.amazonaws.com/prod/user/account/update?id=xxxxx&username=xxxxx
//
// Sample API Params:
//  id=ADMN-878495
//  username=passwordabc
//
// Sample API Payload:
// 	{
// 		"address": "Långbro, Stockholm",
// 		"mobile_number": "0586-4404205"
// 	}
//
// Sample API Response:
// 	{
// 	  "id": "ADMN-878495",
// 	  "user_type": "ADMIN",
// 	  "first_name": "Emily",
// 	  "last_name": "Davis",
// 	  "username": "emilydavis",
// 	  "password": "passwordabc",
// 	  "address": "Långbro, Stockholm",
// 	  "email": "emilydavis@example.com",
// 	  "mobile_number": "0586-4404205",
// 	  "date_created": "1687849585"
// 	}
func handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var (
		user           = new(schema.User)
		id_query       = request.QueryStringParameters["id"]
		username_query = request.QueryStringParameters["username"]
	)

	err := user.IsEmptyPayload(request.Body)
	if err != nil {
		return api.StatusBadRequest(err)
	}

	// Unmarshal the received JSON-encoded data
	err = utility.ParseJSON([]byte(request.Body), user)
	if err != nil {
		user.Error(err, "JSONError", "failed to unmarshal the JSON-encoded data",
			utility.KVP{Key: "payload", Value: request.Body})
		return api.StatusInternalServerError(err)
	}

	// Fetch the existing user account record
	accounts, err := query.GetUserAccountRecords(ctx, id_query, username_query)
	if err != nil {
		user.Error(err, "DynamoDBError", "failed to fetch the user account record")
		return api.StatusInternalServerError(err)
	}

	account := accounts[0]
	if account == (schema.User{}) {
		err := errors.New("the account you're trying to update is non-existent")
		user.Error(err, "APIError", "the account does not exist")

		return api.StatusBadRequest(err)
	}

	// Create a composite key that has both the partition/primary key
	// and the sort key of the item.
	var compositKey = map[string]types.AttributeValue{
		"id":       &types.AttributeValueMemberS{Value: id_query},
		"username": &types.AttributeValueMemberS{Value: username_query},
	}

	// Construct the update builder
	account = validate.UpdateUserAccountFields(*user, account)
	var update = expression.Set(expression.Name("first_name"), expression.Value(account.FirstName)).
		Set(expression.Name("last_name"), expression.Value(account.LastName)).
		Set(expression.Name("address"), expression.Value(account.Address)).
		Set(expression.Name("email"), expression.Value(account.Email)).
		Set(expression.Name("mobile_number"), expression.Value(account.MobileNumber))

	result, err := query.UpdateUserAcccount(ctx, compositKey, update)
	if err != nil {
		account.Error(err, "DynamoDBError", "failed to update the user account record")
		return api.StatusInternalServerError(err)
	}

	return api.StatusOK(result)
}
