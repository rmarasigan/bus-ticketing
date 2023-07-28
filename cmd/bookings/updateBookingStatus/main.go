package main

import (
	"context"
	"encoding/json"
	"errors"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rmarasigan/bus-ticketing/api"
	"github.com/rmarasigan/bus-ticketing/api/schema"
	"github.com/rmarasigan/bus-ticketing/internal/app/query"
	"github.com/rmarasigan/bus-ticketing/internal/app/validate"
	awswrapper "github.com/rmarasigan/bus-ticketing/internal/aws_wrapper"
	"github.com/rmarasigan/bus-ticketing/internal/utility"
)

func main() {
	lambda.Start(handler)
}

// It receives the Amazon API Gateway event record data as input, validates the
// request query and body, updates the booking record and responds with a 200 OK
// HTTP Status without body.
//
// Method: POST
//
// Endpoint: https://{api_id}.execute-api.{region}.amazonaws.com/prod/bus-route/update?id=xxxxx&bus_route_id=xxxxx
//
// Sample API Params:
//  id=bd866a7e-34cd-4ea1-8411-5351a6b76ffd
//  bus_route_id=RTBRTC15001900884691
//
// Sample API Payload:
// 	{
// 	  "status": "CONFIRMED",
// 	  "seat_number": "23,24,25"
// 	}
func handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var (
		booking       schema.Bookings
		eventbus      = os.Getenv("EVENT_BUS")
		id_query      = request.QueryStringParameters["id"]
		routeId_query = request.QueryStringParameters["bus_route_id"]
	)

	err := booking.IsEmptyPayload(request.Body)
	if err != nil {
		return api.StatusBadRequest(err)
	}

	// Check if the EventBridge Event Bus is configured
	if eventbus == "" {
		err := errors.New("eventbridge EVENT_BUS environment variable is not set")
		booking.Error(err, "EventBridgeError", "eventbridge EVENT_BUS is not configured on the environment")

		return api.StatusInternalServerError(err)
	}

	// Unmarshal the received JSON-encoded data
	err = utility.ParseJSON([]byte(request.Body), &booking)
	if err != nil {
		booking.Error(err, "JSONError", "failed to unmarshal the JSON-encoded data", utility.KVP{Key: "payload", Value: request.Body})
		return api.StatusInternalServerError(err)
	}

	// ********************************************************************* //
	// ***************** Fetch and validate booking record ***************** //
	// ********************************************************************* //
	// 1. Fetch the existing booking record
	records, err := query.GetBookingRecords(ctx, id_query, routeId_query)
	if err != nil {
		booking.Error(err, "DynamoDBError", "failed to fetch the booking record")
		return api.StatusInternalServerError(err)
	}

	// 2. Check if there is an existing record
	record := records[0]
	if record == (schema.Bookings{}) {
		err := errors.New("the booking record you're trying to update is non-existent")
		booking.Error(err, "APIError", "the booking record does not exist")

		return api.StatusBadRequest(err)
	}

	// 3. Return an error if thay want to re-confirm the booking
	// that was already cancelled.
	if booking.Status == booking.Status.Confirmed() && *record.IsCancelled {
		err := errors.New("this booking has been cancelled and cannot be re-confirmed")
		booking.Error(err, "APIError", "booking confirmation failed")

		return api.StatusBadRequest(err)
	}

	// 4. If the booking status is cancelled, set the "is_cancelled"
	// field automatically to "true".
	if booking.Status == booking.Status.Cancelled() {
		flag := true
		record.IsCancelled = &flag
		booking.IsCancelled = &flag
	}
	record = validate.UpdateBookingFields(booking, record)

	// 5. Check if the booking status is valid or not
	err = record.IsValidStatus()
	if err != nil {
		record.Error(err, "APIError", "the booking status is invalid")
		return api.StatusBadRequest(err)
	}

	// 6. Validate if the booking status is a valid event source
	eventSource, err := record.EventSource()
	if err != nil {
		record.Error(err, "EventBridgeError", "incorrect event source of booking")
		return api.StatusBadRequest(err)
	}

	// 7. Check if it is a cancelled booking and validate if the
	// required fields are present.
	err = record.IsBookingCancelled()
	if err != nil {
		record.Error(err, "APIError", "the booking cancellation details are not set")
		return api.StatusBadRequest(err)
	}

	// ********************************************************************* //
	// ******************** Send events to the EventBus ******************** //
	// ********************************************************************* //
	detail, err := json.Marshal(&record)
	if err != nil {
		booking.Error(err, "JSONError", "failed to marshal booking object")
		return api.StatusInternalServerError(err)
	}

	// Send event to the configured EventBridge Event Bus and specific source
	err = awswrapper.EventBridgePutEvents(ctx, string(detail), eventSource, eventbus)
	if err != nil {
		record.Error(err, "EventBridgeError", "failed to send events to the EventBus", utility.KVP{Key: "source", Value: eventSource})
		return api.StatusInternalServerError(err)
	}

	return api.StatusOKWithoutBody()
}
