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
// request query, fetches the booking record(s), and responds with a 200
// OK HTTP Status.
//
// Method: GET
//
// Endpoint: https://{api_id}.execute-api.{region}.amazonaws.com/prod/bookings/get
//
// Sample API Params:
//  id=bd866a7e-34cd-4ea1-8411-5351a6b76ffd
//  bus_route_id=RTBRTC15001900884691
//
// Sample API Response:
// 	[
// 	  {
// 	    "id": "bd866a7e-34cd-4ea1-8411-5351a6b76ffd",
// 	    "user_id": "ADMN-878495",
// 	    "bus_id": "BCBSCMPN-884690",
// 	    "bus_route_id": "RTBRTC15001900884691",
// 	    "status": "PENDING",
// 	    "seat_number": "23,24,25,26",
// 	    "travel_date": "2023-07-06 19:30",
// 	    "date_created": "2023-07-05 07:48:26",
// 	    "cancelled": {
// 	      "id": "",
// 	      "booking_id": "",
// 	      "reason": "",
// 	      "cancelled_by": "",
// 	      "date_cancelled": ""
// 	    },
// 	    "timestamp": "2023-07-01 10:30"
// 	  }
// 	]
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
