package main

import (
	"context"
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rmarasigan/bus-ticketing/api"
	"github.com/rmarasigan/bus-ticketing/api/schema"
	"github.com/rmarasigan/bus-ticketing/internal/app/query"
	"github.com/rmarasigan/bus-ticketing/internal/utility"
)

func main() {
	lambda.Start(handler)
}

// It receives the Amazon API Gateway event record data as input, validates the
// request query, fetches the user account record/information, and responds with
// a 200 OK HTTP Status.
//
// Endpoint:
//  https://{api_id}.execute-api.{region}.amazonaws.com/prod/user/account/get?id=xxxxx&username=xxxxx
//
// Sample API Params:
//  id=CSTMR-855048
// 	username=j.doe
//
// Sample API Response:
// 	{
// 		"id": "CSTMR-855048",
// 		"user_type": "CUSTOMER",
// 		"first_name": "John",
// 		"last_name": "Doe",
// 		"username": "j.doe",
// 		"address": "Långbro, Stockholm",
// 		"email": "j.doe@outlook.com",
// 		"mobile_number": "0586-4404205"
// 	}
func handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var (
		id_query       = request.QueryStringParameters["id"]
		username_query = request.QueryStringParameters["username"]
	)

	// Fetch the existing user account record/information
	account, err := query.GetUserAccount(ctx, id_query, username_query)
	if err != nil {
		utility.Error(err, "DynamoDBError", "failed to fetch the user account record", utility.KVP{Key: "username", Value: username_query})
		return api.StatusInternalServerError()
	}

	if account == (schema.User{}) {
		err := errors.New("the account you're trying to fetch is non-existent")
		utility.Error(err, "APIError", "the account does not exist")

		return api.StatusBadRequest(err)
	}

	return api.StatusOK(account)
}
