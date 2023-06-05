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
// request query, fetches the bus line record/information, and responds with a 200
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

	// Check whether the request queries are present
	if id_query == "" {
		err := errors.New("'id' parameter is not set")
		utility.Error(err, "APIError", "'id' is not implemented")

		return api.StatusBadRequest(err)
	}

	if name_query == "" {
		err := errors.New("'name' parameter is not set")
		utility.Error(err, "APIError", "'name' is not implemented")

		return api.StatusBadRequest(err)
	}

	// Fetch the existing bus line record/information
	bus, err := query.GetBusLine(ctx, id_query, name_query)
	if err != nil {
		utility.Error(err, "DynamoDBError", "failed to fetch the bus line information/record", utility.KVP{Key: "id", Value: id_query}, utility.KVP{Key: "name", Value: name_query})
		return api.StatusInternalServerError()
	}

	if bus == (schema.Bus{}) {
		err := errors.New("the bus line information you're trying to fetch is non-existent")
		utility.Error(err, "APIError", "the bus line does not exist", utility.KVP{Key: "id", Value: id_query}, utility.KVP{Key: "name", Value: name_query})

		return api.StatusBadRequest(err)
	}

	return api.StatusOK(bus)
}
