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
// request query, fetches the bus line record(s), and responds with a 200
// OK HTTP Status.
//
// Endpoint:
//  https://{api_id}.execute-api.{region}.amazonaws.com/prod/bus/get?name=xxxxxx&id=xxxxxx
//
// Sample API Params:
//  id=RLBSW-856996
// 	name=Thunder Rail Bus Line
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
// 	}
func handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var (
		id_query   = request.QueryStringParameters["id"]
		name_query = request.QueryStringParameters["name"]
	)

	// Fetch the existing bus line record
	busList, err := query.GetBusLineRecords(ctx, id_query, name_query)
	if err != nil {
		utility.Error(err, "DynamoDBError", "failed to fetch the bus line record", utility.KVP{Key: "id", Value: id_query},
			utility.KVP{Key: "name", Value: name_query})

		return api.StatusInternalServerError(err)
	}

	if len(busList) == 0 {
		return api.StatusOK(api.Message{Custom: "no record(s) found"})
	}

	return api.StatusOK(busList)
}
