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
// request query, fetches the bus route record(s), and responds with a 200 OK HTTP
// Status.
//
// Method: GET
//
// Endpoint: https://{api_id}.execute-api.{region}.amazonaws.com/prod/bus-route/get
//
// Sample API Params:
//  bus_id=SNRSBSS-875011
//  id=RTRTC15001900877753
//
// Sample API Response:
// 	[
// 	  {
// 	    "rate": 120,
// 	    "active": true,
// 	    "currency_code": "PHP",
// 	    "id": "RTRTC15001900877753",
// 	    "bus_id": "SNRSBSS-875011",
// 	    "bus_unit_id": "SNRSBSSBUS002",
// 	    "departure_time": "15:00",
// 	    "arrival_time": "17:00",
// 	    "from_route": "Route A",
// 	    "to_route": "Route B"
// 	  }
// 	]
func handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var (
		id_query    = request.QueryStringParameters["id"]
		busId_query = request.QueryStringParameters["bus_id"]
	)

	// Fetch the existing bus route record
	routes, err := query.GetBusRouteRecords(ctx, id_query, busId_query)
	if err != nil {
		utility.Error(err, "DynamoDBError", "failed to fetch the bus route record", utility.KVP{Key: "id", Value: id_query},
			utility.KVP{Key: "bus_id", Value: busId_query})

		return api.StatusInternalServerError(err)
	}

	if len(routes) == 0 {
		return api.StatusOK(api.Message{Custom: "no record(s) found"})
	}

	return api.StatusOK(routes)
}
