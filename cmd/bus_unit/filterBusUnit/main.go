package main

import (
	"context"
	"errors"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rmarasigan/bus-ticketing/api"
	"github.com/rmarasigan/bus-ticketing/internal/app/query"
	"github.com/rmarasigan/bus-ticketing/internal/utility"
)

func main() {
	lambda.Start(handler)
}

// It receives the Amazon API Gateway event record as input, fetches the
// bus unit's record, and responds with a 200 OK HTTP Status.
//
// Endpoint:
//  https://{api_id}.execute-api.{region}.amazonaws.com/prod/bus-unit/search?bus_id=xxxxxx
//
// Sample API Params:
//  bus_id=RLBSW-856996
//
// Sample API Response:
// 	{
// 		"bus_id": "RLBSW-856996",
// 		"code": "RLBSWV1_0606",
// 		"active": true,
// 		"min_capacity": 40,
// 		"max_capacity": 50,
//		"date_created": "1687501761"
// 	}
func handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var (
		active       *bool
		code_query   = request.QueryStringParameters["code"]
		busId_query  = request.QueryStringParameters["bus_id"]
		active_query = request.QueryStringParameters["active"]
	)

	// Convert the 'active' query into a boolean value
	// if it is present in the request parameters.
	if active_query != "" {
		value, err := strconv.ParseBool(active_query)
		if err != nil {
			utility.Error(err, "StrConvError", "failed to convert 'active' string to a boolean value", utility.KVP{Key: "active_query", Value: active_query})
			return api.StatusBadRequest(errors.New("invalid 'active' parameter value"))
		}

		active = &value
	}

	listOfBusUnit, err := query.FilterBusUnit(ctx, code_query, busId_query, active)
	if err != nil {
		utility.Error(err, "DynamoDBError", "failed to filter the bus unit")
		return api.StatusInternalServerError(err)
	}

	return api.StatusOK(listOfBusUnit)
}
