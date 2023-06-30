package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rmarasigan/bus-ticketing/api/schema"
	"github.com/rmarasigan/bus-ticketing/internal/app/query"
	"github.com/rmarasigan/bus-ticketing/internal/utility"
)

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, event events.SQSEvent) error {
	var (
		booking schema.Bookings
		records = event.Records
	)

	if len(records) == 0 {
		utility.Info("SQSEvent", "no records found")
		return nil
	}

	for _, record := range records {
		// Unmarshal the event message
		err := utility.ParseJSON([]byte(record.Body), &booking)
		if err != nil {
			booking.Error(err, "JSONError", "failed to unmarshal the JSON-encoded data", utility.KVP{Key: "payload", Value: record.Body})
			return err
		}

		// Set default values of the booking record
		booking.SetValues()

		// Inserts a new booking record to the DynamoDB
		err = query.CreateBooking(ctx, booking)
		if err != nil {
			booking.Error(err, "DynamoDBError", "failed to create a new booking record")
			return err
		}
	}

	return nil
}
