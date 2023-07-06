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
		id_query         = request.QueryStringParameters["id"]
		busRouteId_query = request.QueryStringParameters["bus_route_id"]
	)

	bookings, err := query.GetBookingRecords(ctx, id_query, busRouteId_query)
	if err != nil {
		utility.Error(err, "DynamoDBError", "failed to fetch the booking record")
		return api.StatusInternalServerError(err)
	}

	if len(bookings) == 0 {
		return api.StatusOK(api.Message{Custom: "no record(s) found"})
	}

	return api.StatusOK(bookings)
}
