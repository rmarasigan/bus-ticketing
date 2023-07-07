package main

import (
	"context"
	"errors"
	"strconv"

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

// It receives the Amazon API Gateway event record as input, fetches the
// bus unit route records, and responds with a 200 OK HTTP Status.
//
// Endpoint:
//  https://{api_id}.execute-api.{region}.amazonaws.com/prod/bus-route/search?bus_id=xxxxxx
//
// Sample API Params:
//  bus_id=SNRSBSS-875011
//
// Sample API Response:
// 	[
//	 {
// 		"id": "RTRTB15001900877732",
// 		"bus_id": "SNRSBSS-875011",
// 		"bus_unit_id": "SNRSBSSBUS002",
// 		"currency_code": "PHP",
// 		"rate": 90,
// 		"active": true,
// 		"departure_time": "15:00",
// 		"arrival_time": "19:00",
// 		"from_route": "Route A",
// 		"to_route": "Route B"
//	 }
// 	]
func handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var (
		active       *bool
		route        schema.BusRouteFilter
		active_query = request.QueryStringParameters["active"]
	)

	// Convert the 'active' query into a boolean value
	// if it is present in the request parameters.
	if active_query != "" {
		value, err := strconv.ParseBool(active_query)
		if err != nil {
			utility.Error(err, "StrConvError", "failed to convert 'active' string to a boolean value", utility.KVP{Key: "active_query", Value: active_query})
			return api.StatusBadRequest(errors.New("invalid 'available' parameter value"))
		}

		active = &value
	}

	route.Active = active
	route.BusID = request.QueryStringParameters["bus_id"]
	route.BusUnitID = request.QueryStringParameters["bus_unit_id"]
	route.Arrival = request.QueryStringParameters["arrival_time"]
	route.Departure = request.QueryStringParameters["departure_time"]
	route.ToRoute = request.QueryStringParameters["to_route"]
	route.FromRoute = request.QueryStringParameters["from_route"]

	listOfBusRoute, err := query.FilterBusRoute(ctx, route)
	if err != nil {
		utility.Error(err, "DynamoDBError", "failed to filter the bus route")
		return api.StatusInternalServerError(err)
	}

	return api.StatusOK(listOfBusRoute)
}
