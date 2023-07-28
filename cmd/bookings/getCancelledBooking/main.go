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
// request query, fetches the cancelled booking record, and responds with a 200
// OK HTTP Status.
//
// Method: GET
//
// Endpoint: https://{api_id}.execute-api.{region}.amazonaws.com/prod/bookings/cancelled/get
//
// Sample API Params:
//  booking_id=ce4e0245-b772-47f8-92fc-0d70cbd511c0
//
// Sample API Response:
// 	[
// 	  {
// 	    "id": "053607ed-3dc6-40a3-aea4-d7e87fd015f6",
// 	    "booking_id": "ce4e0245-b772-47f8-92fc-0d70cbd511c0",
// 	    "reason": "sample reason",
// 	    "cancelled_by": "ADMN-878495",
// 	    "date_cancelled": "2023-07-05 04:16:41"
// 	  }
// 	]
func handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var bookingId_query = request.QueryStringParameters["booking_id"]

	cancelledBookings, err := query.GetCancelledBookingRecords(ctx, bookingId_query)
	if err != nil {
		utility.Error(err, "DynamoDBError", "failed to fetch the cancelled booking record")
		return api.StatusInternalServerError(err)
	}

	if len(cancelledBookings) == 0 {
		return api.StatusOK(api.Message{Custom: "no record(s) found"})
	}

	return api.StatusOK(cancelledBookings)
}
