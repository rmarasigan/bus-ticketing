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
// Endpoint:
//  https://{api_id}.execute-api.{region}.amazonaws.com/prod/bus/search?name=xxxxxx&company=xxxxxx
//
// Sample API Params:
// 	name=Thunder Rail Bus Line
//  company=Rail Bus Way
//
// Sample API Response:
// 	{
// 		"id": "RLBSW-856996",
// 		"name": "Thunder Rail Bus Line",
// 		"owner": "Thando Oyibo Emmett",
// 		"email": "thando.emmet@outlook.com",
// 		"address": "1986 Bogisich Junctions, Hamillhaven, Kansas",
// 		"company": "Rail Bus Way",
// 		"mobile_number": "+1-335-908-1432"
// 		"date_created": "1685699666"
// 	}
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
