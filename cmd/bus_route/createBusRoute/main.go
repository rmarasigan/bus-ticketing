package main

import (
	"context"
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
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
// request body, saves the validated request body to the DynamoDB Table, and
// responds with a 200 OK HTTP Status.
//
// Endpoint:
//  https://{api_id}.execute-api.{region}.amazonaws.com/prod/bus-route/create
//
// Sample API Payload:
// 	{
// 	  "rate": 120,
// 	  "available": true,
// 	  "currency_code": "PHP",
// 	  "bus_id": "BCBSCMPN-875011",
// 	  "bus_unit_id": "BCBSCMPNBUS001",
// 	  "departure_time": "15:00",
// 	  "arrival_time": "17:00",
// 	  "from_route": "Route A",
// 	  "to_route": "Route B"
// 	}
func handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var route schema.BusRoute

	// Unmarshal the received JSON-encoded data
	err := utility.ParseJSON([]byte(request.Body), &route)
	if err != nil {
		route.Error(err, "JSONError", "failed to unmarshal the JSON-encoded data", utility.KVP{Key: "payload", Value: request.Body})
		return api.StatusInternalServerError(err)
	}

	routeExist, err := validate.IsBusRouteExisting(ctx, route.SetFilter())
	if err != nil {
		route.Error(err, "IsBusRouteExisting", "failed to validate bus route if it exist")
		return api.StatusInternalServerError(err)
	}

	// If the bus route exists, stop the execution
	if routeExist {
		err := errors.New("already existing bus route")
		utility.Info("BusRouteExisting", err.Error(), utility.KVP{Key: "route", Value: route})
		return api.StatusBadRequest(err)
	}

	// Set default values of the bus route information
	route.SetValues()

	// Inserts a new bus route record to the DynamoDB
	err = query.CreateBusRoute(ctx, route)
	if err != nil {
		route.Error(err, "DynamoDBError", "failed to create a new bus route record")
		return api.StatusInternalServerError(err)
	}

	return api.StatusOKWithoutBody()
}
