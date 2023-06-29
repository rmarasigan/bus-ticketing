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
// request query and body, updates the bus route record and responds with a 200
// OK HTTP Status.
//
// Endpoint:
//  https://{api_id}.execute-api.{region}.amazonaws.com/prod/bus-route/update?id=xxxxx&bus_id=xxxxx
//
// Sample API Params:
//  id=RTRTB15001900880101
// 	bus_id=SNRSBSS-875011
//
// Sample API Payload:
// 	{
// 		"active": false
// 	}
//
// Sample API Response:
// 	{
// 		"id": "RTRTB15001900880101",
// 		"bus_id": "SNRSBSS-875011",
// 		"bus_unit_id": "SNRSBSSBUS002",
// 		"currency_code": "PHP",
// 		"rate": 90,
// 		"active": false,
// 		"departure_time": "15:00",
// 		"arrival_time": "19:00",
// 		"from_route": "Route A",
// 		"to_route": "Route B",
// 		"date_created": "1688010114"
// 	}
func handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var (
		route       schema.BusRoute
		id_query    = request.QueryStringParameters["id"]
		busId_query = request.QueryStringParameters["bus_id"]
	)

	err := route.IsEmptyPayload(request.Body)
	if err != nil {
		return api.StatusBadRequest(err)
	}

	// Unmarshal the received JSON-encoded data
	err = utility.ParseJSON([]byte(request.Body), &route)
	if err != nil {
		route.Error(err, "JSONError", "failed to unmarshal the JSON-encoded data",
			utility.KVP{Key: "payload", Value: request.Body})

		return api.StatusInternalServerError(err)
	}

	// Fetch the existing bus route record
	busRoute, err := query.GetBusRoute(ctx, id_query, busId_query)
	if err != nil {
		route.Error(err, "DynamoDBError", "failed to fetch the bus route record")
		return api.StatusInternalServerError(err)
	}

	if busRoute == (schema.BusRoute{}) {
		err := errors.New("the bus route you're trying to update is non-existent")
		route.Error(err, "APIError", "the bus route does not exist")

		return api.StatusBadRequest(err)
	}

	// Create a composite key that has both the partition/primary key
	// and the sort key of the item.
	var compositeKey = map[string]types.AttributeValue{
		"id":     &types.AttributeValueMemberS{Value: id_query},
		"bus_id": &types.AttributeValueMemberS{Value: busId_query},
	}

	// Construct the update builder
	busRoute = validate.UpdateBusRouteFields(route, busRoute)
	var update = expression.Set(expression.Name("currency_code"), expression.Value(busRoute.Currency)).
		Set(expression.Name("rate"), expression.Value(busRoute.Rate)).
		Set(expression.Name("active"), expression.Value(busRoute.Active)).
		Set(expression.Name("departure_time"), expression.Value(busRoute.DepartureTime)).
		Set(expression.Name("arrival_time"), expression.Value(busRoute.ArrivalTime)).
		Set(expression.Name("from_route"), expression.Value(busRoute.FromRoute)).
		Set(expression.Name("to_route"), expression.Value(busRoute.ToRoute))

	result, err := query.UpdateBusRoute(ctx, compositeKey, update)
	if err != nil {
		busRoute.Error(err, "DynamoDBError", "failed to update the bus route record")
		return api.StatusInternalServerError(err)
	}

	return api.StatusOK(result)
}
