package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"github.com/rmarasigan/bus-ticketing/api/schema"
	"github.com/rmarasigan/bus-ticketing/internal/app/config"
	"github.com/rmarasigan/bus-ticketing/internal/app/email/template"
	"github.com/rmarasigan/bus-ticketing/internal/app/query"
	"github.com/rmarasigan/bus-ticketing/internal/app/validate"
	"github.com/rmarasigan/bus-ticketing/internal/utility"
)

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, event events.CloudWatchEvent) error {
	var (
		detail  = event.Detail
		booking schema.Bookings
	)

	// Unmarshal the received JSON-encoded event data
	err := utility.ParseJSON([]byte(detail), &booking)
	if err != nil {
		booking.Error(err, "JSONError", "failed to unmarshal the JSON-encoded event data", utility.KVP{Key: "event", Value: event})
		return err
	}

	// ********************************************************************* //
	// ********************* Update the booking record ********************* //
	// ********************************************************************* //
	// Create a composite key that has both the partition/primary key
	// and the sort key of the item.
	var bookingCompositeKey = map[string]types.AttributeValue{
		"id":           &types.AttributeValueMemberS{Value: booking.ID},
		"bus_route_id": &types.AttributeValueMemberS{Value: booking.BusRouteID},
	}

	// Construct the update builder for booking
	var updateBooking = expression.Set(expression.Name("status"), expression.Value(booking.Status)).
		Set(expression.Name("is_cancelled"), expression.Value(booking.IsCancelled)).
		Set(expression.Name("date_confirmed"), expression.Value(""))

	bookingResult, err := query.UpdateBooking(ctx, bookingCompositeKey, updateBooking)
	if err != nil {
		booking.Error(err, "DynamoDBError", "failed to update the booking record")
		return err
	}

	// ********************************************************************* //
	// ************** Update/add the cancelled booking record ************** //
	// ********************************************************************* //
	booking.Cancelled.ID = uuid.NewString()
	booking.Cancelled.BookingID = booking.ID
	booking.Cancelled.DateCancelled = time.Now().Format("2006-01-02 15:04:05")

	// Check if the record exist
	cancelledBookingExists, err := validate.IsCancelledBookingExists(ctx, booking.ID)
	if err != nil {
		booking.Error(err, "IsCancelledBookingExists", "failed to validate cancelled booking if it exist")
		return err
	}

	if !cancelledBookingExists {
		// Create a partition/primary key of the item.
		var cancelledBookingKey = map[string]types.AttributeValue{
			"booking_id": &types.AttributeValueMemberS{Value: booking.ID},
		}

		// Construct the update builder for cancelled booking.
		var updateCancelledBooking = expression.Set(expression.Name("id"), expression.Value(booking.Cancelled.ID)).
			Set(expression.Name("reason"), expression.Value(booking.Cancelled.Reason)).
			Set(expression.Name("cancelled_by"), expression.Value(booking.Cancelled.CancelledBy)).
			Set(expression.Name("date_cancelled"), expression.Value(booking.Cancelled.DateCancelled))

		_, err := query.RecordBookingCancelled(ctx, cancelledBookingKey, updateCancelledBooking)
		if err != nil {
			booking.Error(err, "DynamoDBError", "failed to record the cancelled booking")
			return err
		}
	}

	// ********************************************************************* //
	// ******************** Sending email to the client ******************** //
	// ********************************************************************* //
	// Fetch email configuration
	email, err := config.GetEmailConfig(ctx)
	if err != nil {
		booking.Error(err, "EmailError", "failed to fetch and set the email configuration")
		return err
	}

	// Fetch the user account record
	user, err := query.GetUserAccountById(ctx, booking.UserID)
	if err != nil {
		booking.Error(err, "DynamoDBError", "failed to fetch the user account")
		return err
	}

	// Fetch the bus route record
	routes, err := query.GetBusRouteRecords(ctx, booking.BusRouteID, booking.BusID)
	if err != nil {
		booking.Error(err, "DynamoDBError", "failed to fetch the bus route record")
		return err
	}
	route := routes[0]

	// Set the email content
	email.Content.To = append(email.Content.To, user.Email)
	email.Content.Subject = fmt.Sprintf("CANCELLED BOOKING: %s to %s [%s]", route.FromRoute, route.ToRoute, booking.TravelDate)

	if strings.HasPrefix(booking.Cancelled.CancelledBy, "ADMN") {
		email.Content.Message = template.CancelledBooking(user, route, booking, email.CustomerSupport)
	} else {
		email.Content.Message = template.CustomerCancelledBooking(user, route, booking, email.CustomerSupport)
	}

	// Send email to the client
	err = email.Send()
	if err != nil {
		booking.Error(err, "EmailError", "failed to send email to client")
		return err
	}
	utility.Info("CancelledBooking", "Successfully cancelled the booking record", utility.KVP{Key: "booking", Value: bookingResult})

	return nil
}
