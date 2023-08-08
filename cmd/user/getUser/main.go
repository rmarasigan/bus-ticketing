package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rmarasigan/bus-ticketing/api"
	"github.com/rmarasigan/bus-ticketing/internal/app/query"
	"github.com/rmarasigan/bus-ticketing/internal/utility"
)

func main() {
	lambda.Start(handler)
}

// It receives the Amazon API Gateway event record data as input, validates the
// request query, fetches the user account record(s), and responds with
// a 200 OK HTTP Status.
//
// Method: GET
//
// Endpoint: https://{api_id}.execute-api.{region}.amazonaws.com/prod/user/account/get
//
// Sample API Params:
//  id=ADMN-878495
//  username=passwordabc
//
// Sample API Response:
// 	[
// 	  {
// 	    "id": "ADMN-878495",
// 	    "user_type": "ADMIN",
// 	    "first_name": "Emily",
// 	    "last_name": "Davis",
// 	    "username": "emilydavis",
// 	    "address": "321 Cedar Road",
// 	    "email": "emilydavis@example.com",
// 	    "mobile_number": "(407) 435-6841"
// 	  }
// 	]
func handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var (
		id_query       = request.QueryStringParameters["id"]
		username_query = request.QueryStringParameters["username"]
	)

	// Fetch the existing user account record
	accounts, err := query.GetUserAccountRecords(ctx, id_query, username_query)
	if err != nil {
		utility.Error(err, "DynamoDBError", "failed to fetch the user account record", utility.KVP{Key: "username", Value: username_query})
		return api.StatusInternalServerError(err)
	}

	if len(accounts) == 0 {
		return api.StatusOK(api.Message{Custom: "no record(s) found"})
	}

	return api.StatusOK(accounts)
}
