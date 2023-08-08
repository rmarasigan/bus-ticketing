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

// It receives the Amazon API Gateway event record data as input, fetches the
// bus line's record, and responds with a 200 OK HTTP Status.
//
// Method: GET
//
// Endpoint: https://{api_id}.execute-api.{region}.amazonaws.com/prod/bus/search?name=xxxxxx&company=xxxxxx
//
// Sample API Params:
//  name=Yellow Sunshine
//  company=Transit Bus Co
//
// Sample API Response:
// 	[
// 	  {
// 	    "id": "TRNSTBSC-875011",
// 	    "name": "Yellow Sunshine",
// 	    "owner": "Melissa Anderson",
// 	    "email": "melissa.anderson@example.com",
// 	    "address": "741 Oak Avenue, Suburb",
// 	    "company": "Transit Bus Co",
// 	    "mobile_number": "999-333-7777",
// 	    "date_created": "1687501112"
// 	  }
// 	]
func handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var (
		name_query    = request.QueryStringParameters["name"]
		company_query = request.QueryStringParameters["company"]
	)

	// Fetch a list of bus line information
	listOfBus, err := query.FilterBusLine(ctx, name_query, company_query)
	if err != nil {
		utility.Error(err, "DynamoDBError", "failed to filter the bus line")
		return api.StatusInternalServerError(err)
	}

	return api.StatusOK(listOfBus)
}
