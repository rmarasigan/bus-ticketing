package main

import (
	"context"
	"errors"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rmarasigan/bus-ticketing/api"
	"github.com/rmarasigan/bus-ticketing/api/schema"
	awswrapper "github.com/rmarasigan/bus-ticketing/internal/aws_wrapper"
	"github.com/rmarasigan/bus-ticketing/internal/utility"
)

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var (
		booking schema.Bookings
		queue   = os.Getenv("BOOKING_QUEUE")
	)

	// Check if the queue is configured
	if queue == "" {
		err := errors.New("sqs BOOKING_QUEUE environment variable is not set")
		booking.Error(err, "SQSError", "sqs BOOKING_QUEUE is not configured on the environment")

		return api.StatusInternalServerError(err)
	}

	// Unmarshal the received JSON-encoded data and check
	// if it is a valid JSON data that we have received.
	err := utility.ParseJSON([]byte(request.Body), &booking)
	if err != nil {
		booking.Error(err, "JSONError", "failed to unmarshal the JSON-encoded data", utility.KVP{Key: "payload", Value: request.Body})
		return api.StatusInternalServerError(err)
	}

	// Validate if the booking status is a valid one.
	err = booking.IsValidStatus()
	if err != nil {
		booking.Error(err, "APIError", "the booking status is invalid")
		return api.StatusBadRequest(err)
	}

	// Send the message to the queue
	err = awswrapper.SQSSendMessage(ctx, queue, request.Body, awswrapper.BOOKING_MSG_GROUP_ID)
	if err != nil {
		booking.Error(err, "SQSError", "failed to send message", utility.KVP{Key: "queue", Value: queue})
		return api.StatusInternalServerError(err)
	}

	return api.StatusOKWithoutBody()
}
