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
// request query, fetches the bus unit record(s), and responds with a 200 OK HTTP
// Status.
//
// Method: GET
//
// Endpoint: https://{api_id}.execute-api.{region}.amazonaws.com/prod/bus-unit/get
//
// Sample API Params:
//  bus_id=BCBSCMPN-875011
//  code=BCBSCMPNBUS002
//
// Sample API Response:
// 	[
// 	  {
// 	    "bus_id": "BCBSCMPN-875011",
// 	    "code": "BCBSCMPNBUS002",
// 	    "active": true,
// 	    "min_capacity": 30,
// 	    "max_capacity": 60
// 	  }
// 	]
func handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var (
		code_query  = request.QueryStringParameters["code"]
		busId_query = request.QueryStringParameters["bus_id"]
	)

	// Fetch the existing bus unit record
	units, err := query.GetBusUnitRecords(ctx, code_query, busId_query)
	if err != nil {
		utility.Error(err, "DynamoDBError", "failed to fetch the bus unit record", utility.KVP{Key: "code", Value: code_query}, utility.KVP{Key: "bus_id", Value: busId_query})

		return api.StatusInternalServerError(err)
	}

	if len(units) == 0 {
		return api.StatusOK(api.Message{Custom: "no record(s) found"})
	}

	return api.StatusOK(units)
}
