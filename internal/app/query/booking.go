package query

import (
	"context"
	"errors"

	"github.com/rmarasigan/bus-ticketing/internal/app/env"
	"github.com/rmarasigan/bus-ticketing/internal/trail"
)

// CreateBooking checks if the DynamoDB Table is configured on the environment, and
// creates a new booking record.
func CreateBooking(ctx context.Context, data interface{}) error {
	var tablename = env.BOOKING_TABLE

	// Check if the DynamoDB Table is configured
	if tablename == "" {
		trail.Error("dynamodb BOOKING_TABLE is not configured on the environment")
		err := errors.New("dynamodb BOOKING_TABLE environment is not set")

		return err
	}

	// Save the Booking record into the DynamoDB Table
	err := InsertItem(ctx, tablename, data)
	if err != nil {
		trail.Error("failed to insert a new booking record")
		return err
	}

	return nil
}
