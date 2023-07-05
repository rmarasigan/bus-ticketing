package main

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/rmarasigan/bus-ticketing/api/schema"
	"github.com/rmarasigan/bus-ticketing/internal/app/config"
	"github.com/rmarasigan/bus-ticketing/internal/app/email/template"
	"github.com/rmarasigan/bus-ticketing/internal/app/query"
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
	var compositeKey = map[string]types.AttributeValue{
		"id":           &types.AttributeValueMemberS{Value: booking.ID},
		"bus_route_id": &types.AttributeValueMemberS{Value: booking.BusRouteID},
	}

	// Construct the update builder
	booking.DateConfirmed = time.Now().Format("2006-01-02 15:04:05")
	var update = expression.Set(expression.Name("status"), expression.Value(booking.Status)).
		Set(expression.Name("date_confirmed"), expression.Value(booking.DateConfirmed)).
		Set(expression.Name("seat_number"), expression.Value(booking.SeatNumber))

	result, err := query.UpdateBooking(ctx, compositeKey, update)
	if err != nil {
		booking.Error(err, "DynamoDBError", "failed to update thte booking record")
		return err
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
	route, err := query.GetBusRoute(ctx, booking.BusRouteID, booking.BusID)
	if err != nil {
		booking.Error(err, "DynamoDBError", "failed to fetch the bus route record")
		return err
	}

	// Set the email content
	email.Content.To = append(email.Content.To, user.Email)
	email.Content.Message = template.ConfirmedBooking(user, route, booking, email.CustomerSupport)
	email.Content.Subject = fmt.Sprintf("BOOKING SCHEDULE: %s to %s [%s]", route.FromRoute, route.ToRoute, booking.TravelDate)

	// Send email to the client
	err = email.Send()
	if err != nil {
		booking.Error(err, "EmailError", "failed to send email to client")
		return err
	}

	utility.Info("ConfirmedBooking", "Successfully confirmed the booking", utility.KVP{Key: "booking", Value: result})

	return nil
}
