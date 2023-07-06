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

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var (
		status_query  = request.QueryStringParameters["status"]
		busId_query   = request.QueryStringParameters["bus_id"]
		routeId_query = request.QueryStringParameters["route_id"]
	)

	bookings, err := query.FilterBookings(ctx, busId_query, routeId_query, status_query)
	if err != nil {
		utility.Error(err, "DynamoDBError", "failed to filter the bookings")
		return api.StatusInternalServerError(err)
	}

	return api.StatusOK(bookings)
}
