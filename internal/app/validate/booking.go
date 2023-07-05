package validate

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/rmarasigan/bus-ticketing/api/schema"
	"github.com/rmarasigan/bus-ticketing/internal/app/env"
	"github.com/rmarasigan/bus-ticketing/internal/app/query"
	"github.com/rmarasigan/bus-ticketing/internal/trail"
)

// UpdateBookingFields validates if the field that are going to be updated
// are empty or not to set its previous value.
//
// Fields that are validated:
//  seat_number, status, reason
//
// Fields that are validated if it is a cancelled booking:
//  reason, cancelled_by
func UpdateBookingFields(booking, old schema.Bookings) schema.Bookings {
	if booking.Status != "" {
		old.Status = booking.Status
	}

	if booking.SeatNumber != "" {
		old.SeatNumber = booking.SeatNumber
	}

	if booking.IsCancelled != nil {
		if booking.Cancelled != (schema.BookingCancelled{}) {
			if booking.Cancelled.Reason != "" {
				old.Cancelled.Reason = booking.Cancelled.Reason
			}

			if booking.Cancelled.CancelledBy != "" {
				old.Cancelled.CancelledBy = booking.Cancelled.CancelledBy
			}
		}
	}

	return old
}

// IsCancelledBookingExists checks if the DynamoDB Table is configured on the environment, and
// returns a boolean and error value to check whether the cancelled booking already exist or not.
func IsCancelledBookingExists(ctx context.Context, bookingId string) (bool, error) {
	var tablename = env.BOOKING_CANCELLED_TABLE

	// Check if the DynamoDB Table is configured
	if tablename == "" {
		trail.Error("dynamodb BOOKING_CANCELLED_TABLE is not configured on the environment")
		err := errors.New("dynamodb BOOKING_CANCELLED_TABLE environment variable is not set")

		return false, err
	}

	// Create a composite key expression
	key := expression.Key("booking_id").Equal(expression.Value(bookingId))

	result, err := query.IsExisting(ctx, tablename, key)
	if err != nil {
		return false, err
	}

	return result, nil
}
